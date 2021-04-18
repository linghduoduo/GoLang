package main

import "fmt"

// main is the entry point for the application.
func main() {
	var x = [3]int
	fmt.Println(x)

	//var x = [3]int{10, 20, 30}
	var x = [12]int{1, 5: 4, 6, 10: 100, 15}
	var x = [...]int{10, 20, 30}

	var x = [...]int{1, 2, 3}
	var y = [3]int{1, 2, 3}
	fmt.Println(x == y) // prints true
	fmt.Println(x)

	var x [2][3]int
	x[0] = [...]int{10, 20, 30}
	fmt.Println(x)

	var x []int
	x = append(x, 10)
	fmt.Println(x)

	var x = []int{1, 2, 3}
	x = append(x, 4)
	fmt.Println(x)

	var x []int
	x = append(x, 5, 6, 7)
	fmt.Println(x)

	y := []int{20, 30, 40}
	x = append(x, y...)
	fmt.Println(x)

	//Example 3-1. Understanding capacity
	var x []int
	fmt.Println(x, len(x), cap(x))
	x = append(x, 10)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 20)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 30)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 40)
	fmt.Println(x, len(x), cap(x))
	x = append(x, 50)
	fmt.Println(x, len(x), cap(x))

	x := make([]int, 5)
	x = append(x, 10)
	fmt.Println(x)

	x := make([]int, 5, 10)
	fmt.Println(x)

	x := make([]int, 0, 10)
	x = append(x, 5,6,7,8)

	//Example 3-3. Declaring a slice with default values
	data := []int{2, 4, 6, 8}
	fmt.Println(data)

	//Example 3-4. Slicing slices
	x := []int{1, 2, 3, 4}
	y := x[:2]
	z := x[1:]
	d := x[1:3]
	e := x[:]
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
	fmt.Println("d:", d)
	fmt.Println("e:", e)


	//Example 3-5. Slices with overlapping storage
	x := []int{1, 2, 3, 4}
	y := x[:2]
	z := x[1:]
	x[1] = 20
	y[0] = 10
	z[1] = 30
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)

	//Example 3-6. append makes overlapping slices more confusing
	x := []int{1, 2, 3, 4}
	y := x[:2]
	fmt.Println(cap(x), cap(y))
	y = append(y, 30)
	fmt.Println("x:", x)
	fmt.Println("y:", y)

	//Example 3-7. Even more confusing slices
	x := make([]int, 0, 5)
	x = append(x, 1, 2, 3, 4)
	y := x[:2]
	z := x[2:]
	fmt.Println(cap(x), cap(y), cap(z))
	y = append(y, 30, 40, 50)
	x = append(x, 60)
	z = append(z, 70)
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)

	//Example 3-8. The full slice expression protects against append
	x := make([]int, 0, 5)
	y := x[:2:2]
	z := x[2:4:4]
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)

	x := []int{1, 2, 3, 4}
	y := make([]int, 4)
	num := copy(y, x)
	fmt.Println(y, num)

	x := []int{1, 2, 3, 4}
	d := [4]int{5, 6, 7, 8}
	y := make([]int, 2)
	copy(y, d[:])
	fmt.Println(y)
	copy(d[:], x)
	fmt.Println(d)

	var a rune    = 'x'
	var s string  = string(a)
	fmt.Println(s)
	var b byte    = 'y'
	var s2 string = string(b)
	fmt.Println(s2)

	teams := map[string][]string {
		"Orcas": []string{"Fred", "Ralph", "Bijou"},
		"Lions": []string{"Sarah", "Peter", "Billie"},
		"Kittens": []string{"Waldo", "Raul", "Ze"},
	}
	fmt.Println(teams)

	//Example 3-10. Using a map
	totalWins := map[string]int{}
	totalWins["Orcas"] = 1
	totalWins["Lions"] = 2
	fmt.Println(totalWins["Orcas"])
	fmt.Println(totalWins["Kittens"])
	totalWins["Kittens"]++
	fmt.Println(totalWins["Kittens"])
	totalWins["Lions"] = 3
	fmt.Println(totalWins["Lions"])

	m := map[string]int{
		"hello": 5,
		"world": 0,
	}
	v, ok := m["hello"]
	fmt.Println(v, ok)

	v, ok = m["world"]
	fmt.Println(v, ok)

	v, ok = m["goodbye"]
	fmt.Println(v, ok)

	m := map[string]int{
		"hello": 5,
		"world": 10,
	}
	delete(m, "hello")
	fmt.Println(m)

	//Example 3-11. Using a map as a set
	intSet := map[int]bool{}
	vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10}
	for _, v := range vals {
		intSet[v] = true
	}
	fmt.Println(len(vals), len(intSet))
	fmt.Println(intSet[5])
	fmt.Println(intSet[500])
	if intSet[100] {
		fmt.Println("100 is in the set")
	}

	intSet := map[int]struct{}{}
	vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10}
	for _, v := range vals {
		intSet[v] = struct{}
	}
	if _, ok := intSet[5]; ok {
		fmt.Println("5 is in the set")
	}

	type firstPerson struct {
		name string
		age  int
	}
	type secondPerson struct {
		name string
		age  int
	}
	type thirdPerson struct {
		age  int
		name string
	}
	type fourthPerson struct {
		firstName string
		age       int
	}
	type fifthPerson struct {
		name          string
		age           int
		favoriteColor string
	}
	type firstPerson struct {
		name string
		age  int
	}
	f := firstPerson{
		name: "Bob",
		age:  50,
	}
	var g struct {
		name string
		age  int
	}

	//compiles -- can use = and == between identical named and anonymous structs
	g = f
	fmt.Println(f == g)
}
