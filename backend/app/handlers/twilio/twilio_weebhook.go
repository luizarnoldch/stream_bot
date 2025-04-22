package twilio

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type TwilioWhatsAppMessage struct {
	SmsMessageSid              string `form:"SmsMessageSid"    json:"SmsMessageSid"`
	NumMedia                   int    `form:"NumMedia"         json:"NumMedia"`
	ProfileName                string `form:"ProfileName"      json:"ProfileName"`
	MessageType                string `form:"MessageType"      json:"MessageType"`
	SmsSid                     string `form:"SmsSid"           json:"SmsSid"`
	WaId                       string `form:"WaId"             json:"WaId"`
	SmsStatus                  string `form:"SmsStatus"        json:"SmsStatus"`
	Body                       string `form:"Body"             json:"Body"`
	To                         string `form:"To"               json:"To"`
	NumSegments                int    `form:"NumSegments"      json:"NumSegments"`
	ReferralNumMedia           int    `form:"ReferralNumMedia" json:"ReferralNumMedia"`
	MessageSid                 string `form:"MessageSid"       json:"MessageSid"`
	AccountSid                 string `form:"AccountSid"       json:"AccountSid"`
	From                       string `form:"From"             json:"From"`
	ApiVersion                 string `form:"ApiVersion"       json:"ApiVersion"`
	TwilioWhatsAppMessageMedia []TwilioWhatsAppMessageMedia
}

type TwilioWhatsAppMessageMedia struct {
	MediaContentType string
	MediaUrl         string
}

func (t *TwilioHandler) TwilioWebhook(c *fiber.Ctx) error {
	var msg TwilioWhatsAppMessage
	if err := c.BodyParser(&msg); err != nil {
		log.Printf("Error parsing form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error processing request: " + err.Error(),
		})
	}

	for i := range msg.NumMedia {
		url := c.FormValue(fmt.Sprintf("MediaUrl%d", i))
		ctype := c.FormValue(fmt.Sprintf("MediaContentType%d", i))
		if url != "" {
			msg.TwilioWhatsAppMessageMedia = append(msg.TwilioWhatsAppMessageMedia, TwilioWhatsAppMessageMedia{
				MediaUrl:         url,
				MediaContentType: ctype,
			})
		}
	}

	var respone string

	switch msg.MessageType {
	case "text":
		log.Printf("📩 Text message from %s: %q\n", msg.From, msg.Body)

		// Procesar mensaje
		formattedText, err := TwilioTextProcessor(t, msg)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing message: " + err.Error(),
			})
		}
		respone = formattedText
	case "image":
		log.Printf("📩 Image message from %s: %q\n", msg.From, msg.Body)

		formattedText, err := TwilioImageProcessor(t, msg)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing message: " + err.Error(),
			})
		}
		respone = formattedText
	case "audio", "video":
		log.Printf("📩 Video and Audio message from %s: %q\n", msg.From, msg.Body)

	default:
		log.Printf("❓ Unknown message type: %s\n", msg.MessageType)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported message type",
		})
	}

	// Enviar respuesta por Twilio
	if _, err := t.app.Twilio.SendWhatsAppMessage(msg.From, respone); err != nil {
		log.Printf("Error sending Twilio message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error sending message: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Message processed successfully",
	})
}
