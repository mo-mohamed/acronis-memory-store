// Package client provides a Go client library for the Memory Store API.
//
// The client supports all core operations for managing strings and lists with TTL:
//   - Set: Store key-value pairs with required TTL
//   - Get: Retrieve values by key
//   - Update: Modify existing key values
//   - Remove: Delete keys
//   - Push: Add items to lists (LPUSH)
//   - Pop: Remove and return items from lists (LPOP)
//
// Basic usage:
//
//	client := client.NewClient("http://localhost:8080")
//	ctx := context.Background()
//
//	// Store a string value with 1-hour TTL
//	err := client.Set(ctx, "user:123", "John Doe", 3600)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Retrieve the value
//	value, err := client.Get(ctx, "user:123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(value) // Output: John Doe
//
//	// Work with lists
//	err = client.Push(ctx, "queue:tasks", "process-order")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	item, err := client.Pop(ctx, "queue:tasks")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(item) // Output: process-order
//
// All operations require proper context for cancellation and timeout handling.
// TTL is required for all Set operations and must be greater than 0.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Acronis Memory Store API client.
// It provides methods to interact with the memory store server
// for managing strings and lists with TTL support.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Acronis Memory Store API client.
// The baseURL should point to the memory store server (e.g., "http://localhost:8080").
//
// Example:
//
//	client := client.NewClient("http://localhost:8080")
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Set stores a key-value pair with the specified TTL in seconds.
// TTL of 0 means no expiration, TTL > 0 means expires after specified seconds.
// The value can be any JSON-serializable type.
//
// Example:
//
//	// Store a string value with 1-hour TTL
//	err := client.Set(ctx, "user:123", "John Doe", 3600)
//
//	// Store a permanent value (no expiration)
//	err := client.Set(ctx, "config:app", "permanent setting", 0)
//
//	// Store a complex object with 30-minute TTL
//	user := map[string]interface{}{
//	    "name": "John Doe",
//	    "age":  30,
//	}
//	err := client.Set(ctx, "user:profile:123", user, 1800)
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	if ttlSeconds < 0 {
		return fmt.Errorf("TTL must be >= 0 (0 = no expiration)")
	}

	req := SetRequest{
		Key:        key,
		Value:      value,
		TTLSeconds: ttlSeconds,
	}

	_, err := c.doRequest(ctx, "POST", "/api/v1/keys", req)
	return err
}

// Get retrieves a value by its key. Returns the value as a string.
// If the key doesn't exist or has expired, returns an error.
//
// Example:
//
//	value, err := client.Get(ctx, "user:123")
//	if err != nil {
//	    if strings.Contains(err.Error(), "Key not found") {
//	        fmt.Println("Key doesn't exist")
//	    } else {
//	        log.Fatal(err)
//	    }
//	}
//	fmt.Println("Value:", value)
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/keys/"+key, nil)
	if err != nil {
		return "", err
	}

	// Parse the response data
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	value, ok := data["value"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected value format")
	}

	return value, nil
}

// Update modifies the value of an existing key. The key must exist.
// This operation preserves the original TTL of the key.
//
// Example:
//
//	// Update existing user
//	err := client.Update(ctx, "user:123", "Jane Doe")
//	if err != nil {
//	    if strings.Contains(err.Error(), "Key not found") {
//	        fmt.Println("Key doesn't exist - use Set instead")
//	    } else {
//	        log.Fatal(err)
//	    }
//	}
func (c *Client) Update(ctx context.Context, key string, value interface{}) error {
	req := UpdateRequest{
		Value: value,
	}

	_, err := c.doRequest(ctx, "PUT", "/api/v1/keys/"+key, req)
	return err
}

// Remove deletes a key and its value from the store.
// If the key doesn't exist, the operation succeeds without error.
//
// Example:
//
//	err := client.Remove(ctx, "user:123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Key removed successfully")
func (c *Client) Remove(ctx context.Context, key string) error {
	_, err := c.doRequest(ctx, "DELETE", "/api/v1/keys/"+key, nil)
	return err
}

// Push adds an item to the front of a list (LPUSH operation).
// If the list doesn't exist, it will be created automatically.
// The item can be any JSON-serializable type.
//
// Example:
//
//	// Add tasks to a queue
//	err := client.Push(ctx, "queue:tasks", "process-order-123")
//	err = client.Push(ctx, "queue:tasks", "send-email-456")
//
//	// Add complex objects
//	task := map[string]interface{}{
//	    "id":   "task-789",
//	    "type": "backup",
//	    "priority": 1,
//	}
//	err = client.Push(ctx, "queue:priority", task)
func (c *Client) Push(ctx context.Context, key string, item interface{}) error {
	req := PushRequest{
		Key:  key,
		Item: item,
	}

	_, err := c.doRequest(ctx, "POST", "/api/v1/lists/push", req)
	return err
}

// Pop removes and returns an item from the front of a list (LPOP operation).
// Returns the item as a string. If the list is empty or doesn't exist,
// returns an error.
//
// Example:
//
//	// Process tasks from queue
//	for {
//	    item, err := client.Pop(ctx, "queue:tasks")
//	    if err != nil {
//	        if strings.Contains(err.Error(), "list is empty") {
//	            fmt.Println("No more tasks to process")
//	            break
//	        }
//	        log.Fatal(err)
//	    }
//	    fmt.Println("Processing:", item)
//	    // Process the item...
//	}
func (c *Client) Pop(ctx context.Context, key string) (string, error) {
	req := PopRequest{
		Key: key,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/lists/pop", req)
	if err != nil {
		return "", err
	}

	// Parse the response data
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	value, ok := data["value"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected value format")
	}

	return value, nil
}

// doRequest performs an HTTP request and handles the response.
// This is an internal method used by all public client methods.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*Response, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !apiResp.Success {
		return &apiResp, fmt.Errorf("API error: %s", apiResp.Error)
	}

	return &apiResp, nil
}
