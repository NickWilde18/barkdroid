package store

import "barkdroid/internal/model"

// Store is the persistence layer for device registration.
type Store interface {
	RegisterDevice(platform, pushProvider, registrationID string) (*model.Device, error)
	GetDeviceByKey(key string) (*model.Device, error)
	ListDevices() ([]*model.Device, error)
	Close() error
}
