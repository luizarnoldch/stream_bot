package twilio

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/cohesion-org/deepseek-go"
	"github.com/go-deepseek/deepseek/request"
	"github.com/luizarnoldch/stream_bot/db"
	"github.com/luizarnoldch/stream_bot/integrations/prompts"
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
	if err != nil || len(conversations) == 0 {
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

	formattedText, err := processOpenAIDataURLsMessages(ctx, t, msg.Body, img_urls, conversation)
	if err != nil {
		log.Printf("Error processing OpenAI data URLs messages: %v", err)
		return "", fmt.Errorf("error al procesar mensajes de OpenAI: %w", err)
	}

	if _, err := t.app.DB.CreateMessage(ctx, formattedText, request.RoleSystem, int32(conversation.ID)); err != nil {
		log.Printf("Error guardando mensaje de sistema: %v", err)
		return "", fmt.Errorf("error al guardar mensaje de sistema: %w", err)
	}

	operation, err := processOperationTextToObjectFormat(t, formattedText)
	if err != nil {
		log.Printf("Error processing operation text to object format: %v", err)
		return "", fmt.Errorf("error al procesar texto de operación a formato de objeto: %w", err)
	}

	payment, err := t.app.DB.CreatePaymentOperation(ctx, operation)
	if err != nil {
		log.Printf("Error creating payment operation: %v", err)
		return "", fmt.Errorf("error al crear operación de pago: %w", err)
	}

	amountFloat64, err := payment.AmountSent.Float64Value()
	if err != nil {
		log.Printf("Error getting amount value: %v", err)
		return "", fmt.Errorf("error al obtener valor del monto: %w", err)
	}
	paymentResponse := fmt.Sprintf("Gracias %s tu pago con el codigo de opearción %s registrado con éxito. El monto enviado es de %.2f %s", msg.ProfileName, payment.OperationNumber, amountFloat64.Float64, payment.Currency)

	return paymentResponse, nil
}

func processOperationTextToObjectFormat(t *TwilioHandler, text string) (db.CreatePaymentOperationParams, error) {

	input := []deepseek.ChatCompletionMessage{
		{
			Role:    deepseek.ChatMessageRoleSystem,
			Content: prompts.PROMPT_STREAM_PAGOS_OBJECT,
		},
		{
			Role:    deepseek.ChatMessageRoleUser,
			Content: text,
		},
	}

	aiResponse, err := t.app.DeepSeek.ChatCompletions(input)
	if err != nil {
		log.Printf("Error from DeepSeek: %v", err)
		return db.CreatePaymentOperationParams{}, fmt.Errorf("error de DeepSeek: %w", err)
	}

	formattedText := formatResponse(aiResponse.Choices[0].Message.Content)
	jsonFormattetdText := formatJSONResponse(formattedText)

	var operation db.CreatePaymentOperationParams
	if err := json.Unmarshal([]byte(jsonFormattetdText), &operation); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return db.CreatePaymentOperationParams{}, fmt.Errorf("error al deserializar JSON: %w", err)
	}

	return operation, nil
}

// func processOperationTextToObjectFormat(ctx context.Context, t *TwilioHandler, text string) (db.CreatePaymentOperationParams, error) {

// 	input := []deepseek.ChatCompletionMessage{
// 		{
// 			Role:    deepseek.ChatMessageRoleSystem,
// 			Content: prompts.PROMPT_STREAM_PAGOS_OBJECT,
// 		},
// 		{
// 			Role:    deepseek.ChatMessageRoleUser,
// 			Content: text,
// 		},
// 	}

// 	aiResponse, err := t.app.DeepSeek.ChatCompletions(input)
// 	if err != nil {
// 		log.Printf("Error from DeepSeek: %v", err)
// 		return db.CreatePaymentOperationParams{}, fmt.Errorf("error de DeepSeek: %w", err)
// 	}

// 	formattedText := formatResponse(aiResponse.Choices[0].Message.Content)

// 	var operation db.CreatePaymentOperationParams
// 	if err := json.Unmarshal([]byte(formattedText), &operation); err != nil {
// 		log.Printf("Error unmarshalling JSON: %v", err)
// 		return db.CreatePaymentOperationParams{}, fmt.Errorf("error al deserializar JSON: %w", err)
// 	}

// }
