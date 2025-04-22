package twilio

import (
	app "github.com/luizarnoldch/stream_bot/types"
)

type TwilioHandler struct {
	app app.HTTPServer
}

func NewTwilioHandler(core app.HTTPServer) *TwilioHandler {
	return &TwilioHandler{
		app: core,
	}
}
