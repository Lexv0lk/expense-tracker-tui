package menu

import (
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/menu/constants"
	"github.com/charmbracelet/bubbles/key"
)

type ActionKeyMap struct {
	Delete key.Binding
	Create key.Binding
}

// ShortHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Delete, km.Create}
}

// FullHelp implements the ActionKeyMap interface.
func (km ActionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Delete, km.Create},
	}
}

// getDefaultKeymap returns a default set of keybindings fot table actions.
func getDefaultKeymap() ActionKeyMap {
	return ActionKeyMap{
		Delete: constants.Keymap.Delete,
		Create: constants.Keymap.Create,
	}
}
