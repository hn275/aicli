package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var InputData string

type Model struct {
	tea.Model
	Input textinput.Model
	err   error
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, tea.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.Input, cmd = m.Input.Update(msg)
	InputData = m.Input.Value()
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("Enter prompt: (Esc to quit)\n%s", m.Input.View())
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter prompt"
	ti.Focus()
	ti.CharLimit = 250
	ti.Width = 100

	return Model{Input: ti, err: nil}
}
