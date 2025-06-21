package store

import "context"

// Store defines the interface for in-memory data structure store
type IStore interface {
	// String operations
	Set(ctx context.Context, key string, value any, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Update(ctx context.Context, key string, value any) error
	Remove(ctx context.Context, key string) error

	// List operations
	Push(ctx context.Context, key string, item any) error
	Pop(ctx context.Context, key string) (string, error)

	// TTL and utility
	StartTTLWorker(ctx context.Context)
	StopTTLWorker()
}
