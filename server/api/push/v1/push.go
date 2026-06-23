// Package v1 defines the request/response types for push API v1.
package v1

// PushBodyReq is GET /:key/:body (Bark-compatible).
type PushBodyReq struct {
	Key  string `v:"required" in:"path" dc:"device key"`
	Body string `v:"required" in:"path" dc:"notification body"`
}

// PushTitleBodyReq is GET /:key/:title/:body (Bark-compatible).
type PushTitleBodyReq struct {
	Key   string `v:"required" in:"path" dc:"device key"`
	Title string `v:"required" in:"path" dc:"notification title"`
	Body  string `v:"required" in:"path" dc:"notification body"`
	URL   string `in:"query" dc:"URL to open on tap"`
}

// PushPostReq is POST /push.
type PushPostReq struct {
	Key   string `v:"required" json:"key" dc:"device key"`
	Title string `json:"title" dc:"notification title"`
	Body  string `v:"required" json:"body" dc:"notification body"`
	URL   string `json:"url" dc:"URL to open on tap"`
}

// PushRes is the response for push delivery.
type PushRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RegisterDeviceReq is POST /register.
type RegisterDeviceReq struct {
	Platform       string `v:"required|in:android,ios" json:"platform" dc:"device platform"`
	PushProvider   string `v:"required|in:jpush,getui,tpns,bark" json:"push_provider" dc:"push provider name"`
	RegistrationID string `v:"required" json:"registration_id" dc:"device registration ID from push SDK"`
}

// RegisterDeviceRes is the response for device registration.
type RegisterDeviceRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *RegisterDeviceResData   `json:"data,omitempty"`
}

// RegisterDeviceResData contains the device credentials returned on registration.
type RegisterDeviceResData struct {
	ID   string `json:"id" dc:"device internal ID"`
	Key  string `json:"key" dc:"Bark-compatible device key"`
	Note string `json:"note" dc:"usage hint"`
}
