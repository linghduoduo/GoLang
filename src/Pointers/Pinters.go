package main

import "fmt"

func main() {
	x := 10
	pointerToX := &x
	fmt.Println(pointerToX)  // prints a memory address
	fmt.Println(*pointerToX) // prints 10

	z := 5 + *pointerToX
	fmt.Println(z)           // prints 15
}