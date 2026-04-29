package notifier

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SendEmail dispatches an email via Gmail's SMTP server using the provided app password.
func SendEmail(sender, password, recipient, title, articleURL string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: 🇮🇹 Alerta Ciudadania!\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	body := buildEmailBody(title, articleURL)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", sender, password, smtpHost)

	address := smtpHost + ":" + smtpPort
	err := smtp.SendMail(address, auth, sender, []string{recipient}, message)
	if err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}

	return nil
}

func buildEmailBody(title, articleURL string) string {
	var sb strings.Builder
	sb.WriteString("¡Hola!\r\n\r\n")
	sb.WriteString("Se detectó un nuevo aviso en el Consulado Italiano de Córdoba:\r\n\r\n")
	sb.WriteString("📋 ")
	sb.WriteString(title)
	sb.WriteString("\r\n🔗 ")
	sb.WriteString(articleURL)
	sb.WriteString("\r\n\r\n")
	sb.WriteString("— cittadinanza-bot")
	return sb.String()
}
