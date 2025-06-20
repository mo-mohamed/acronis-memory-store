package memory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mo-mohamed/acronis-memory-store/internal/store/memory"
)

func BenchmarkSet(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "value", 0)
	}
}

func BenchmarkSetWithTTL(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "value", 60)
	}
}

func BenchmarkGet(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	// Pre-populate store
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		store.Get(ctx, key)
	}
}

func BenchmarkConcurrentGet(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	// Pre-populate store
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "value", 0)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i%1000)
			store.Get(ctx, key)
			i++
		}
	})
}

func BenchmarkConcurrentSet(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i)
			store.Set(ctx, key, "value", 0)
			i++
		}
	})
}

func BenchmarkScalability(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			store := memory.NewMemoryStore()
			defer store.StopTTLWorker()
			ctx := context.Background()

			// Pre-populate
			for i := 0; i < size; i++ {
				key := fmt.Sprintf("key_%d", i)
				store.Set(ctx, key, "value", 0)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key_%d", i%size)
				store.Get(ctx, key)
			}
		})
	}
}
