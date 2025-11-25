## Purpose

Create performance benchmarks for rendering operations to establish baseline metrics, identify optimization opportunities, and prevent performance regressions.

## Acceptance Criteria

- [ ] Benchmark tests for style rendering (<1ms for simple styles)
- [ ] Benchmark tests for box rendering (<10ms for complex compositions)
- [ ] Benchmark tests for text alignment and padding
- [ ] Benchmark tests for ANSI code generation
- [ ] Benchmark results documented with baseline metrics
- [ ] Performance budgets enforced in CI (warn on >10% regression)
- [ ] Profiling data collected for hot paths

## Technical Approach

**Benchmark Structure**:
```go
func BenchmarkStyle_Render(b *testing.B) {
    s := style.New().Bold().Foreground(color.Red)
    text := "Hello, World!"
    for i := 0; i < b.N; i++ {
        s.Render(text)
    }
}
```

**Benchmark Categories**:
1. Style operations (ANSI code generation)
2. Text rendering (alignment, padding, wrapping)
3. Border rendering (box drawing, titles)
4. Complex compositions (nested boxes, dashboard)
5. Memory allocations (identify allocation hotspots)

**Performance Targets**:
- Simple style rendering: <1ms (target: <100μs)
- Box rendering: <10ms (target: <1ms)
- Text alignment: <500μs
- ANSI code generation: <50μs
- Dashboard composition: <50ms

**Profiling**:
```bash
go test -bench=. -cpuprofile=cpu.out
go tool pprof cpu.out
```

**Files to Create/Modify**:
- style/style_bench_test.go (new)
- render/render_bench_test.go (new)
- border/border_bench_test.go (new)
- layout/layout_bench_test.go (new)
- docs/performance.md (benchm ark results)

**Dependencies**:
- Go testing package (built-in benchmarks)

## Testing Strategy

**Benchmark Execution**:
```bash
go test -bench=. -benchmem ./...
go test -bench=BenchmarkStyle -count=5 -benchtime=3s
```

**Memory Profiling**:
- Track allocations per operation
- Identify string concatenation hotspots
- Optimize with strings.Builder
- Reduce interface{} boxing

**Regression Detection**:
- Run benchmarks on every PR
- Compare with baseline metrics
- Fail CI if >10% performance regression
- Use benchstat for statistical comparison

**Optimization Strategy**:
1. Profile to identify hot paths
2. Optimize string building (use strings.Builder)
3. Cache ANSI sequences where possible
4. Minimize allocations in tight loops
5. Re-benchmark after optimization

## Notes

**Benchmark Example**:
```go
func BenchmarkBorder_Box(b *testing.B) {
    content := "Sample content for box rendering"
    for i := 0; i < b.N; i++ {
        border.Box(content, border.Single, 40, 10)
    }
}

func BenchmarkRender_Alignment(b *testing.B) {
    text := "Center this text"
    for i := 0; i < b.N; i++ {
        render.Align(text, render.Center, 50)
    }
}
```

**Interpreting Results**:
```
BenchmarkStyle_Render-8    1000000    1234 ns/op    128 B/op    4 allocs/op
```
- 1234 ns/op: nanoseconds per operation (target: <1000)
- 128 B/op: bytes allocated per operation (minimize)
- 4 allocs/op: number of allocations (reduce if possible)

**Performance Budget**:
- Simple operations: <1ms
- Complex operations: <10ms
- Memory allocations: minimize to <10 per operation
- No allocations in hot paths where possible

**Tools**:
- go test -bench: Built-in benchmarking
- pprof: CPU and memory profiling
- benchstat: Statistical comparison of benchmark runs
- flamegraphs: Visualize performance bottlenecks

