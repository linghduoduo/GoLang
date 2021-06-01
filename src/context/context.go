package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func logic(ctx context.Context, info string) (string, error) {
	// do some interesting stuff here
	return "", nil
}

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		// wrap the context with stuff -- we'll see how soon!
		req = req.WithContext(ctx)
		handler.ServeHTTP(rw, req)
	})
}

func handler(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := req.ParseForm()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	data := req.FormValue("data")
	result, err := logic(ctx, data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write([]byte(result))
}

type ServiceCaller struct {
	client *http.Client
}

//func (sc ServiceCaller) callAnotherService(ctx context.Context, data string)
//(string, error) {
//req, err := http.NewRequest(http.MethodGet,
//"http://example.com?data="+data, nil)
//if err != nil {
//return "", err
//}
//req = req.WithContext(ctx)
//resp, err := sc.client.Do(req)
//if err != nil {
//return "", err
//}
//defer resp.Body.Close()
//if resp.StatusCode != http.StatusOK {
//return "", fmt.Errorf("Unexpected status code %d",
//resp.StatusCode)
//}
//// do the rest of the stuff to process the response
//id, err := processResponse(resp.Body)
//return id, err
//}]

//func longRunningThingManager(ctx context.Context, data string) (string, error) {
//	type wrapper struct {
//		result string
//		err    error
//	}
//	ch := make(chan wrapper, 1)
//	go func() {
//		// do the long running thing
//		result, err := longRunningThing(ctx, data)
//		ch <- wrapper{result, err}
//	}()
//	select {
//	case data := <-ch:
//		return data.result, data.err
//	case <-ctx.Done():
//		return "", ctx.Err()
//	}
//}

func main() {
	fmt.Println("Hello")

	//ss := slowServer()
	//defer ss.Close()
	//fs := fastServer()
	//defer fs.Close()

	//ctx := context.Background()
	//callBoth(ctx, os.Args[1], ss.URL, fs.URL)

	ctx := context.Background()
	parent, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	child, cancel2 := context.WithTimeout(parent, 3*time.Second)
	defer cancel2()
	start := time.Now()
	<-child.Done()
	end := time.Now()
	fmt.Println(end.Sub(start))


}
