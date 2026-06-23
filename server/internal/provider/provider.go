package provider

// PushMessage represents a notification to be delivered.
type PushMessage struct {
	Title    string // notification title
	Body     string // notification body
	URL      string // optional URL to open on tap
	DeviceID string // registration ID, device token, or Bark key
}

// Provider abstracts a push delivery backend (JPush, Getui, Bark, etc.).
type Provider interface {
	Name() string
	Push(msg *PushMessage) error
}
