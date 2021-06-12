package main

import (
	"fmt"
)

func stringp(s string) *string {
	return &s
}

type person struct {
	FirstName  string
	MiddleName *string
	LastName   string
}

func failedUpdate(px *int) {
	x2 := 20
	px = &x2
}

func update(px *int) {
	*px = 20
}

//func MakeFoo() (Foo, error) {
//	f := Foo{
//		Field1: "val",
//		Field2: 20,
//	}
//	return f, nil
//}

func main() {
	//A pointer type is a type that represents a pointer. It is written with a * before a type name. A pointer type can be based on any type
	//x := 10
	//pointerToX := &x
	//fmt.Println(pointerToX)  // prints a memory address
	//fmt.Println(*pointerToX) // prints 10
	//
	//z := 5 + *pointerToX
	//fmt.Println(z)           // prints 15

	//x := "hello"
	//pointerToX := &x
	//fmt.Println(*pointerToX)

	//x := 10
	//var pointerToX *int
	//pointerToX = &x
	//fmt.Println(*pointerToX)

	//var x = new(int)
	//fmt.Println(x == nil) // prints false
	//fmt.Println(*x)       // prints 0

	p := person{
		FirstName:  "Pat",
		MiddleName: stringp("Perry"), // This works
		LastName:   "Peterson",
	}

	fmt.Println(p.FirstName)
	fmt.Println(*p.MiddleName)
	fmt.Println(p.LastName)

	x := 10
	failedUpdate(&x)
	fmt.Println(x) // prints 10
	update(&x)
	fmt.Println(x) // prints 20

	//The Unmarshal function populates a variable from a slice of bytes containing JSON.
	//f := struct {
	//	Name string `json:"name"`
	//	Age int `json:"age"`
	//}
	//err := json.Unmarshal([]byte(`{"name": "Bob", "age": 30}`), &f)

	//r = open_resource()
	//while r.has_data() {
	//	data_chunk = r.next_chunk()
	//	process(data_chunk)
	//}
	//close(r)

	//file, err := os.Open(fileName)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//data := make([]byte, 100)
	//for {
	//	count, err := file.Read(data)
	//	if err != nil {
	//		return err
	//	}
	//	if count == 0 {
	//		return nil
	//	}
	//	process(data[:count])
	//}

}