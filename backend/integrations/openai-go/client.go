package opeani_go

import (
	"context"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client openai.Client
}

func NewOpenAIClient(apiKey string) OpenAIClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))

	return OpenAIClient{
		client: client,
	}
}

func (d *OpenAIClient) ChatCompletions(messages []openai.ChatCompletionMessageParamUnion) (*openai.ChatCompletion, error) {
	chatReq := openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4oMini,
		Messages: messages,
	}

	chatResp, err := d.client.Chat.Completions.New(context.Background(), chatReq)
	if err != nil {
		log.Println("Error =>", err)
		return &openai.ChatCompletion{}, err
	}
	return chatResp, nil
}
