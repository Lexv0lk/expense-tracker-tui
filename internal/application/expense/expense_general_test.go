package expense

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense/mocks"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestAddExpense(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name            string
		storageFn       func(t *testing.T, expense domain.Expense) domain.ExpenseStorage
		expectedExpense domain.Expense
		expectedErr     error
	}

	testCases := []testCase{
		{
			name: "Successful addition",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
					{Id: 2, Description: "Lunch", Amount: 12.0},
				}

				resExpenses := make([]domain.Expense, len(currentExpenses))
				copy(resExpenses, currentExpenses)
				resExpenses = append(resExpenses, expense)

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          3,
				Description: "Dinner",
				Amount:      20.0,
				SpentAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			name: "Empty storage load",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{}
				resExpenses := []domain.Expense{expense}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          1,
				Description: "Groceries",
				Amount:      45.0,
			},
			expectedErr: nil,
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedExpense: domain.Expense{},
			expectedErr:     assert.AnError,
		},
		{
			name: "Save error",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Subscription", Amount: 9.99},
				}
				resExpenses := []domain.Expense{
					{Id: 1, Description: "Subscription", Amount: 9.99},
					expense,
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          2,
				Description: "Book",
				Amount:      15.0,
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t, tt.expectedExpense)
			result, err := addExpense(mockStorage, tt.expectedExpense.SpentAt, tt.expectedExpense.Description, tt.expectedExpense.Amount)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExpense, result)
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name            string
		storageFn       func(t *testing.T, expense domain.Expense) domain.ExpenseStorage
		expectedExpense domain.Expense
		expectedErr     error
	}

	testCases := []testCase{
		{
			name: "Successful update",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5, SpentAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Lunch", Amount: 12.0, SpentAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)},
				}

				resExpenses := []domain.Expense{
					currentExpenses[0],
					expense,
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          2,
				Description: "Brunch",
				Amount:      15.0,
				SpentAt:     time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			name: "Expense not found",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expectedExpense: domain.Expense{Id: 2},
			expectedErr:     &ExpenseNotFoundError{ID: 2},
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedExpense: domain.Expense{Id: 1},
			expectedErr:     assert.AnError,
		},
		{
			name: "Save error",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Subscription", Amount: 9.99},
					{Id: 2, Description: "Book", Amount: 15.0},
				}
				resExpenses := []domain.Expense{
					currentExpenses[0],
					expense,
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          2,
				Description: "E-Book",
				Amount:      10.0,
				SpentAt:     time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
			},
			expectedErr: assert.AnError,
		},
		{
			name: "Update to same values",
			storageFn: func(t *testing.T, expense domain.Expense) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Groceries", Amount: 50.0, SpentAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)},
				}
				resExpenses := []domain.Expense{
					expense,
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expectedExpense: domain.Expense{
				Id:          1,
				Description: "Groceries",
				Amount:      50.0,
				SpentAt:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t, tt.expectedExpense)
			result, err := updateExpense(mockStorage, tt.expectedExpense.Id, tt.expectedExpense.Description, tt.expectedExpense.Amount, tt.expectedExpense.SpentAt)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExpense, result)
			}
		})
	}
}

func TestDeleteExpense(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name        string
		storageFn   func(t *testing.T) domain.ExpenseStorage
		expenseId   int
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Successful deletion",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
					{Id: 2, Description: "Lunch", Amount: 12.0},
					{Id: 3, Description: "Dinner", Amount: 20.0},
					{Id: 4, Description: "Snacks", Amount: 5.0},
				}

				resExpenses := []domain.Expense{
					currentExpenses[0],
					currentExpenses[2],
					currentExpenses[3],
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expenseId:   2,
			expectedErr: nil,
		},
		{
			name: "Expense not found",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
					{Id: 2, Description: "Lunch", Amount: 12.0},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expenseId:   54,
			expectedErr: &ExpenseNotFoundError{ID: 54},
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expenseId:   1,
			expectedErr: assert.AnError,
		},
		{
			name: "Save error",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Subscription", Amount: 9.99},
					{Id: 2, Description: "Book", Amount: 15.0},
				}
				resExpenses := []domain.Expense{
					currentExpenses[0],
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(assert.AnError).Times(1).After(firstCall)

				return result
			},
			expenseId:   2,
			expectedErr: assert.AnError,
		},
		{
			name: "Delete the only expense",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Groceries", Amount: 50.0},
				}
				resExpenses := []domain.Expense{}

				result := mocks.NewMockExpenseStorage(ctrl)
				firstCall := result.EXPECT().Load().Return(currentExpenses, nil).Times(1)
				result.EXPECT().Save(gomock.Eq(resExpenses)).Return(nil).Times(1).After(firstCall)

				return result
			},
			expenseId:   1,
			expectedErr: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t)
			err := deleteExpense(mockStorage, tt.expenseId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetExpense(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name            string
		storageFn       func(t *testing.T) domain.ExpenseStorage
		expenseId       int
		expectedExpense domain.Expense
		expectedErr     error
	}

	testCases := []testCase{
		{
			name: "Successful retrieval",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
					{Id: 2, Description: "Lunch", Amount: 12.0},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expenseId: 2,
			expectedExpense: domain.Expense{
				Id:          2,
				Description: "Lunch",
				Amount:      12.0,
			},
			expectedErr: nil,
		},
		{
			name: "Expense not found",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expenseId:       54,
			expectedExpense: domain.Expense{},
			expectedErr:     &ExpenseNotFoundError{ID: 54},
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expenseId:       1,
			expectedExpense: domain.Expense{},
			expectedErr:     assert.AnError,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t)
			result, err := getExpense(mockStorage, tt.expenseId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExpense, result)
			}
		})
	}
}

func TestExpenseFileStorage_Load(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name           string
		storageFn      func(t *testing.T) domain.ExpenseStorage
		expectedAmount float64
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "Successful summary calculation",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5},
					{Id: 2, Description: "Lunch", Amount: 12.0},
					{Id: 3, Description: "Dinner", Amount: 20.0},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expectedAmount: 35.5,
			expectedErr:    nil,
		},
		{
			name: "Empty expenses list",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			expectedAmount: 0,
			expectedErr:    nil,
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			expectedAmount: 0,
			expectedErr:    assert.AnError,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t)
			result, err := getAllExpensesSummary(mockStorage)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAmount, result)
			}
		})
	}
}

func TestGetExpensesSummary(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type testCase struct {
		name           string
		storageFn      func(t *testing.T) domain.ExpenseStorage
		year           int
		month          time.Month
		expectedAmount float64
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "Successful monthly summary calculation",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5, SpentAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Lunch", Amount: 12.0, SpentAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
					{Id: 3, Description: "Dinner", Amount: 20.0, SpentAt: time.Date(2024, 2, 1, 19, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			year:           2024,
			month:          time.January,
			expectedAmount: 15.5,
		},
		{
			name: "No expenses for the month",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5, SpentAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Lunch", Amount: 12.0, SpentAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			year:           2024,
			month:          time.February,
			expectedAmount: 0,
			expectedErr:    nil,
		},
		{
			name: "Load error",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(nil, assert.AnError).Times(1)

				return result
			},
			year:           2024,
			month:          time.January,
			expectedAmount: 0,
			expectedErr:    assert.AnError,
		},
		{
			name: "Expenses in different years",
			storageFn: func(t *testing.T) domain.ExpenseStorage {
				t.Helper()

				currentExpenses := []domain.Expense{
					{Id: 1, Description: "Coffee", Amount: 3.5, SpentAt: time.Date(2023, 12, 31, 9, 0, 0, 0, time.UTC)},
					{Id: 2, Description: "Lunch", Amount: 12.0, SpentAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
					{Id: 3, Description: "Dinner", Amount: 20.0, SpentAt: time.Date(2024, 1, 20, 19, 0, 0, 0, time.UTC)},
					{Id: 4, Description: "Snacks", Amount: 5.0, SpentAt: time.Date(2024, 2, 1, 15, 0, 0, 0, time.UTC)},
				}

				result := mocks.NewMockExpenseStorage(ctrl)
				result.EXPECT().Load().Return(currentExpenses, nil).Times(1)

				return result
			},
			year:           2024,
			month:          time.January,
			expectedAmount: 32.0,
			expectedErr:    nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tt.storageFn(t)
			result, err := getExpensesSummary(mockStorage, tt.year, tt.month)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAmount, result)
			}
		})
	}
}
