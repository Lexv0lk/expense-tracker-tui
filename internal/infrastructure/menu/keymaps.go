package menu

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/key"
)

type ActionKeyMap struct {
	Delete key.Binding
	Create key.Binding
	Quit   key.Binding
	Edit   key.Binding
	GetSum key.Binding
	Filter key.Binding
	Export key.Binding
}

type NavigationKeyMap struct {
	Back    key.Binding
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
}

// ShortHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Create, km.Delete, km.Edit, km.Filter, km.GetSum, km.Export, km.Quit}
}

// FullHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Create, km.Delete, km.Edit, km.Filter, km.GetSum, km.Export, km.Quit},
	}
}

// ShortHelp implements the NavigationKeyMap interface.
func (km NavigationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Confirm, km.Up, km.Down, km.Back}
}

// FullHelp implements the NavigationKeyMap interface.
func (km NavigationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Confirm, km.Up, km.Down, km.Back},
	}
}

// getActionKeymap returns a default set of keybindings fot table actions.
func getActionKeymap() ActionKeyMap {
	return ActionKeyMap{
		Delete: constants.Keymap.Delete,
		Create: constants.Keymap.Create,
		Quit:   constants.Keymap.Quit,
		Edit:   constants.Keymap.Enter,
		GetSum: key.NewBinding(key.WithKeys("s"),
			key.WithHelp("s", "summary")),
		Filter: key.NewBinding(key.WithKeys("ctrl+f"),
			key.WithHelp("ctrl+f", "filter")),
		Export: key.NewBinding(key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "export to csv")),
	}
}

// getNavigationKeymap returns a default set of keybindings for navigation actions.
func getNavigationKeymap() NavigationKeyMap {
	return NavigationKeyMap{
		Back: constants.Keymap.Back,
		Up:   constants.Keymap.Up,
		Down: constants.Keymap.Down,
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}
}
