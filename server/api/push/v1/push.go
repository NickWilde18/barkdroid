package v1

import "github.com/gogf/gf/v2/frame/g"

// PushBodyReq GET /:key/:body
type PushBodyReq struct {
	g.Meta `path:"/:key/:body" method:"get" summary:"推送通知（仅正文）"`
	Key    string `v:"required" in:"path" dc:"设备密钥"`
	Body   string `v:"required" in:"path" dc:"通知正文"`
}

// PushTitleBodyReq GET /:key/:title/:body
type PushTitleBodyReq struct {
	g.Meta `path:"/:key/:title/:body" method:"get" summary:"推送通知（含标题）"`
	Key    string `v:"required" in:"path" dc:"设备密钥"`
	Title  string `v:"required" in:"path" dc:"通知标题"`
	Body   string `v:"required" in:"path" dc:"通知正文"`
	URL    string `in:"query" dc:"点击通知打开的 URL"`
}

// PushPostReq POST /push
type PushPostReq struct {
	g.Meta `path:"/push" method:"post" summary:"推送通知（JSON）"`
	Key    string `v:"required" json:"key" dc:"设备密钥"`
	Title  string `json:"title" dc:"通知标题"`
	Body   string `v:"required" json:"body" dc:"通知正文"`
	URL    string `json:"url" dc:"点击通知打开的 URL"`
}

// PushRes 推送响应
type PushRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RegisterDeviceReq POST /register
type RegisterDeviceReq struct {
	g.Meta         `path:"/register" method:"post" summary:"注册设备"`
	Platform       string `v:"required|in:android,ios" json:"platform" dc:"设备平台"`
	PushProvider   string `v:"required|in:jpush,getui,tpns,bark" json:"push_provider" dc:"推送服务商"`
	RegistrationID string `v:"required" json:"registration_id" dc:"推送 SDK 返回的注册 ID"`
}

// RegisterDeviceRes 注册响应
type RegisterDeviceRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *RegisterDeviceResData `json:"data,omitempty"`
}

// RegisterDeviceResData 注册返回的密钥信息
type RegisterDeviceResData struct {
	ID   string `json:"id" dc:"设备内部 ID"`
	Key  string `json:"key" dc:"Bark 兼容密钥"`
	Note string `json:"note" dc:"使用提示"`
}
