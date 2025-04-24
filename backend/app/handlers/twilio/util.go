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

	"github.com/cohesion-org/deepseek-go"
	"github.com/luizarnoldch/stream_bot/db"
	"github.com/luizarnoldch/stream_bot/integrations/prompts"
	"github.com/openai/openai-go"
)

func formatResponse(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	return strings.ReplaceAll(text, "\r\n", " ")
}

func formatJSONResponse(response string) string {
	cleaned := strings.ReplaceAll(response, "`", "")
	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(cleaned[start : end+1])
	}
	return strings.TrimSpace(cleaned)
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

	if _, err = t.app.DB.CreateMessage(ctx, prompts.PROMPT_STREAM_VENTAS_BOT, deepseek.ChatMessageRoleSystem, newConversation.ID); err != nil {
		log.Printf("Error creating system message: %v", err)
	}

	return newUser, newConversation, nil
}

func buildDeepSeekInput(messages []db.AiMessage) []deepseek.ChatCompletionMessage {
	input := make([]deepseek.ChatCompletionMessage, 0, len(messages))
	for _, msg := range messages {
		input = append(input, deepseek.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return input
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
	var msgs []openai.ChatCompletionMessageParamUnion

	if text != "" {
		if _, err := t.app.DB.CreateMessage(ctx, text, deepseek.ChatMessageRoleUser, int32(conv.ID)); err != nil {
			log.Printf("Error guardando texto: %v", err)
		}
		msgs = append(msgs, openai.UserMessage(text))
	}

	for _, url := range dataURLs {
		parts := []openai.ChatCompletionContentPartUnionParam{
			openai.TextContentPart(prompts.PROMPT_STREAM_PAGOS_BOT),
			openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL:    url,
				Detail: string(openai.ImageFileDeltaDetailLow),
			}),
		}

		msgs = append(msgs, openai.ChatCompletionMessageParamUnion{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfArrayOfContentParts: parts,
				},
			},
		})
	}

	response, err := t.app.OpenAI.ChatCompletions(msgs)
	if err != nil {
		log.Printf("OpenAI error: %v", err)
		return "", fmt.Errorf("error de OpenAI: %w", err)
	}

	result := formatResponse(response.Choices[0].Message.Content)
	if _, err := t.app.DB.CreateMessage(ctx, result, deepseek.ChatMessageRoleSystem, int32(conv.ID)); err != nil {
		log.Printf("Error guardando respuesta: %v", err)
		return "", fmt.Errorf("error al guardar mensaje de sistema: %w", err)
	}
	return result, nil
}
