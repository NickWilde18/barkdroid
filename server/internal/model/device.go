package model

import "time"

// Device represents a registered push device.
type Device struct {
	ID             string    `json:"id"`
	Key            string    `json:"key"`             // Bark-compatible key (short unique identifier)
	Platform       string    `json:"platform"`        // "android" or "ios"
	PushProvider   string    `json:"push_provider"`   // "jpush", "getui", "tpns", "bark"
	RegistrationID string    `json:"registration_id"` // JPush registrationId or APNs device token
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RegisterInput is the payload for device registration.
type RegisterInput struct {
	Platform       string `json:"platform" v:"required|in:android,ios"`
	PushProvider   string `json:"push_provider" v:"required|in:jpush,getui,tpns,bark"`
	RegistrationID string `json:"registration_id" v:"required"`
}

// PushInput is the payload for POST /push.
type PushInput struct {
	Key   string `json:"key" v:"required"`
	Title string `json:"title"`
	Body  string `json:"body" v:"required"`
	URL   string `json:"url"`
}
