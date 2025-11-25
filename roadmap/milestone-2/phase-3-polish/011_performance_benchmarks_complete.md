## Purpose

Measure and document the performance characteristics of the copy-on-write builder pattern. Establish baseline benchmarks for copy overhead, method chaining, and allocation patterns to guide future optimizations.

## Acceptance Criteria

- [ ] Benchmarks for single method calls (<500ns per call)
- [ ] Benchmarks for method chaining (linear scaling)
- [ ] Benchmarks for style copying (full struct copy)
- [ ] Memory allocation profiling (bytes per operation)
- [ ] Comparison: shallow copy vs deep copy
- [ ] Benchmark results documented in comments
- [ ] All benchmarks run with -benchmem flag
- [ ] Performance budget: <500ns per method call, <10 allocs per chain

## Technical Approach

Create a dedicated `benchmark_test.go` file with comprehensive performance benchmarks. Measure time and allocations for various builder API usage patterns.

**Benchmark Scenarios**:

1. **Single Method Call**:
```go
func BenchmarkStyle_SingleMethodCall(b *testing.B) {
    s := NewStyle()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = s.Bold(true)
    }
}
// Expected: <100ns/op, 1 alloc/op (bool pointer)
```

2. **Method Chaining - Short**:
```go
func BenchmarkStyle_ShortChain(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = NewStyle().
            Bold(true).
            Foreground(Red).
            Padding(2)
    }
}
// Expected: <500ns/op, 5-6 allocs/op
```

3. **Method Chaining - Long**:
```go
func BenchmarkStyle_LongChain(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = NewStyle().
            Bold(true).
            Italic(true).
            Underline(true).
            Foreground(Blue).
            Background(White).
            Width(100).
            Height(30).
            Padding(2, 4).
            Margin(1).
            Border(Rounded).
            BorderForeground(Gray).
            Align(Center).
            AlignVertical(Middle)
    }
}
// Expected: <2000ns/op, 15-20 allocs/op
```

4. **Style Copy Overhead**:
```go
func BenchmarkStyle_Copy(b *testing.B) {
    s := NewStyle().
        Bold(true).
        Foreground(Red).
        Width(80).
        Padding(2).
        Border(Rounded)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = s.Italic(true)  // Forces copy
    }
}
// Expected: <100ns/op (shallow copy of ~32 pointers)
```

5. **CSS Shorthand**:
```go
func BenchmarkStyle_CSSShorthand(b *testing.B) {
    b.Run("Padding1", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = NewStyle().Padding(2)
        }
    })
    b.Run("Padding2", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = NewStyle().Padding(1, 2)
        }
    })
    b.Run("Padding4", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = NewStyle().Padding(1, 2, 3, 4)
        }
    })
}
// Expected: Padding(1) ~500ns, Padding(2) ~700ns, Padding(4) ~1000ns
```

6. **Branching Styles**:
```go
func BenchmarkStyle_Branching(b *testing.B) {
    baseStyle := NewStyle().
        Padding(1, 3).
        Border(Rounded).
        Align(Center)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = baseStyle.Foreground(Blue).Background(White)
        _ = baseStyle.Foreground(Red).Background(White)
        _ = baseStyle.Foreground(Green).Background(White)
    }
}
// Expected: <1000ns/op for 3 branches
```

7. **Allocation Breakdown**:
```go
func BenchmarkStyle_AllocationBreakdown(b *testing.B) {
    b.Run("TextAttribute", func(b *testing.B) {
        s := NewStyle()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = s.Bold(true)
        }
        // Allocates: 1 bool pointer
    })

    b.Run("Color", func(b *testing.B) {
        s := NewStyle()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = s.Foreground(Red)
        }
        // Allocates: 1 Color pointer
    })

    b.Run("Dimension", func(b *testing.B) {
        s := NewStyle()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = s.Width(80)
        }
        // Allocates: 1 int pointer
    })
}
```

8. **Comparison: Mutable vs Immutable**:
```go
// Document the trade-off
func BenchmarkStyle_MutableHypothetical(b *testing.B) {
    // Hypothetical mutable implementation (for comparison)
    // Would be 0 allocs but loses immutability guarantees
    b.Skip("Hypothetical comparison - mutable would be faster but unsafe")
}
```

**Benchmark Analysis**:
```bash
go test -bench=. -benchmem -benchtime=5s > benchmarks.txt
go test -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
# Analyze allocation hotspots
```

**Files to Create/Modify**:
- benchmark_test.go (create new file)

**Dependencies**:
- All builder methods from tasks 002-007

## Testing Strategy

**Benchmark Goals**:
- Single method call: <500ns/op, 1-2 allocs/op
- Short chain (3 methods): <500ns/op, 5-6 allocs/op
- Long chain (15 methods): <2000ns/op, 20-25 allocs/op
- CSS shorthand: Linear scaling with arg count

**Performance Analysis**:
- Run benchmarks on clean system (no background load)
- Compare results across Go versions
- Profile allocations with pprof
- Identify optimization opportunities (if needed)

**Documentation**:
- Add benchmark results to godoc
- Document performance budget in README
- Explain copy overhead trade-offs

## Notes

- Benchmark results vary by hardware - document test environment
- Shallow copy overhead is acceptable for immutability guarantees
- If performance issues found, consider arena allocation in later milestone
- Most real-world usage won't chain 100+ methods
- Profile allocations with `go test -bench=. -benchmem -memprofile=mem.prof`
- Use `benchstat` tool to compare benchmark runs statistically


