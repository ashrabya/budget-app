package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ashrabya/budget-app/internal/models"
	"github.com/ashrabya/budget-app/internal/storage"
	"github.com/ashrabya/budget-app/internal/uuid"
)

type Handler struct {
	store *storage.Store
}

func New(store *storage.Store) *Handler {
	return &Handler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Transactions
	mux.HandleFunc("/api/transactions", h.handleTransactions)
	mux.HandleFunc("/api/transactions/", h.handleTransaction)

	// Budgets
	mux.HandleFunc("/api/budgets", h.handleBudgets)
	mux.HandleFunc("/api/budgets/", h.handleBudget)

	// Summary
	mux.HandleFunc("/api/summary", h.handleSummary)

	// Categories
	mux.HandleFunc("/api/categories", h.handleCategories)
}

// --- Transactions ---

func (h *Handler) handleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.store.ListTransactions())
	case http.MethodPost:
		var t models.Transaction
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		t.ID = uuid.New()
		t.CreatedAt = time.Now().UTC()
		if t.Date.IsZero() {
			t.Date = time.Now().UTC()
		}
		if t.Amount <= 0 {
			writeError(w, http.StatusBadRequest, "amount must be positive")
			return
		}
		if t.Type != models.Income && t.Type != models.Expense {
			writeError(w, http.StatusBadRequest, "type must be 'income' or 'expense'")
			return
		}
		if err := h.store.AddTransaction(&t); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, t)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		t, err := h.store.GetTransaction(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, t)
	case http.MethodPut:
		var t models.Transaction
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		t.ID = id
		if err := h.store.UpdateTransaction(&t); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, t)
	case http.MethodDelete:
		if err := h.store.DeleteTransaction(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// --- Budgets ---

func (h *Handler) handleBudgets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.store.ListBudgets())
	case http.MethodPost:
		var b models.Budget
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		b.ID = uuid.New()
		b.CreatedAt = time.Now().UTC()
		b.Period = "monthly"
		if b.Limit <= 0 {
			writeError(w, http.StatusBadRequest, "limit must be positive")
			return
		}
		if err := h.store.AddBudget(&b); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, b)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleBudget(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/budgets/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	if r.Method == http.MethodDelete {
		if err := h.store.DeleteBudget(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
		return
	}
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

// --- Summary ---

func (h *Handler) handleSummary(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	var month time.Time
	if monthStr == "" {
		month = time.Now().UTC()
	} else {
		var err error
		month, err = time.Parse("2006-01", monthStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid month format, use YYYY-MM: %v", err))
			return
		}
	}
	writeJSON(w, http.StatusOK, h.store.GetSummary(month))
}

// --- Categories ---

func (h *Handler) handleCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"expense": models.ExpenseCategories,
		"income":  models.IncomeCategories,
	})
}
