package twilio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func TwilioTextProcessor(t *TwilioHandler, msg TwilioWhatsAppMessage) (string, error) {
	ctx := context.Background()

	// User handling
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

	// Existing user conversation
	conversations, err := t.app.DB.ListConversationsByUser(ctx, int32(user.ID))
	if err != nil || len(conversations) == 0 {
		log.Printf("Error getting conversations: %v", err)
		return "", fmt.Errorf("error al obtener conversaciones: %w", err)
	}

	return processTextConversation(ctx, t, msg, conversations[0])
}
