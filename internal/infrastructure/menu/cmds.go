package menu

import tea "github.com/charmbracelet/bubbletea"

type backMsg struct{}
type addMsg struct{}
type errorMsg struct {
	error
	sourceBack tea.Cmd
}
type editMsg struct {
	id int
}
type summaryMsg struct{}

func backToTableCmd() tea.Cmd {
	return func() tea.Msg {
		return backMsg{}
	}
}

func goToAddCmd() tea.Cmd {
	return func() tea.Msg {
		return addMsg{}
	}
}

func goToEditCmd(id int) tea.Cmd {
	return func() tea.Msg {
		return editMsg{id}
	}
}

func goToSummaryCmd() tea.Cmd {
	return func() tea.Msg {
		return summaryMsg{}
	}
}

func errorCmd(err error, sourceBack tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return errorMsg{err, sourceBack}
	}
}
