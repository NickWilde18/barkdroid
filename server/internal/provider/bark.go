package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BarkConfig holds the configuration for iOS Bark forwarding.
type BarkConfig struct {
	BaseURL string // e.g. https://api.day.app
	Key     string // Bark device key
}

// BarkProvider forwards push messages to an iOS Bark server.
type BarkProvider struct {
	cfg    BarkConfig
	client *http.Client
}

// NewBark creates a BarkProvider.
func NewBark(cfg BarkConfig) *BarkProvider {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return &BarkProvider{
		cfg:    BarkConfig{BaseURL: baseURL, Key: cfg.Key},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns "bark".
func (p *BarkProvider) Name() string { return "bark" }

// Push forwards a message to a Bark server.
// Uses the Bark URL scheme: /{key}/{title}/{body}?url=...
func (p *BarkProvider) Push(msg *PushMessage) error {
	// Build Bark URL: /{key}/{title}/{body}
	title := url.PathEscape(msg.Title)
	body := url.PathEscape(msg.Body)

	u := fmt.Sprintf("%s/%s/%s/%s", p.cfg.BaseURL, p.cfg.Key, title, body)
	if msg.URL != "" {
		u += "?url=" + url.QueryEscape(msg.URL)
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("create bark request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("bark http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark returned HTTP %d", resp.StatusCode)
	}

	return nil
}
