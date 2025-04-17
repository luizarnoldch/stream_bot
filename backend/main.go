package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
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

var (
	TWILIO_ACCOUNT_SID     = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN      = os.Getenv("TWILIO_AUTH_TOKEN")
	TWILIO_WHATSAPP_NUMBER = os.Getenv("TWILIO_WHATSAPP_NUMBER")
	PORT                   = os.Getenv("PORT")
)

func sendSMS(c *fiber.Ctx) error {
	rawBody := c.Body()
	log.Println("Raw Request Body:", string(rawBody))

	// Parse form into struct
	var msg WhatsAppMessage
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing form: " + err.Error(),
		})
	}

	// Log parsed struct
	log.Printf("Parsed Struct: %+v\n", msg)

	// Preparar el cuerpo para el POST al endpoint /chat
	postBody, _ := json.Marshal(map[string]string{
		"prompt": msg.Body,
	})

	// Hacer el POST a http://localhost:4000/chat
	resp, err := http.Post("http://localhost:4000/chat", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al hacer el POST: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	// Leer la respuesta del servidor /chat
	var chatResponse struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al decodificar la respuesta: " + err.Error(),
		})
	}

	// Log de la respuesta del chat
	log.Printf("Respuesta del Chat: %+v\n", chatResponse)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(msg.From)
	params.SetFrom(TWILIO_WHATSAPP_NUMBER)
	params.SetBody(chatResponse.Text)

	respTwilio, err := client.Api.CreateMessage(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "sent",
		"sid":     *respTwilio.Sid,
		"message": chatResponse.Text, // Enviar el mensaje devuelto de la API /chat
	})
}

func main() {
	// Load .env if available
	_ = godotenv.Load()

	app := fiber.New()

	app.Post("/send-sms", sendSMS)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", PORT)))
}
