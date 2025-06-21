package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/mo-mohamed/acronis-memory-store/internal/store/memory"
)

func TestSet(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	t.Run("basic set", func(t *testing.T) {
		err := store.Set(ctx, "key1", "value1", 60)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}
	})

	t.Run("set with TTL", func(t *testing.T) {
		err := store.Set(ctx, "key2", "value2", 60)
		if err != nil {
			t.Errorf("Set with TTL failed: %v", err)
		}
	})

	t.Run("set with negative TTL", func(t *testing.T) {
		err := store.Set(ctx, "key3", "value3", -1)
		if err == nil {
			t.Error("Expected error for negative TTL")
		}
		if err != nil && err.Error() != "invalid TTL value" {
			t.Errorf("Expected 'invalid TTL value', got %v", err)
		}
	})

	t.Run("set with zero TTL (no expiration)", func(t *testing.T) {
		err := store.Set(ctx, "key4", "value4", 0)
		if err != nil {
			t.Errorf("Set with zero TTL should be valid: %v", err)
		}

		value, err := store.Get(ctx, "key4")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if value != "value4" {
			t.Errorf("Expected 'value4', got %q", value)
		}

		time.Sleep(100 * time.Millisecond)
		value, err = store.Get(ctx, "key4")
		if err != nil {
			t.Errorf("Get after wait failed: %v", err)
		}
		if value != "value4" {
			t.Errorf("Expected 'value4', got %q", value)
		}
	})

	t.Run("set complex type", func(t *testing.T) {
		err := store.Set(ctx, "key5", map[string]int{"count": 42}, 60)
		if err != nil {
			t.Errorf("Set complex type failed: %v", err)
		}
	})
}

func TestGet(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	store.Set(ctx, "existing", "value", 60)
	store.Set(ctx, "complex", map[string]int{"num": 42}, 60)

	t.Run("get existing key", func(t *testing.T) {
		value, err := store.Get(ctx, "existing")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if value != "value" {
			t.Errorf("Expected 'value', got %q", value)
		}
	})

	t.Run("get non-existing key", func(t *testing.T) {
		_, err := store.Get(ctx, "nonexistent")
		if err == nil {
			t.Error("Expected error for non-existing key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("get complex type", func(t *testing.T) {
		value, err := store.Get(ctx, "complex")
		if err != nil {
			t.Errorf("Get complex type failed: %v", err)
		}
		expected := `{"num":42}`
		if value != expected {
			t.Errorf("Expected %q, got %q", expected, value)
		}
	})
}

func TestUpdate(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	store.Set(ctx, "existing", "original", 60)

	t.Run("update existing key", func(t *testing.T) {
		err := store.Update(ctx, "existing", "updated")
		if err != nil {
			t.Errorf("Update failed: %v", err)
		}

		value, err := store.Get(ctx, "existing")
		if err != nil {
			t.Errorf("Get after update failed: %v", err)
		}
		if value != "updated" {
			t.Errorf("Expected 'updated', got %q", value)
		}
	})

	t.Run("update non-existing key", func(t *testing.T) {
		err := store.Update(ctx, "nonexistent", "value")
		if err == nil {
			t.Error("Expected error for non-existing key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("update preserves TTL", func(t *testing.T) {
		store.Set(ctx, "ttl_key", "original", 10)

		err := store.Update(ctx, "ttl_key", "updated")
		if err != nil {
			t.Errorf("Update failed: %v", err)
		}

		value, err := store.Get(ctx, "ttl_key")
		if err != nil {
			t.Errorf("Get after update failed: %v", err)
		}
		if value != "updated" {
			t.Errorf("Expected 'updated', got %q", value)
		}
	})
}

func TestRemove(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	store.Set(ctx, "to_remove", "value", 60)

	t.Run("remove existing key", func(t *testing.T) {
		err := store.Remove(ctx, "to_remove")
		if err != nil {
			t.Errorf("Remove failed: %v", err)
		}

		_, err = store.Get(ctx, "to_remove")
		if err == nil {
			t.Error("Expected error after removing key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("remove non-existing key", func(t *testing.T) {
		err := store.Remove(ctx, "nonexistent")
		if err == nil {
			t.Error("Expected error for non-existing key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})
}

func TestPush(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	t.Run("push to new list", func(t *testing.T) {
		err := store.Push(ctx, "list1", "item1")
		if err != nil {
			t.Errorf("Push to new list failed: %v", err)
		}
	})

	t.Run("push to existing list", func(t *testing.T) {
		list := "list2"
		store.Push(ctx, list, "first")
		err := store.Push(ctx, list, "second")
		if err != nil {
			t.Errorf("Push to existing list failed: %v", err)
		}
		_, err = store.Pop(ctx, list)
		if err != nil {
			t.Errorf("Pop from existing list failed: %v", err)
		}
	})

	t.Run("push to string key", func(t *testing.T) {
		key := "string_key"
		store.Set(ctx, key, "value", 60)
		err := store.Push(ctx, key, "item")
		if err == nil {
			t.Error("Expected error when pushing to string key")
		}
		if err != nil && err.Error() != "operation not supported for this data type" {
			t.Errorf("Expected 'operation not supported for this data type', got %v", err)
		}
	})

	t.Run("push complex type", func(t *testing.T) {
		err := store.Push(ctx, "list3", map[string]int{"id": 1})
		if err != nil {
			t.Errorf("Push complex type failed: %v", err)
		}
	})

	t.Run("push multiple items and verify order", func(t *testing.T) {
		store.Push(ctx, "order_test", "first")
		store.Push(ctx, "order_test", "second")
		store.Push(ctx, "order_test", "third")

		item, err := store.Pop(ctx, "order_test")
		if err != nil {
			t.Errorf("Pop failed: %v", err)
		}
		if item != "third" {
			t.Errorf("Expected 'third', got %q", item)
		}
	})
}

func TestPop(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	store.Push(ctx, "list1", "first")
	store.Push(ctx, "list1", "second")

	t.Run("pop from list", func(t *testing.T) {
		item, err := store.Pop(ctx, "list1")
		if err != nil {
			t.Errorf("Pop failed: %v", err)
		}
		if item != "second" {
			t.Errorf("Expected 'second', got %q", item)
		}
	})

	t.Run("pop from empty list", func(t *testing.T) {
		list := "list1"
		store.Pop(ctx, list)

		_, err := store.Pop(ctx, list)
		if err == nil {
			t.Error("Expected error when popping from empty list")
		}
		if err != nil && err.Error() != "list is empty" {
			t.Errorf("Expected 'list is empty', got %v", err)
		}
	})

	t.Run("pop from non-existing key", func(t *testing.T) {
		_, err := store.Pop(ctx, "nonexistent")
		if err == nil {
			t.Error("Expected error for non-existing key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("pop from string key", func(t *testing.T) {
		key := "string_key_pop"
		store.Set(ctx, key, "value", 60)
		_, err := store.Pop(ctx, key)
		if err == nil {
			t.Error("Expected error when popping from string key")
		}
		if err != nil && err.Error() != "operation not supported for this data type" {
			t.Errorf("Expected 'operation not supported for this data type', got %v", err)
		}
	})
}

func TestTTLExpiration(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	t.Run("string TTL expiration", func(t *testing.T) {
		// Set key with 1 second as TTL
		err := store.Set(ctx, "expire_string", "temporary", 1)
		if err != nil {
			t.Errorf("Set with TTL failed: %v", err)
		}

		// Check if the key exists
		value, err := store.Get(ctx, "expire_string")
		if err != nil {
			t.Errorf("Get before expiration failed: %v", err)
		}
		if value != "temporary" {
			t.Errorf("Expected 'temporary', got %q", value)
		}

		// Wait for expiration + cleanup
		time.Sleep(2500 * time.Millisecond)

		// Check if the key expired
		_, err = store.Get(ctx, "expire_string")
		if err == nil {
			t.Error("Expected error for expired key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("list TTL expiration", func(t *testing.T) {
		list := "expire_list"

		store.Push(ctx, list, "item1")

		time.Sleep(1500 * time.Millisecond)

		err := store.Push(ctx, list, "new_item")
		if err != nil {
			t.Errorf("Push to expired list failed: %v", err)
		}
	})
}

func TestConcurrentOperations(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	t.Run("concurrent reads and writes for strings", func(t *testing.T) {
		done := make(chan bool, 2)
		key := "concurrent"

		// set values
		go func() {
			for i := 0; i < 100; i++ {
				store.Set(ctx, key, "value", 60)
			}
			done <- true
		}()

		// get values
		go func() {
			for i := 0; i < 100; i++ {
				store.Get(ctx, key)
			}
			done <- true
		}()

		<-done
		<-done
	})

	t.Run("concurrent list operations for lists", func(t *testing.T) {
		done := make(chan bool, 2)
		list := "concurrent"

		// set values from a list
		go func() {
			for i := 0; i < 50; i++ {
				store.Push(ctx, list, "item")
			}
			done <- true
		}()

		// get values from a list
		go func() {
			for i := 0; i < 50; i++ {
				store.Pop(ctx, list)
			}
			done <- true
		}()

		<-done
		<-done
	})
}

func TestStringify(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	tests := []struct {
		name     string
		input    any
		expected string
		hasError bool
	}{
		{"string", "hello", "hello", false},
		{"int", 42, "42", false},
		{"bool", true, "true", false},
		{"slice", []string{"a", "b"}, `["a","b"]`, false},
		{"map", map[string]int{"key": 1}, `{"key":1}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.Stringify(tt.input)

			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLazyExpiration(t *testing.T) {
	store := memory.NewMemoryStore()
	defer store.StopTTLWorker()
	ctx := context.Background()

	t.Run("Get deletes expired keys", func(t *testing.T) {
		err := store.Set(ctx, "lazy_expire", "value", 1)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		_, err = store.Get(ctx, "lazy_expire")
		if err == nil {
			t.Error("Expected error for expired key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}

		_, err = store.Get(ctx, "lazy_expire")
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found' on second access, got %v", err)
		}
	})

	t.Run("Update deletes expired keys", func(t *testing.T) {
		err := store.Set(ctx, "lazy_expire_update", "value", 1)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		err = store.Update(ctx, "lazy_expire_update", "new_value")
		if err == nil {
			t.Error("Expected error for expired key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})

	t.Run("Pop deletes expired keys", func(t *testing.T) {
		err := store.Set(ctx, "lazy_expire_pop", "value", 1)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		_, err = store.Pop(ctx, "lazy_expire_pop")
		if err == nil {
			t.Error("Expected error for expired key")
		}
		if err != nil && err.Error() != "key not found" {
			t.Errorf("Expected 'key not found', got %v", err)
		}
	})
}
