package ui

// most of this initial content was copied from the lib example
// https://github.com/charmbracelet/bubbletea/blob/master/examples/chat/main.go

import (
	"fmt"
	"strings"

	"github.com/bp-chat/bp-tui/commands"
	"github.com/bp-chat/bp-tui/commands/calls"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Chat struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	receipStyle lipgloss.Style
	err         error
	sendAction  func(string)
}

func New(send func(string)) Chat {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Chat{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		receipStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		err:         nil,
		sendAction:  send,
	}
}

func (chat Chat) Init() tea.Cmd {
	return nil
}

func (m Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			textMsg := m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+textMsg)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.sendAction(textMsg)
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
		break

		// We handle errors just like any other message
		//this probably worked in a previous version
	case error:
		m.err = msg
		return m, nil
	case *commands.Command:
		cmsg := calls.FromCommand(msg)
		m.messages = append(m.messages, m.receipStyle.Render("One: ")+cmsg.Message)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, tea.Batch(tiCmd, vpCmd)
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Chat) View() string {
	var s string
	s = fmt.Sprintf(
		"%s\n\n%s\n\n",
		m.viewport.View(),
		m.textarea.View(),
	)
	return s
}
