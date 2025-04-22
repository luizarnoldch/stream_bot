package types

import (
	"github.com/luizarnoldch/stream_bot/config"
	"github.com/luizarnoldch/stream_bot/db"

	"github.com/luizarnoldch/stream_bot/integrations/deepseek"
	"github.com/luizarnoldch/stream_bot/integrations/openai"
	"github.com/luizarnoldch/stream_bot/integrations/twilio"
)

type HTTPServer struct {
	DB       *db.Queries
	Config   *config.CONFIG
	DeepSeek *deepseek.DeepSeekClient
	OpenAI   *openai.OpenAIClient
	Twilio   *twilio.TwilioClient
}

func NewHTTPServer(db *db.Queries, config *config.CONFIG, deepSeek *deepseek.DeepSeekClient, openai *openai.OpenAIClient, twilio *twilio.TwilioClient) *HTTPServer {
	return &HTTPServer{
		DB:       db,
		Config:   config,
		DeepSeek: deepSeek,
		OpenAI:   openai,
		Twilio:   twilio,
	}
}
