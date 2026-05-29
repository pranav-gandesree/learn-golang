package mathutils

import "fmt"

func Add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func Double(n int) int {
	return multiply(n, 2)
}

func Arrays() {
	var intArray [3]int32 = [3]int32{4, 5, 6}
	fmt.Println("Integer Array:", intArray)

	fmt.Println("addresses of intArray[0], intArray[1], intArray[2]:", &intArray[0], &intArray[1], &intArray[2])
}
