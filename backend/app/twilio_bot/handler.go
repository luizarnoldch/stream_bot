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
	rawBody := c.Body()
	log.Println("Raw Request Body:", string(rawBody))
	log.Println("[SendSMS] Starting request processing")

	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[SendSMS] Error loading config: %v", err)
	}
	log.Println("[SendSMS] Config loaded successfully")

	// Parse incoming form data
	var msg WhatsAppMessage
	if err := c.BodyParser(&msg); err != nil {
		log.Printf("[SendSMS] Error parsing form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al parsear formulario: " + err.Error(),
		})
	}
	log.Printf("[SendSMS] Parsed WhatsApp message: %+v", msg)

	// Prepare DeepSeek client and request
	deepSeekClient := deepseek_service.NewDeepSeekClient(env.MICRO.DEEPSEEK_API_KEY)
	input := []*request.Message{
		{Role: request.RoleSystem, Content: deepseek_service.PROMPT_STREAM_BOT_YOUTUBE},
		{Role: request.RoleUser, Content: msg.Body},
	}
	log.Printf("[SendSMS] Sending to DeepSeek: %s", msg.Body)

	ai_response, err := deepSeekClient.ChatCompletionsRequest(input)
	if err != nil {
		log.Printf("[SendSMS] Error from DeepSeek: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener respuesta de DeepSeek: " + err.Error(),
		})
	}
	log.Printf("[SendSMS] Received DeepSeek response: %s", ai_response.Choices[0].Message.Content)

	// Format the AI response
	formattedText := strings.ReplaceAll(ai_response.Choices[0].Message.Content, "\n", " ")
	formattedText = strings.ReplaceAll(formattedText, "\r\n", " ")
	log.Printf("[SendSMS] Formatted text for WhatsApp: %s", formattedText)

	// Prepare Twilio client and send message
	twilClient := twilio_service.NewTwilioClient(
		env.MICRO.TWILIO.ACCOUNT_SID,
		env.MICRO.TWILIO.AUTH_TOKEN,
		env.MICRO.TWILIO.WHATSAPP_NUMBER,
	)
	log.Printf("[SendSMS] Sending WhatsApp message to %s", msg.From)

	twilio_response, err := twilClient.SendWhatsAppMessage(msg.From, formattedText)
	if err != nil {
		log.Printf("[SendSMS] Error from Twilio: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener respuesta de Twilo: " + err.Error(),
		})
	}
	log.Printf("[SendSMS] Twilio response: %+v", twilio_response)

	// Final response
	log.Println("[SendSMS] Successfully sent response back to caller")
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": twilio_response,
	})
}
