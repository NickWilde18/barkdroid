package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"barkdroid/internal/model"

	_ "modernc.org/sqlite"
)

// SQLiteStore implements Store backed by a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) the SQLite database and ensures the schema exists.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Enable WAL mode for better concurrent reads
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable wal: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS devices (
			id              TEXT PRIMARY KEY,
			key             TEXT UNIQUE NOT NULL,
			platform        TEXT NOT NULL DEFAULT 'android',
			push_provider   TEXT NOT NULL,
			registration_id TEXT NOT NULL,
			created_at      DATETIME NOT NULL DEFAULT (datetime('now')),
			updated_at      DATETIME NOT NULL DEFAULT (datetime('now'))
		);
		CREATE INDEX IF NOT EXISTS idx_devices_key ON devices(key);
	`)
	return err
}

// RegisterDevice inserts or upserts a device and returns it with a generated key.
func (s *SQLiteStore) RegisterDevice(platform, pushProvider, registrationID string) (*model.Device, error) {
	// Check if this registration_id already exists (upsert)
	existing, err := s.deviceByRegistration(registrationID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Update timestamp and return existing key
		_, err = s.db.Exec(
			`UPDATE devices SET updated_at = datetime('now') WHERE id = ?`,
			existing.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("update device: %w", err)
		}
		existing.UpdatedAt = time.Now()
		return existing, nil
	}

	// Generate a new device with a random key
	id := newID()
	key := newKey()

	now := time.Now()
	_, err = s.db.Exec(
		`INSERT INTO devices (id, key, platform, push_provider, registration_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, key, platform, pushProvider, registrationID, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}

	return &model.Device{
		ID:             id,
		Key:            key,
		Platform:       platform,
		PushProvider:   pushProvider,
		RegistrationID: registrationID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// GetDeviceByKey looks up a device by its Bark-compatible key.
func (s *SQLiteStore) GetDeviceByKey(key string) (*model.Device, error) {
	dev, err := s.scanDevice(`SELECT id, key, platform, push_provider, registration_id, created_at, updated_at FROM devices WHERE key = ?`, key)
	if err != nil {
		return nil, err
	}
	if dev == nil {
		return nil, fmt.Errorf("device not found: %s", key)
	}
	return dev, nil
}

// ListDevices returns all registered devices, newest first.
func (s *SQLiteStore) ListDevices() ([]*model.Device, error) {
	rows, err := s.db.Query(
		`SELECT id, key, platform, push_provider, registration_id, created_at, updated_at
		 FROM devices ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()

	var devices []*model.Device
	for rows.Next() {
		dev, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}
	return devices, rows.Err()
}

// Close closes the database.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) deviceByRegistration(regID string) (*model.Device, error) {
	return s.scanDevice(
		`SELECT id, key, platform, push_provider, registration_id, created_at, updated_at FROM devices WHERE registration_id = ?`,
		regID,
	)
}

func (s *SQLiteStore) scanDevice(query string, args ...interface{}) (*model.Device, error) {
	row := s.db.QueryRow(query, args...)
	return scanRow(row)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanRow(sc scanner) (*model.Device, error) {
	var d model.Device
	err := sc.Scan(&d.ID, &d.Key, &d.Platform, &d.PushProvider, &d.RegistrationID, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
