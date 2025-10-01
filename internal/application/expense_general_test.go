package application

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNextExpenseId(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		expenses []domain.Expense
		expected int
	}

	testCases := []testCase{
		{
			name:     "No expenses",
			expenses: []domain.Expense{},
			expected: 1,
		},
		{
			name: "Some expenses",
			expenses: []domain.Expense{
				{Id: 1},
				{Id: 2},
				{Id: 3},
			},
			expected: 4,
		},
		{
			name: "Non-sequential IDs",
			expenses: []domain.Expense{
				{Id: 1},
				{Id: 3},
				{Id: 5},
			},
			expected: 6,
		},
		{
			name: "Single expense",
			expenses: []domain.Expense{
				{Id: 42},
			},
			expected: 43,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := getNextExpenseId(tt.expenses)
			assert.Equal(t, tt.expected, result)
		})
	}
}
