package models

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Category string

const (
	// Expense categories
	CategoryFood          Category = "Food & Dining"
	CategoryTransport     Category = "Transport"
	CategoryHousing       Category = "Housing"
	CategoryEntertainment Category = "Entertainment"
	CategoryHealth        Category = "Health"
	CategoryShopping      Category = "Shopping"
	CategoryUtilities     Category = "Utilities"
	CategoryOther         Category = "Other"

	// Income categories
	CategorySalary     Category = "Salary"
	CategoryFreelance  Category = "Freelance"
	CategoryInvestment Category = "Investment"
	CategoryGift       Category = "Gift"
	CategoryBonus      Category = "Bonus"
)

type Transaction struct {
	ID          string          `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Category    Category        `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
}

type Budget struct {
	ID        string    `json:"id"`
	Category  Category  `json:"category"`
	Limit     float64   `json:"limit"`
	Period    string    `json:"period"` // "monthly"
	CreatedAt time.Time `json:"created_at"`
}

type Summary struct {
	TotalIncome   float64            `json:"total_income"`
	TotalExpenses float64            `json:"total_expenses"`
	Balance       float64            `json:"balance"`
	ByCategoryExp map[string]float64 `json:"by_category_expenses"`
	ByCategoryInc map[string]float64 `json:"by_category_income"`
	BudgetStatus  []BudgetStatus     `json:"budget_status"`
}

type BudgetStatus struct {
	Budget  Budget  `json:"budget"`
	Spent   float64 `json:"spent"`
	Percent float64 `json:"percent"`
}

var ExpenseCategories = []Category{
	CategoryFood,
	CategoryTransport,
	CategoryHousing,
	CategoryEntertainment,
	CategoryHealth,
	CategoryShopping,
	CategoryUtilities,
	CategoryOther,
}

var IncomeCategories = []Category{
	CategorySalary,
	CategoryFreelance,
	CategoryInvestment,
	CategoryGift,
	CategoryBonus,
}
