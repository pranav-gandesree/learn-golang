package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// =========================================================
	// 1. WHAT IS A GOROUTINE?
	// =========================================================
	//
	// A goroutine is a lightweight, independently-scheduled function call.
	// You launch one with the `go` keyword: `go someFunction()`.
	//
	// Goroutines are NOT OS threads — they're much cheaper. You can run
	// thousands or millions of them. The Go runtime multiplexes goroutines
	// onto a small pool of real OS threads automatically.
	//
	// THE CATCH: when main() returns, the WHOLE PROGRAM EXITS — including
	// any goroutines still running. So you must explicitly WAIT for them
	// to finish before main returns. That's what sync.WaitGroup is for.

	// =========================================================
	// 2. ONE GOROUTINE + WaitGroup
	// =========================================================
	//
	// sync.WaitGroup is a counter that lets the main goroutine wait
	// for a set of other goroutines to finish.
	//
	//   wg.Add(n)  — increment counter by n  (call BEFORE launching)
	//   wg.Done()  — decrement counter by 1  (each goroutine calls when done)
	//   wg.Wait()  — block until counter reaches 0

	var wg sync.WaitGroup

	wg.Add(1) // about to launch 1 goroutine — increment counter
	go func() {
		defer wg.Done() // GUARANTEES Done() runs even if the function panics
		fmt.Println("[g1] starting work")
		time.Sleep(100 * time.Millisecond) // pretend to do work
		fmt.Println("[g1] finished")
	}()

	fmt.Println("main: waiting for g1...")
	wg.Wait() // blocks until g1 calls Done()
	fmt.Println("main: g1 has finished, continuing")

	// =========================================================
	// 3. MULTIPLE GOROUTINES RUNNING IN PARALLEL
	// =========================================================
	//
	// Three goroutines, each sleeping for different durations.
	// Total wall-clock time will be ~150ms (the longest), NOT
	// 50+100+150=300ms — because they run CONCURRENTLY.

	fmt.Println("\nlaunching 3 workers in parallel:")
	start := time.Now()
	wg.Add(3)

	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond)
		fmt.Println("  [A] done (slept 150ms)")
	}()
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  [B] done (slept 100ms)")
	}()
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		fmt.Println("  [C] done (slept 50ms)")
	}()

	wg.Wait()
	fmt.Printf("all 3 workers finished in %v (NOT 300ms — they ran in parallel)\n", time.Since(start))

	// =========================================================
	// 4. LOOP + GOROUTINES — the closure-capture gotcha
	// =========================================================
	//
	// Common pattern: launch one goroutine per loop iteration.
	// There's a classic trap. The SAFE pattern is to pass the
	// loop variable as an ARGUMENT to the goroutine function,
	// so each goroutine gets its own copy.

	fmt.Println("\nlaunching 5 numbered workers (passing i as argument):")
	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go func(n int) { // n is a fresh parameter for each goroutine
			defer wg.Done()
			fmt.Printf("  worker #%d running\n", n)
		}(i) // ← pass i at LAUNCH time, capturing its current value
	}
	wg.Wait()

	// If you wrote `go func() { fmt.Println(i) }()` instead (capturing i
	// by closure), in OLD Go (< 1.22) every goroutine would print "5" —
	// they'd all share the same `i` variable, which becomes 5 by the time
	// they actually run. Go 1.22+ fixed this for `for` loops, but
	// passing-as-argument works on EVERY Go version. Always use it.

	// =========================================================
	// 5. ORDER IS NON-DETERMINISTIC
	// =========================================================
	//
	// The order of "worker #N running" prints above will vary between
	// runs. The Go runtime schedules goroutines however it pleases —
	// never depend on a particular order without explicit synchronization.

	fmt.Println("\nrun the program multiple times — the worker order will change.")
}
