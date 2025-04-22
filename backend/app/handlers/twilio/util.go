package twilio

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-deepseek/deepseek/request"
	"github.com/luizarnoldch/stream_bot/db"
	"github.com/luizarnoldch/stream_bot/integrations/prompts"
	"github.com/sashabaranov/go-openai"
)

func formatResponse(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.ReplaceAll(text, "\r\n", " ")
}

func handleUserCreation(ctx context.Context, t *TwilioHandler, phone string, profileName string) (db.AuthUser, db.AiConversation, error) {
	newUser, err := t.app.DB.CreateUser(ctx, profileName, phone)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return db.AuthUser{}, db.AiConversation{}, fmt.Errorf("error al crear usuario: %w", err)
	}

	newConversation, err := t.app.DB.CreateConversation(ctx, "Conversacion de "+profileName, newUser.ID)
	if err != nil {
		log.Printf("Error creating conversation: %v", err)
		return newUser, db.AiConversation{}, fmt.Errorf("error al crear conversación: %w", err)
	}

	if _, err = t.app.DB.CreateMessage(ctx, prompts.PROMPT_STREAM_VENTAS_BOT, request.RoleSystem, newConversation.ID); err != nil {
		log.Printf("Error creating system message: %v", err)
	}

	return newUser, newConversation, nil
}

func buildDeepSeekInput(messages []db.AiMessage) []*request.Message {
	input := make([]*request.Message, 0, len(messages))
	for _, msg := range messages {
		input = append(input, &request.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return input
}

func processTextConversation(ctx context.Context, t *TwilioHandler, msg TwilioWhatsAppMessage, conversation db.AiConversation) (string, error) {
	_, err := t.app.DB.CreateMessage(ctx, msg.Body, request.RoleUser, int32(conversation.ID))
	if err != nil {
		log.Printf("Error saving %s message: %v", request.RoleUser, err)
		return "", fmt.Errorf("error al guardar mensaje de %s: %w", request.RoleUser, err)
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
	_, err = t.app.DB.CreateMessage(ctx, formattedText, request.RoleSystem, int32(conversation.ID))
	if err != nil {
		log.Printf("Error saving %s message: %v", request.RoleUser, err)
		return "", fmt.Errorf("error al guardar mensaje de %s: %w", request.RoleUser, err)
	}

	return formattedText, nil
}

func downloadPNG(t *TwilioHandler, dir, baseName, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creando petición: %w", err)
	}

	accountSid := t.app.Config.MICRO.TWILIO.ACCOUNT_SID
	authToken := t.app.Config.MICRO.TWILIO.AUTH_TOKEN
	req.SetBasicAuth(accountSid, authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en petición HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta no OK: %s", resp.Status)
	}

	ext := ".png"
	filename := baseName + ext
	fullpath := filepath.Join(dir, filename)

	out, err := os.Create(fullpath)
	if err != nil {
		return "", fmt.Errorf("error creando archivo: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("error guardando archivo: %w", err)
	}

	return fullpath, nil
}

func prepareImageDataURLs(t *TwilioHandler, msg TwilioWhatsAppMessage) ([]string, error) {
	dir := "downloads"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Error creando dir %s: %v", dir, err)
		return nil, fmt.Errorf("no se pudo crear directorio: %w", err)
	}

	var urls []string
	for idx, media := range msg.TwilioWhatsAppMessageMedia {
		file := fmt.Sprintf("%s_%d", msg.SmsMessageSid, idx)
		img_path, err := downloadPNG(t, dir, file, media.MediaUrl)
		if err != nil {
			log.Printf("Error descargando media %d: %v", idx, err)
			continue
		}

		img_bytes, err := os.ReadFile(img_path)
		if err != nil {
			log.Printf("Error leyendo %s: %v", img_path, err)
			continue
		}
		encoded := base64.StdEncoding.EncodeToString(img_bytes)
		urls = append(urls, fmt.Sprintf("data:image/png;base64,%s", encoded))
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("no se guardó ninguna imagen")
	}
	return urls, nil
}

func processOpenAIDataURLsMessages(ctx context.Context, t *TwilioHandler, text string, dataURLs []string, conv db.AiConversation) (string, error) {
	var msgs []openai.ChatCompletionMessage

	if text != "" {
		if _, err := t.app.DB.CreateMessage(ctx, text, request.RoleUser, int32(conv.ID)); err != nil {
			log.Printf("Error guardando texto: %v", err)
		}
		msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: text})
	}

	for _, url := range dataURLs {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{Type: openai.ChatMessagePartTypeText, Text: prompts.PROMPT_STREAM_PAGOS_BOT},
				{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{URL: url, Detail: "low"}},
			},
		})
	}

	response, err := t.app.OpenAI.ChatCompletions(msgs)
	if err != nil {
		log.Printf("OpenAI error: %v", err)
		return "", fmt.Errorf("error de OpenAI: %w", err)
	}

	result := formatResponse(response.Choices[0].Message.Content)
	if _, err := t.app.DB.CreateMessage(ctx, result, request.RoleSystem, int32(conv.ID)); err != nil {
		log.Printf("Error guardando respuesta: %v", err)
		return "", fmt.Errorf("error al guardar mensaje de sistema: %w", err)
	}
	return result, nil
}
