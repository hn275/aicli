package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var key string

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	key = os.Getenv("OPENAI_API_KEY")
}

func ChatRequest(prompt string) (string, error) {
	var client http.Client
	req := OpenAIRequest{
		Model: GPT35_TURBO,
		Messages: []RequestMessage{
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

	request, _ := http.NewRequest("POST", Chat_URL, buf)
	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("content-type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO: get error message from chat gpt here
		return "", fmt.Errorf("fetched failed: %v", resp.StatusCode)
	}

	var response OpenAIResponse
	json.NewDecoder(resp.Body).Decode(&response)

	result := response.Choices[0]
	return result.Message.Content, nil
}
