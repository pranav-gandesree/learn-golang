# Mutexes — Mutual Exclusion Locks

> See [mutexes/main.go](../mutexes/main.go) for a runnable demo.

## Why mutexes?

When multiple goroutines READ AND WRITE the same memory without coordination, you get a **data race** — unpredictable values, lost updates, corruption.

Even `counter++` is not atomic — it's three CPU steps:
1. read counter from memory into a register
2. add 1 to the register
3. write the register back to memory

If two goroutines interleave these steps, both can read the same old value, both write the same new value — one update silently overwrites the other.

Demo output from [mutexes/main.go](../mutexes/main.go):
```
UNSAFE counter: expected 100000, got  72646 (lost 27354 updates)
SAFE counter:   expected 100000, got 100000
```

## sync.Mutex — basic locking

A lock that only one goroutine can hold at a time:

```go
import "sync"

var mu sync.Mutex

mu.Lock()
defer mu.Unlock()
// ... critical section ...
```

| Method | Behavior |
|---|---|
| `mu.Lock()` | If no one holds it → take it. If someone does → BLOCK until released. |
| `mu.Unlock()` | Release. A waiting goroutine can now take it. |

Code between `Lock` and `Unlock` is the **critical section** — exactly one goroutine can be inside it at a time.

## Always `defer mu.Unlock()`

Without `defer`, a panic or early `return` leaves the lock held forever → every other goroutine that tries to `Lock` blocks forever. **This is THE most common mutex bug.**

```go
func update() {
    mu.Lock()
    defer mu.Unlock()       // first thing after Lock

    // ... whatever ... (returns and panics are now safe)
}
```

## sync.RWMutex — many readers, one writer

If reads are frequent and writes are rare, `RWMutex` lets readers proceed concurrently:

```go
var rw sync.RWMutex

// Reading — multiple readers can hold simultaneously
rw.RLock()
defer rw.RUnlock()
v := cache[key]

// Writing — exclusive
rw.Lock()
defer rw.Unlock()
cache[key] = newValue
```

| Method | Use | Concurrency |
|---|---|---|
| `RLock` / `RUnlock` | Reading | Many readers OK simultaneously |
| `Lock` / `Unlock` | Writing | Exclusive — blocks all readers and writers |

## Mutex vs RWMutex vs atomic — which to use

| Pattern | Use |
|---|---|
| Single int / pointer | `sync/atomic` (fastest) |
| Roughly balanced reads/writes | `sync.Mutex` (RWMutex has overhead) |
| Reads dominate writes (>10×) | `sync.RWMutex` |
| Map/struct with mixed access | `sync.Mutex` |

Don't reach for `RWMutex` by default — it's actually slower than `Mutex` in low-contention scenarios because of its bookkeeping.

## The race detector

`go run -race ./mutexes` instruments the binary to detect races at runtime:

```
==================
WARNING: DATA RACE
Read at 0x00c000116038 by goroutine 9:
      .../mutexes/main.go:39 +0x88

Previous write at 0x00c000116038 by goroutine 7:
      .../mutexes/main.go:39 +0x98
==================
Found 2 data race(s)
```

It tells you the exact line, which goroutines collided, and where they were launched.

**Always run tests with `-race` in CI.** Slower (~2-20×) but catches bugs you'd never see by eye.

```bash
go test -race ./...
go run  -race ./mutexes
```

## Five pitfalls

### 1. Never copy a mutex

```go
// BAD — each function call gets its own independent mutex
func work(mu sync.Mutex) { ... }

// GOOD — share via pointer
func work(mu *sync.Mutex) { ... }
```

A copied mutex is a SECOND, INDEPENDENT lock — the original and the copy don't see each other. `go vet` catches most cases.

### 2. Always `defer mu.Unlock()`

See above. Without `defer`, any panic or early `return` deadlocks the rest of the program.

### 3. sync.Mutex is NOT reentrant

```go
mu.Lock()
mu.Lock()                    // ← DEADLOCKS against itself
```

If you're tempted to recursively lock, restructure: extract the locked work into a helper that *expects* the lock to be held, and call it without re-locking.

### 4. Keep critical sections small

```go
// BAD — holds lock during slow network call
mu.Lock()
data, _ := http.Get(url)
cache[key] = data
mu.Unlock()

// GOOD — do slow work FIRST, then take lock just to write
data, _ := http.Get(url)
mu.Lock()
cache[key] = data
mu.Unlock()
```

Locks should protect data, not entire workflows.

### 5. For a single int/pointer, use `sync/atomic`

```go
import "sync/atomic"

var counter atomic.Int64
counter.Add(1)
v := counter.Load()
```

Atomic operations are single CPU instructions — way faster than `Lock`/`Unlock` overhead. Use them for counters, flags, single pointers. For maps, slices, structs → use a mutex.

## Quick reference

```go
import "sync"

// Plain mutex
var mu sync.Mutex
mu.Lock(); defer mu.Unlock()

// Read-write mutex
var rw sync.RWMutex
rw.RLock();  defer rw.RUnlock()     // readers
rw.Lock();   defer rw.Unlock()      // writers

// Atomic counter (no mutex needed)
var n atomic.Int64
n.Add(1)
n.Load()

// Race detection
// go test -race ./...
// go run -race ./...
```

## Debugging data races

1. **Always `go test -race` in CI.** Hard build break on race detection.
2. **Run programs locally with `-race`** when you suspect concurrency bugs.
3. The race report points at exact lines — fix the read/write pair, re-run, repeat.
4. If a race appears intermittently, run the test in a loop: `go test -race -count=100 ./...`.
