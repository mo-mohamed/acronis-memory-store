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

func BenchmarkUpdate(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "initial_value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		store.Update(ctx, key, "updated_value")
	}
}

func BenchmarkRemove(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		store.Set(ctx, key, "value", 0)
		store.Remove(ctx, key)
	}
}

func BenchmarkPush(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list_%d", i%100)
		item := fmt.Sprintf("item_%d", i)
		store.Push(ctx, key, item)
	}
}

func BenchmarkPop(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("list_%d", i)
		for j := 0; j < 1000; j++ {
			item := fmt.Sprintf("item_%d_%d", i, j)
			store.Push(ctx, key, item)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list_%d", i%100)
		store.Pop(ctx, key)
	}
}

func BenchmarkConcurrentGet(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

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

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i)
			store.Set(ctx, key, "value", 0)
			i++
		}
	})
}

func BenchmarkConcurrentPush(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("list_%d", i%100)
			item := fmt.Sprintf("item_%d", i)
			store.Push(ctx, key, item)
			i++
		}
	})
}

func BenchmarkConcurrentPop(b *testing.B) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("list_%d", i)
		for j := 0; j < 10000; j++ {
			item := fmt.Sprintf("item_%d_%d", i, j)
			store.Push(ctx, key, item)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("list_%d", i%100)
			store.Pop(ctx, key)
			i++
		}
	})
}
