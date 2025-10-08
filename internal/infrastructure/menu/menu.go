package menu

import (
	"fmt"
	"github.com/Lexv0lk/expense-tracker-tui/internal/application/expense"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	tableState state = iota
	addState
	msgState
	sumState
)

type MainModel struct {
	currentState state
	table        tableModel
	addInput     changeFormModel
	summaryModel summaryInfoModel
	msgView      msgModel
}

func InitialModel() (tea.Model, error) {
	table, err := getNewTableModel()
	if err != nil {
		return nil, err
	}

	return MainModel{
		currentState: tableState,
		table:        table.(tableModel),
	}, nil
}

func (m MainModel) Init() tea.Cmd { return nil }

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case backMsg:
		allExpensesSum, err := expense.GetAllExpensesSummary()

		if err != nil {
			return m, errorCmd(fmt.Errorf("Error when calculating expenses sum: %w", err), backToTableCmd())
		}

		m.table.expensesSum = allExpensesSum
		m.currentState = tableState
	case addMsg:
		newAddInput, _ := getAddingModel()
		cAddInput, ok := newAddInput.(changeFormModel)

		if !ok {
			panic("Failed assertion to changeFormModel")
		}

		m.addInput = cAddInput
		m.currentState = addState
	case editMsg:
		exp, err := expense.GetExpense(msg.id)
		if err != nil {
			return m, errorCmd(err, backToTableCmd())
		}

		newEditInput, _ := getChangeModel(exp)
		cEditInput, ok := newEditInput.(changeFormModel)

		if !ok {
			panic("Failed assertion to changeFormModel")
		}

		m.addInput = cEditInput
		m.currentState = addState
	case summaryMsg:
		newSummaryModel := newSummaryInfoModel()
		m.summaryModel = newSummaryModel
		m.currentState = sumState
	case errorMsg:
		newMsgModel := getNewMsgModel(msg.error.Error(), msg.sourceBack)
		cMsgModel, ok := newMsgModel.(msgModel)

		if !ok {
			panic("Failed assertion to msgModel")
		}

		m.msgView = cMsgModel
		m.currentState = msgState
	}

	switch m.currentState {
	case tableState:
		newTable, newCmd := m.table.Update(msg)
		tableModel, ok := newTable.(tableModel)

		if !ok {
			panic("Failed assertion to tableModel")
		}

		m.table = tableModel
		cmd = newCmd
	case addState:
		newAdd, newCmd := m.addInput.Update(msg)
		addModel, ok := newAdd.(changeFormModel)

		if !ok {
			panic("Failed assertion to changeFormModel")
		}

		m.addInput = addModel
		cmd = newCmd
	case sumState:
		newSummary, newCmd := m.summaryModel.Update(msg)
		summaryModel, ok := newSummary.(summaryInfoModel)

		if !ok {
			panic("Failed assertion to summaryInfoModel")
		}

		m.summaryModel = summaryModel
		cmd = newCmd
	case msgState:
		newMsg, newCmd := m.msgView.Update(msg)
		msgModel, ok := newMsg.(msgModel)

		if !ok {
			panic("Failed assertion to msgModel")
		}

		m.msgView = msgModel
		cmd = newCmd
	}

	return m, cmd
}

func (m MainModel) View() string {
	switch m.currentState {
	case tableState:
		return m.table.View()
	case addState:
		return m.addInput.View()
	case msgState:
		return m.msgView.View()
	case sumState:
		return m.summaryModel.View()
	default:
		panic("Unknown menu state")
	}
}
