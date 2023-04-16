package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var InputData string

// var wg sync.WaitGroup

func InitialModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "..."
	textInput.Focus()
	textInput.CharLimit = 250
	textInput.Width = 100

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return model{
		input:     textInput,
		output:    "",
		err:       nil,
		spinner:   sp,
		isLoading: false,
	}
}

type model struct {
	tea.Model
	input     textinput.Model
	output    string
	err       error
	spinner   spinner.Model
	isLoading bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.isLoading = true
			// MockRequest(&wg, &m)
			return m, m.spinner.Tick
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

	m.input, cmd = m.input.Update(msg)
	InputData = m.input.Value()
	return m, cmd
}

func (m model) View() string {
	if m.isLoading {
		return m.spinner.View()
	}
	if m.output != "" {
		return m.output
	}

	return fmt.Sprintf("Enter prompt: (Esc to quit)\n%s", m.input.View())
}
