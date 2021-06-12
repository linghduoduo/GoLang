package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//The Done Channel Pattern
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

//Using a Cancel Function to Terminate a Goroutine
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

//When to Use Buffered and Unbuffered Channels
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

//Backpressure
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

//Turning Off a case in a select
// in and in2 are channels, done is a done channel.
//for {
//select {
//case v, ok := <-in:
//if !ok {
//in = nil // the case will never succeed again!
//continue
//}
//// process the v that was read from in
//case v, ok := <-in2:
//if !ok {
//in2 = nil // the case will never succeed again!
//continue
//}
//// process the v that was read from in2
//case <-done:
//return
//}
//}

//How to Time Out Code
func timeLimit() (int, error) {
	var result int
	var err error
	done := make(chan struct{})
	go func() {
		result, err = fmt.Println("Do something")
		close(done)
	}()
	select {
	case <-done:
		return result, err
	case <-time.After(2 * time.Second):
		return 0, errors.New("work timed out")
	}
}

func processAndGather(in <-chan int, processor func(int) int, num int) []int {
	out := make(chan int, num)
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for v := range in {
				out <- processor(v)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	var result []int
	for v := range out {
		result = append(result, v)
	}
	return result
}


type SlowComplicatedParser interface {
	Parse(string) string
}

//Running Code Exactly Once
var parser SlowComplicatedParser
var once sync.Once

//func Parse(dataToParse string) string {
//	once.Do(func() {
//		parser = initParser()
//	})
//	return parser.Parse(dataToParse)
//}
//
//func initParser() SlowComplicatedParser {
//	// do all sorts of setup and loading here
//	fmt.Println("do all sorts of setup and loading here")
//}


//Putting Our Concurrent Tools Together
//type processor struct {
//	outA chan AOut
//	outB chan BOut
//	outC chan COut
//	inC  chan CIn
//	errs chan error
//}
//
//func GatherAndProcess(ctx context.Context, data Input) (COut, error) {
//	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
//	defer cancel()
//	p := processor{
//		outA: make(chan AOut, 1),
//		outB: make(chan BOut, 1),
//		inC:  make(chan CIn, 1),
//		outC: make(chan COut, 1),
//		errs: make(chan error, 2),
//	}
//	p.launch(ctx, data)
//	inputC, err := p.waitForAB(ctx)
//	if err != nil {
//		return COut{}, err
//	}
//	p.inC <- inputC
//	out, err := p.waitForC(ctx)
//	return out, err
//}
//
//func (p *processor) launch(ctx context.Context, data Input) {
//	go func() {
//		aOut, err := getResultA(ctx, data.A)
//		if err != nil {
//			p.errs <- err
//			return
//		}
//		p.outA <- aOut
//	}()
//	go func() {
//		bOut, err := getResultB(ctx, data.B)
//		if err != nil {
//			p.errs <- err
//			return
//		}
//		p.outB <- bOut
//	}()
//	go func() {
//		select {
//		case <-ctx.Done():
//			return
//		case inputC := <-p.inC:
//			cOut, err := getResultC(ctx, inputC)
//			if err != nil {
//				p.errs <- err
//				return
//			}
//			p.outC <- cOut
//		}
//	}()
//}
//
//func (p *processor) waitForAB(ctx context.Context) (CIn, error) {
//	var inputC CIn
//	count := 0
//	for count < 2 {
//		select {
//		case a := <-p.outA:
//			inputC.A = a
//			count++
//		case b := <-p.outB:
//			inputC.B = b
//			count++
//		case err := <-p.errs:
//			return CIn{}, err
//		case <-ctx.Done():
//			return CIn{}, ctx.Err()
//		}
//	}
//	return inputC, nil
//}
//
//func (p *processor) waitForC(ctx context.Context) (COut, error) {
//	select {
//	case out := <-p.outC:
//		return out, nil
//	case err := <-p.errs:
//		return COut{}, err
//	case <-ctx.Done():
//		return COut{}, ctx.Err()
//	}
//}


//When to Use Mutexes Instead of Channels
func scoreboardManager(in <-chan func(map[string]int), done <-chan struct{}) {
	scoreboard := map[string]int{}
	for {
		select {
		case <-done:
			return
		case f := <-in:
			f(scoreboard)
		}
	}
}

type ChannelScoreboardManager chan func(map[string]int)

func NewChannelScoreboardManager() (ChannelScoreboardManager, func()) {
	ch := make(ChannelScoreboardManager)
	done := make(chan struct{})
	go scoreboardManager(ch, done)
	return ch, func() {
		close(done)
	}
}

func (csm ChannelScoreboardManager) Update(name string, val int) {
	csm <- func(m map[string]int) {
		m[name] = val
	}
}

func (csm ChannelScoreboardManager) Read(name string) (int, bool) {
	var out int
	var ok bool
	done := make(chan struct{})
	csm <- func(m map[string]int) {
		out, ok = m[name]
		close(done)
	}
	<-done
	return out, ok
}

type MutexScoreboardManager struct {
	l          sync.RWMutex
	scoreboard map[string]int
}

func NewMutexScoreboardManager() *MutexScoreboardManager {
	return &MutexScoreboardManager{
		scoreboard: map[string]int{},
	}
}

func (msm *MutexScoreboardManager) Update(name string, val int) {
	msm.l.Lock()
	defer msm.l.Unlock()
	msm.scoreboard[name] = val
}

func (msm *MutexScoreboardManager) Read(name string) (int, bool) {
	msm.l.RLock()
	defer msm.l.RUnlock()
	val, ok := msm.scoreboard[name]
	return val, ok
}

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

	//Using WaitGroups
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		fmt.Println("doThing1()")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("doThing2()")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("doThing3()")
	}()
	wg.Wait()

}