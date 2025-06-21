# Acronis In Memory Store
## Requirements 

Need to implement simple in-memory data structure store. See Redis for example.

Required data strutures:

Strings
Lists
Required operations:

- Get
- Set
- Update
- Remove
- Push for lists
- Pop for lists


Required features:

- Keys with a limited TTL
- Go client API library
- HTTP REST API

Add unit tests for Go API and integration tests for REST API (without full coverage, just for example). Provide REST API specs with examples, client library API docs and deployment docs (for Docker).

Optional features:

- Data persistence 
- Perfomance tests 
- Authentication

## Application Setup

### Prerequisites
- Go 1.21.1 or higher

### Installation & Running


#### Running the Application Locally

1. **Clone the repository**
```bash
git clone https://github.com/mo-mohamed/acronis-memory-store
cd acronis-memory-store
```

2. **Build and run the server**
```bash
go run cmd/server/main.go
```

3. **Server will start on port 8080**
```
starting server on port 8080
```

4. **Custom port (optional)**
```bash
PORT=3000 go run cmd/server/main.go
```

#### Running the Application in Docker
```bash
docker compose up
```

## Documentation

- [API Documentation](./API.md) - Detailed API specifications with examples
- [Benchamring Results](./benchmarks.md) - Testing and building instructions
