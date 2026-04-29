package monitor

import (
	"cittadinanza-bot/internal/config"
	"cittadinanza-bot/internal/detector"
	"cittadinanza-bot/internal/notifier"
	"cittadinanza-bot/internal/parser"
	"cittadinanza-bot/internal/scraper"
	"cittadinanza-bot/internal/storage"
	"fmt"
	"log"
	"time"
)

// Run executes a single monitoring cycle: fetch -> parse -> detect -> notify -> log
func Run(cfg *config.Config, db *storage.DB) {
	log.Printf("Starting check cycle for %s", cfg.TargetURL)

	result := storage.CheckResult{
		CheckedAt: time.Now(),
	}

	// 1. Fetch HTML
	rawHTML, err := scraper.FetchHTML(cfg.TargetURL)
	if err != nil {
		log.Printf("ERROR fetching HTML: %v", err)
		result.Error = fmt.Sprintf("fetch: %v", err)
		_ = db.LogCheck(result)
		return
	}

	// 2. Parse Articles
	articles, err := parser.ParseArticles(rawHTML)
	if err != nil {
		log.Printf("ERROR parsing HTML: %v", err)
		result.Error = fmt.Sprintf("parse: %v", err)
		_ = db.LogCheck(result)
		return
	}
	result.ArticlesFound = len(articles)

	// 3. Process Articles
	for _, article := range articles {
		// Dedup by URL
		isNew, err := db.IsNew(article.URL)
		if err != nil {
			log.Printf("ERROR checking if article is new (%s): %v", article.URL, err)
			continue
		}

		if !isNew {
			continue // Already seen this URL
		}

		result.NewArticles++
		notified := false

		// 4. Keyword Detection
		if detector.HasKeyword(article.Title, cfg.Keywords) {
			log.Printf("KEYWORD MATCH! Title: %q URL: %q", article.Title, article.URL)
			
			// 5. Notify Channels
			if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
				err := notifier.SendTelegram(cfg.TelegramToken, cfg.TelegramChatID, article.Title, article.URL)
				if err != nil {
					log.Printf("ERROR sending Telegram notification: %v", err)
				} else {
					log.Printf("Notification sent successfully via Telegram")
					notified = true
				}
			}

			if cfg.EmailEnabled {
				err := notifier.SendEmail(cfg.EmailSender, cfg.EmailAppPassword, cfg.EmailRecipient, article.Title, article.URL)
				if err != nil {
					log.Printf("ERROR sending Email notification: %v", err)
				} else {
					log.Printf("Notification sent successfully via Email")
					notified = true
				}
			}

			if notified {
				result.AlertsSent++
			}
		} else {
			log.Printf("New article seen (no keyword match): %q", article.Title)
		}

		// 6. Mark as seen (whether matched or not)
		if err := db.MarkAsSeen(article, notified); err != nil {
			log.Printf("ERROR marking article as seen (%s): %v", article.URL, err)
		}
	}

	// 7. Log cycle
	if err := db.LogCheck(result); err != nil {
		log.Printf("ERROR logging check result: %v", err)
	}

	log.Printf("Check cycle completed. Found: %d, New: %d, Alerts Sent: %d", 
		result.ArticlesFound, result.NewArticles, result.AlertsSent)
}
