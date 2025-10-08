package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
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

	allExpenses    []domain.Expense
	expensesToShow []domain.Expense

	//filter
	filterEnabled bool
	filterInput   textinput.Model
}

func getNewTableModel() (tea.Model, error) {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Category", Width: 15},
		{Title: "Description", Width: 30},
		{Title: "Amount", Width: 10},
		{Title: "Date", Width: 12},
	}

	allExpenses, err := expense.GetAllExpenses()
	if err != nil {
		return nil, fmt.Errorf("Error getting all expenses: %w", err)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(lo.Map(allExpenses, getRow)),
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

	f := textinput.New()
	f.Placeholder = "Type to filter by category"
	f.Prompt = "Filter: "
	f.Width = 20
	f.CharLimit = 0

	return tableModel{
		table:          t,
		help:           help.New(),
		actionsKeyMap:  getActionKeymap(),
		filterEnabled:  false,
		filterInput:    f,
		allExpenses:    allExpenses,
		expensesToShow: allExpenses,
	}, nil
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filterEnabled {
			switch {
			case key.Matches(msg, constants.Keymap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.actionsKeyMap.Filter):
				m.filterEnabled = false
				m.filterInput.Blur()
				m.table.Focus()
			}
		} else {
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

					err = m.UpdateExpenses()
					if err != nil {
						return m, errorCmd(err, backToTableCmd())
					}
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
			case key.Matches(msg, m.actionsKeyMap.Filter):
				m.filterEnabled = true
				m.table.Blur()
				m.filterInput.SetValue("")
				m.filterInput.Focus()
			}
		}
	case backMsg:
		err := m.UpdateExpenses()
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}
	}

	m.UpdateShowData()
	m.table, cmd = m.table.Update(msg)

	if m.filterEnabled {
		newFilterInput, filterCmd := m.filterInput.Update(msg)
		cmd = tea.Batch(cmd, filterCmd)
		m.filterInput = newFilterInput
	}

	return m, cmd
}

func (m tableModel) View() string {
	var sb strings.Builder

	if m.filterEnabled {
		sb.WriteString(m.filterInput.View() + "\n\n")
	}

	sb.WriteString(tableStyle.Render(m.table.View() + "\n"))
	sb.WriteString("\n" + fmt.Sprintf("Total spent: %.2f", m.expensesSum) + "\n")
	sb.WriteString(m.help.View(m.actionsKeyMap))

	return sb.String()
}

func (m *tableModel) UpdateExpenses() error {
	allExpenses, err := expense.GetAllExpenses()
	if err != nil {
		return fmt.Errorf("error getting all expenses: %w", err)
	}

	m.allExpenses = allExpenses
	return nil
}

func (m *tableModel) UpdateShowData() {
	if m.filterEnabled {
		filterVal := m.filterInput.Value()
		filtered := lo.Filter(m.allExpenses, func(expense domain.Expense, _ int) bool {
			return strings.Contains(strings.ToLower(expense.Category), strings.ToLower(filterVal))
		})
		m.expensesToShow = filtered
	} else {
		m.expensesToShow = m.allExpenses
	}

	m.expensesSum = lo.SumBy(m.expensesToShow, func(expense domain.Expense) float64 {
		return expense.Amount
	})
	m.table.SetRows(lo.Map(m.expensesToShow, getRow))
}

func getRow(expense domain.Expense, _ int) table.Row {
	return table.Row{
		strconv.Itoa(expense.Id),
		expense.Category,
		expense.Description,
		fmt.Sprintf("%.2f", expense.Amount),
		expense.SpentAt.Format("2006-01-02"),
	}
}
