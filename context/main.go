package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {
	// Derive a context that auto-cancels 2 seconds from now.
	// req.Context() is the parent (cancels if the client disconnects);
	// the child also cancels on its own after 2s, whichever comes first.
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel() // ALWAYS call cancel to release the timer; defer = run on return

	fmt.Println("server: hello handler started")
	defer fmt.Println("server: hello handler ended")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "hello\n")
	case <-ctx.Done():

		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func main() {
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":3000", nil)
}
