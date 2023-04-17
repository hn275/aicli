package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hn275/aicli/openai"
	"github.com/joho/godotenv"
)

type aiResponse string

func (m model) renderOutput() string {
	out := ""
	for _, v := range m.output {
		col := 69
		if v.sender == gpt {
			col = 123
		}

		out += fmt.Sprintf("  %s: %s\n", withColor(col, v.sender), v.content)
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

// DBG
func dbg(a string) tea.Msg {
	time.Sleep(time.Second)
	return aiResponse(a)

}
func fetch(prompt string) tea.Cmd {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	key := os.Getenv("OPENAI_API_KEY")

	return func() tea.Msg {
		return dbg(prompt)
		var client http.Client
		req := openai.OpenAIRequest{
			Model: openai.GPT35_TURBO,
			Messages: []openai.RequestMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		}

		bodyJson, err := json.Marshal(&req)
		if err != nil {
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(bodyJson)

		request, _ := http.NewRequest("POST", openai.Chat_URL, buf)
		request.Header.Add("authorization", fmt.Sprintf("Bearer %s", key))
		request.Header.Add("content-type", "application/json")

		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Sprintf("Fetched failed: %v", resp.StatusCode)
		}

		var response openai.OpenAIResponse
		json.NewDecoder(resp.Body).Decode(&response)

		result := response.Choices[0]
		return aiResponse(result.Message.Content)
	}
}

func withColor(color int, content string) string {
	if color < 0 || color > 256 {
		panic("color out of range: 0 - 256")
	}

	col := fmt.Sprintf("%d", color)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render(content)
}
