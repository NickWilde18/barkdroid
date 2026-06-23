package push

import (
	"log"

	"barkdroid/internal/provider"
	"barkdroid/internal/store"

	"github.com/gogf/gf/v2/net/ghttp"
)

// push 是私有 helper，处理推送核心逻辑。
func push(r *ghttp.Request, st store.Store, providers map[string]provider.Provider, key, title, body, url string) {
	if key == "" || body == "" {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "key and body are required",
		})
		return
	}

	dev, err := st.GetDeviceByKey(key)
	if err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	prov, ok := providers[dev.PushProvider]
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
