package ui

import (
	"fmt"

	bp "github.com/bp-chat/bp-tui/client"
	"github.com/bp-chat/bp-tui/commands"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type greeter struct {
	config    bp.Config
	textInput textinput.Model
	err       error
	info      string
}

func newGreeter(config bp.Config) greeter {
	ti := textinput.New()
	ti.Placeholder = "Heisenberg"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20

	return greeter{
		config:    config,
		textInput: ti,
		err:       nil,
		info:      "",
	}
}

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

func (m greeter) Init() tea.Cmd {
	return textinput.Blink
}

func (m greeter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case connectionFailedMsg:
		m.textInput.Reset()
		m.info = "could not connect to the server"
		m.err = msg.err
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.info = "connecting..."
			return m, connect(m.config, m.textInput.Value())
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m greeter) View() string {
	return fmt.Sprintf(
		"What’s your name?\n\n%s\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
		m.info,
	) + "\n"
}
