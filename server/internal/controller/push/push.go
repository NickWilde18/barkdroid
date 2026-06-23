package push

import (
	"log"

	"barkdroid/internal/model"
	"barkdroid/internal/provider"
	"barkdroid/internal/store"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Controller handles Bark-compatible HTTP requests.
type Controller struct {
	store     store.Store
	providers map[string]provider.Provider
}

// New creates a Controller.
func New(st store.Store, providers map[string]provider.Provider) *Controller {
	return &Controller{
		store:     st,
		providers: providers,
	}
}

func (c *Controller) PushBody(r *ghttp.Request) {
	key := r.Get("key").String()
	body := r.Get("body").String()
	c.push(r, key, "", body, "")
}

// PushTitleBody handles GET /:key/:title/:body
func (c *Controller) PushTitleBody(r *ghttp.Request) {
	key := r.Get("key").String()
	title := r.Get("title").String()
	body := r.Get("body").String()
	url := r.Get("url").String()
	c.push(r, key, title, body, url)
}

// PushPost handles POST /push
func (c *Controller) PushPost(r *ghttp.Request) {
	var input model.PushInput
	if err := r.Parse(&input); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "bad request: " + err.Error(),
		})
		return
	}
	c.push(r, input.Key, input.Title, input.Body, input.URL)
}

// RegisterDevice handles POST /register
func (c *Controller) RegisterDevice(r *ghttp.Request) {
	var input model.RegisterInput
	if err := r.Parse(&input); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "bad request: " + err.Error(),
		})
		return
	}

	dev, err := c.store.RegisterDevice(input.Platform, input.PushProvider, input.RegistrationID)
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

func (c *Controller) push(r *ghttp.Request, key, title, body, url string) {
	if key == "" || body == "" {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "key and body are required",
		})
		return
	}

	dev, err := c.store.GetDeviceByKey(key)
	if err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	prov, ok := c.providers[dev.PushProvider]
	if !ok {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "provider '" + dev.PushProvider + "' not configured on server",
		})
		return
	}

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
