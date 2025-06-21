package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mo-mohamed/acronis-memory-store/internal/store"
)

type Handler struct {
	store store.IStore
}

func NewHandler(s store.IStore) *Handler {
	return &Handler{store: s}
}

// SetHandler handles SET operations
// POST /api/v1/keys
func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Key == "" {
		h.writeError(w, http.StatusBadRequest, "Key is required")
		return
	}

	if req.TTLSeconds <= 0 {
		h.writeError(w, http.StatusBadRequest, "TTL is required and must be greater than 0")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Set(ctx, req.Key, req.Value, req.TTLSeconds); err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set key: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Key set successfully"})
}

// GetHandler handles GET operations
// GET /api/v1/keys/{key}
func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := r.URL.Path[len("/api/v1/keys/"):]
	if key == "" {
		h.writeError(w, http.StatusBadRequest, "Key is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := h.store.Get(ctx, key)
	if err != nil {
		if err.Error() == "key not found" {
			h.writeError(w, http.StatusNotFound, "Key not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get key: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"key": key, "value": value})
}

// UpdateHandler handles UPDATE operations
// PUT /api/v1/keys/{key}
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := r.URL.Path[len("/api/v1/keys/"):]
	if key == "" {
		h.writeError(w, http.StatusBadRequest, "Key is required")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Update(ctx, key, req.Value); err != nil {
		if err.Error() == "key not found" {
			h.writeError(w, http.StatusNotFound, "Key not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update key: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Key updated successfully"})
}

// RemoveHandler handles DELETE operations
// DELETE /api/v1/keys/{key}
func (h *Handler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	key := r.URL.Path[len("/api/v1/keys/"):]
	if key == "" {
		h.writeError(w, http.StatusBadRequest, "Key is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Remove(ctx, key); err != nil {
		if err.Error() == "key not found" {
			h.writeError(w, http.StatusNotFound, "Key not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to remove key: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Key removed successfully"})
}

// PushHandler handles PUSH operations for lists
// POST /api/v1/lists/push
func (h *Handler) PushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Push(ctx, req.Key, req.Item); err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to push item: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Item pushed successfully"})
}

// PopHandler handles POP operations for lists
// POST /api/v1/lists/pop
func (h *Handler) PopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req PopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := h.store.Pop(ctx, req.Key)
	if err != nil {
		if err.Error() == "key not found" {
			h.writeError(w, http.StatusNotFound, "Key not found")
			return
		}
		if err.Error() == "empty list" {
			h.writeError(w, http.StatusBadRequest, "List is empty")
			return
		}
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to pop item: %v", err))
		return
	}

	h.writeSuccess(w, map[string]string{"key": req.Key, "value": value})
}

// SetupRoutes sets up all the HTTP routes
func (h *Handler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/keys", h.SetHandler)
	// This is for GET, PUT and DELETE
	mux.HandleFunc("/api/v1/keys/", h.keyOperation)

	mux.HandleFunc("/api/v1/lists/push", h.PushHandler)
	mux.HandleFunc("/api/v1/lists/pop", h.PopHandler)

	return mux
}

// keyOperation handles GET, PUT and DELETE operations for keys as the request path is the same.
func (h *Handler) keyOperation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetHandler(w, r)
	case http.MethodPut:
		h.UpdateHandler(w, r)
	case http.MethodDelete:
		h.RemoveHandler(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// writeJSON is a helper function to write JSON responses
func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeError is a helper function to write error responses
func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// writeSuccess is a helper function to write success responses
func (h *Handler) writeSuccess(w http.ResponseWriter, data any) {
	h.writeJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}
