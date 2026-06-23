package store

import (
	"os"
	"testing"
)

func TestSQLiteStore_RegisterAndGet(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer s.Close()

	// Register a device
	dev, err := s.RegisterDevice("android", "jpush", "regid-abc123")
	if err != nil {
		t.Fatalf("RegisterDevice: %v", err)
	}
	if dev.Key == "" || dev.ID == "" {
		t.Errorf("expected key and id, got key=%q id=%q", dev.Key, dev.ID)
	}
	if len(dev.Key) != 8 {
		t.Errorf("expected 8-char key, got %d chars: %q", len(dev.Key), dev.Key)
	}
	if dev.Platform != "android" || dev.PushProvider != "jpush" || dev.RegistrationID != "regid-abc123" {
		t.Errorf("unexpected device fields: %+v", dev)
	}

	// Get by key
	got, err := s.GetDeviceByKey(dev.Key)
	if err != nil {
		t.Fatalf("GetDeviceByKey: %v", err)
	}
	if got.ID != dev.ID || got.Key != dev.Key {
		t.Errorf("GetDeviceByKey mismatch: want %+v, got %+v", dev, got)
	}

	// Get non-existent key
	_, err = s.GetDeviceByKey("nonexist")
	if err == nil {
		t.Error("expected error for non-existent key")
	}
}

func TestSQLiteStore_DuplicateRegistration(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer s.Close()

	// First registration
	dev1, err := s.RegisterDevice("android", "jpush", "regid-same")
	if err != nil {
		t.Fatalf("RegisterDevice 1: %v", err)
	}

	// Same registration_id should return the same device
	dev2, err := s.RegisterDevice("android", "jpush", "regid-same")
	if err != nil {
		t.Fatalf("RegisterDevice 2: %v", err)
	}

	if dev1.Key != dev2.Key {
		t.Errorf("duplicate registration should return same key: %q vs %q", dev1.Key, dev2.Key)
	}
	if dev1.ID != dev2.ID {
		t.Errorf("duplicate registration should return same id: %q vs %q", dev1.ID, dev2.ID)
	}
}

func TestSQLiteStore_ListDevices(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer s.Close()

	// Empty list
	devs, err := s.ListDevices()
	if err != nil {
		t.Fatalf("ListDevices empty: %v", err)
	}
	if len(devs) != 0 {
		t.Errorf("expected empty list, got %d devices", len(devs))
	}

	// Add two devices
	s.RegisterDevice("android", "jpush", "regid-1")
	s.RegisterDevice("ios", "bark", "token-xyz")

	devs, err = s.ListDevices()
	if err != nil {
		t.Fatalf("ListDevices after add: %v", err)
	}
	if len(devs) != 2 {
		t.Errorf("expected 2 devices, got %d", len(devs))
	}
}

func TestSQLiteStore_FilePersistence(t *testing.T) {
	path := "/tmp/barkdroid_test.db"
	os.Remove(path) // clean up from previous runs

	s, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	s.RegisterDevice("android", "jpush", "regid-persist")
	s.Close()

	// Reopen — data should survive
	s2, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore reopen: %v", err)
	}
	defer s2.Close()

	devs, err := s2.ListDevices()
	if err != nil {
		t.Fatalf("ListDevices after reopen: %v", err)
	}
	if len(devs) != 1 {
		t.Errorf("expected 1 device after reopen, got %d", len(devs))
	}

	os.Remove(path)
}
