package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
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

func add(i int, j int) int { return i + j }

func sub(i int, j int) int { return i - j }

func mul(i int, j int) int { return i * j }

func div2(i int, j int) int { return i / j }

var opMap = map[string]func(int, int) int{
	"+": add,
	"-": sub,
	"*": mul,
	"/": div2,
}

//Returning Functions from Functions
func makeMult(base int) func(int) int {
	return func(factor int) int {
		return base * factor
	}
}
func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string)(err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
	return err
	}
	defer func() {
	if err == nil {
	err = tx.Commit()
	}
	if err != nil {
	tx.Rollback()
	}
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
	if err != nil {
	return err
	}
	// use tx to do more database inserts here
	return nil
}

func getFile(name string) (*os.File, func(), error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	return file, func() {
		file.Close()
	}, err
}
//f, closer, err := getFile(os.Args[1])
//if err != nil {
//log.Fatal(err)
//}
//defer closer()

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

	expressions := [][]string{
		[]string{"2", "+", "3"},
		[]string{"2", "-", "3"},
		[]string{"2", "*", "3"},
		[]string{"2", "/", "3"},
		[]string{"2", "%", "3"},
		[]string{"two", "+", "three"},
		[]string{"5"},
	}

	for _, expression := range expressions {
		if len(expression) != 3 {
			fmt.Println("invalid expression:", expression)
			continue
		}
		//the strconv.Atoi function in the standard library to convert a string to an int.
		p1, err := strconv.Atoi(expression[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		op := expression[1]
		opFunc, ok := opMap[op]
		if !ok {
			fmt.Println("unsupported operator:", op)
			continue
		}
		p2, err := strconv.Atoi(expression[2])
		if err != nil {
			fmt.Println(err)
			continue
		}
		result := opFunc(p1, p2)
		fmt.Println(result)
	}

	for i := 0; i < 5; i++ {
		func(j int) {
			fmt.Println("printing", j, "from inside of an anonymous function")
		}(i)
	}

	type Person struct {
		FirstName string
		LastName  string
		Age       int
	}

	//Passing Functions as Parameters
	people := []Person{
		{"Pat", "Patterson", 37},
		{"Tracy", "Bobbert", 23},
		{"Fred", "Fredson", 18},
	}
	fmt.Println(people)

	// sort by last name
	sort.Slice(people, func(i int, j int) bool {
		return people[i].LastName < people[j].LastName
	})
	fmt.Println(people)

	// sort by age
	sort.Slice(people, func(i int, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println(people)

	//Returning Functions from Functions
	twoBase := makeMult(2)
	threeBase := makeMult(3)
	for i := 0; i < 3; i++ {
		fmt.Println(twoBase(i), threeBase(i))
	}

	//Go Is Call By Value

}

//defer
//func main() {
//	if len(os.Args) < 2 {
//		log.Fatal("no file specified")
//	}
//	f, err := os.Open(os.Args[1])
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//	data := make([]byte, 2048)
//	for {
//		count, err := f.Read(data)
//		os.Stdout.Write(data[:count])
//		if err != nil {
//			if err != io.EOF {
//				log.Fatal(err)
//			}
//			break
//		}
//	}
//}

