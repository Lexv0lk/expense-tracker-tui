package expense

import "fmt"

type ExpenseNotFoundError struct {
	ID int
}

func (e *ExpenseNotFoundError) Error() string {
	return fmt.Sprintf("Expense with ID %d not found", e.ID)
}

func (e *ExpenseNotFoundError) Is(target error) bool {
	cErr, ok := target.(*ExpenseNotFoundError)

	if !ok {
		return false
	}

	return cErr.ID == e.ID
}
