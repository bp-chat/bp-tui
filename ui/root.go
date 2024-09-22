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

type vec2 struct {
	x int
	y int
}

type Model struct {
	users       []user
	client      *bp.Client
	activeModel tea.Model
}

type tickMsg time.Time

type connectionFailedMsg struct{ err error }
type userConnectedMsg struct {
	connectedClient bp.Client
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var windowSize vec2

func New(config bp.Config) Model {
	var client *bp.Client
	u := []user{
		{name: "user 1", canChat: false},
		{name: "user 2", canChat: false},
		{name: "user 3", canChat: false},
	}
	greeter := newGreeter(config)
	return Model{
		users:       u,
		client:      client,
		activeModel: greeter,
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
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		windowSize.x = msg.Width - h
		windowSize.y = msg.Height - v
	case userConnectedMsg:
		m.client = &msg.connectedClient
		m.activeModel = newUserList(m.client, &m.users, windowSize)
		return m, cmd
	}
	m.activeModel, cmd = m.activeModel.Update(message)
	return m, cmd
}

func (m Model) View() string {
	return m.activeModel.View()
}
