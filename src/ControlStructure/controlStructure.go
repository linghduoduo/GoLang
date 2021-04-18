package main

import (
	"fmt"
	"math/rand"
)

func main() {

	//Example 4-1. Shadowing variables
	//x := 10
	//if x > 5 {
	//	fmt.Println(x)
	//	x := 5
	//	fmt.Println(x)
	//}
	//fmt.Println(x)

	//Example 4-2. Shadowing with multiple assignment
	//x := 10
	//if x > 5 {
	//	x, y := 5, 20
	//	fmt.Println(x, y)
	//}
	//fmt.Println(x)

	//Example 4-3. Shadowing package names
	//x := 10
	//fmt.Println(x)
	//fmt := "oops"
	//fmt.Println(fmt)

	//Example 4-5. if and else
	n := rand.Intn(10)
	if n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big:", n)
	} else {
		fmt.Println("That's a good number:", n)
	}

	//Example 4-6. Scoping a variable to an if statement
	if n := rand.Intn(10); n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big:", n)
	} else {
		fmt.Println("That's a good number:", n)
	}

	//Example 4-7. Out of scope…
	if n := rand.Intn(10); n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big:", n)
	} else {
		fmt.Println("That's a good number:", n)
	}
	fmt.Println(n)

	//Example 4-8. A complete for statement
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	//Example 4-9. A condition-only for statement
	i := 1
	for i < 100 {
		fmt.Println(i)
		i = i * 2
	}

	//Example 4-11. Confusing code
	for i := 1; i <= 100; i++ {
		if i%3 == 0 {
			if i%5 == 0 {
				fmt.Println("FizzBuzz")
			} else {
				fmt.Println("Fizz")
			}
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}

	//Example 4-12. Using continue to make code clearer
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
			continue
		}
		if i%3 == 0 {
			fmt.Println("Fizz")
			continue
		}
		if i%5 == 0 {
			fmt.Println("Buzz")
			continue
		}
		fmt.Println(i)
	}

	//Example 4-13. The for-range loop
	//evenVals := []int{2, 4, 6, 8, 10, 12}
	//for i, v := range evenVals {
	//	fmt.Println(i, v)
	//}

	//Example 4-14. Ignoring the key in a for-range loop
	//evenVals := []int{2, 4, 6, 8, 10, 12}
	//for _, v := range evenVals {
	//	fmt.Println(v)
	//}

	uniqueNames := map[string]bool{"Fred": true, "Raul": true, "Wilma": true}
	for k := range uniqueNames {
		fmt.Println(k)
	}

	//Example 4-15. Map iteration order varies
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}
	for i := 0; i < 3; i++ {
		fmt.Println("Loop", i)
		for k, v := range m {
			fmt.Println(k, v)
		}
	}

	//Example 4-16. Iterating over strings
	samples := []string{"hello", "apple_π!"}
	for _, sample := range samples {
		for i, r := range sample {
			fmt.Println(i, r, string(r))
		}
		fmt.Println()
	}

	//Example 4-17. Modifying the value doesn’t modify the source
	//evenVals := []int{2, 4, 6, 8, 10, 12}
	//for _, v := range evenVals {
	//	v *= 2
	//}
	//fmt.Println(evenVals)

	//main(){
	//	samples := []string{"hello", "apple_π!"}
	//outer:
	//	for _, sample := range samples {
	//		for i, r := range sample {
	//			fmt.Println(i, r, string(r))
	//			if r == 'l' {
	//				continue outer
	//			}
	//		}
	//		fmt.Println()
	//	}

	//evenVals := []int{2, 4, 6, 8, 10}
	//for i, v := range evenVals {
	//	if i == 0 {
	//		continue
	//	}
	//	if i == len(evenVals)-2 {
	//		break
	//	}
	//	fmt.Println(i, v)
	//}

	evenVals := []int{2, 4, 6, 8, 10}
	for i := 1; i < len(evenVals)-1; i++ {
		fmt.Println(i, evenVals[i])
	}

	//Example 4-19. The switch statement
	words := []string{"a", "cow", "smile", "gopher", "octopus", "anthropologist"}
	for _, word := range words {
		switch size := len(word); size {
		case 1, 2, 3, 4:
			fmt.Println(word, "is a short word!")
		case 5:
			wordLen := len(word)
			fmt.Println(word, "is exactly the right length:", wordLen)
		case 6, 7, 8, 9:
		default:
			fmt.Println(word, "is a long word!")
		}
	}


	for i := 0; i < 10; i++ {
		switch {
		case i%2 == 0:
			fmt.Println(i, "is even")
		case i%3 == 0:
			fmt.Println(i, "is divisible by 3 but not 2")
		case i%7 == 0:
			fmt.Println("exit the loop!")
			break
		default:
			fmt.Println(i, "is boring")
		}
	}

	//Example 4-21. The blank switch
	words2 := []string{"hi", "salutations", "hello"}
	for _, word := range words2 {
		switch wordLen := len(word); {
		case wordLen < 5:
			fmt.Println(word, "is a short word!")
		case wordLen > 10:
			fmt.Println(word, "is a long word!")
		default:
			fmt.Println(word, "is exactly the right length.")
		}
	}

	aa := rand.Int()
	switch {
	case aa == 2:
		fmt.Println("a is 2")
	case aa == 3:
		fmt.Println("a is 3")
	case aa == 4:
		fmt.Println("a is 4")
	default:
		fmt.Println("a is ", aa)
	}

	switch aa {
	case 2:
		fmt.Println("a is 2")
	case 3:
		fmt.Println("a is 3")
	case 4:
		fmt.Println("a is 4")
	default:
		fmt.Println("a is ", aa)
	}

	//Example 4-22. Rewriting if/else with a blank switch
	switch n := rand.Intn(10); {
	case n == 0:
		fmt.Println("That's too low")
	case n > 5:
		fmt.Println("That's too big:", n)
	default:
		fmt.Println("That's a good number:", n)
	}
	
	//Example 4-23. Go’s goto has rules
	//func main() {
	//	a := 10
	//	goto skip
	//	b := 20
	//skip:
	//	c := 30
	//	fmt.Println(a, b, c)
	//	if c > a {
	//		goto inner
	//	}
	//	if a < b {
	//	inner:
	//		fmt.Println("a is less than b")
	//	}
	//}

}
