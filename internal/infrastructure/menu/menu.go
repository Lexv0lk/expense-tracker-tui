package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	tableState state = iota
	addState
)

type MainModel struct {
	currentState state
	table        tableModel
	addInput     addModel
}

func InitialModel() (tea.Model, error) {
	table, err := getNewTableModel()

	if err != nil {
		return nil, err
	}

	addInput, err := getAdditionModel()

	if err != nil {
		return nil, err
	}

	return MainModel{
		currentState: tableState,
		table:        table.(tableModel),
		addInput:     addInput.(addModel),
	}, nil
}

func (m MainModel) Init() tea.Cmd { return nil }

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case backMsg:
		m.currentState = tableState
	case addMsg:
		newAddInput, _ := getAdditionModel()
		cAddInput, ok := newAddInput.(addModel)

		if !ok {
			panic("Failed assertion to addModel")
		}

		m.addInput = cAddInput
		m.currentState = addState
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
		addModel, ok := newAdd.(addModel)

		if !ok {
			panic("Failed assertion to addModel")
		}

		m.addInput = addModel
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
	default:
		panic("Unknown menu state")
	}
}
