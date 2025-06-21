# In Memory Store - Performance Benchmarks

## Test Environment

- **Operating System**: macOS 14.3.0 (Darwin 24.3.0)
- **Architecture**: Apple Silicon (ARM64)
- **CPU Cores**: 11
- **Memory**: 18 GB
- **Go Version**: go1.21.1 darwin/arm64
- **Test Iterations**: 20,000 per benchmark

## Benchmark Results

### Core Operations Performance

| Operation | Time per Op | Memory per Op | Allocs per Op |
|-----------|-------------|---------------|---------------|
| **Set (Permanent)** | 339.2 ns | 348 B | 2 allocs |
| **Set (with TTL)** | 307.9 ns | 348 B | 2 allocs |
| **Get** | 113.6 ns | 13 B | 1 alloc |
| **Update** | 135.0 ns | 13 B | 1 alloc |
| **Remove** | 130.8 ns | 20 B | 1 alloc |
| **Push** | 583.0 ns | 1,766 B | 5 allocs |
| **Pop** | 74.21 ns | 8 B | 1 alloc |

### Concurrent Operations Performance

| Operation | Time per Op | Memory per Op | Allocs per Op |
|-----------|-------------|---------------|---------------|
| **Get (Concurrent)** | 103.3 ns | 13 B | 1 alloc |
| **Set (Concurrent)** | 200.4 ns | 56 B | 1 alloc |
| **Push (Concurrent)** | 598.5 ns | 1,763 B | 5 allocs |
| **Pop (Concurrent)** | 221.5 ns | 8 B | 1 alloc |