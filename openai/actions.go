package openai

import (
	"bytes"
	"encoding/json"
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

func ChatRequest(prompt string) (*OpenAIResponse, *OpenAIError) {
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

	body, err := json.Marshal(&req)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(body)

	request, _ := http.NewRequest("POST", Chat_URL, buf)
	request.Header.Add("authorization", "Bearer"+key)
	request.Header.Add("content-type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var err OpenAIError
		if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
			panic(err)
		}
		return nil, &err
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		panic(err)
	}

	return &response, nil
}
