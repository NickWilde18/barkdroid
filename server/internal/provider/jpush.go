package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// JPushConfig holds the credentials for JPush.
type JPushConfig struct {
	AppKey       string
	MasterSecret string
}

// JPushProvider delivers notifications via JPush REST API.
type JPushProvider struct {
	cfg    JPushConfig
	client *http.Client
}

// NewJPush creates a JPushProvider.
func NewJPush(cfg JPushConfig) *JPushProvider {
	return &JPushProvider{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns "jpush".
func (p *JPushProvider) Name() string { return "jpush" }

// Push sends a push notification via JPush v3 Push API.
// Docs: https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push
func (p *JPushProvider) Push(msg *PushMessage) error {
	payload := map[string]interface{}{
		"platform": "android",
		"audience": map[string]interface{}{
			"registration_id": []string{msg.DeviceID},
		},
		"notification": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":  msg.Body,
				"title":  msg.Title,
				"extras": jpushExtras(msg.URL),
			},
		},
		"options": map[string]interface{}{
			"apns_production": true,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal jpush payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.jpush.cn/v3/push", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create jpush request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+p.basicAuth())

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("jpush http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("jpush returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Read and check the JPush response for sendno/msg_id
	var jresp jpushResponse
	if err := json.NewDecoder(resp.Body).Decode(&jresp); err != nil {
		// Non-fatal: push was accepted (HTTP 200), just couldn't parse response
		return nil
	}
	if jresp.Error.Code != 0 {
		return fmt.Errorf("jpush error: code=%d msg=%s", jresp.Error.Code, jresp.Error.Message)
	}

	return nil
}

func (p *JPushProvider) basicAuth() string {
	auth := p.cfg.AppKey + ":" + p.cfg.MasterSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func jpushExtras(url string) map[string]interface{} {
	extras := make(map[string]interface{})
	if url != "" {
		extras["url"] = url
	}
	return extras
}

type jpushResponse struct {
	SendNo string `json:"sendno"`
	MsgID  string `json:"msg_id"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
