package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const telegramAPIBase = "https://api.telegram.org/bot%s/sendMessage"

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// SendTelegram posts a Telegram message to the given chat using the bot token.
func SendTelegram(token, chatID, title, articleURL string) error {
	message := buildMessage(title, articleURL)

	payload := telegramPayload{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling telegram payload: %w", err)
	}

	url := fmt.Sprintf(telegramAPIBase, token)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sending telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// buildMessage formats the Telegram notification text.
func buildMessage(title, articleURL string) string {
	return fmt.Sprintf(
		"🇮🇹 *Alerta Ciudadanía*\n\nSe detectó un nuevo aviso en el Consulado Italiano de Córdoba:\n\n📋 %s\n🔗 %s",
		escapeMarkdown(title),
		articleURL,
	)
}

// escapeMarkdown escapes special Markdown characters in text to avoid
// Telegram parse errors when titles contain symbols like – or ().
func escapeMarkdown(s string) string {
	replacer := []struct{ old, new string }{
		{"_", "\\_"},
		{"*", "\\*"},
		{"[", "\\["},
		{"`", "\\`"},
	}
	for _, r := range replacer {
		s = replaceAll(s, r.old, r.new)
	}
	return s
}

func replaceAll(s, old, new string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result = append(result, new...)
			i += len(old)
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}
