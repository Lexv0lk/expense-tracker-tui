package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

var (
	tableStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

type tableModel struct {
	table         table.Model
	help          help.Model
	actionsKeyMap ActionKeyMap
	expensesSum   float64
}

func getNewTableModel() (tea.Model, error) {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Description", Width: 30},
		{Title: "Amount", Width: 10},
		{Title: "Date", Width: 12},
	}
	rows, err := getExpensesRows()

	if err != nil {
		return nil, fmt.Errorf("error getting expenses rows: %w", err)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	t.Focus()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("170")).
		Bold(false)
	t.SetStyles(s)

	return tableModel{
		table:         t,
		help:          help.New(),
		actionsKeyMap: getActionKeymap(),
	}, nil
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Delete):
			if len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()

				id, err := strconv.Atoi(selectedRow[0])
				if err != nil {
					return m, errorCmd(err, backToTableCmd())
				}

				err = expense.DeleteExpense(id)
				if err != nil {
					return m, errorCmd(err, backToTableCmd())
				}

				rows, err := getExpensesRows()
				if err != nil {
					return m, errorCmd(err, backToTableCmd())
				}

				allExpensesSum, err := expense.GetAllExpensesSummary()
				if err != nil {
					return m, errorCmd(fmt.Errorf("Error when calculating expenses sum: %w", err), backToTableCmd())
				}

				m.expensesSum = allExpensesSum
				m.table.SetRows(rows)
			}
		case key.Matches(msg, constants.Keymap.Create):
			return m, goToAddCmd()
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.actionsKeyMap.GetSum):
			return m, goToSummaryCmd()
		case key.Matches(msg, constants.Keymap.Enter):
			if len(m.table.Rows()) > 0 {
				selectedRow := m.table.SelectedRow()

				id, err := strconv.Atoi(selectedRow[0])
				if err != nil {
					return m, errorCmd(err, backToTableCmd())
				}

				return m, goToEditCmd(id)
			}
		}
	case backMsg:
		rows, err := getExpensesRows()
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		m.table.SetRows(rows)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	var sb strings.Builder
	sb.WriteString(tableStyle.Render(m.table.View() + "\n"))
	sb.WriteString("\n" + fmt.Sprintf("Total spent: %.2f", m.expensesSum) + "\n")
	sb.WriteString(m.help.View(m.actionsKeyMap))

	return sb.String()
}

func getExpensesRows() ([]table.Row, error) {
	allExpenses, err := expense.GetAllExpenses()

	if err != nil {
		return nil, fmt.Errorf("error getting all expenses: %w", err)
	}

	return lo.Map(allExpenses, getRow), nil
}

func getRow(expense domain.Expense, _ int) table.Row {
	return table.Row{
		strconv.Itoa(expense.Id),
		expense.Description,
		fmt.Sprintf("%.2f", expense.Amount),
		expense.SpentAt.Format("2006-01-02"),
	}
}
