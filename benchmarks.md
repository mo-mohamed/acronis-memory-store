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
| **Set (Permanent)** | 749.1 ns | 348 B | 2 allocs |
| **Set (with TTL)** | 447.4 ns | 347 B | 2 allocs |
| **Get** | 155.5 ns | 13 B | 1 alloc |
| **Get (Concurrent)** | 132.7 ns | 14 B | 1 alloc |
| **Set (Concurrent)** | 365.7 ns | 56 B | 1 alloc |
