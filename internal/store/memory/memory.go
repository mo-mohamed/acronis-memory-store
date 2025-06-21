// Package memory implements an in memory key value store. It also supports keys cleanup based on TTL.
package memory

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrTypeMismatch  = errors.New("operation not supported for this data type")
	ErrInvalidTTL    = errors.New("invalid TTL value")
	ErrEmptyList     = errors.New("list is empty")
	ErrMarshalFailed = errors.New("failed to marshal value to JSON")
)

type MemoryStore struct {
	mu        sync.RWMutex
	data      map[string]Value
	ttlCtx    context.Context
	ttlCancel context.CancelFunc
}

// NewMemoryStore initializes a new in memory store.
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		data:      make(map[string]Value),
		ttlCtx:    nil,
		ttlCancel: nil,
	}

	// Start the bakground worker to clean expired keys
	s.ttlCtx, s.ttlCancel = context.WithCancel(context.Background())
	s.doStartTTLWorker()

	return s
}

// StartTTLWorker starts a background worker to clean up keys that are expired
func (s *MemoryStore) StartTTLWorker(ctx context.Context) {
	s.mu.Lock()
	if s.ttlCancel != nil {
		s.ttlCancel()
	}

	// Create a cancel context.
	s.ttlCtx, s.ttlCancel = context.WithCancel(context.Background())
	s.mu.Unlock()

	s.doStartTTLWorker()
}

// StopTTLWorker stops the clean up worker from removing expired keys
func (s *MemoryStore) StopTTLWorker() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ttlCancel != nil {
		s.ttlCancel()
		s.ttlCancel = nil
		s.ttlCtx = nil
	}
}

// Set sets a key with a value and optional ttl (0 = no TTL)
func (s *MemoryStore) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	if ttlSeconds < 0 {
		return ErrInvalidTTL
	}

	stringValue, err := s.Stringify(value)
	if err != nil {
		return ErrMarshalFailed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var ttl time.Time
	if ttlSeconds > 0 {
		ttl = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}
	// If ttlSeconds == 0, ttl remains zero (no expiration)

	s.data[key] = Value{Val: stringValue, TTL: ttl, IsList: false}
	return nil
}

// Get gets a value from the store
func (s *MemoryStore) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()

	v, ok := s.data[key]
	if !ok {
		s.mu.RUnlock()
		return "", ErrKeyNotFound
	}

	// key exists and is not expired or doesn't have a TTL, return the value
	if v.TTL.IsZero() || time.Now().Before(v.TTL) {
		if v.IsList {
			s.mu.RUnlock()
			return "", ErrTypeMismatch
		}
		result := v.Val
		s.mu.RUnlock()
		return result, nil
	}

	// key is expired, lazy delete it
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok = s.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	if !v.TTL.IsZero() && time.Now().After(v.TTL) {
		delete(s.data, key)
		return "", ErrKeyNotFound
	}

	if v.IsList {
		return "", ErrTypeMismatch
	}

	return v.Val, nil
}

// Update updates a value in the store
func (s *MemoryStore) Update(ctx context.Context, key string, value any) error {
	stringValue, err := s.Stringify(value)
	if err != nil {
		return ErrMarshalFailed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.data[key]
	if !exists {
		return ErrKeyNotFound
	}

	// If expired, delete it and return key not found
	if !v.TTL.IsZero() && time.Now().After(v.TTL) {
		delete(s.data, key)
		return ErrKeyNotFound
	}

	if v.IsList {
		return ErrTypeMismatch
	}

	v.Val = stringValue
	s.data[key] = v
	return nil
}

// Remove deletes a key from the store
func (s *MemoryStore) Remove(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return ErrKeyNotFound
	}

	delete(s.data, key)
	return nil
}

// List operations
func (s *MemoryStore) Push(ctx context.Context, key string, item any) error {
	stringItem, err := s.Stringify(item)
	if err != nil {
		return ErrMarshalFailed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	v := s.data[key]

	// If the key doesn't exist, or exists but expired, then create a new list
	if _, exists := s.data[key]; !exists || (!v.TTL.IsZero() && time.Now().After(v.TTL)) {
		// If key exists but is expired, lazy delete it first
		if _, exists := s.data[key]; exists && (!v.TTL.IsZero() && time.Now().After(v.TTL)) {
			delete(s.data, key)
		}
		v = Value{IsList: true, List: []string{}}
	}

	if !v.IsList {
		return ErrTypeMismatch
	}

	v.List = append([]string{stringItem}, v.List...)
	s.data[key] = v
	return nil
}

// Pop takes a value from the list
func (s *MemoryStore) Pop(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.data[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	// If expired, lazy delete it
	if !v.TTL.IsZero() && time.Now().After(v.TTL) {
		delete(s.data, key)
		return "", ErrKeyNotFound
	}

	if !v.IsList {
		return "", ErrTypeMismatch
	}

	if len(v.List) == 0 {
		return "", ErrEmptyList
	}

	item := v.List[0]
	v.List = v.List[1:]
	s.data[key] = v
	return item, nil
}

// Stringify converts any value to string
func (s *MemoryStore) Stringify(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

// doStartTTLWorker starts the actual TTL cleanup worker
func (s *MemoryStore) doStartTTLWorker() {
	s.mu.RLock()
	ctx := s.ttlCtx
	s.mu.RUnlock()

	if ctx == nil {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				for k, v := range s.data {
					select {
					case <-ctx.Done(): // Here we check again to avoid race condition, giving a chance to listen to the context done signal
						s.mu.Unlock()
						return
					default:
						if !v.TTL.IsZero() && time.Now().After(v.TTL) {
							delete(s.data, k)
						}
					}
				}
				s.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}
