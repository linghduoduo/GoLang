package adder

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_addNumbers(t *testing.T) {
	result := addNumbers(2,3)
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
	}
}

var testTime time.Time

func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")
	testTime = time.Now()
	exitVal := m.Run()
	fmt.Println("Clean up stuff after tests here")
	os.Exit(exitVal)
}

func TestFirst(t *testing.T) {
	fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
	fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}


// createFile is a helper function called from multiple tests
func createFile(t *testing.T) (string, error) {
	f, err := os.Create("tempFile")
	if err != nil {
		return "", err
	}
	// write some data to f
	t.Cleanup(func() {
		os.Remove(f.Name())
	})
	return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
	fName, err := createFile(t)
	if err != nil {
		t.Fatal(err)
	}
	// do testing, don't worry about cleanup
}
