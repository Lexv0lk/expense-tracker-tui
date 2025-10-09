package menu

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const continueMsg = "Press any key to continue..."

type msgModel struct {
	viewport viewport.Model
	backCmd  tea.Cmd
	err      error
}

func newMsgModel(content string, backCmd tea.Cmd) msgModel {
	const width = 4 + 30 + 10 + 12 + 15
	const height = 7

	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	// We need to adjust the width of the glamour render from our main width
	// to account for a few things:
	//
	//  * The viewport border width
	//  * The viewport padding
	//  * The viewport margins
	//  * The gutter glamour applies to the left side of the content
	//
	const glamourGutter = 2
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return msgModel{
			err:     fmt.Errorf("Error creating glamour renderer: %w", err),
			backCmd: backCmd,
		}
	}

	str, err := renderer.Render(content)
	if err != nil {
		return msgModel{
			err:     fmt.Errorf("Error rendering content: %w", err),
			backCmd: backCmd,
		}
	}

	vp.SetContent(str)

	return msgModel{
		viewport: vp,
		err:      nil,
		backCmd:  backCmd,
	}
}

func (m msgModel) Init() tea.Cmd { return nil }

func (m msgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, m.backCmd
	}

	return m, nil
}

func (m msgModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error rendering message: %v\n%s", m.err, continueMsg)
	}

	return m.viewport.View() + "\n" + continueMsg
}
