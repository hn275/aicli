package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hn275/aicli/openai"
)

type aiResponse string

func (m model) renderOutput() string {
	out := ""
	for _, v := range m.output {
		col := 69
		if v.sender == gpt {
			col = 123
		}

		out += fmt.Sprintf("%s:\n%s\n\n", withColor(col, v.sender), v.content)
	}

	return out
}

func (m model) fetchAI() (model, tea.Cmd) {
	if m.input.Value() == "" || !m.allowInput {
		return m, nil
	}

	m.output = append(m.output, chatMessage{you, m.input.Value()})

	m.isLoading = true
	m.allowInput = false
	prompt := m.input.Value()
	m.input.Reset()
	cmd := tea.Batch(m.spinner.Tick, fetch(prompt))
	return m, cmd
}

//lint:ignore U1000 for debugging and developing, so you don't use up your token
func dbg(a string) tea.Msg {
	time.Sleep(time.Second)
	return aiResponse(a)

}

func fetch(prompt string) tea.Cmd {
	return func() tea.Msg {
		// return dbg(prompt)
		result, err := openai.ChatRequest(prompt)
		if err != nil {
			return aiResponse(err.Error.Message)
		}
		return aiResponse(result.Choices[0].Message.Content)
	}
}

func withColor(color int, content string) string {
	if color < 0 || color > 256 {
		panic("color out of range: 0 - 256")
	}

	col := fmt.Sprintf("%d", color)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render(content)
}
