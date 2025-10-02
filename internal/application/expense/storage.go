package expense

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/files"
)

type expenseFileStorage struct {
}

var defaultExpenseStorage domain.ExpenseStorage = &expenseFileStorage{}

func (t *expenseFileStorage) Save(tasks []domain.Expense) error {
	return files.SaveToFile(tasks)
}

func (t *expenseFileStorage) Load() ([]domain.Expense, error) {
	return files.GetFromFile[[]domain.Expense]()
}
