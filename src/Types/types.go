package main

import (
	"errors"
	"fmt"
	"net/http"
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

type Employee struct {
	Name         string
	ID           string
}

func (e Employee) Description() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.ID)
}

type Manager struct {
	Employee
	Reports []Employee
}

func (m Manager) FindNewEmployees() []Employee {
	// do business logic
	return m.Reports
}

type Inner struct {
	X int
}

type Outer struct {
	Inner
	X int
}

//type Inner struct {
//	A int
//}

func (i Inner) IntPrinter(val int) string {
	return fmt.Sprintf("Inner: %d", val)
}

//func (i Inner) Double() string {
//	return i.IntPrinter(i.A * 2)
//}
//
//type Outer struct {
//	Inner
//	S string
//}
//
//func (o Outer) IntPrinter(val int) string {
//	return fmt.Sprintf("Outer: %d", val)
//}

type LogicProvider struct {}

//func (lp LogicProvider) Process(data string) string {
//	// business logic
//}

//Interfaces specify what callers need. The client code defines the interface to specify what functionality it requires.
//type Logic interface {
//	Process(data string) string
//}
//
//type Client struct{
//	L Logic
//}
//
//func(c Client) Program() {
//	// get data from somewhere
//	c.L.Process(data)
//}

//func process(r io.Reader) error


// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
//func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
//	// If the reader has a WriteTo method, use it to do the copy.
//	// Avoids an allocation and a copy.
//	if wt, ok := src.(WriterTo); ok {
//		return wt.WriteTo(dst)
//	}
//	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
//	if rt, ok := dst.(ReaderFrom); ok {
//		return rt.ReadFrom(src)
//	}
//	// function continues...
//}

//func ctxDriverStmtExec(ctx context.Context, si driver.Stmt,
//	nvdargs []driver.NamedValue) (driver.Result, error) {
//	if siCtx, is := si.(driver.StmtExecContext); is {
//		return siCtx.ExecContext(ctx, nvdargs)
//	}
//	// fallback code is here
//}

//func walkTree(t *treeNode) (int, error) {
//	switch val := t.val.(type) {
//	case nil:
//		return 0, errors.New("invalid expression")
//	case number:
//		// we know that t.val is of type number, so return the
//		// int value
//		return int(val), nil
//	case operator:
//		// we know that t.val is of type operator, so
//		// find the values of the left and right children, then
//		// call the process() method on operator to return the
//		// result of processing their values.
//		left, err := walkTree(t.lchild)
//		if err != nil {
//			return 0, err
//		}
//		right, err := walkTree(t.rchild)
//		if err != nil {
//			return 0, err
//		}
//		return val.process(left, right), nil
//	default:
//		// if a new treeVal type is defined, but walkTree wasn't updated
//		// to process it, this detects it
//		return 0, errors.New("unknown node type")
//	}
//}

func LogOutput(message string) {
	fmt.Println(message)
}

type SimpleDataStore struct {
	userData map[string]string
}

func (sds SimpleDataStore) UserNameForID(userID string) (string, bool) {
	name, ok := sds.userData[userID]
	return name, ok
}

func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[string]string{
			"1": "Fred",
			"2": "Mary",
			"3": "Pat",
		},
	}
}

type DataStore interface {
	UserNameForID(userID string) (string, bool)
}

type Logger interface {
	Log(message string)
}

type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

type SimpleLogic struct {
	l  Logger
	ds DataStore
}

func (sl SimpleLogic) SayHello(userID string) (string, error) {
	sl.l.Log("in SayHello for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Hello, " + name, nil
}

func (sl SimpleLogic) SayGoodbye(userID string) (string, error) {
	sl.l.Log("in SayGoodbye for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Goodbye, " + name, nil
}

func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{
		l:    l,
		ds: ds,
	}
}

type Logic interface {
	SayHello(userID string) (string, error)
}

type Controller struct {
	l     Logger
	logic Logic
}

func (c Controller) HandleGreeting(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In SayHello")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func NewController(l Logger, logic Logic) Controller {
	return Controller{
		l:     l,
		logic: logic,
	}
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

	type MailCategory int

	const (
		Uncategorized MailCategory = iota
		Personal
		Spam
		Social
		Advertisements
	)

	m := Manager{
		Employee: Employee{
			Name:         "Bob Bobson",
			ID:             "12345",
		},
		Reports: []Employee{},
	}
	fmt.Println(m.ID)            // prints 12345
	fmt.Println(m.Description()) // prints Bob Bobson (12345)

	o := Outer{
		Inner: Inner{
			X: 10,
		},
		X: 20,
	}
	fmt.Println(o.X)       // prints 20
	fmt.Println(o.Inner.X) // prints 10

	//c := Client{
	//	L: LogicProvider{},
	//}
	//c.Program()

	//r, err := os.Open(fileName)
	//if err != nil {
	//	return err
	//}
	//defer r.Close()
	//return process(r)
	//return nil

	//r, err := os.Open(fileName)
	//if err != nil {
	//	return err
	//}
	//defer r.Close()
	//gz, err = gzip.NewReader(r)
	//if err != nil {
	//	return err
	//}
	//defer gz.Close()
	//return process(gz)

	//Interfaces and nil
	var s *string
	fmt.Println(s == nil) // prints true
	var i interface{}
	fmt.Println(i == nil) // prints true
	i = s
	fmt.Println(i == nil) // prints false

	//One common use of the empty interface is as a placeholder for data of uncertain schema that’s read from an external source, like a JSON file
	// one set of braces for the interface{} type,
	// the other to instantiate an instance of the map
	//data := map[string]interface{}{}
	//contents, err := ioutil.ReadFile("testdata/sample.json")
	//if err != nil {
	//	return err
	//}
	//defer contents.Close()
	//json.Unmarshal(contents, &data)
	// the contents are now in the data map

	//Another use of interface{} is as a way to store a value in a user-created data structure.
	type LinkedList struct {
		Value interface{}
		Next    *LinkedList
	}

	//func (ll *LinkedList) Insert(pos int, val interface{}) *LinkedList {
	//if ll == nil || pos == 0 {
	//return &LinkedList{
	//Value: val,
	//Next:    ll,
	//}
	//}
	//ll.Next = ll.Next.Insert(pos-1, val)
	//return ll
	//}
	//
	//Type Assertions and Type Switches

	type MyInt int
	var ii interface{}
	var mine MyInt = 20
	ii = mine
	i2 := ii.(MyInt)
	fmt.Println(i2 + 1)


	//l := LoggerAdapter(LogOutput)
	//ds := NewSimpleDataStore()
	//logic := NewSimpleLogic(l, ds)
	//cc := NewController(l, logic)
	//http.HandleFunc("/hello", cc.SayHello)
	//http.ListenAndServe(":8080", nil)
}
