package menu

import tea "github.com/charmbracelet/bubbletea"

type backMsg struct{}
type addMsg struct{}
type errorMsg struct{ error }

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

func errorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg{err}
	}
}
