package main

import (
	"cittadinanza-bot/internal/config"
	"cittadinanza-bot/internal/monitor"
	"cittadinanza-bot/internal/storage"
	"flag"
	"log"
	"os"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.Printf("Starting cittadinanza-bot...")

	// 1. Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config. Interval: %d minutes", cfg.CheckIntervalMinutes)

	// 2. Open Storage
	db, err := storage.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open storage: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to SQLite database at %s", cfg.DBPath)

	// 3. Run a single monitoring loop and exit (GitHub Actions approach)
	log.Printf("Executing single monitoring cycle...")
	monitor.Run(cfg, db)
	log.Printf("Run finished successfully.")
}
