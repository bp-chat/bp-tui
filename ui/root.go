package ui

import (
	"time"

	bp "github.com/bp-chat/bp-tui/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	foregroundColor       lipgloss.Color = lipgloss.Color("#ffffff")
	foregroundAccentColor lipgloss.Color = lipgloss.Color("#ff0000")
)

type Model struct {
	client bp.Client
	chat   Chat
}

type tickMsg time.Time

func New(client bp.Client) Model {
	return Model{
		client: client,
		chat:   NewChat(client),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

		return m, nil
	}
	return m.chat.Update(message)
}

func (m Model) View() string {
	return m.chat.View()
}
