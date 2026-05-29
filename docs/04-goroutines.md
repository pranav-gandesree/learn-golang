# Goroutines & WaitGroups

> See [goroutines/main.go](../goroutines/main.go) for a runnable demo.

## What's a goroutine?

A lightweight, independently-scheduled function call. Launched with the `go` keyword:

```go
go someFunction()             // returns IMMEDIATELY
go func() { ... }()           // anonymous function variant
```

- ~2KB to start (vs ~1MB for an OS thread)
- Multiplexed onto a small pool of OS threads by the Go runtime
- You can run millions of them

**Critical:** when `main()` returns, the WHOLE PROGRAM EXITS — including any goroutines still running. You must explicitly wait for them.

## Concurrent vs parallel

- **Concurrent** = structured so multiple things can be in progress at the same time
- **Parallel** = multiple things actually executing simultaneously on different CPU cores

Goroutines are always concurrent. Whether they're truly *parallel* depends on the workload and hardware. For I/O-bound work (network, files, sleep), even one CPU core handles many goroutines efficiently.

## sync.WaitGroup — the synchronization primitive

A counter with three methods:

| Method | Effect |
|---|---|
| `wg.Add(n)` | counter += n |
| `wg.Done()` | counter -= 1 |
| `wg.Wait()` | blocks until counter == 0 |

## The canonical pattern

```go
import "sync"

var wg sync.WaitGroup

wg.Add(1)                       // BEFORE launching
go func() {
    defer wg.Done()             // FIRST LINE — guarantees Done() runs
    // ... work ...
}()

wg.Wait()                       // blocks until all goroutines call Done()
```

## Three rules to avoid 90% of bugs

1. **`Add` BEFORE `go`.** If `Add(1)` is inside the goroutine, `Wait()` might see counter=0 and return before the goroutine has even started.
2. **`defer wg.Done()` as the first line of the goroutine.** Guarantees `Done` runs even if the function panics.
3. **Every `Add(n)` must match exactly `n` `Done()` calls.** Too few → deadlock. Too many → "negative counter" panic.

## Loop + goroutines — pass arguments

```go
for i := 0; i < n; i++ {
    wg.Add(1)
    go func(x int) {            // pass i as a parameter
        defer wg.Done()
        process(x)
    }(i)                        // ← pass i AT LAUNCH
}
wg.Wait()
```

**Why:** in Go < 1.22, capturing `i` by closure meant all goroutines shared the same variable and saw its final value. Go 1.22+ fixed this for `for` loops, but the pass-as-argument pattern works on EVERY Go version. Always use it.

## Multiple goroutines run concurrently

```go
wg.Add(3)
go func() { defer wg.Done(); time.Sleep(150 * time.Millisecond) }()
go func() { defer wg.Done(); time.Sleep(100 * time.Millisecond) }()
go func() { defer wg.Done(); time.Sleep(50  * time.Millisecond) }()
wg.Wait()
// Total wall-clock: ~150ms (the longest), NOT 300ms — they ran in parallel
```

## Order is NON-DETERMINISTIC

Goroutines run in whatever order the runtime schedules them. Two runs of the same program → potentially different orders. Never depend on a specific order without explicit synchronization (channels, mutexes, more WaitGroups).

## Don't use `time.Sleep` as a synchronization tool

```go
go doStuff()
time.Sleep(time.Second)          // BAD — guessing how long doStuff takes
```

- If `doStuff` takes longer than the sleep → main exits with work unfinished
- If shorter → wasted time

Use `WaitGroup` (or channels) — they return the *instant* work is done.

## Go 1.25+ shorthand: `wg.Go`

```go
// Pre-1.25
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()

// 1.25+
wg.Go(func() {
    doWork()
})
```

Does the same thing under the hood. Most existing Go code uses the explicit form.

## Quick reference

```go
import "sync"

var wg sync.WaitGroup
wg.Add(n)                      // before launching n goroutines
go func() {
    defer wg.Done()            // first line, always
    // ... work ...
}()
wg.Wait()                      // block until all Done() called

// Loop pattern
for _, item := range items {
    wg.Add(1)
    go func(x Item) {          // pass as argument
        defer wg.Done()
        process(x)
    }(item)
}
wg.Wait()
```

## Common pitfalls

| Symptom | Cause |
|---|---|
| Program exits silently with no output | Forgot to `Wait()` for goroutines |
| `panic: sync: negative WaitGroup counter` | Called `Done()` more times than `Add()` |
| `fatal error: all goroutines are asleep - deadlock!` | Called `Wait()` when no `Done()` will ever fire |
| All goroutines print the same value | Captured loop variable by closure (Go < 1.22) — pass as argument |
