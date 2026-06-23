package handler

import (
	"log"

	"barkdroid/internal/model"
	"barkdroid/internal/provider"
	"barkdroid/internal/store"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Handler handles Bark-compatible HTTP requests.
type Handler struct {
	store     store.Store
	providers map[string]provider.Provider // provider name → implementation
}

// New creates a Handler.
func New(st store.Store, providers map[string]provider.Provider) *Handler {
	return &Handler{
		store:     st,
		providers: providers,
	}
}

// PushBody handles GET /:key/:body
// Compatible with: curl https://api.day.app/{key}/{body}
func (h *Handler) PushBody(r *ghttp.Request) {
	key := r.Get("key").String()
	body := r.Get("body").String()
	h.push(r, key, "", body, "")
}

// PushTitleBody handles GET /:key/:title/:body
// Compatible with: curl https://api.day.app/{key}/{title}/{body}
func (h *Handler) PushTitleBody(r *ghttp.Request) {
	key := r.Get("key").String()
	title := r.Get("title").String()
	body := r.Get("body").String()
	url := r.Get("url").String() // optional query param
	h.push(r, key, title, body, url)
}

// PushPost handles POST /push
// JSON body: {"key": "...", "title": "...", "body": "...", "url": "..."}
func (h *Handler) PushPost(r *ghttp.Request) {
	var input model.PushInput
	if err := r.Parse(&input); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "bad request: " + err.Error(),
		})
		return
	}
	h.push(r, input.Key, input.Title, input.Body, input.URL)
}

// RegisterDevice handles POST /register
// JSON body: {"platform": "android", "push_provider": "jpush", "registration_id": "..."}
func (h *Handler) RegisterDevice(r *ghttp.Request) {
	var input model.RegisterInput
	if err := r.Parse(&input); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "bad request: " + err.Error(),
		})
		return
	}

	dev, err := h.store.RegisterDevice(input.Platform, input.PushProvider, input.RegistrationID)
	if err != nil {
		log.Printf("[ERROR] register device: %v", err)
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "register failed",
		})
		return
	}

	r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "ok",
		Data: map[string]interface{}{
			"key":  dev.Key,
			"id":   dev.ID,
			"note": "Use this key to send push: GET /{key}/{title}/{body}",
		},
	})
}

// push is the shared push logic for all Bark-compatible endpoints.
func (h *Handler) push(r *ghttp.Request, key, title, body, url string) {
	if key == "" || body == "" {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "key and body are required",
		})
		return
	}

	// Look up the device by its Bark-compatible key
	dev, err := h.store.GetDeviceByKey(key)
	if err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	// Find the configured provider for this device
	prov, ok := h.providers[dev.PushProvider]
	if !ok {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "provider '" + dev.PushProvider + "' not configured on server",
		})
		return
	}

	// Deliver the push
	if err := prov.Push(&provider.PushMessage{
		Title:    title,
		Body:     body,
		URL:      url,
		DeviceID: dev.RegistrationID,
	}); err != nil {
		log.Printf("[ERROR] push via %s failed: %v", dev.PushProvider, err)
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "push delivery failed: " + err.Error(),
		})
		return
	}

	r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "ok",
	})
}
