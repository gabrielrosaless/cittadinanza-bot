package parser

import (
	"cittadinanza-bot/internal/storage"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// ParseArticles extracts all news articles (title + URL) from the consulate's
// news page HTML. Each article is identified by an <h5> tag containing an <a>.
func ParseArticles(rawHTML string) ([]storage.Article, error) {
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var articles []storage.Article
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if isH5(n) {
			if a := findAnchor(n); a != nil {
				title := strings.TrimSpace(extractText(a))
				href := getAttr(a, "href")
				if title != "" && href != "" {
					articles = append(articles, storage.Article{
						Title: title,
						URL:   normalizeURL(href),
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return articles, nil
}

// isH5 returns true if the node is an <h5> element.
func isH5(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "h5"
}

// findAnchor returns the first <a> element found within n, or nil.
func findAnchor(n *html.Node) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			return c
		}
		if found := findAnchor(c); found != nil {
			return found
		}
	}
	return nil
}

// extractText returns the concatenated visible text of all descendant text nodes.
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractText(c))
	}
	return sb.String()
}

// getAttr returns the value of the given attribute from an element node.
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// normalizeURL ensures the URL is absolute. The consulate site uses absolute
// URLs in article links, but this guard handles any relative ones just in case.
func normalizeURL(href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	const base = "https://conscordoba.esteri.it"
	return base + href
}
