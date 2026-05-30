package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	g := Greeter{Name: "world"}  // make a Greeter (fill its Name field)
	fmt.Fprintln(w, g.Message()) // call its method, write the result to the response
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// Greeter is a struct — it holds ONLY data (fields). No methods inside.
type Greeter struct {
	Name string
}

// Message is a method ON Greeter. It lives OUTSIDE the struct.
// The "(g Greeter)" part is the receiver: it says "this method belongs to Greeter,
// and inside here I can refer to that Greeter as g".
func (g Greeter) Message() string {
	return "hello " + g.Name
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}
