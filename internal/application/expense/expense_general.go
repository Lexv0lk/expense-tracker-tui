package expense

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/samber/lo"
	"time"
)

func addExpense(storage domain.ExpenseStorage, spentTime time.Time, description string, category string, amount float64) (domain.Expense, error) {
	expenses, err := storage.Load()

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error loading expenses: %w", err)
	}

	newExpense := domain.Expense{
		Id:          getNextExpenseId(expenses),
		Description: description,
		Category:    category,
		Amount:      amount,
		SpentAt:     spentTime,
	}

	expenses = append(expenses, newExpense)
	err = storage.Save(expenses)

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error saving expenses: %w", err)
	}

	return newExpense, nil
}

func updateExpense(storage domain.ExpenseStorage, id int, description string, category string, amount float64, spentAt time.Time) (domain.Expense, error) {
	expenses, err := storage.Load()

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error loading expenses: %w", err)
	}

	var updatedExpense domain.Expense
	found := false

	for i := range expenses {
		if expenses[i].Id == id {
			expenses[i].Description = description
			expenses[i].Amount = amount
			expenses[i].SpentAt = spentAt
			expenses[i].Category = category
			updatedExpense = expenses[i]
			found = true
			break
		}
	}

	if !found {
		return domain.Expense{}, &ExpenseNotFoundError{ID: id}
	}

	err = storage.Save(expenses)

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error saving expenses: %w", err)
	}

	return updatedExpense, nil
}

func deleteExpense(storage domain.ExpenseStorage, id int) error {
	expenses, err := storage.Load()

	if err != nil {
		return fmt.Errorf("Error loading expenses: %w", err)
	}

	_, index, found := lo.FindIndexOf(expenses, func(e domain.Expense) bool {
		return e.Id == id
	})

	if !found {
		return &ExpenseNotFoundError{ID: id}
	}

	expenses = append(expenses[:index], expenses[index+1:]...)
	err = storage.Save(expenses)

	if err != nil {
		return fmt.Errorf("Error saving expenses: %w", err)
	}

	return nil
}

func getExpense(storage domain.ExpenseStorage, id int) (domain.Expense, error) {
	expenses, err := storage.Load()

	if err != nil {
		return domain.Expense{}, fmt.Errorf("Error loading expenses: %w", err)
	}

	expense, _, found := lo.FindIndexOf(expenses, func(e domain.Expense) bool {
		return e.Id == id
	})

	if !found {
		return domain.Expense{}, &ExpenseNotFoundError{ID: id}
	}

	return expense, nil
}

func getAllExpensesSummary(storage domain.ExpenseStorage) (float64, error) {
	expenses, err := storage.Load()

	if err != nil {
		return 0, fmt.Errorf("Error loading expenses: %w", err)
	}

	return lo.SumBy(expenses, func(e domain.Expense) float64 {
		return e.Amount
	}), nil
}

func getExpensesSummary(storage domain.ExpenseStorage, year int, month time.Month) (float64, error) {
	expenses, err := storage.Load()

	if err != nil {
		return 0, fmt.Errorf("Error loading expenses: %w", err)
	}

	filteredExpenses := lo.Filter(expenses, func(e domain.Expense, _ int) bool {
		return e.SpentAt.Year() == year && e.SpentAt.Month() == month
	})

	return lo.SumBy(filteredExpenses, func(e domain.Expense) float64 {
		return e.Amount
	}), nil
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
