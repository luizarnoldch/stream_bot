package handlers

import (
	"github.com/luizarnoldch/stream_bot/app/handlers/twilio"
	"github.com/luizarnoldch/stream_bot/types"
)

type Handlers struct {
	Twilio twilio.TwilioHandler
}

func NewHandlers(app types.HTTPServer) *Handlers {
	return &Handlers{
		Twilio: *twilio.NewTwilioHandler(app),
	}
}
