package openai

import (
	"context"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) OpenAIClient {
	client := openai.NewClient(apiKey)

	return OpenAIClient{
		client: client,
	}
}

func (d *OpenAIClient) ChatCompletions(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	chatReq := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini20240718,
		Stream:   false,
		Messages: messages,
	}

	chatResp, err := d.client.CreateChatCompletion(context.Background(), chatReq)
	if err != nil {
		log.Println("Error =>", err)
		return openai.ChatCompletionResponse{}, err
	}
	return chatResp, nil
}
