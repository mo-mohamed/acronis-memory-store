// Package client provides data structures and models for the Memory Store API client.
package client

// Response represents the standard API response structure returned by all endpoints.
// It contains a success flag, optional data payload, and optional error message.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SetRequest represents the request payload for SET operations.
// It contains the key to store, the value to associate with the key,
// and the TTL in seconds.
type SetRequest struct {
	Key        string `json:"key"`
	Value      any    `json:"value"`
	TTLSeconds int    `json:"ttl_seconds"`
}

// UpdateRequest represents the request payload for UPDATE operations.
// It contains only the new value to update an existing key with.
// The key is specified in the URL path.
type UpdateRequest struct {
	Value any `json:"value"`
}

// PushRequest represents the request payload for PUSH operations on lists.
// It contains the list key and the item to add to the front of the list.
type PushRequest struct {
	Key  string `json:"key"`
	Item any    `json:"item"`
}

// PopRequest represents the request payload for POP operations on lists.
// It contains only the list key from which to remove and return the front item.
type PopRequest struct {
	Key string `json:"key"`
}
