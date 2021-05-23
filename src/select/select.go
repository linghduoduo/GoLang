package main

import (
	"errors"
	"fmt"
	"time"
)

func searchData(s string, searchers []func(string) []string) []string {
	done := make(chan struct{})
	result := make(chan []string)
	for _, searcher := range searchers {
		go func(searcher func(string) []string) {
			select {
			case result <- searcher(s):
			case <-done:
			}
		}(searcher)
	}
	r := <-result
	close(done)
	return r
}

func countTo(max int) (<-chan int, func()) {
	ch := make(chan int)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}
	go func() {
		for i := 0; i < max; i++ {
			select {
			case <-done:
				return
			case ch<-i:

			}
		}
		close(ch)
	}()
	return ch, cancel
}

func process(val int) int {
	fmt.Println(val)
	return val*2
}

func processChannel(ch chan int) []int {
	const conc = 10
	results := make(chan int, conc)
	for i := 0; i < conc; i++ {
		go func() {
			v := <- ch
			results <- process(v)
		}()
	}
	var out []int
	for i := 0; i < conc; i++ {
		out = append(out, <-results)
	}
	return out
}

type PressureGauge struct {
	ch chan struct{}
}

func New(limit int) *PressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}
	return &PressureGauge{
		ch: ch,
	}
}

func (pg *PressureGauge) Process(f func()) error {
	select {
	case <-pg.ch:
		f()
		pg.ch <- struct{}{}
		return nil
	default:
		return errors.New("no more capacity")
	}
}

func doThingThatShouldBeLimited() string {
	time.Sleep(2 * time.Second)
	return "done"
}

//func timeLimit() (int, error) {
//	var result int
//	var err error
//	done := make(chan struct{})
//	go func() {
//		result, err = doSomeWork()
//		close(done)
//	}()
//	select {
//	case <-done:
//		return result, err
//	case <-time.After(2 * time.Second):
//		return 0, errors.New("work timed out")
//	}
//}


func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		v := 1
		ch1 <- v
		v2 := <-ch2
		fmt.Println(v, v2)
	}()
	v := 2
	var v2 int
	select {
	case ch2 <- v:
	case v2 = <-ch1:
	}
	fmt.Println(v, v2)

	select {
	case v := <-ch1:
		fmt.Println("read from ch:", v)
	default:
		fmt.Println("no value written to ch")
	}

	a := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(a))
	for _, v := range a {
		go func(val int) {
			ch <- val * 2
		}(v)
	}
	for i := 0; i < len(a); i++ {
		fmt.Println(<-ch)
	}

	ch3, cancel := countTo(10)
	for i := range ch3 {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel()

	//pg := New(10)
	//http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
	//	err := pg.Process(func() {
	//		w.Write([]byte(doThingThatShouldBeLimited()))
	//	})
	//	if err != nil {
	//		w.WriteHeader(http.StatusTooManyRequests)
	//		w.Write([]byte("Too many requests"))
	//	}
	//})
	//http.ListenAndServe(":8080", nil)

	//for {
	//	select {
	//	case v, ok := <-in:
	//		if !ok {
	//			in = nil // the case will never succeed again!
	//			continue
	//		}
	//		// process the v that was read from in
	//	case v, ok := <-in2:
	//		if !ok {
	//			in2 = nil // the case will never succeed again!
	//			continue
	//		}
	//		// process the v that was read from in2
	//	case <-done:
	//		return
	//	}
	//}

}