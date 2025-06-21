# Acronis In Memory Store API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Response Format
All API responses follow below standardized JSON structure:

**Success Response:**
```json
{
  "success": true,
  "data": {
    // Response information here
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message description"
}
```

## Content Type
All requests that include a body must use:
```
Content-Type: application/json
```

---

## Key-Value Operations

### 1. Set Key-Value Pair

Store a key-value pair with TTL.

**Endpoint:** `POST /api/v1/keys`

**Request Body:**
```json
{
  "key": "string (required)",
  "value": "any (required)",
  "ttl_seconds": "integer (required)"
}
```

**Parameters:**
- `key` (string, required): The key to store
- `value` (any, required): The value to store (can be string, number, object, etc.)
- `ttl_seconds` (integer, required): Time to live in seconds

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/keys \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "value": "my user",
    "ttl_seconds": 3600
  }'
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Key set successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `500 Internal Server Error`: Server error during operation

---

### 2. Get Value by Key

Retrieve a value by its key.

**Endpoint:** `GET /api/v1/keys/{key}`

**Path Parameters:**
- `key` (string, required): The key to retrieve

**Example Request:**
```bash
curl http://localhost:8080/api/v1/keys/user:123
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "key": "user:123",
    "value": "my user"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Key parameter is missing
- `404 Not Found`: Key does not exist or has expired
- `500 Internal Server Error`: Server error during operation

---

### 3. Update Key Value

Update the value of an existing key.

**Endpoint:** `PUT /api/v1/keys/{key}`

**Path Parameters:**
- `key` (string, required): The key to update

**Request Body:**
```json
{
  "value": "any (required)"
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/keys/user:123 \
  -H "Content-Type: application/json" \
  -d '{
    "value": "my user"
  }'
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Key updated successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing value field
- `404 Not Found`: Key does not exist or has expired
- `500 Internal Server Error`: Server error during operation

---

### 4. Delete Key

Remove a key and its value from the store.

**Endpoint:** `DELETE /api/v1/keys/{key}`

**Path Parameters:**
- `key` (string, required): The key to delete

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/keys/user:123
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Key removed successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Key parameter is missing
- `404 Not Found`: Key does not exist
- `500 Internal Server Error`: Server error during operation

---

## List Operations

### 5. Push Item to List (LPUSH)

Add an item to the front of a list. If the list doesn't exist, it will be created.

**Endpoint:** `POST /api/v1/lists/push`

**Request Body:**
```json
{
  "key": "string (required)",
  "item": "any (required)"
}
```

**Parameters:**
- `key` (string, required): The list key
- `item` (any, required): The item to add to the front of the list

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/lists/push \
  -H "Content-Type: application/json" \
  -d '{
    "key": "queue:tasks",
    "item": "my item"
  }'
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Item pushed successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `500 Internal Server Error`: Server error during operation

---

### 6. Pop Item from List (LPOP)

Remove and return an item from the front of a list.

**Endpoint:** `POST /api/v1/lists/pop`

**Request Body:**
```json
{
  "key": "string (required)"
}
```

**Parameters:**
- `key` (string, required): The list key

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/lists/pop \
  -H "Content-Type: application/json" \
  -d '{
    "key": "queue:tasks"
  }'
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "key": "queue:tasks",
    "value": "my item"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON, missing key field, or list is empty
- `404 Not Found`: Key does not exist
- `500 Internal Server Error`: Server error during operation

---

## HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 400 | Bad Request - Invalid request format or parameters |
| 404 | Not Found - Requested resource does not exist |
| 405 | Method Not Allowed - HTTP method not supported for endpoint |
| 500 | Internal Server Error - Server encountered an error |

---

## Error Messages

### Common Error Messages

| Error Message | Description | Status Code |
|---------------|-------------|-------------|
| "Key is required" | The key parameter is missing or empty | 400 |
| "Key not found" | The requested key does not exist or has expired | 404 |
| "Invalid JSON payload" | The request body contains invalid JSON | 400 |
| "Method not allowed" | The HTTP method is not supported for this endpoint | 405 |
| "List is empty" | Attempted to pop from an empty list | 400 |
| "Failed to set key: ..." | Server error during set operation | 500 |
| "Failed to get key: ..." | Server error during get operation | 500 |
| "Failed to update key: ..." | Server error during update operation | 500 |
| "Failed to remove key: ..." | Server error during remove operation | 500 |
| "Failed to push item: ..." | Server error during push operation | 500 |
| "Failed to pop item: ..." | Server error during pop operation | 500 |

---