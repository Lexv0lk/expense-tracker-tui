package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
	"strings"
	"time"
)

const summaryInfoTitle = "Summary Info"

var months = map[string]time.Month{
	"january":   time.January,
	"february":  time.February,
	"march":     time.March,
	"april":     time.April,
	"may":       time.May,
	"june":      time.June,
	"july":      time.July,
	"august":    time.August,
	"september": time.September,
	"october":   time.October,
	"november":  time.November,
	"december":  time.December,
}

type summaryInfoModel struct {
	focusIndex int
	inputs     []textinput.Model

	//data
	currentSummary float64

	//help
	helpModel      help.Model
	navigationKeys NavigationKeyMap

	//err
	summaryErr error
}

func (m summaryInfoModel) Init() tea.Cmd { return nil }

func (m summaryInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.navigationKeys.Down, m.navigationKeys.Up, constants.Keymap.Enter):
			if key.Matches(msg, constants.Keymap.Enter) {
				if m.focusIndex == len(m.inputs) {
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

					year, err := strconv.Atoi(m.inputs[0].Value())
					if err != nil {
						return m, errorCmd(err, goToSummaryCmd())
					}

					monthStr := m.inputs[1].Value()
					monthStr = strings.ToLower(strings.TrimSpace(monthStr))
					month, ok := months[monthStr]
					if !ok {
						return m, errorCmd(fmt.Errorf("invalid month: %s", monthStr), goToSummaryCmd())
					}

					summary, err := expense.GetMonthlyExpensesSummary(year, month)
					if err != nil {
						m.summaryErr = fmt.Errorf("Eror fetching summary: %w", err)
					} else {
						m.summaryErr = nil
					}

					m.currentSummary = summary
					return m, nil
				} else {
					m.focusIndex++
				}
			} else {
				if key.Matches(msg, m.navigationKeys.Up) {
					m.focusIndex--
				} else if key.Matches(msg, m.navigationKeys.Down) {
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
		case key.Matches(msg, m.navigationKeys.Back):
			return m, backToTableCmd()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m summaryInfoModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(summaryInfoTitle) + "\n\n")

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

	if m.summaryErr != nil {
		b.WriteString(m.summaryErr.Error())
	} else {
		b.WriteString(fmt.Sprintf("Total spent for chosen month: %.2f\n\n", m.currentSummary))
	}

	b.WriteString(m.helpModel.View(m.navigationKeys))

	return b.String()
}

func newSummaryInfoModel() (tea.Model, error) {
	m := summaryInfoModel{
		inputs:         make([]textinput.Model, 2),
		helpModel:      help.New(),
		navigationKeys: getNavigationKeymap(),
	}

	var t textinput.Model

	// Year input
	t = textinput.New()
	t.Placeholder = "2023"
	t.Prompt = "Year: "
	t.Width = 20
	t.Validate = validateYear
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.SetValue(strconv.Itoa(time.Now().Year()))
	m.inputs[0] = t

	// Month input
	t = textinput.New()
	t.Placeholder = "January"
	t.Prompt = "Month: "
	t.Width = 20
	t.Validate = validateMonth
	t.PromptStyle = blurredStyle
	t.TextStyle = blurredStyle
	t.SetValue(time.Now().Month().String())
	m.inputs[1] = t

	m.inputs[0].Focus()

	return m, nil
}

func (m summaryInfoModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func validateYear(s string) error {
	if s == "" {
		return fmt.Errorf("year cannot be empty")
	}

	_, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("year must be a number")
	}

	return nil
}

func validateMonth(s string) error {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	_, ok := months[s]
	if !ok {
		return fmt.Errorf("invalid month: %s", s)
	}

	return nil
}
