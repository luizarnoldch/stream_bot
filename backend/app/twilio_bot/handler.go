package twiliobot

import (
	"log"
	"strings"

	"github.com/go-deepseek/deepseek/request"
	"github.com/gofiber/fiber/v2"
	"github.com/luizarnoldch/stream_bot/config"
	deepseek_service "github.com/luizarnoldch/stream_bot/services/deepseek"
	twilio_service "github.com/luizarnoldch/stream_bot/services/twilo"
)

type WhatsAppMessage struct {
	SmsMessageSid    string `form:"SmsMessageSid"`
	NumMedia         string `form:"NumMedia"`
	ProfileName      string `form:"ProfileName"`
	MessageType      string `form:"MessageType"`
	SmsSid           string `form:"SmsSid"`
	WaId             string `form:"WaId"`
	SmsStatus        string `form:"SmsStatus"`
	Body             string `form:"Body"`
	To               string `form:"To"`
	NumSegments      string `form:"NumSegments"`
	ReferralNumMedia string `form:"ReferralNumMedia"`
	MessageSid       string `form:"MessageSid"`
	AccountSid       string `form:"AccountSid"`
	From             string `form:"From"`
	ApiVersion       string `form:"ApiVersion"`
}

func SendSMS(c *fiber.Ctx) error {
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Parse incoming form data
	var msg WhatsAppMessage
	if err := c.BodyParser(&msg); err != nil {
		log.Printf("Error parsing form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al parsear formulario: " + err.Error(),
		})
	}
	deepSeekClient := deepseek_service.NewDeepSeekClient(env.MICRO.DEEPSEEK_API_KEY)
	input := []*request.Message{
		{Role: request.RoleSystem, Content: deepseek_service.PROMPT_STREAM_BOT_YOUTUBE},
		{Role: request.RoleUser, Content: msg.Body},
	}

	ai_response, err := deepSeekClient.ChatCompletionsRequest(input)
	if err != nil {
		log.Printf("Error from DeepSeek: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener respuesta de DeepSeek: " + err.Error(),
		})
	}

	// Format the AI response
	formattedText := strings.ReplaceAll(ai_response.Choices[0].Message.Content, "\n", " ")
	formattedText = strings.ReplaceAll(formattedText, "\r\n", " ")

	// Prepare Twilio client and send message
	twilClient := twilio_service.NewTwilioClient(
		env.MICRO.TWILIO.ACCOUNT_SID,
		env.MICRO.TWILIO.AUTH_TOKEN,
		env.MICRO.TWILIO.WHATSAPP_NUMBER,
	)

	twilio_response, err := twilClient.SendWhatsAppMessage(msg.From, formattedText)
	if err != nil {
		log.Printf("Error from Twilio: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener respuesta de Twilo: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": twilio_response,
	})
}
