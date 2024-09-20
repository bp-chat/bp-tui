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
	client      *bp.Client
	activeModel tea.Model
	users       []user
}

type tickMsg time.Time

type connectionFailedMsg struct{ err error }
type userConnectedMsg struct {
	connectedClient bp.Client
}

func New(config bp.Config) Model {
	u := []user{
		{name: "user 1"},
		{name: "user 2"},
		{name: "user 3"},
	}
	var client *bp.Client
	greeter := newGreeter(config)
	return Model{
		client:      client,
		activeModel: greeter,
		users:       u,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case userConnectedMsg:
		m.client = &msg.connectedClient
		m.activeModel = newUserList(m.client, m.users)
		return m, nil
	}
	m.activeModel, cmd = m.activeModel.Update(message)
	return m, cmd
}

func (m Model) View() string {
	return m.activeModel.View()
}
