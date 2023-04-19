package openai

const Chat_URL = "https://api.openai.com/v1/chat/completions"

const GPT35_TURBO = "gpt-3.5-turbo"

type OpenAIError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type OpenAIRequest struct {
	Model    string           `json:"model"`
	Messages []RequestMessage `json:"messages"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Usage   ResponseUsage    `json:"usage"`
	Choices []ResponseChoice `json:"choices"`
}

type ResponseChoice struct {
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Index        int64           `json:"index"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}
