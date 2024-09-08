package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type greeter struct {
	textInput textinput.Model
	err       error
}

func newGreeter() greeter {
	ti := textinput.New()
	ti.Placeholder = "Heisenberg"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20

	return greeter{
		textInput: ti,
		err:       nil,
	}
}

func (m greeter) Init() tea.Cmd {
	return textinput.Blink
}

func (m greeter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m greeter) View() string {
	return fmt.Sprintf(
		"What’s your name?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
