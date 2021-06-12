package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

func seedRand() *rand.Rand {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		panic("cannot seed with cryptographic random number generator")
	}
	r := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
	return r
}

type Foo struct {
	x int
	S string
}

func (f Foo) Hello() string {
	return "hello"
}

func (f Foo) goodbye() string {
	return "goodbye"
}

type Bar = Foo
func MakeBar() Bar {
	bar := Bar{
		x: 20,
		S: "Hello",
	}
	var f Foo = bar
	fmt.Println(f.Hello())
	fmt.Println(f.goodbye())
	return bar
}

func main() {
	fmt.Println(seedRand())

	testFoo := Foo{
		100,
		"Test",
	}
	fmt.Println(testFoo.x)
	fmt.Println(testFoo.S)
	fmt.Println("-------")

	fmt.Println(MakeBar())
	fmt.Println("-------")

	fmt.Println(MakeBar().S)
	fmt.Println(MakeBar().goodbye())
}
