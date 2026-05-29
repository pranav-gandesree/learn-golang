package main

import "fmt"

func main() {
	intArr := [...]int32{1, 2, 3}
	fmt.Println(intArr)

	var intSlice []int32 = []int32{4, 5, 6}
	fmt.Printf("The length is %v with capacity %v", len(intSlice), cap(intSlice))
	intSlice = append(intSlice, 7)
	fmt.Printf("\nThe length is %v with capacity %v\n", len(intSlice), cap(intSlice))
	fmt.Println(intSlice[3])

	var intSlice2 []int32 = []int32{8, 8}
	intSlice = append(intSlice, intSlice2...)

	fmt.Printf("\nThe length is %v with capacity %v after slice 2 \n", len(intSlice), cap(intSlice))

	var intSlice3 []int32 = make([]int32, 5, 111)
	fmt.Printf("\nThe length is %v with capacity %v in slice 3 \n", len(intSlice3), cap(intSlice3))
	fmt.Println("intSlice3 values:", intSlice3)

	intSlice = append(intSlice, intSlice3...)
	fmt.Printf("\nintSlice after appending intSlice3: %v\n", intSlice)
	fmt.Printf("intSlice now: len=%d cap=%d\n", len(intSlice), cap(intSlice))
}
