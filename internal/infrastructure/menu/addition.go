package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
	"time"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type addModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func (m addModel) Init() tea.Cmd { return nil }

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

					//TODO: add view for showing errors
					amount, _ := strconv.ParseFloat(m.inputs[1].Value(), 64)
					date, _ := time.Parse("2006-01-02", m.inputs[2].Value())
					expense.AddExpense(m.inputs[0].Value(), amount, date)
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

func (m addModel) View() string {
	var b strings.Builder

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

	return b.String()
}

func getAdditionModel() (tea.Model, error) {
	m := addModel{
		inputs: make([]textinput.Model, 3),
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
			t.Placeholder = "Amount"
			t.SetValue("0")
			t.Validate = validateAmount
		case 2:
			t.Placeholder = "YYYY-MM-DD"
			t.Validate = validateDate
			t.SetValue(time.Now().Format("2006-01-02"))
		}

		m.inputs[i] = t
	}

	return m, nil
}

func (m addModel) updateInputs(msg tea.Msg) tea.Cmd {
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
