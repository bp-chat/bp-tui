package ui

import (
	"time"

	bp "github.com/bp-chat/bp-tui/client"
	"github.com/bp-chat/bp-tui/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	foregroundColor       lipgloss.Color = lipgloss.Color("#ffffff")
	foregroundAccentColor lipgloss.Color = lipgloss.Color("#ff0000")
)

func connect(config bp.Config, name string) tea.Cmd {
	return func() tea.Msg {

		var username commands.UserName
		copy(username[:], name[:])
		eu := bp.EphemeralUser{
			Name: username,
			Keys: bp.CreateKeys(),
		}
		conn, err := bp.Connect(config.Host)
		if err != nil {
			return connectionFailedMsg{err}
		}
		client := bp.New(eu, conn)
		err = client.RefreshKeys()
		if err != nil {

			return connectionFailedMsg{err: nil}
		}
		return userConnectedMsg{
			connectedClient: client,
		}
	}
}

type Model struct {
	users       userList
	client      *bp.Client
	chat        Chat
	activeModel tea.Model
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
	return Model{
		client:      client,
		chat:        NewChat(client),
		users:       newUserList(client, u),
		activeModel: newGreeter(config),
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
	case userConnectedMsg:
		m.client = &msg.connectedClient
		m.activeModel = m.users
		return m, nil
	}
	return m.activeModel.Update(message)
}

func (m Model) View() string {
	return m.activeModel.View()
}
