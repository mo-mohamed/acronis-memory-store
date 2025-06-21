package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mo-mohamed/acronis-memory-store/pkg/client"
)

func TestClient_Set(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	err := c.Set(ctx, "test_key", "test_value", 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_Set_TTLValidation(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	// zero TTL
	err := c.Set(ctx, "test_key", "test_value", 0)
	if err == nil {
		t.Error("Expected error for zero TTL, got nil")
	}
	if !strings.Contains(err.Error(), "TTL is required and must be greater than 0") {
		t.Errorf("Expected TTL validation error, got %v", err)
	}

	// negative TTL
	err = c.Set(ctx, "test_key", "test_value", -5)
	if err == nil {
		t.Error("Expected error for negative TTL, got nil")
	}
	if !strings.Contains(err.Error(), "TTL is required and must be greater than 0") {
		t.Errorf("Expected TTL validation error, got %v", err)
	}
}

func TestClient_Get(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	value, err := c.Get(ctx, "test_key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}
}

func TestClient_Get_NotFound(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	_, err := c.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}

	if !strings.Contains(err.Error(), "Key not found") {
		t.Errorf("Expected 'Key not found' error, got %v", err)
	}
}

func TestClient_Update(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	err := c.Update(ctx, "test_key", "updated_value")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_Remove(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	err := c.Remove(ctx, "test_key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_Push(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	err := c.Push(ctx, "test_list", "test_item")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_Pop(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	value, err := c.Pop(ctx, "test_list")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if value != "test_item" {
		t.Errorf("Expected 'test_item', got %s", value)
	}
}

func TestClient_Pop_EmptyList(t *testing.T) {
	server := mockServer()
	defer server.Close()

	c := client.NewClient(server.URL)
	ctx := context.Background()

	_, err := c.Pop(ctx, "empty_list")
	if err == nil {
		t.Error("Expected error for empty list, got nil")
	}

	if !strings.Contains(err.Error(), "List is empty") {
		t.Errorf("Expected 'List is empty' error, got %v", err)
	}
}

// mockServer mimics the memory store API
func mockServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]any{
			"success": true,
			"data":    map[string]string{"message": "Key set successfully"},
		}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/api/v1/keys/", func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/api/v1/keys/")

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			if key == "nonexistent" {
				response := map[string]any{
					"success": false,
					"error":   "Key not found",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := map[string]any{
				"success": true,
				"data": map[string]string{
					"key":   key,
					"value": "test_value",
				},
			}
			json.NewEncoder(w).Encode(response)

		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			response := map[string]any{
				"success": true,
				"data":    map[string]string{"message": "Key updated successfully"},
			}
			json.NewEncoder(w).Encode(response)

		case http.MethodDelete:
			w.Header().Set("Content-Type", "application/json")
			response := map[string]any{
				"success": true,
				"data":    map[string]string{"message": "Key removed successfully"},
			}
			json.NewEncoder(w).Encode(response)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/lists/push", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]any{
			"success": true,
			"data":    map[string]string{"message": "Item pushed successfully"},
		}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/api/v1/lists/pop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)
		key := req["key"].(string)

		w.Header().Set("Content-Type", "application/json")

		if key == "empty_list" {
			response := map[string]any{
				"success": false,
				"error":   "List is empty",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := map[string]any{
			"success": true,
			"data": map[string]string{
				"key":   key,
				"value": "test_item",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(mux)
}
