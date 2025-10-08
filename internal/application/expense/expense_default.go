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

func UpdateExpense(id int, description string, amount float64, spentAt time.Time) (domain.Expense, error) {
	return updateExpense(defaultExpenseStorage, id, description, amount, spentAt)
}

func GetExpense(id int) (domain.Expense, error) {
	return getExpense(defaultExpenseStorage, id)
}

func GetAllExpenses() ([]domain.Expense, error) {
	return defaultExpenseStorage.Load()
}

func GetAllExpensesSummary() (float64, error) {
	return getAllExpensesSummary(defaultExpenseStorage)
}

func GetMonthlyExpensesSummary(year int, month time.Month) (float64, error) {
	return getExpensesSummary(defaultExpenseStorage, year, month)
}
