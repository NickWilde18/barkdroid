package push

import (
	"log"

	"barkdroid/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

// RegisterDevice 处理 POST /register —— 注册设备，返回 Bark 兼容密钥。
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
