package push

import "github.com/gogf/gf/v2/net/ghttp"

// PushBody 处理 GET /:key/:body —— Bark 兼容：仅正文推送。
func (c *Controller) PushBody(r *ghttp.Request) {
	key := r.Get("key").String()
	body := r.Get("body").String()
	push(r, c.store, c.providers, key, "", body, "")
}
