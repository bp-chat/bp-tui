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
	users       userList
	client      bp.Client
	chat        Chat
	activeModel tea.Model
}

type tickMsg time.Time

func New(config bp.Config, client bp.Client) Model {
	u := []user{
		{name: "user 1"},
		{name: "user 2"},
		{name: "user 3"},
	}
	return Model{
		client:      client,
		chat:        NewChat(client),
		users:       newUserList(client, u),
		activeModel: newGreeter(),
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
	return m.activeModel.Update(message)
}

func (m Model) View() string {
	return m.activeModel.View()
}
