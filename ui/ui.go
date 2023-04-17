package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tea.Model
	onInput   bool
	input     textinput.Model
	output    aiResponse
	err       error
	spinner   spinner.Model
	isLoading bool
}

func NewModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "Hello AI"
	textInput.Focus()
	textInput.CharLimit = 250
	textInput.Width = 100

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	return model{
		onInput:   true,
		input:     textInput,
		output:    "",
		err:       nil,
		spinner:   sp,
		isLoading: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case aiResponse:
		return m.renderResponse(msg)

	case tea.KeyMsg:
		k := msg.String()
		if k == "r" && !m.onInput {
			m.onInput = true
			return m, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			return m.fetchAI()

		case tea.KeyEsc:
			return m, tea.Quit

		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case error:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}

	if m.onInput {
		m.input, _ = m.input.Update(msg)
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {

	if m.isLoading {
		return m.spinner.View()
	}

	if !m.onInput {
		out := fmt.Sprintf("%s\n\nr to restart\nescape to quit", m.output)
		return out
	}

	if m.onInput {
		return fmt.Sprintf(
			"Enter prompt:\n%s\n%s",
			m.input.View(), "(esc to quit)",
		)
	}

	return ""
}
