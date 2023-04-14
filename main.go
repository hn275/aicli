package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hn275/aicli/openai"
	"github.com/joho/godotenv"
)

var key string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	key = os.Getenv("OPENAI_API_KEY")
}

func main() {
	var client http.Client
	req := openai.OpenAIRequest{
		Model: openai.GPT35_TURBO,
		Messages: []openai.RequestMessage{
			{Role: "user", Content: "Give me code for a drawer React component with TypeScript, no explanation needed"},
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

	var response openai.OpenAIResponse
	json.NewDecoder(resp.Body).Decode(&response)

	choices := response.Choices
	for _, v := range choices {
		log.Println(v.Message.Content)
	}

}
