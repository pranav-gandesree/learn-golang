package main

import "fmt"

func main() {
	s := []int{}
	fmt.Printf("start: len=%d cap=%d\n", len(s), cap(s))

	for i := 1; i <= 20; i++ {
		s = append(s, i)
		fmt.Printf("after append %2d: len=%d cap=%d\n", i, len(s), cap(s))
	}
}
