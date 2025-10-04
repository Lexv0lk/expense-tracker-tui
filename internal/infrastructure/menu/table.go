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
)

var (
	tableStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

type tableModel struct {
	table         table.Model
	help          help.Model
	actionsKeyMap ActionKeyMap
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
		actionsKeyMap: getDefaultKeymap(),
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
				//TODO: add error handling and confirmation
				selectedRow := m.table.SelectedRow()
				id, _ := strconv.Atoi(selectedRow[0])
				_ = expense.DeleteExpense(id)
				rows, _ := getExpensesRows()
				m.table.SetRows(rows)
			}
		case key.Matches(msg, constants.Keymap.Create):
			return m, goToAddCmd()
		case key.Matches(msg, constants.Keymap.Back):
			return m, tea.Quit
		}
	case backMsg:
		//TODO: add view for showing errors
		rows, _ := getExpensesRows()
		m.table.SetRows(rows)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return tableStyle.Render(m.table.View()+"\n") + "\n" + m.help.View(m.actionsKeyMap)
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
