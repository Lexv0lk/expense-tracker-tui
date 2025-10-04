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
}

type NavigationKeyMap struct {
	Back key.Binding
	Up   key.Binding
	Down key.Binding
}

// ShortHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Delete, km.Create, km.Quit, km.Edit}
}

// FullHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Delete, km.Create, km.Quit, km.Edit},
	}
}

// ShortHelp implements the NavigationKeyMap interface.
func (km NavigationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Back, km.Up, km.Down}
}

// FullHelp implements the NavigationKeyMap interface.
func (km NavigationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Back, km.Up, km.Down},
	}
}

// getActionKeymap returns a default set of keybindings fot table actions.
func getActionKeymap() ActionKeyMap {
	return ActionKeyMap{
		Delete: constants.Keymap.Delete,
		Create: constants.Keymap.Create,
		Quit:   constants.Keymap.Quit,
		Edit:   constants.Keymap.Enter,
	}
}

// getNavigationKeymap returns a default set of keybindings for navigation actions.
func getNavigationKeymap() NavigationKeyMap {
	return NavigationKeyMap{
		Back: constants.Keymap.Back,
		Up:   constants.Keymap.Up,
		Down: constants.Keymap.Down,
	}
}
