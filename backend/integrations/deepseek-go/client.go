package deepseek

import (
	"context"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
)

type DeepSeekClient struct {
	client *deepseek.Client
}

func NewDeepSeekClient(apiKey string) DeepSeekClient {
	client := deepseek.NewClient(apiKey)

	return DeepSeekClient{
		client: client,
	}
}

func (d *DeepSeekClient) ChatCompletions(messages []deepseek.ChatCompletionMessage) (*deepseek.ChatCompletionResponse, error) {
	chatReq := &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: messages,
	}

	chatResp, err := d.client.CreateChatCompletion(context.Background(), chatReq)
	if err != nil {
		log.Println("Error =>", err)
		return &deepseek.ChatCompletionResponse{}, err
	}
	return chatResp, nil
}