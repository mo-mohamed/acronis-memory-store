package store

import "context"

// Store defines the interface for in memory data structure store
type IStore interface {
	Set(ctx context.Context, key string, value any, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Update(ctx context.Context, key string, value any) error
	Remove(ctx context.Context, key string) error
	Push(ctx context.Context, key string, item any) error
	Pop(ctx context.Context, key string) (string, error)
	StartTTLWorker(ctx context.Context)
	StopTTLWorker()
}
