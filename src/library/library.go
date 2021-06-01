package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	
}
