package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/domain"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
	"time"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Underline(true)

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

const changeFormTitle = "Expense Change Form"

type changeFormModel struct {
	focusIndex int
	inputs     []textinput.Model
	editingId  *int

	//help
	helpModel      help.Model
	navigationKeys NavigationKeyMap
}

func (m changeFormModel) Init() tea.Cmd { return nil }

func (m changeFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Enter, constants.Keymap.Up, constants.Keymap.Down):
			curInd := m.focusIndex

			if key.Matches(msg, constants.Keymap.Enter) {
				if curInd == len(m.inputs) {
					hasError := false

					for i := range m.inputs {
						err := m.inputs[i].Err

						if err != nil {
							m.inputs[i].SetValue(err.Error())
							hasError = true
						}
					}

					if hasError {
						return m, nil
					}

					description := m.inputs[0].Value()
					if description == "" {
						description = "-"
					}

					category := m.inputs[1].Value()
					if category == "" {
						category = "-"
					}

					amount, err := strconv.ParseFloat(m.inputs[2].Value(), 64)
					if err != nil {
						return m, errorCmd(err, goToAddCmd())
					}

					date, err := time.Parse("2006-01-02", m.inputs[3].Value())
					if err != nil {
						return m, errorCmd(err, goToAddCmd())
					}

					if m.editingId == nil {
						_, err = expense.AddExpense(description, category, amount, date)
						if err != nil {
							return m, errorCmd(err, goToAddCmd())
						}
					} else {
						_, err = expense.UpdateExpense(*m.editingId, description, category, amount, date)
						if err != nil {
							return m, errorCmd(err, goToAddCmd())
						}
					}

					return m, backToTableCmd()
				} else {
					m.focusIndex++
				}
			} else {
				if key.Matches(msg, constants.Keymap.Up) {
					m.focusIndex--
				} else if key.Matches(msg, constants.Keymap.Down) {
					m.focusIndex++
				}
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = blurredStyle
					m.inputs[i].TextStyle = blurredStyle
				}
			}

			return m, tea.Batch(cmds...)
		case key.Matches(msg, constants.Keymap.Back):
			return m, backToTableCmd()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m changeFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(changeFormTitle) + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(m.helpModel.View(m.navigationKeys))

	return b.String()
}

func getAddingModel() (tea.Model, error) {
	m := changeFormModel{
		inputs:         make([]textinput.Model, 4),
		helpModel:      help.New(),
		navigationKeys: getNavigationKeymap(),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Width = 100
		t.CharLimit = 0

		switch i {
		case 0:
			t.Placeholder = "Description"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Category"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 2:
			t.Placeholder = "Amount"
			t.SetValue("0")
			t.Validate = validateAmount
		case 3:
			t.Placeholder = "YYYY-MM-DD"
			t.Validate = validateDate
			t.SetValue(time.Now().Format("2006-01-02"))
		}

		m.inputs[i] = t
	}

	return m, nil
}

func getChangeModel(existingExpense domain.Expense) (tea.Model, error) {
	m, err := getAddingModel()

	if err != nil {
		return nil, err
	}

	cModel := m.(changeFormModel)
	cModel.inputs[0].SetValue(existingExpense.Description)
	cModel.inputs[1].SetValue(existingExpense.Category)
	cModel.inputs[2].SetValue(fmt.Sprintf("%.2f", existingExpense.Amount))
	cModel.inputs[3].SetValue(existingExpense.SpentAt.Format("2006-01-02"))

	id := existingExpense.Id
	cModel.editingId = &id

	return cModel, nil
}

func (m changeFormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func validateDate(dateStr string) error {
	if dateStr == "" {
		return fmt.Errorf("date cannot be empty")
	}

	_, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		return fmt.Errorf("date should be in YYYY-MM-DD format")
	}

	return nil
}

func validateAmount(amountStr string) error {
	if amountStr == "" {
		return fmt.Errorf("amount cannot be empty")
	}

	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)

	if err == nil && amount <= 0 {
		return fmt.Errorf("amount should be a positive number")
	}

	if err != nil {
		return fmt.Errorf("invalid amount format")
	}

	return nil
}
