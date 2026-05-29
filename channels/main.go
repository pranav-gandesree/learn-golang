package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ==========================================================
	// 1. WHAT IS A CHANNEL?
	// ==========================================================
	//
	// A channel is a TYPED CONDUIT for sending values between
	// goroutines. Think of it as a thread-safe pipe.
	//
	//   ch := make(chan int)    creates a channel of int
	//   ch <- 42                SEND  42 into ch     (arrow points INTO ch)
	//   v := <-ch               RECEIVE from ch      (arrow points OUT of ch)
	//
	// Channels are the Go-idiomatic way for goroutines to COMMUNICATE.
	// The Go proverb: "Don't communicate by sharing memory;
	//                  share memory by communicating."

	// ==========================================================
	// 2. UNBUFFERED CHANNELS — synchronous handoff
	// ==========================================================
	//
	// An unbuffered channel has NO storage. A send BLOCKS until
	// another goroutine is ready to receive, and vice versa.
	// It's a hand-off, not a mailbox.

	fmt.Println("--- 1. Unbuffered channel ---")
	ch := make(chan int) // capacity 0 = unbuffered

	go func() {
		fmt.Println("  goroutine: about to send 42")
		ch <- 42 // BLOCKS here until main reads
		fmt.Println("  goroutine: sent 42 — receiver took it")
	}()

	time.Sleep(50 * time.Millisecond) // let the goroutine start
	fmt.Println("  main: about to receive")
	val := <-ch
	fmt.Println("  main: got", val)

	// ==========================================================
	// 3. BUFFERED CHANNELS — async until full
	// ==========================================================
	//
	// `make(chan T, N)` creates a channel with a buffer of size N.
	// Sends DON'T block until the buffer is FULL.
	// Receives DON'T block until the buffer is EMPTY.

	fmt.Println("\n--- 2. Buffered channel ---")
	buf := make(chan string, 3)

	buf <- "first"  // doesn't block — buffer has room
	buf <- "second" // doesn't block
	buf <- "third"  // doesn't block — buffer now FULL
	// buf <- "fourth"  // would BLOCK forever (no receiver)

	fmt.Println("  buffer length:", len(buf), "of capacity", cap(buf))

	fmt.Println("  received:", <-buf) // first  (FIFO)
	fmt.Println("  received:", <-buf) // second
	fmt.Println("  received:", <-buf) // third

	// ==========================================================
	// 4. CLOSING A CHANNEL + range
	// ==========================================================
	//
	// close(ch) marks a channel as "no more values will be sent."
	// Receivers can use `for v := range ch` to loop until close.
	//
	// Rules:
	//   - Only the SENDER should close. NEVER close from receiver side.
	//   - Sending to a closed channel PANICS.
	//   - Receiving from a closed channel returns the ZERO VALUE
	//     immediately. Use `v, ok := <-ch` to detect close (ok=false).

	fmt.Println("\n--- 3. close() + range ---")
	nums := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			nums <- i * 10
		}
		close(nums) // signals "I'm done sending"
	}()

	for v := range nums { // loop exits automatically when channel is closed
		fmt.Println("  received:", v)
	}

	// ==========================================================
	// 5. CHANNEL DIRECTION (send-only / receive-only types)
	// ==========================================================
	//
	// You can restrict a channel parameter to send-only or
	// receive-only. Compiler enforces this — it documents intent
	// and prevents bugs.
	//
	//   chan<- T   send-only   (arrow points INTO chan)
	//   <-chan T   receive-only (arrow points OUT of chan)

	fmt.Println("\n--- 4. Send-only / receive-only ---")
	pipe := make(chan int, 3)

	go produce(pipe)
	consume(pipe)

	// ==========================================================
	// 6. SELECT — wait on MULTIPLE channels
	// ==========================================================
	//
	// select is to channels what switch is to values. It blocks
	// until ONE of its cases can proceed, then runs that case.
	// If multiple are ready, one is chosen at random.

	fmt.Println("\n--- 5. select ---")
	chA := make(chan string)
	chB := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		chA <- "from A"
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		chB <- "from B"
	}()

	// Receive whichever arrives first, twice
	for i := 0; i < 2; i++ {
		select {
		case msg := <-chA:
			fmt.Println("  got:", msg)
		case msg := <-chB:
			fmt.Println("  got:", msg)
		case <-time.After(200 * time.Millisecond):
			fmt.Println("  timed out")
		}
	}

	// ==========================================================
	// 7. WORKER POOL — the classic concurrency pattern
	// ==========================================================
	//
	// Fan out work to N worker goroutines, collect results back.
	// Two channels: one for jobs, one for results.

	fmt.Println("\n--- 6. Worker pool (3 workers, 5 jobs) ---")
	const numWorkers = 3
	const numJobs = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start N workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // tell workers there are no more jobs

	// Closer goroutine: closes results once all workers are done.
	// This lets the main goroutine use `for r := range results`.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for r := range results {
		fmt.Println("  result:", r)
	}

	// ==========================================================
	// 8. COMMON PITFALLS
	// ==========================================================
	//
	// a) DEADLOCK: send/receive on a channel that no one else
	//    is reading/writing. Go detects this at runtime:
	//      fatal error: all goroutines are asleep - deadlock!
	//
	// b) PANIC: sending to a CLOSED channel panics.
	//      panic: send on closed channel
	//
	// c) PANIC: closing a channel twice panics.
	//
	// d) NIL CHANNEL: send/receive on a nil channel BLOCKS FOREVER
	//    (this is sometimes useful in select — disables a case).
	//
	// e) ONLY THE SENDER should close. The receiver doesn't know
	//    if other senders are still active.
	//
	// f) UNBUFFERED ≠ "buffer of 1". Unbuffered is a synchronous
	//    handshake; buffered (even size 1) has temporary storage.
}

// produce sends 3 ints into ch and closes it.
// The `chan<- int` type means this function can ONLY SEND.
// Calling `<-ch` here would be a compile error.
func produce(ch chan<- int) {
	for i := 1; i <= 3; i++ {
		ch <- i * 100
	}
	close(ch)
}

// consume receives from ch until it's closed.
// The `<-chan int` type means this function can ONLY RECEIVE.
// Calling `ch <- v` or `close(ch)` here would be a compile error.
func consume(ch <-chan int) {
	for v := range ch {
		fmt.Println("  consumed:", v)
	}
}

// worker reads jobs from `jobs`, processes them, and writes
// results to `results`. Calls wg.Done() when the jobs channel closes.
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(20 * time.Millisecond) // pretend to do work
		results <- j * j
		_ = id // (would normally log "worker %d processed job %d")
	}
}
