package main

import "fmt"

func loopsMain() {
	fmt.Println("Loops in Go")
	sum := 0
	for i := 0; i < 5; i++ {
		sum += i
		fmt.Printf("For loop iteration: %d, Sum: %d\n", i, sum)
	}

}
