package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mo-mohamed/acronis-memory-store/internal/store"
	"github.com/mo-mohamed/acronis-memory-store/internal/store/memory"
)

func TestHandler_SetAndGet(t *testing.T) {
	var memoryStore store.IStore = memory.NewMemoryStore()
	defer memoryStore.StopTTLWorker()

	handler := NewHandler(memoryStore)
	mux := handler.SetupRoutes()

	setPayload := SetRequest{
		Key:        "test_key",
		Value:      "test_value",
		TTLSeconds: 60,
	}
	payloadBytes, _ := json.Marshal(setPayload)

	req := httptest.NewRequest("POST", "/api/v1/keys", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var setResponse Response
	if err := json.NewDecoder(w.Body).Decode(&setResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !setResponse.Success {
		t.Errorf("Expected success=true, got %v", setResponse.Success)
	}

	req = httptest.NewRequest("GET", "/api/v1/keys/test_key", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var getResponse Response
	if err := json.NewDecoder(w.Body).Decode(&getResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !getResponse.Success {
		t.Errorf("Expected success=true, got %v", getResponse.Success)
	}

	data := getResponse.Data.(map[string]any)
	if data["value"] != "test_value" {
		t.Errorf("Expected value 'test_value', got %v", data["value"])
	}
}

func TestHandler_ListOperations(t *testing.T) {
	var memoryStore store.IStore = memory.NewMemoryStore()
	defer memoryStore.StopTTLWorker()

	handler := NewHandler(memoryStore)
	mux := handler.SetupRoutes()

	pushPayload := PushRequest{
		Key:  "test_list",
		Item: "item1",
	}
	payloadBytes, _ := json.Marshal(pushPayload)

	req := httptest.NewRequest("POST", "/api/v1/lists/push", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	popPayload := PopRequest{
		Key: "test_list",
	}
	payloadBytes, _ = json.Marshal(popPayload)

	req = httptest.NewRequest("POST", "/api/v1/lists/pop", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var popResponse Response
	if err := json.NewDecoder(w.Body).Decode(&popResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !popResponse.Success {
		t.Errorf("Expected success=true, got %v", popResponse.Success)
	}

	data := popResponse.Data.(map[string]any)
	if data["value"] != "item1" {
		t.Errorf("Expected value 'item1', got %v", data["value"])
	}
}

func TestHandler_UpdateAndRemove(t *testing.T) {
	var memoryStore store.IStore = memory.NewMemoryStore()
	defer memoryStore.StopTTLWorker()

	handler := NewHandler(memoryStore)
	mux := handler.SetupRoutes()

	setPayload := SetRequest{
		Key:        "update_key",
		Value:      "initial_value",
		TTLSeconds: 60,
	}
	payloadBytes, _ := json.Marshal(setPayload)

	req := httptest.NewRequest("POST", "/api/v1/keys", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	updatePayload := UpdateRequest{
		Value: "updated_value",
	}
	payloadBytes, _ = json.Marshal(updatePayload)

	req = httptest.NewRequest("PUT", "/api/v1/keys/update_key", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	req = httptest.NewRequest("DELETE", "/api/v1/keys/update_key", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_ErrorCases(t *testing.T) {
	var memoryStore store.IStore = memory.NewMemoryStore()
	handler := NewHandler(memoryStore)
	mux := handler.SetupRoutes()

	req := httptest.NewRequest("GET", "/api/v1/keys/nonexistent", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", "/api/v1/keys", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	req = httptest.NewRequest("PATCH", "/api/v1/keys/test", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandler_TTLValidation(t *testing.T) {
	var memoryStore store.IStore = memory.NewMemoryStore()
	defer memoryStore.StopTTLWorker()
	handler := NewHandler(memoryStore)
	mux := handler.SetupRoutes()

	t.Run("set with zero TTL", func(t *testing.T) {
		setPayload := SetRequest{
			Key:        "zero_ttl_key",
			Value:      "test_value",
			TTLSeconds: 0,
		}
		payloadBytes, _ := json.Marshal(setPayload)

		req := httptest.NewRequest("POST", "/api/v1/keys", bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var response Response
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success=false for zero TTL")
		}

		if response.Error != "TTL is required and must be greater than 0" {
			t.Errorf("Expected TTL error message, got %s", response.Error)
		}
	})

	t.Run("set with negative TTL", func(t *testing.T) {
		setPayload := SetRequest{
			Key:        "negative_ttl_key",
			Value:      "test_value",
			TTLSeconds: -5,
		}
		payloadBytes, _ := json.Marshal(setPayload)

		req := httptest.NewRequest("POST", "/api/v1/keys", bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var response Response
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success=false for negative TTL")
		}

		if response.Error != "TTL is required and must be greater than 0" {
			t.Errorf("Expected TTL error message, got %s", response.Error)
		}
	})
}
