package domain

import "time"

type Expense struct {
	Id          int
	CreatedAt   time.Time
	Description string
	Amount      float64
}

type ExpenseStorage interface {
	Save(expenses []Expense) error
	Load() ([]Expense, error)
}
