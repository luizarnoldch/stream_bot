package twilio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/cohesion-org/deepseek-go"
	"github.com/luizarnoldch/stream_bot/db"
)

func TwilioTextProcessor(t *TwilioHandler, msg TwilioWhatsAppMessage) (string, error) {
	ctx := context.Background()

	user, err := t.app.DB.GetUserByPhone(ctx, msg.From)
	if errors.Is(err, sql.ErrNoRows) {
		_, conversation, err := handleUserCreation(ctx, t, msg.From, msg.ProfileName)
		if err != nil {
			return "", err
		}
		return processTextConversation(ctx, t, msg, conversation)
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return "", fmt.Errorf("error al obtener usuario: %w", err)
	}

	conversations, err := t.app.DB.ListConversationsByUser(ctx, int32(user.ID))
	if err != nil || len(conversations) == 0 {
		log.Printf("Error getting conversations: %v", err)
		return "", fmt.Errorf("error al obtener conversaciones: %w", err)
	}

	return processTextConversation(ctx, t, msg, conversations[0])
}

func processTextConversation(ctx context.Context, t *TwilioHandler, msg TwilioWhatsAppMessage, conversation db.AiConversation) (string, error) {
	_, err := t.app.DB.CreateMessage(ctx, msg.Body, deepseek.ChatMessageRoleUser, int32(conversation.ID))
	if err != nil {
		log.Printf("Error saving %s message: %v", deepseek.ChatMessageRoleUser, err)
		return "", fmt.Errorf("error al guardar mensaje de %s: %w", deepseek.ChatMessageRoleUser, err)
	}

	messages, err := t.app.DB.ListMessagesByConversation(ctx, conversation.ID)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return "", fmt.Errorf("error al obtener mensajes: %w", err)
	}

	aiResponse, err := t.app.DeepSeek.ChatCompletions(buildDeepSeekInput(messages))
	if err != nil {
		log.Printf("Error from DeepSeek: %v", err)
		return "", fmt.Errorf("error de DeepSeek: %w", err)
	}

	formattedText := formatResponse(aiResponse.Choices[0].Message.Content)
	_, err = t.app.DB.CreateMessage(ctx, formattedText, deepseek.ChatMessageRoleSystem, int32(conversation.ID))
	if err != nil {
		log.Printf("Error saving %s message: %v", deepseek.ChatMessageRoleSystem, err)
		return "", fmt.Errorf("error al guardar mensaje de %s: %w", deepseek.ChatMessageRoleSystem, err)
	}
	
	return formattedText, nil
}
