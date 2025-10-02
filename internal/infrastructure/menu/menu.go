package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"io"
	"strconv"
	"strings"
	"time"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	tableStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

type mode int

const (
	modeMenu mode = iota
	modeAdd
	modeList
)

type addStep int

const (
	stepWhat addStep = iota
	stepAmount
	stepDate
)

type menuItem struct {
	name string
	mode mode
}

func (i menuItem) Title() string       { return i.name }
func (i menuItem) Description() string { return "" }
func (i menuItem) FilterValue() string { return i.name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	mode     mode
	expenses []domain.Expense

	menu   list.Model
	choice menuItem

	// adding steps
	step      addStep
	inputWhat textinput.Model
	inputAmt  textinput.Model
	inputDate textinput.Model

	//table
	table table.Model
}

func InitialModel() Model {
	items := []list.Item{
		menuItem{"Add Expense", modeAdd},
		menuItem{"List Expenses", modeList},
	}

	const defaultWidth = 20
	menu := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	menu.Title = "Menu"
	menu.SetShowStatusBar(false)
	menu.SetFilteringEnabled(false)
	menu.Styles.Title = titleStyle
	menu.Styles.PaginationStyle = paginationStyle
	menu.Styles.HelpStyle = helpStyle

	tiWhat := textinput.New()
	tiWhat.Placeholder = "Some product"
	tiWhat.CharLimit = 156
	tiWhat.Width = 20

	tiAmt := textinput.New()
	tiAmt.Placeholder = "100"
	tiAmt.CharLimit = 156
	tiAmt.Width = 20

	tiDate := textinput.New()
	tiDate.Placeholder = "YYYY-MM-DD"
	tiDate.SetValue(time.Now().Format("2006-01-02"))
	tiDate.CharLimit = 50
	tiDate.Width = 20

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Description", Width: 30},
		{Title: "Amount", Width: 10},
		{Title: "Date", Width: 12},
	}
	rows := getExpensesRows()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

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

	return Model{
		mode:      modeMenu,
		expenses:  []domain.Expense{},
		menu:      menu,
		inputWhat: tiWhat,
		inputAmt:  tiAmt,
		inputDate: tiDate,
		table:     t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeMenu:
		return m.updateMenu(msg)

	case modeAdd:
		switch m.step {
		case stepWhat:
			var cmd tea.Cmd
			m.inputWhat, cmd = m.inputWhat.Update(msg)
			if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
				m.step = stepAmount
				m.inputAmt.Focus()
				m.inputWhat.Blur()
			}
			return m, cmd

		case stepAmount:
			var cmd tea.Cmd
			m.inputAmt, cmd = m.inputAmt.Update(msg)
			if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
				m.step = stepDate
				m.inputDate.Focus()
				m.inputAmt.Blur()
			}
			return m, cmd

		case stepDate:
			var cmd tea.Cmd
			m.inputDate, cmd = m.inputDate.Update(msg)
			if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
				// TODO: error handling on EACH step
				amt, _ := strconv.ParseFloat(m.inputAmt.Value(), 64)
				dt, err := time.Parse("2006-01-02", m.inputDate.Value())
				if err != nil {
					dt = time.Now()
				}

				expense.AddExpense(m.inputWhat.Value(), amt, dt)
				m.table.SetRows(getExpensesRows())

				// очистить поля и вернуться в меню
				m.inputWhat.SetValue("")
				m.inputAmt.SetValue("")
				m.inputDate.SetValue(time.Now().Format("2006-01-02"))
				m.inputWhat.Blur()
				m.inputAmt.Blur()
				m.inputDate.Blur()

				m.mode = modeMenu
			}
			return m, cmd
		}

	case modeList:
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.mode = modeMenu
				return m, nil
			case "q", "ctrl+c":
				m.mode = modeMenu
				return m, nil //TODO: change to tea.Quit
			}
		}
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func getExpensesRows() []table.Row {
	allExpenses, _ := expense.GetAllExpenses()
	return lo.Map(allExpenses, getRow)
}

func getRow(expense domain.Expense, _ int) table.Row {
	return table.Row{
		strconv.Itoa(expense.Id),
		expense.Description,
		fmt.Sprintf("%.2f", expense.Amount),
		expense.SpentAt.Format("2006-01-02"),
	}
}

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.menu.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.Type; keypress {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			cItem, ok := m.menu.SelectedItem().(menuItem)
			if ok {
				m.mode = cItem.mode

				if m.mode == modeAdd {
					m.inputWhat.Focus()
				} else if m.mode == modeList {
					m.table.Focus()
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.mode {
	case modeMenu:
		return "\n" + m.menu.View()

	case modeAdd:
		switch m.step {
		case stepWhat:
			return lipgloss.NewStyle().Bold(true).Render("Добавление траты\n\nЧто купил:") + "\n" + m.inputWhat.View() + "\n"
		case stepAmount:
			return lipgloss.NewStyle().Bold(true).Render("Добавление траты\n\nСумма:") + "\n" + m.inputAmt.View() + "\n"
		case stepDate:
			return lipgloss.NewStyle().Bold(true).Render("Добавление траты\n\nДата (YYYY-MM-DD):") + "\n" + m.inputDate.View() + "\n"
		}

	case modeList:
		return tableStyle.Render(m.table.View() + "\n")
	}

	return ""
}
