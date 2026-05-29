package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ==========================================================
	// 1. WHY MUTEXES? — THE DATA RACE PROBLEM
	// ==========================================================
	//
	// When multiple goroutines READ AND WRITE the same memory
	// without coordination, they "race" each other and the result
	// is unpredictable. This is a DATA RACE.
	//
	// Concretely: `counter++` is NOT a single atomic operation.
	// At the CPU level it takes three steps:
	//   1. read counter from memory into a register
	//   2. add 1 to the register
	//   3. write the register back to memory
	//
	// If two goroutines interleave these steps, one increment can
	// be LOST — both read the same old value, both write the same
	// new value, one update silently overwrites the other.

	// --- Demo: 1000 goroutines each increment a counter 100 times.
	// --- EXPECTED result: 100,000. ACTUAL (without locking): less.

	var unsafeCounter int
	var wg sync.WaitGroup

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				unsafeCounter++ // ← THREE STEPS, NOT ATOMIC
			}
		}()
	}
	wg.Wait()
	fmt.Printf("UNSAFE counter: expected 100000, got %6d (lost %d updates)\n",
		unsafeCounter, 100000-unsafeCounter)

	// Try running this with `go run -race ./mutexes` — Go's built-in
	// race detector will print a detailed report showing exactly which
	// lines raced. It's the BEST way to find data races in real code.

	// ==========================================================
	// 2. FIX IT WITH sync.Mutex
	// ==========================================================
	//
	// A mutex (MUTual EXclusion lock) lets only ONE goroutine into
	// a "critical section" at a time. Two methods:
	//   mu.Lock()    — wait until no one holds the lock, then take it
	//   mu.Unlock()  — release the lock so someone else can take it
	//
	// PATTERN: always `defer mu.Unlock()` right after `mu.Lock()`,
	// so the lock is released even if the function panics. Forgetting
	// to unlock = lock held forever = every other goroutine deadlocks.

	var safeCounter int
	var mu sync.Mutex

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				safeCounter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("SAFE counter:   expected 100000, got %6d\n", safeCounter)

	// ==========================================================
	// 3. sync.RWMutex — MANY READERS, ONE WRITER
	// ==========================================================
	//
	// If your data is read OFTEN but written RARELY, a regular Mutex
	// is wasteful — readers don't interfere with each other but with
	// a Mutex they'd still queue up behind each other.
	//
	// sync.RWMutex has FOUR methods:
	//   rw.Lock()    / rw.Unlock()    — exclusive WRITER lock
	//   rw.RLock()   / rw.RUnlock()   — shared READER lock
	//
	// Many goroutines can hold RLock SIMULTANEOUSLY. But Lock blocks
	// until ALL readers AND writers are done, and prevents new ones
	// from starting.

	fmt.Println("\n--- RWMutex demo: 5 readers + 1 writer ---")
	cache := map[string]string{"hello": "world"}
	var rw sync.RWMutex

	wg.Add(6)
	for i := 1; i <= 5; i++ {
		go func(id int) {
			defer wg.Done()
			rw.RLock() // shared read lock — many can hold this at once
			defer rw.RUnlock()
			val := cache["hello"]
			fmt.Printf("  reader %d sees hello=%q\n", id, val)
			time.Sleep(20 * time.Millisecond)
		}(i)
	}
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond) // give readers a head start
		rw.Lock()                         // exclusive write lock
		defer rw.Unlock()
		cache["hello"] = "updated"
		fmt.Println("  writer changed hello -> \"updated\"")
	}()
	wg.Wait()

	// ==========================================================
	// 4. COMMON PITFALLS
	// ==========================================================
	//
	// a) NEVER COPY A MUTEX. Always share via pointer.
	//    A copied mutex is a SECOND, INDEPENDENT lock — the original
	//    and the copy don't see each other. `go vet` warns about this.
	//
	// b) ALWAYS `defer mu.Unlock()`.
	//    Without it, a panic or early `return` leaves the lock held
	//    forever — every other goroutine that tries to Lock blocks
	//    forever.
	//
	// c) sync.Mutex is NOT REENTRANT.
	//    If the same goroutine calls Lock() twice without Unlock(),
	//    it deadlocks against ITSELF. Don't recursively lock.
	//
	// d) KEEP CRITICAL SECTIONS SMALL.
	//    Don't make network calls or file I/O while holding a lock —
	//    you'll serialize the whole program behind that slow thing.
	//
	// e) FOR A SINGLE INT, sync/atomic IS FASTER.
	//      import "sync/atomic"
	//      var counter int64
	//      atomic.AddInt64(&counter, 1)
	//    For anything richer (map, struct), use a mutex.
}
