package expense

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"time"
)

func AddExpense(description string, amount float64, spentTime time.Time) (domain.Expense, error) {
	return addExpense(defaultExpenseStorage, spentTime, description, amount)
}

func DeleteExpense(id int) error {
	return deleteExpense(defaultExpenseStorage, id)
}

func GetAllExpenses() ([]domain.Expense, error) {
	return defaultExpenseStorage.Load()
}
