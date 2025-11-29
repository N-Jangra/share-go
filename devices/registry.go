package devices

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"time"
)

// Device represents a receiving device registered by a user.
type Device struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RegisteredAt time.Time `json:"registeredAt"`
	LastSeen     time.Time `json:"lastSeen"`
}

// PendingTransfer represents a transfer that should be delivered to a device.
type PendingTransfer struct {
	TransferID string    `json:"transferId"`
	Token      string    `json:"token"`
	Files      []PendingFile `json:"files"`
	SentAt     time.Time `json:"sentAt"`
}

// PendingFile is a lightweight descriptor for a file being shared via devices API.
type PendingFile struct {
	Name string `json:"name"`
	Mime string `json:"mime"`
	Size int64  `json:"size"`
}

var (
	// ErrDeviceNotFound indicates an unknown device id.
	ErrDeviceNotFound = errors.New("device not found")
)

type deviceState struct {
	info    *Device
	pending *PendingTransfer
}

// Registry keeps track of registered devices and pending transfers.
type Registry struct {
	mu       sync.RWMutex
	devices  map[string]*deviceState
	maxCount int
}

// NewRegistry creates an empty registry.
func NewRegistry(maxDevices int) *Registry {
	return &Registry{
		devices:  make(map[string]*deviceState),
		maxCount: maxDevices,
	}
}

// Register inserts a new device with the provided friendly name.
func (r *Registry) Register(name string) (*Device, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("device name cannot be empty")
	}
	if len(name) > 40 {
		name = name[:40]
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.registerLocked(name)
}

func (r *Registry) registerLocked(name string) (*Device, error) {
	if r.maxCount > 0 && len(r.devices) >= r.maxCount {
		return nil, errors.New("device registry is full")
	}

	now := time.Now().UTC()
	id := randomString(12)
	device := &Device{
		ID:           id,
		Name:         name,
		RegisteredAt: now,
		LastSeen:     now,
	}
	r.devices[id] = &deviceState{info: device}
	return device, nil
}

// Update mutates an existing device name and touch timestamp.
func (r *Registry) Update(id, name string) (*Device, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("device name cannot be empty")
	}
	if len(name) > 40 {
		name = name[:40]
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.devices[id]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	state.info.Name = name
	state.info.LastSeen = time.Now().UTC()
	return state.info, nil
}

// Upsert updates a device if the id exists, otherwise registers a new device.
func (r *Registry) Upsert(id, name string) (*Device, error) {
	if id != "" {
		device, err := r.Update(id, name)
		if err == nil || !errors.Is(err, ErrDeviceNotFound) {
			return device, err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	return r.registerLocked(name)
}

// List returns a copy of all known devices.
func (r *Registry) List() []*Device {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Device, 0, len(r.devices))
	for _, state := range r.devices {
		deviceCopy := *state.info
		out = append(out, &deviceCopy)
	}
	return out
}

// Notify assigns a transfer to the given device.
func (r *Registry) Notify(deviceID string, pending *PendingTransfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.devices[deviceID]
	if !ok {
		return ErrDeviceNotFound
	}
	state.pending = pending
	state.info.LastSeen = time.Now().UTC()
	return nil
}

// Pending returns the pending transfer for the device, if any.
func (r *Registry) Pending(deviceID string) (*PendingTransfer, error) {
	r.mu.RLock()
	state, ok := r.devices[deviceID]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrDeviceNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	state.info.LastSeen = time.Now().UTC()
	if state.pending == nil {
		return nil, nil
	}
	copy := *state.pending
	return &copy, nil
}

// Clear removes the pending transfer for the device if ids match.
func (r *Registry) Clear(deviceID, transferID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.devices[deviceID]
	if !ok || state.pending == nil {
		return
	}
	if transferID == "" || state.pending.TransferID == transferID {
		state.pending = nil
	}
}

// ClearByTransfer removes any pending notifications that reference the transfer id.
func (r *Registry) ClearByTransfer(transferID string) {
	if transferID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, state := range r.devices {
		if state.pending != nil && state.pending.TransferID == transferID {
			state.pending = nil
		}
	}
}

func randomString(length int) string {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)[:length]
}
