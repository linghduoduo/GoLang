package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type NotHowReaderIsDefined interface {
	Read() (p []byte, err error)
}

func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048)
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, nil, err
	}
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

type Closer interface {
	Close() error
}

type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{r}
}

type HelloHandler struct{}

func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

var securityMsg = []byte("You didn't give the secret password\n")
func RequestTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("request time for %s: %v", r.URL.Path, end.Sub(start))
	})
}
func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func main() {
	//s := "The quick brown fox jumped over the lazy dog"
	//sr := strings.NewReader(s)
	//counts, err := countLetters(sr)
	//if err != nil {
	//	return
	//}
	//fmt.Println(counts)
	//
	//r, closer, err := buildGZipReader("my_data.txt.gz")
	//if err != nil {
	//	return
	//}
	//defer closer()
	//counts2, err := countLetters(r)
	//if err != nil {
	//	return
	//}
	//fmt.Println(counts2)
	//
	//fileName := "my_data.txt"
	//f, err := os.Open(fileName)
	//if err != nil {
	//	return
	//}
	//defer f.Close()

	d := 2 * time.Hour + 30 * time.Minute // d is of type time.Duration
	fmt.Println(d)

	t, err := time.Parse("2006-02-01 15:04:05 -0700", "2016-13-03 00:00:00 +0000")
	if err != nil {
		return
	}
	fmt.Println(t.Format("January 2, 2006 at 3:04:05PM MST"))

	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type Order struct {
		ID            string        `json:"id"`
		DateOrdered time.Time `json:"date_ordered"`
		CustomerID    string        `json:"customer_id"`
		Items         []Item        `json:"items"`
	}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	toFile := Person {
		Name: "Fred",
		Age:  40,
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "sample-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())
	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		panic(err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	//dec := json.NewDecoder(strings.NewReader(data))
	//for dec.More() {
	//	err := dec.Decode(&t)
	//	if err != nil {
	//		panic(err)
	//	}
	//	// process t
	//}

	//var b bytes.Buffer
	//enc := json.NewEncoder(&b)
	//for _, input := range allInputs {
	//	t := process(input)
	//	err = enc.Encode(t)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//out := b.String()

	type RFC822ZTime struct {
		time.Time
	}

	//func (rt RFC822ZTime) MarshalJSON() ([]byte, error) {
	//	out := rt.Time.Format(time.RFC822Z)
	//	return []byte(`"` + out + `"`), nil
	//}
	//
	//func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error {
	//	if string(b) == "null" {
	//	return nil
	//}
	//	t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
	//	if err != nil {
	//	return err
	//}
	//	*rt = RFC822ZTime{t}
	//	return nil
	//}

	//type Order struct {
	//	ID          string      `json:"id"`
	//	DateOrdered RFC822ZTime `json:"date_ordered"`
	//	CustomerID  string      `json:"customer_id"`
	//	Items       []Item      `json:"items"`
	//}

	a := make(map[int]string)
	a[1] = "asdf"
	a[-1] = "qwer"
	fmt.Println("Initial:     ",a)

	stuff, err := json.Marshal(a)
	fmt.Println("Serialized:  ", string(stuff), err)

	b := make(map[int]string)
	err = json.Unmarshal(stuff, &b)
	fmt.Println("Deserialized:", b, err)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-My-Client", "Learning Go")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}
	fmt.Println(res.Header.Get("Content-Type"))
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data)

	type ResponseWriter interface {
		Header() http.Header
		Write([]byte) (int, error)
		WriteHeader(statusCode int)
	}

	//s := http.Server{
	//	Addr:         ":8080",
	//	ReadTimeout:  30 * time.Second,
	//	WriteTimeout: 90 * time.Second,
	//	IdleTimeout:  120 * time.Second,
	//	Handler:      HelloHandler{},
	//}
	//err := s.ListenAndServe()
	//if err != nil {
	//	if err != http.ErrServerClosed {
	//		panic(err)
	//	}
	//}

	//mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello!\n"))
	//})

	//person := http.NewServeMux()
	//person.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("greetings!\n"))
	//})
	//dog := http.NewServeMux()
	//dog.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("good puppy!\n"))
	//})
	//mux := http.NewServeMux()
	//mux.Handle("/person/", http.StripPrefix("/person", person))
	//mux.Handle("/dog/", http.StripPrefix("/dog", dog))

	//terribleSecurity := TerribleSecurityProvider("GOPHER")
	//
	//mux.Handle("/hello", terribleSecurity(RequestTimer(
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("Hello!\n"))
	//	}))))

	//helloHandler := func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello!\n"))
	//}
	//chain := alice.New(terribleSecurity, RequestTimer).ThenFunc(helloHandler)
	//mux.Handle("/hello", chain)
}
