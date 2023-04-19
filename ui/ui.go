package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	gpt = "GPT"
	you = "YOU"
)

type chatMessage struct {
	sender  string
	content string
}

type model struct {
	tea.Model
	allowInput bool
	input      textinput.Model
	output     []chatMessage
	err        error
	spinner    spinner.Model
	isLoading  bool
}

func NewModel() model {
	textInput := textinput.New()
	textInput.Placeholder = "Hello AI"
	textInput.Focus()
	textInput.CharLimit = 250
	textInput.Width = 100

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor))

	return model{
		allowInput: true,
		input:      textInput,
		output:     []chatMessage{},
		err:        nil,
		spinner:    sp,
		isLoading:  false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case aiResponse:
		response := withWordWrap(string(msg))
		m.output = append(m.output, chatMessage{gpt, response})
		m.isLoading = false
		m.allowInput = true
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.fetchAI()

		case tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlC:
			os.Exit(0)
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case error:
		m.err = msg
		return m, nil
	}

	if m.allowInput {
		m.input, _ = m.input.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	out := []string{
		withColor(themeColor, "AIcli\n"),
		m.renderOutput(),
	}

	if m.isLoading {
		out = append(out, m.spinner.View()+" uno momento...")
	} else {
		c := fmt.Sprintf("(%d/%d)", len(m.input.Value()), m.input.CharLimit)
		c = withColor(themeColor, c)
		count := fmt.Sprintf("\n%s%s", c, m.input.View())
		out = append(out, count)
	}

	return strings.Join(out, "\n")
}
