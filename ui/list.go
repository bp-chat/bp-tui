package ui

import (
	bp "github.com/bp-chat/bp-tui/client"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type userList struct {
	client bp.Client
	list   list.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type user struct {
	name string
}

func (i user) Title() string       { return i.name }
func (i user) Description() string { return "" }
func (i user) FilterValue() string { return i.name }

func mapUsers(users []user) []list.Item {
	result := make([]list.Item, len(users))
	for _, u := range users {
		result = append(result, u)
	}
	return result
}

func newList(client bp.Client, users []user) userList {
	return userList{
		client: client,
		list:   list.New(mapUsers(users), list.NewDefaultDelegate(), 0, 0),
	}
}

func (m userList) Init() tea.Cmd {
	return nil
}

func (m userList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m userList) View() string {
	return docStyle.Render(m.list.View())
}
