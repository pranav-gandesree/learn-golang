# Channels

> See [channels/main.go](../channels/main.go) for a runnable demo.

## What is a channel?

A **channel** is a typed conduit for sending values between goroutines. A thread-safe pipe.

```go
ch := make(chan int)    // create
ch <- 42                // SEND (arrow points INTO ch)
v := <-ch               // RECEIVE (arrow points OUT of ch)
```

The Go proverb: **"Don't communicate by sharing memory; share memory by communicating."** Channels are the idiomatic Go way for goroutines to coordinate. Mutexes guard shared state; channels pass values.

## Unbuffered vs buffered

| Kind | Create | Send blocks if... | Receive blocks if... |
|---|---|---|---|
| **Unbuffered** | `make(chan T)` | no receiver ready | no sender ready |
| **Buffered**  | `make(chan T, N)` | buffer is FULL | buffer is EMPTY |

### Unbuffered = synchronous handshake

```go
ch := make(chan int)
go func() { ch <- 42 }()    // blocks until main receives
v := <-ch                   // both sides meet
```

Sender and receiver rendezvous — neither proceeds until both are ready.

### Buffered = mailbox with capacity

```go
ch := make(chan string, 3)
ch <- "a"   // no block
ch <- "b"   // no block
ch <- "c"   // no block — buffer full
// ch <- "d"   // would BLOCK forever (no receiver)

len(ch)      // 3 — current items in buffer
cap(ch)      // 3 — buffer size
```

FIFO order. Receives don't block until the buffer is empty.

**Unbuffered ≠ "buffer of 1".** Unbuffered is a synchronous handshake; size-1 buffered has temporary storage that decouples sender and receiver in time.

## Closing a channel + range

```go
close(ch)                    // signals "no more values will be sent"
for v := range ch { ... }    // loops until ch is closed
v, ok := <-ch                // ok=false after channel is closed and drained
```

### Rules of close

1. **Only the SENDER should close.** A receiver doesn't know if other senders are still active.
2. **Sending to a closed channel PANICS** (`send on closed channel`).
3. **Closing a closed channel PANICS.**
4. **Receiving from a closed channel returns the ZERO VALUE immediately.** Use comma-ok (`v, ok := <-ch`) to detect.
5. **You don't have to close every channel** — only when receivers need to know "no more values coming" (e.g., for `range` loops).

## Channel direction (send-only, receive-only)

```go
func produce(ch chan<- int) { ch <- 42 }    // send-only
func consume(ch <-chan int) { v := <-ch }   // receive-only
```

| Type | Means |
|---|---|
| `chan T`   | bidirectional (the default when you `make`) |
| `chan<- T` | send-only (arrow points INTO chan) |
| `<-chan T` | receive-only (arrow points OUT of chan) |

A regular `chan T` implicitly converts to either direction when passed as an argument. This is purely a compile-time constraint — useful for documenting intent and preventing bugs (e.g., a "producer" function that mistakenly tries to receive).

## `select` — wait on multiple channels

`select` is to channels what `switch` is to values. Blocks until ONE case can proceed; if multiple are ready, picks randomly.

```go
select {
case v := <-ch1:
    // got something from ch1
case v := <-ch2:
    // got something from ch2
case ch3 <- 42:
    // sent 42 into ch3
case <-time.After(time.Second):
    // 1-second timeout
default:
    // ran immediately because no case was ready (non-blocking select)
}
```

Common uses:
- **Timeout**: pair a real receive with `<-time.After(...)`
- **Non-blocking try**: add `default:` so it returns instantly if nothing is ready
- **Multiplexing**: react to whichever of several channels fires first

## Worker pool pattern

The classic Go concurrency pattern. Fan work out to N workers, fan results back in.

```go
const numWorkers = 3
const numJobs    = 5

jobs    := make(chan int, numJobs)
results := make(chan int, numJobs)
var wg sync.WaitGroup

// Spawn workers
for w := 1; w <= numWorkers; w++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for j := range jobs {
            results <- j * j
        }
    }()
}

// Send jobs, then signal "no more jobs"
for j := 1; j <= numJobs; j++ {
    jobs <- j
}
close(jobs)

// Once all workers done, close results so main's range exits
go func() {
    wg.Wait()
    close(results)
}()

// Collect
for r := range results {
    fmt.Println("result:", r)
}
```

Two channels, one WaitGroup. The "closer goroutine" pattern (`go func() { wg.Wait(); close(results) }()`) is the standard way to close a result channel once all senders are done.

## Common pitfalls

| Symptom | Cause |
|---|---|
| `fatal error: all goroutines are asleep - deadlock!` | Send/receive with no counterpart |
| `panic: send on closed channel` | Sender wrote after `close()` |
| `panic: close of closed channel` | `close()` called twice |
| Goroutine hangs silently | Send/receive on a `nil` channel — blocks FOREVER |
| `range` loop never exits | Forgot to close the channel |
| Result channel range exits early | Closed BEFORE all workers finished sending |

### The nil-channel trick

```go
var ch chan int    // nil channel
// ch <- 1         // BLOCKS forever
// <-ch            // BLOCKS forever
```

This sounds useless, but in a `select`, a nil case is effectively *disabled* — useful for dynamically turning off branches:

```go
select {
case v := <-ch1:
    // ...
case v := <-ch2:    // if you set ch2 = nil, this case is ignored
    // ...
}
```

## Quick reference

```go
// Create
ch := make(chan T)        // unbuffered
ch := make(chan T, N)     // buffered, capacity N

// Send / receive
ch <- v                   // send
v := <-ch                 // receive
v, ok := <-ch             // ok=false if closed and drained

// Close (sender side only!)
close(ch)

// Range
for v := range ch { ... } // until ch is closed

// Direction
chan<- T                  // send-only
<-chan T                  // receive-only

// Multiplexing
select {
case v := <-ch1: ...
case ch2 <- v:   ...
case <-time.After(d): ...
default: ...              // non-blocking
}

// Stats
len(ch)   // items currently buffered
cap(ch)   // buffer capacity (0 for unbuffered)
```

## When to use goroutines and channels

### The single most important question

Before reaching for goroutines or channels, ask:

> **"What am I waiting for?"**

- Waiting for the **network or disk** → goroutines almost always help.
- Waiting for the **CPU** to finish math → goroutines help ONLY if you have more cores than you're using.
- Not waiting at all (everything is fast and sequential) → don't add concurrency. Synchronous code is shorter, easier to debug, and easier to reason about.

And before reaching for a channel:

> **"Am I moving a value between goroutines, or signaling an event?"**

- Yes → channel.
- No → probably a mutex (or just a variable).

### When goroutines are a good fit

| Scenario | Why |
|---|---|
| Multiple network/API calls in parallel | Each goroutine waits on I/O; total time = slowest call, not sum |
| HTTP server handling requests | `net/http` already gives each request its own goroutine |
| Background work (logging, metrics, cleanup) | Runs alongside main flow without blocking it |
| Producer → consumer pipelines | Each stage progresses independently |
| Fan-out over a worker pool | Saturate available CPU/IO without per-item overhead |

### When goroutines are a BAD fit

| Anti-pattern | Why |
|---|---|
| Pure CPU work with thousands of goroutines | Scheduling overhead > benefit; cap at `runtime.NumCPU()` |
| Sequential dependencies (B needs A's result) | No parallelism possible — adds complexity for no speedup |
| Tiny per-iteration work | Goroutine setup cost exceeds the work itself |
| "It feels more Go-y" | Premature concurrency = the new premature optimization |

### When channels are a good fit

| Scenario | Pattern |
|---|---|
| Worker pool | One `jobs` channel, one `results` channel |
| Signal "done" / "stop" | `close(done)` — many goroutines can wait on it |
| Streaming pipeline | Channels connect each stage; data flows, never held in full |
| Timeout / cancellation | `select` with `<-time.After(d)` |
| Rate limiting / throttling | Buffered channel as token bucket |
| Pool of reusable resources | Buffered channel of connections / buffers |

### When channels are a BAD fit (use a mutex instead)

| Anti-pattern | Better tool |
|---|---|
| Single goroutine touches the data | No concurrency primitive at all — plain variable |
| Shared counter / map / cache | `sync.Mutex` (or `sync/atomic` for a single int) |
| Two-step "send then immediately receive" with no other goroutine involved | Just a function call |

### Channels vs mutexes — picking by problem shape

| Channels are best for... | Mutexes are best for... |
|---|---|
| Passing OWNERSHIP of values between goroutines | Protecting shared MUTABLE state in place |
| Coordinating WORK (job queues, pipelines, fan-out) | Counters, caches, structs many goroutines touch |
| Signaling events ("done", "stop", "ready") | Tight critical sections |
| Implementing timeouts and cancellation | Single-int hot paths (often atomics beat both) |

Channels usually express *flow* more naturally; mutexes usually express *state*. Both are first-class in Go — use the one that makes the code read clearly.

### Decision tree

```
1. Do I actually need concurrency?
   ├─ No (single-threaded is fine, fast, simple)  → just write sync code
   └─ Yes → go to 2.

2. What kind of concurrency?
   ├─ One thing happens in the background      → 1 goroutine, no channel
   ├─ N independent things, gather all results → goroutines + WaitGroup
   ├─ Many goroutines update a shared variable → goroutines + Mutex/atomic
   ├─ Producer(s) → consumer(s) flow            → goroutines + channel
   ├─ Need to time-out / cancel / multiplex     → channel + select
   └─ Wait for the first of N events            → channel + select
```

### Real-world scenarios at a glance

| Scenario | Tool |
|---|---|
| HTTP server handling requests | Goroutines (auto, from `net/http`) |
| Hitting 50 microservices in parallel and aggregating | Goroutines + `sync.WaitGroup` (or `errgroup`) |
| Background log flusher every 5 seconds | 1 goroutine + `time.Ticker` |
| Processing a queue with N workers | Goroutines + 2 channels (jobs, results) |
| Cache shared between request handlers | `sync.Map` or a mutex-protected map (NO channel) |
| Database connection pool | Buffered channel of connections |
| "Finish within X seconds OR give up" | Channel + `select` with `time.After` |
| Stream a 10GB file through several transformations | Pipeline of goroutines connected by channels |
| Counter across goroutines | `sync/atomic` (NOT channels) |
| Single function does CPU math | Plain code — no concurrency |

### Common mistakes to avoid

1. **Goroutine without a way to wait for it** → silent data loss when main exits.
2. **Goroutines for tiny work** → overhead exceeds benefit.
3. **Channels where a mutex would do** → harder to read, more state, more potential deadlocks.
4. **Mutexes where a channel would do** → tangled signaling, hard-to-spot races.
5. **`time.Sleep` instead of synchronizing** → fragile, slow, often wrong.
6. **Unbounded goroutines** (`go foo()` in a loop with no limit) → can spawn thousands and exhaust memory/connections. Use a worker pool or a semaphore channel.

### Beginner's rule

> **Default to no concurrency. Add it only when you can answer "what specific wait am I eliminating?" in one sentence.**

That single discipline prevents 90% of the concurrency bugs people write while learning Go. Concurrency is a tool, not a goal — synchronous code that works is better than concurrent code that's subtly broken.
