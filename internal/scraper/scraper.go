package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxRetries    = 3
	retryWait     = 5 * time.Second
	requestTimeout = 15 * time.Second
	userAgent     = "Mozilla/5.0 (compatible; cittadinanza-bot/1.0)"
)

// FetchHTML downloads the HTML content of the given URL.
// It retries up to maxRetries times on failure, waiting retryWait between attempts.
func FetchHTML(url string) (string, error) {
	client := &http.Client{Timeout: requestTimeout}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		html, err := doRequest(client, url)
		if err == nil {
			return html, nil
		}
		lastErr = err
		if attempt < maxRetries {
			time.Sleep(retryWait)
		}
	}

	return "", fmt.Errorf("fetching %q after %d attempts: %w", url, maxRetries, lastErr)
}

func doRequest(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "es-AR,es;q=0.9,it;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	return string(body), nil
}
