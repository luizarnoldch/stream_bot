package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/luizarnoldch/stream_bot/app/handlers"
	"github.com/luizarnoldch/stream_bot/config"
	"github.com/luizarnoldch/stream_bot/db"
	"github.com/luizarnoldch/stream_bot/types"

	deepseek_service "github.com/luizarnoldch/stream_bot/integrations/deepseek-go"
	openai_service "github.com/luizarnoldch/stream_bot/integrations/openai-go"
	twilio_service "github.com/luizarnoldch/stream_bot/integrations/twilio"
)

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	psqlClient := db.GetPSQLClient(env.GetConnString())
	deepSeekClient := deepseek_service.NewDeepSeekClient(env.MICRO.DEEPSEEK_API_KEY)
	openAIClient := openai_service.NewOpenAIClient(env.MICRO.OPENAI_API_KEY)
	twilClient := twilio_service.NewTwilioClient(
		env.MICRO.TWILIO.ACCOUNT_SID,
		env.MICRO.TWILIO.AUTH_TOKEN,
		env.MICRO.TWILIO.WHATSAPP_NUMBER,
		env.MICRO.TWILIO.MESSAGING_SERVICE_SID,
	)

	newHttpServer := types.NewHTTPServer(
		psqlClient,
		env,
		&deepSeekClient,
		&openAIClient,
		&twilClient,
	)

	router := handlers.NewHandlers(*newHttpServer)

	app := fiber.New()
	app.Post("/webhook", router.Twilio.TwilioWebhook)

	log.Println("Server started on port " + env.MICRO.API.PORT)
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.MICRO.API.HOST, env.MICRO.API.PORT)))
}
