package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ashrabya/budget-app/internal/models"
)

type Store struct {
	mu           sync.RWMutex
	transactions map[string]*models.Transaction
	budgets      map[string]*models.Budget
	dataFile     string
}

type persistData struct {
	Transactions []*models.Transaction `json:"transactions"`
	Budgets      []*models.Budget      `json:"budgets"`
}

func NewStore(dataFile string) (*Store, error) {
	s := &Store{
		transactions: make(map[string]*models.Transaction),
		budgets:      make(map[string]*models.Budget),
		dataFile:     dataFile,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading data: %w", err)
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.dataFile)
	if err != nil {
		return err
	}
	var pd persistData
	if err := json.Unmarshal(data, &pd); err != nil {
		return err
	}
	for _, t := range pd.Transactions {
		s.transactions[t.ID] = t
	}
	for _, b := range pd.Budgets {
		s.budgets[b.ID] = b
	}
	return nil
}

func (s *Store) save() error {
	pd := persistData{}
	for _, t := range s.transactions {
		pd.Transactions = append(pd.Transactions, t)
	}
	for _, b := range s.budgets {
		pd.Budgets = append(pd.Budgets, b)
	}
	data, err := json.MarshalIndent(pd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dataFile, data, 0644)
}

// Transactions

func (s *Store) AddTransaction(t *models.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[t.ID] = t
	return s.save()
}

func (s *Store) GetTransaction(id string) (*models.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.transactions[id]
	if !ok {
		return nil, fmt.Errorf("transaction %s not found", id)
	}
	return t, nil
}

func (s *Store) ListTransactions() []*models.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Transaction, 0, len(s.transactions))
	for _, t := range s.transactions {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.After(list[j].Date)
	})
	return list
}

func (s *Store) DeleteTransaction(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transactions[id]; !ok {
		return fmt.Errorf("transaction %s not found", id)
	}
	delete(s.transactions, id)
	return s.save()
}

func (s *Store) UpdateTransaction(t *models.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transactions[t.ID]; !ok {
		return fmt.Errorf("transaction %s not found", t.ID)
	}
	s.transactions[t.ID] = t
	return s.save()
}

// Budgets

func (s *Store) AddBudget(b *models.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.budgets[b.ID] = b
	return s.save()
}

func (s *Store) ListBudgets() []*models.Budget {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Budget, 0, len(s.budgets))
	for _, b := range s.budgets {
		list = append(list, b)
	}
	return list
}

func (s *Store) DeleteBudget(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.budgets[id]; !ok {
		return fmt.Errorf("budget %s not found", id)
	}
	delete(s.budgets, id)
	return s.save()
}

// Summary

func (s *Store) GetSummary(month time.Time) *models.Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	summary := &models.Summary{
		ByCategoryExp: make(map[string]float64),
		ByCategoryInc: make(map[string]float64),
	}

	// Filter to current month
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	for _, t := range s.transactions {
		if t.Date.Before(start) || t.Date.After(end) {
			continue
		}
		if t.Type == models.Income {
			summary.TotalIncome += t.Amount
			summary.ByCategoryInc[string(t.Category)] += t.Amount
		} else {
			summary.TotalExpenses += t.Amount
			summary.ByCategoryExp[string(t.Category)] += t.Amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpenses

	// Budget status
	for _, b := range s.budgets {
		spent := summary.ByCategoryExp[string(b.Category)]
		pct := 0.0
		if b.Limit > 0 {
			pct = (spent / b.Limit) * 100
		}
		summary.BudgetStatus = append(summary.BudgetStatus, models.BudgetStatus{
			Budget:  *b,
			Spent:   spent,
			Percent: pct,
		})
	}

	return summary
}
