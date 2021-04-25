package main

import (
	"errors"
	"fmt"
	"os"
)

func div(numerator int, denominator int) int {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

//Example 5-1. Using a struct to simulate named parameters
type MyFuncOpts struct {
	FirstName string
	LastName string
	Age int
}
func MyFunc(opts MyFuncOpts) {
	fmt.Println(opts.FirstName)
	fmt.Println(opts.LastName)
	fmt.Println(opts.Age)
}

//variadic parameters  - Since the variadic parameter is converted to a slice, you can supply a slice as the input. However, you must put three dots (…) after the variable or slice literal.
func addTo(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func divAndRemainder(numerator int, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, numerator % denominator, nil
}

//Example 5-2. Multiple return values in Python are destructured tuples
func divAndRemainder2(numerator int, denominator int) (result int, remainder int, err error) {
	if denominator == 0 {
		err = errors.New("cannot divide by zero")
		return result, remainder, err
	}
	result, remainder = numerator/denominator, numerator%denominator
	return result, remainder, err
}

func divAndRemainder3(numerator, denominator int) (result int, remainder int, err error) {
	// assign some values
	result, remainder = 20, 30
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, numerator % denominator, nil
}

func divAndRemainder4(numerator, denominator int) (result int, remainder int, err error) {
	if denominator == 0 {
		err = errors.New("cannot divide by zero")
		return
	}
	result, remainder = numerator/denominator, numerator%denominator
	return
}

func main() {
	result := div(5, 2)
	fmt.Println(result)

	MyFunc(MyFuncOpts {
		LastName: "Patel",
		Age: 50,
	})

	MyFunc(MyFuncOpts {
		FirstName: "Joe",
		LastName: "Smith",
	})

	fmt.Println(addTo(3))
	fmt.Println(addTo(3, 2))
	fmt.Println(addTo(3, 2, 4, 6, 8))
	a := []int{4, 3}
	fmt.Println(addTo(3, a...))
	fmt.Println(addTo(3, []int{1, 2, 3, 4, 5}...))

	result, remainder, err := divAndRemainder(5, 2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result, remainder)

	//x, y, z := divAndRemainder2(5, 2)
	//fmt.Println(x, y, z)

	//x, y, z := divAndRemainder3(5, 2)
	//fmt.Println(x, y, z)

	x, y, z := divAndRemainder4(5, 2)
	fmt.Println(x, y, z)
}