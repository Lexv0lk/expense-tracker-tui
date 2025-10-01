//go:generate mockgen -destination=mocks/files.go -package=mocks io WriteCloser,ReadCloser
package files

import (
	"encoding/json"
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/files/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type TestCase struct {
		name          string
		writeCloserFn func(t *testing.T, expenses []domain.Expense) io.WriteCloser
		expenses      []domain.Expense
		expectedErr   error
	}

	tests := []TestCase{
		{
			name: "Successful Save",
			writeCloserFn: func(t *testing.T, expenses []domain.Expense) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(expenses, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			expenses: []domain.Expense{
				{Id: 1, Description: "Expense 1", Amount: 10},
				{Id: 2, Description: "Expense 2", Amount: 20},
				{Id: 3, Description: "Expense 3", Amount: 30},
			},
			expectedErr: nil,
		},
		{
			name: "Error writer",
			writeCloserFn: func(t *testing.T, expenses []domain.Expense) io.WriteCloser {
				t.Helper()

				testErr := fmt.Errorf("test error")

				result := mocks.NewMockWriteCloser(ctrl)
				result.EXPECT().Write(gomock.Any()).Return(0, testErr).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			expenses:    []domain.Expense{},
			expectedErr: fmt.Errorf("test error"),
		},
		{
			name: "Empty Expense List",
			writeCloserFn: func(t *testing.T, expenses []domain.Expense) io.WriteCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(expenses, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks.NewMockWriteCloser(ctrl)
				firstCall := result.EXPECT().Write(gomock.Eq(correctJson)).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			expenses:    []domain.Expense{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := saveToFile(tt.writeCloserFn(t, tt.expenses), tt.expenses)

			if tt.expectedErr != nil {
				assert.EqualError(err, tt.expectedErr.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGetFromFile(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type TestCase struct {
		name             string
		readCloserFn     func(t *testing.T, expenses []domain.Expense) io.ReadCloser
		expectedExpenses []domain.Expense
		expectedErr      error
	}

	tests := []TestCase{
		{
			name: "Successful Read",
			readCloserFn: func(t *testing.T, expenses []domain.Expense) io.ReadCloser {
				t.Helper()

				correctJson, _ := json.MarshalIndent(expenses, "", "  ")
				correctJson = append(correctJson, '\n')

				result := mocks.NewMockReadCloser(ctrl)
				firstCall := result.EXPECT().Read(gomock.Any()).DoAndReturn(
					func(p []byte) (n int, err error) {
						copy(p, correctJson)
						return len(correctJson), nil
					}).Times(1)
				result.EXPECT().Close().Times(1).After(firstCall)

				return result
			},
			expectedExpenses: []domain.Expense{
				{Id: 1, Description: "Expense 1", Amount: 10},
				{Id: 2, Description: "Expense 2", Amount: 20},
				{Id: 3, Description: "Expense 3", Amount: 30},
			},
			expectedErr: nil,
		},
		{
			name: "Error reader",
			readCloserFn: func(t *testing.T, expenses []domain.Expense) io.ReadCloser {
				t.Helper()

				testErr := fmt.Errorf("test error")

				result := mocks.NewMockReadCloser(ctrl)
				result.EXPECT().Read(gomock.Any()).Return(0, testErr).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			expectedExpenses: nil,
			expectedErr:      fmt.Errorf("test error"),
		},
		{
			name: "No err if empty file",
			readCloserFn: func(t *testing.T, expenses []domain.Expense) io.ReadCloser {
				t.Helper()

				result := mocks.NewMockReadCloser(ctrl)
				result.EXPECT().Read(gomock.Any()).Return(0, io.EOF).Times(1)
				result.EXPECT().Close().Times(1)

				return result
			},
			expectedExpenses: []domain.Expense{},
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			expenses, err := getFromFile[[]domain.Expense](tt.readCloserFn(t, tt.expectedExpenses))

			if tt.expectedErr != nil {
				assert.EqualError(err, tt.expectedErr.Error())
			} else {
				assert.EqualValues(tt.expectedExpenses, expenses)
				assert.NoError(err)
			}
		})
	}
}
