package main

import (
	"fmt"
	"time"
)

type Person struct {
	FirstName string
	LastName string
	Age int
}

type Score int
type Converter func(string)Score
type TeamScores map[string]Score

// Method declarations look just like function declarations, with one addition: the receiver specification.
func (p Person) String() string {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}

type Counter struct {
	total             int
	lastUpdated time.Time
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

func doUpdateWrong(c Counter) {
	c.Increment()
	fmt.Println("in doUpdateWrong:", c.String())
}

func doUpdateRight(c *Counter) {
	c.Increment()
	fmt.Println("in doUpdateRight:", c.String())
}

type IntTree struct {
	val         int
	left, right *IntTree
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		return &IntTree{val: val}
	}
	if val < it.val {
		it.left = it.left.Insert(val)
	} else if val > it.val {
		it.right = it.right.Insert(val)
	}
	return it
}

func (it *IntTree) Contains(val int) bool {
	switch {
	case it == nil:
		return false
	case val < it.val:
		return it.left.Contains(val)
	case val > it.val:
		return it.right.Contains(val)
	default:
		return true
	}
}

type Adder struct {
	start int
}

func (a Adder) AddTo(val int) int {
	return a.start + val
}

func main() {

	p := Person {
		FirstName: "Fred",
		LastName:"Fredson",
		Age: 52,
	}
	output := p.String()

	fmt.Println(output)

	//var c Counter
	//fmt.Println(c.String())
	//c.Increment()
	//fmt.Println(c.String())

	var c Counter
	doUpdateWrong(c)
	fmt.Println("in main:", c.String())
	doUpdateRight(&c)
	fmt.Println("in main:", c.String())

	var it *IntTree
	it = it.Insert(5)
	it = it.Insert(3)
	it = it.Insert(10)
	it = it.Insert(2)
	fmt.Println(it.Contains(2))  // true
	fmt.Println(it.Contains(12)) // false

	//Methods Are Functions Too
	myAdder := Adder{start: 10}
	fmt.Println(myAdder.AddTo(5)) // prints 15

	f1 := myAdder.AddTo
	fmt.Println(f1(10))

	f2 := Adder.AddTo
	fmt.Println(f2(myAdder, 15))  // prints 25

	//Any time your logic depends on values that are configured at startup or changed while your program is running, those values should be stored in a struct and that logic should be implemented as a method. If your logic only depends on the input parameters, then it should be a function.
	// assigning untyped constants is valid
	//var i int = 300
	//var s Score = 100
	//var hs HighScore = 200
	//hs = s                  // compilation error!
	//s = i                   // compilation error!
	//s = Score(i)            // ok
	//hs = HighScore(s)       // ok
	

}
