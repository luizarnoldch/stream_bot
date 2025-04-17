package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	tb "github.com/luizarnoldch/stream_bot/app/twilio_bot"
	"github.com/luizarnoldch/stream_bot/config"
)

func Server() {
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := fiber.New()
	app.Post("/send-sms", tb.SendSMS)

	log.Println("Server started on port " + env.MICRO.API.PORT)
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.MICRO.API.HOST, env.MICRO.API.PORT)))
}
