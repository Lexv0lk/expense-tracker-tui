package menu

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	tableState state = iota
	changeState
	msgState
	sumState
)

type MainModel struct {
	currentState state
	models       map[state]tea.Model
}

func InitialModel() (tea.Model, error) {
	tableModel, err := newTableModel()
	if err != nil {
		return nil, err
	}

	return MainModel{
		currentState: tableState,
		models: map[state]tea.Model{
			tableState:  tableModel,
			changeState: changeFormModel{},
			sumState:    summaryInfoModel{},
			msgState:    msgModel{},
		},
	}, nil
}

func (m MainModel) Init() tea.Cmd { return nil }

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case backMsg:
		newTable, err := m.models[tableState].(tableModel).UpdateExpenses()
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		m.models[tableState] = newTable
		m.currentState = tableState
	case addMsg:
		newAddInput, err := newAdditionModel()
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		m.models[changeState] = newAddInput
		m.currentState = changeState
	case editMsg:
		exp, err := expense.GetExpense(msg.id)
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		newEditInput, err := newChangeModel(exp)
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		m.models[changeState] = newEditInput
		m.currentState = changeState
	case summaryMsg:
		newSummaryModel, err := newSummaryInfoModel()
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		m.models[sumState] = newSummaryModel
		m.currentState = sumState
	case errorMsg:
		m.models[msgState] = newMsgModel(msg.error.Error(), msg.sourceBack)
		m.currentState = msgState
	case infoMsg:
		m.models[msgState] = newMsgModel(msg.message, msg.sourceBack)
		m.currentState = msgState
	}

	m.models[m.currentState], cmd = m.models[m.currentState].Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	return m.models[m.currentState].View()
}
