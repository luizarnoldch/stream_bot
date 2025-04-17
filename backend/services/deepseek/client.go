package deepseek

import (
	"context"
	"fmt"
	"log"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
	"github.com/go-deepseek/deepseek/response"
)

type DeepSeekClient struct {
	client deepseek.Client
}

func NewDeepSeekClient(apiKey string) DeepSeekClient {
	client, err := deepseek.NewClient(apiKey)
	if err != nil {
		log.Fatalln("Error creating DeepSeek client:", err)
	}

	return DeepSeekClient{
		client: client,
	}
}

func (d *DeepSeekClient) ChatCompletionsRequest(messages []*request.Message) (*response.ChatCompletionsResponse, error) {
	chatReq := &request.ChatCompletionsRequest{
		Model:    deepseek.DEEPSEEK_CHAT_MODEL,
		Stream:   false,
		Messages: messages,
	}

	chatResp, err := d.client.CallChatCompletionsChat(context.Background(), chatReq)
	if err != nil {
		fmt.Println("Error =>", err)
		return &response.ChatCompletionsResponse{}, err
	}
	return chatResp, nil
}
