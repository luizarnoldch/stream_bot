package twilio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-deepseek/deepseek/request"
	"github.com/luizarnoldch/stream_bot/db"
)

func TwilioImageProcessor(t *TwilioHandler, msg TwilioWhatsAppMessage) (string, error) {
	ctx := context.Background()

	user, err := t.app.DB.GetUserByPhone(ctx, msg.From)
	if errors.Is(err, sql.ErrNoRows) {
		_, conversation, err := handleUserCreation(ctx, t, msg.From, msg.ProfileName)
		if err != nil {
			return "", err
		}
		return processImageConversation(ctx, t, msg, conversation)
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return "", fmt.Errorf("error al obtener usuario: %w", err)
	}

	conversations, err := t.app.DB.ListConversationsByUser(ctx, int32(user.ID))
	if err != nil {
		log.Printf("Error getting conversations: %v", err)
		return "", fmt.Errorf("error al obtener conversaciones: %w", err)
	}

	return processImageConversation(ctx, t, msg, conversations[0])
}

func processImageConversation(ctx context.Context, t *TwilioHandler, msg TwilioWhatsAppMessage, conversation db.AiConversation) (string, error) {
	img_urls, err := prepareImageDataURLs(t, msg)
	if err != nil {
		log.Printf("Error preparing image data URLs: %v", err)
		return "", fmt.Errorf("error al preparar URL de imagen: %w", err)
	}

	formattedAIResponseText, err := processOpenAIDataURLsMessages(ctx, t, msg.Body, img_urls, conversation)
	if err != nil {
		log.Printf("Error processing OpenAI data URLs messages: %v", err)
		return "", fmt.Errorf("error al procesar mensajes de OpenAI: %w", err)
	}
	
	if _, err := t.app.DB.CreateMessage(ctx, formattedAIResponseText, request.RoleSystem, int32(conversation.ID)); err != nil {
		log.Printf("Error guardando mensaje de sistema: %v", err)
		return "", fmt.Errorf("error al guardar mensaje de sistema: %w", err)
	}

	return formattedAIResponseText, nil
}
