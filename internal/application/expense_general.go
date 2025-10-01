package application

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/samber/lo"
	"time"
)

func addExpense(storage domain.ExpenseStorage, now func() time.Time, description string, amount float64) (domain.Expense, error) {
	expenses, err := storage.Load()

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error loading expenses: %w", err)
	}

	newExpense := domain.Expense{
		Id:          getNextExpenseId(expenses),
		Description: description,
		Amount:      amount,
		CreatedAt:   now(),
	}

	expenses = append(expenses, newExpense)
	err = storage.Save(expenses)

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error saving expenses: %w", err)
	}

	return newExpense, nil
}

func getNextExpenseId(expenses []domain.Expense) int {
	if len(expenses) == 0 {
		return 1
	}

	maxId := lo.MaxBy(expenses, func(a, b domain.Expense) bool {
		return a.Id > b.Id
	}).Id

	return maxId + 1
}
