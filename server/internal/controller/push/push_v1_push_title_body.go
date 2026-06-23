package push

import "github.com/gogf/gf/v2/net/ghttp"

// PushTitleBody 处理 GET /:key/:title/:body —— Bark 兼容：带标题推送。
func (c *Controller) PushTitleBody(r *ghttp.Request) {
	key := r.Get("key").String()
	title := r.Get("title").String()
	body := r.Get("body").String()
	url := r.Get("url").String()
	push(r, c.store, c.providers, key, title, body, url)
}
