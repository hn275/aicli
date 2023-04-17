package ui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hn275/aicli/openai"
	"github.com/joho/godotenv"
)

type aiResponse string

func (m model) renderResponse(response aiResponse) (model, tea.Cmd) {
	var cmd tea.Cmd
	m.output = response
	m.isLoading = false
	m.input.Reset()

	return m, cmd

}

func (m model) fetchAI() (model, tea.Cmd) {
	if m.input.Value() == "" || !m.onInput {
		return m, nil
	}

	m.isLoading = true
	m.onInput = false
	cmd := tea.Batch(m.spinner.Tick, fetch(&m))
	return m, cmd
}

func fetch(m *model) tea.Cmd {
	if err := godotenv.Load(); err != nil {
		m.err = errors.New("[ERROR] unable to load env file")
		return nil
	}

	key := os.Getenv("OPENAI_API_KEY")

	return func() tea.Msg {
		time.Sleep(time.Second)
		return aiResponse(m.input.Value())
		var client http.Client
		req := openai.OpenAIRequest{
			Model: openai.GPT35_TURBO,
			Messages: []openai.RequestMessage{
				{
					Role:    "user",
					Content: m.input.Value(),
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
