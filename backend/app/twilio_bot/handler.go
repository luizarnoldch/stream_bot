package twiliobot

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/go-deepseek/deepseek/request"
	"github.com/gofiber/fiber/v2"
	"github.com/luizarnoldch/stream_bot/config"
	"github.com/luizarnoldch/stream_bot/db"

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
	var msg WhatsAppMessage
	if err := c.BodyParser(&msg); err != nil {
		log.Printf("Error parsing form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al parsear formulario: " + err.Error(),
		})
	}

	ctx := context.Background()
	env, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	psqlClient := db.GetPSQLClient(env.GetConnString())
	deepSeekClient := deepseek_service.NewDeepSeekClient(env.MICRO.DEEPSEEK_API_KEY)
	twilClient := twilio_service.NewTwilioClient(
		env.MICRO.TWILIO.ACCOUNT_SID,
		env.MICRO.TWILIO.AUTH_TOKEN,
		env.MICRO.TWILIO.WHATSAPP_NUMBER,
	)

	var (
		user         db.AuthUser
		conversation db.AiConversation
	)

	// Manejo de usuario existente/nuevo
	user, err = psqlClient.GetUserByPhone(ctx, msg.From)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Crear nuevo usuario
			newUser, err := psqlClient.CreateUser(ctx, msg.ProfileName, msg.From)
			if err != nil {
				log.Printf("Error creating user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error al crear usuario: " + err.Error(),
				})
			}

			// Crear nueva conversación
			newConversation, err := psqlClient.CreateConversation(ctx, "Conversacion de "+msg.ProfileName, newUser.ID)
			if err != nil {
				log.Printf("Error creating conversation: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error al crear conversación: " + err.Error(),
				})
			}

			user = newUser
			conversation = newConversation

			// Crear mensaje de sistema inicial
			_, err = psqlClient.CreateMessage(ctx, deepseek_service.PROMPT_STREAM_BOT_YOUTUBE, "system", newConversation.ID)
			if err != nil {
				log.Printf("Error creating system message: %v", err)
			}
		} else {
			log.Printf("Error getting user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al obtener usuario: " + err.Error(),
			})
		}
	} else {
		// Obtener conversación existente
		conversations, err := psqlClient.ListConversationsByUser(ctx, user.ID)
		if err != nil || len(conversations) == 0 {
			log.Printf("Error getting conversations: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al obtener conversaciones",
			})
		}
		conversation = conversations[0]
	}

	// Guardar mensaje del usuario
	_, err = psqlClient.CreateMessage(ctx, msg.Body, "user", conversation.ID)
	if err != nil {
		log.Printf("Error saving user message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar mensaje: " + err.Error(),
		})
	}

	// Obtener historial de mensajes
	messages, err := psqlClient.ListMessagesByConversation(ctx, conversation.ID)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener mensajes: " + err.Error(),
		})
	}

	// Construir input para DeepSeek
	input := make([]*request.Message, 0, len(messages))
	for _, msg := range messages {
		input = append(input, &request.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Obtener respuesta de AI
	aiResponse, err := deepSeekClient.ChatCompletionsRequest(input)
	if err != nil {
		log.Printf("Error from DeepSeek: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener respuesta de DeepSeek: " + err.Error(),
		})
	}

	// Formatear y guardar respuesta
	formattedText := strings.ReplaceAll(aiResponse.Choices[0].Message.Content, "\n", " ")
	formattedText = strings.ReplaceAll(formattedText, "\r\n", " ")

	_, err = psqlClient.CreateMessage(ctx, formattedText, "assistant", conversation.ID)
	if err != nil {
		log.Printf("Error saving AI message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar respuesta: " + err.Error(),
		})
	}

	// Enviar por Twilio
	_, err = twilClient.SendWhatsAppMessage(msg.From, formattedText)
	if err != nil {
		log.Printf("Error from Twilio: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al enviar mensaje: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Mensaje procesado correctamente",
	})
}
