package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// The WebSocket library. It does the messy parts for us:
	// the "upgrade" handshake, splitting data into frames, etc.
	"github.com/coder/websocket"
)

// wsHandler is a NORMAL http handler at first.
// Every WebSocket connection STARTS LIFE as a plain HTTP GET request that carries
// special headers (Connection: Upgrade, Upgrade: websocket). websocket.Accept reads
// those headers, replies "ok let's upgrade", and from then on this same TCP socket is
// a two-way WebSocket connection instead of a one-shot request/response.
func wsHandler(w http.ResponseWriter, r *http.Request) {

	// Accept performs the handshake. If it succeeds we get back a *websocket.Conn:
	// the live connection we read from and write to.
	//
	// InsecureSkipVerify: true turns OFF the browser-origin check. A real server would
	// instead set OriginPatterns to whitelist allowed origins. For local learning this
	// is fine — just don't ship it.
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		// Handshake failed (e.g. a normal browser visit with no upgrade headers).
		// Accept already wrote an error response, so we just stop.
		fmt.Println("server: accept failed:", err)
		return
	}
	// CloseNow is the "just hang up" close — safe to defer so the socket is always
	// released when this function returns, no matter how we exit the loop below.
	defer c.CloseNow()

	fmt.Println("server: a client connected")
	defer fmt.Println("server: a client disconnected")

	// The connection can stay open a long time, but we don't want a single client to
	// hang forever. This context cancels the whole conversation after 10 minutes.
	// (Same context.WithTimeout pattern you used in ../context/main.go — it just wraps
	// a longer-lived connection here instead of one request.)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	// THE ECHO LOOP.
	// A WebSocket is a stream of messages, so we loop: read one, write it straight back.
	for {
		// Read blocks until the client sends a message (or the connection ends).
		// It returns:
		//   typ  - the message kind: websocket.MessageText or websocket.MessageBinary
		//   data - the raw bytes of that message
		//   err  - non-nil when the client closed the connection or something broke
		typ, data, err := c.Read(ctx)
		if err != nil {
			// Most common reason: the client closed the tab / went away.
			// That's normal — we just leave the loop, and the deferred closes run.
			fmt.Println("server: read ended:", err)
			return
		}

		fmt.Printf("server: got %q\n", data)

		// Write sends the SAME bytes back with the SAME message type. That's the "echo".
		// Write also blocks until the message is handed off (or fails).
		if err := c.Write(ctx, typ, data); err != nil {
			fmt.Println("server: write failed:", err)
			return
		}
	}
}

// homeHandler serves the chat web page. This is a PLAIN http handler — no websockets.
// The browser loads this HTML first; the JavaScript inside it is what then opens the
// websocket back to /ws. So one server serves both: the page AND the live connection.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// ServeFile reads index.html from disk and writes it as the response.
	// Note: the path is relative to where you RUN the program from (the project root),
	// which is why it's "websocket/index.html" and not just "index.html".
	http.ServeFile(w, r, "websocket/index.html")
}

func main() {
	http.HandleFunc("/", homeHandler) // the web page
	http.HandleFunc("/ws", wsHandler) // the websocket endpoint

	fmt.Println("server: listening on http://localhost:8090  (websocket at ws://localhost:8090/ws)")
	http.ListenAndServe(":8090", nil)
}
