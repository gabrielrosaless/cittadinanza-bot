package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all runtime configuration for the bot.
type Config struct {
	// Telegram overrides
	TelegramToken       string   `json:"telegram_token"`
	TelegramChatID      string   `json:"telegram_chat_id"`

	// Email fields
	EmailEnabled        bool     `json:"email_enabled"`
	EmailSender         string   `json:"email_sender"`
	EmailAppPassword    string   `json:"email_app_password"`
	EmailRecipient      string   `json:"email_recipient"`

	// General
	TargetURL           string   `json:"target_url"`
	CheckIntervalMinutes int     `json:"check_interval_minutes"`
	DBPath              string   `json:"db_path"`
	Keywords            []string `json:"keywords"`
}

// Load reads and validates the config from the given JSON file path.
func Load(path string) (*Config, error) {
	// Try to open the file but do not fail hard immediately. We might be in a
	// GitHub Actions environment without a config.json present.
	f, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("opening config file %q: %w", path, err)
	}
	if f != nil {
		defer f.Close()
	}

	var cfg Config
	// It's acceptable if the config file doesn't exist (e.g., in GitHub Actions)
	if f != nil {
		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, fmt.Errorf("decoding config JSON: %w", err)
		}
	}

	// Override from OS environment variables (GitHub Secrets)
	if env := os.Getenv("TELEGRAM_TOKEN"); env != "" {
		cfg.TelegramToken = env
	}
	if env := os.Getenv("TELEGRAM_CHAT_ID"); env != "" {
		cfg.TelegramChatID = env
	}
	if env := os.Getenv("EMAIL_ENABLED"); env == "true" {
		cfg.EmailEnabled = true
	}
	if env := os.Getenv("EMAIL_SENDER"); env != "" {
		cfg.EmailSender = env
	}
	if env := os.Getenv("EMAIL_APP_PASSWORD"); env != "" {
		cfg.EmailAppPassword = env
	}
	if env := os.Getenv("EMAIL_RECIPIENT"); env != "" {
		cfg.EmailRecipient = env
	}
	if env := os.Getenv("TARGET_URL"); env != "" {
		cfg.TargetURL = env
	}

	// Supply defaults if empty
	if cfg.TargetURL == "" {
		cfg.TargetURL = "https://conscordoba.esteri.it/es/news/"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "bot.db"
	}
	if len(cfg.Keywords) == 0 {
		cfg.Keywords = []string{"turnos", "cittadinanza", "prenotazioni", "apertura", "ciudadania"}
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	hasTelegram := (c.TelegramToken != "" && c.TelegramToken != "YOUR_BOT_TOKEN_HERE") &&
		(c.TelegramChatID != "" && c.TelegramChatID != "YOUR_CHAT_ID_HERE")

	hasEmail := c.EmailEnabled

	if !hasTelegram && !hasEmail {
		return fmt.Errorf("you must configure at least one notification method (Telegram or Email)")
	}

	if hasEmail {
		if c.EmailSender == "" || c.EmailSender == "your.email@gmail.com" {
			return fmt.Errorf("email_sender is required when email is enabled")
		}
		if c.EmailAppPassword == "" || c.EmailAppPassword == "abcd efgh ijkl mnop" {
			return fmt.Errorf("email_app_password is required when email is enabled")
		}
		if c.EmailRecipient == "" || c.EmailRecipient == "your.email@gmail.com" {
			return fmt.Errorf("email_recipient is required when email is enabled")
		}
	}

	if c.TargetURL == "" {
		return fmt.Errorf("target_url is required")
	}
	if c.CheckIntervalMinutes <= 0 {
		return fmt.Errorf("check_interval_minutes must be greater than 0")
	}
	if c.DBPath == "" {
		c.DBPath = "bot.db"
	}
	if len(c.Keywords) == 0 {
		return fmt.Errorf("at least one keyword is required")
	}
	return nil
}
