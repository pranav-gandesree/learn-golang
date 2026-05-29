package main

import "fmt"

func main() {
	fmt.Println("practice: channels")
	ch := make(chan int) // UNBUFFERED channel — sends BLOCK until someone is ready to receive

	// WHY A GOROUTINE HERE? Because this channel is unbuffered.
	//
	// If we wrote `writetochannel(ch)` (no `go`), main would call
	// the function, which would try `ch <- 0` on the FIRST iteration.
	// That send blocks until SOMEONE receives... but the only
	// receiver is the `for i := range ch` loop BELOW — which never
	// gets to run because main is stuck on the first send.
	//
	//   → DEADLOCK: "all goroutines are asleep"
	//
	// By writing `go writetochannel(ch)`, we run the sender in a
	// SEPARATE goroutine. The sender blocks on its first send,
	// but main keeps going, reaches the `for i := range ch` loop,
	// receives the value, the sender unblocks, sends the next, etc.
	//
	// Channels need TWO goroutines to flow:
	//   one to send, one to receive — running in parallel.
	go writetochannel(ch)

	// receiver loop — runs in MAIN goroutine.
	// `range ch` keeps receiving until the channel is CLOSED.
	for i := range ch {
		fmt.Println(i)
	}

	// EXCEPTION: if you make the channel BUFFERED with enough room
	// for all sends, e.g. `ch := make(chan int, 5)`, then you COULD
	// write to it from main without a goroutine — all 5 sends would
	// fit in the buffer without blocking. But once the buffer fills,
	// you're back to needing a separate goroutine to drain it.
}

// writetochannel runs as a GOROUTINE (launched with `go` in main).
// Each `ch <- value` blocks until main's range loop receives the
// previous value — sender and receiver hand off one value at a time.
func writetochannel(ch chan int) {
	defer close(ch) // signals "no more values" so main's `range ch` loop exits

	for i := 0; i < 5; i++ {
		ch <- i + 5 // BLOCKS here until main reads — that's the handshake
	}
	// close(ch) // same as the defer above; only need ONE close
}
