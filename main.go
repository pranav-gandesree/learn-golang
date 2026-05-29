package main

import (
	"fmt"

	"learning/mathutils"
)

func main() {
	helloMain()
	loopsMain()

	fmt.Println("Add(2, 3) =", mathutils.Add(2, 3))
	fmt.Println("Double(5) =", mathutils.Double(5))

	mathutils.Arrays()
	fmt.Println("Called mathutils.Arrays()")
}
