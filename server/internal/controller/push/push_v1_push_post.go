package push

import (
	"barkdroid/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

// PushPost 处理 POST /push —— JSON 格式推送。
func (c *Controller) PushPost(r *ghttp.Request) {
	var input model.PushInput
	if err := r.Parse(&input); err != nil {
		r.Response.WriteJsonExit(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "bad request: " + err.Error(),
		})
		return
	}
	push(r, c.store, c.providers, input.Key, input.Title, input.Body, input.URL)
}
