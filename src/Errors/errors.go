package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
)

func calcRemainderAndMod(numerator, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("denominator is 0")
	}
	return numerator / denominator, numerator % denominator, nil
}

//Use Strings for Simple Errors
func doubleEven(i int) (int, error) {
	if i % 2 != 0 {
		return 0, errors.New("only even numbers are processed")
	}
	return i * 2, nil
}

//func doubleEven(i int) (int, error) {
//	if i % 2 != 0 {
//		return 0, fmt.Errorf("%d isn't an even number", i)
//	}
//	return i * 2, nil
//}

type Status int

const (
	InvalidLogin Status = iota + 1
	NotFound
)

//type StatusErr struct {
//	Status    Status
//	Message string
//}

//func (se StatusErr) Error() string {
//	return se.Message
//}

//func LoginAndGetData(uid, pwd, file string) ([]byte, error) {
//	err := login(uid, pwd)
//	if err != nil {
//		return nil, StatusErr{
//			Status:    InvalidLogin,
//			Message: fmt.Sprintf("invalid credentials for user %s", uid),
//		}
//	}
//	data, err := getData(file)
//	if err != nil {
//		return nil, StatusErr{
//			Status:    NotFound,
//			Message: fmt.Sprintf("file %s not found", file),
//		}
//	}
//	return data, nil
//}

//func GenerateError(flag bool) error {
//	var genErr StatusErr
//	if flag {
//		genErr = StatusErr{
//			Status: NotFound,
//		}
//	}
//	return genErr
//}

//func fileChecker(name string) error {
//	f, err := os.Open(name)
//	if err != nil {
//		return fmt.Errorf("in fileChecker: %w", err)
//	}
//	f.Close()
//	return nil
//}

//type StatusErr struct {
//	Status Status
//	Message string
//	err error
//	//In some cases, expecting
//	}
//
//func (se StatusErr) Error() string {
//	return se.Message
//}
//
//func (se StatusError) Unwrap() error {
//	return se.err
//}
//
//func LoginAndGetData(uid, pwd, file string) ([]byte, error) {
//	err := login(uid,pwd)
//	if err != nil {
//		return nil, StatusErr {
//			Status: InvalidLogin,
//			Message: fmt.Sprintf("invalid credentials for user %s",uid),
//			Err: err,
//		}
//	}
//	data, err := getData(file)
//	if err != nil {
//		return nil, StatusErr {
//			Status: NotFound,
//			Message: fmt.Sprintf("file %s not found",file),
//			Err: err,
//		}
//	}
//	return data, nil
//}

type MyErr struct {
	Codes []int
}

func (me MyErr) Error() string {
	return fmt.Sprintf("codes: %v", me.Codes)
}

func (me MyErr) Is(target error) bool {
	if me3, ok := target.(MyErr); ok {
		return reflect.DeepEqual(me,me3)
	}
	return false
}

func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("in fileChecker: %w", err)
	}
	f.Close()
	return nil
}

//type ResourceErr struct {
//	Resource     string
//	Code         int
//}
//
//func (re ResourceErr) Is(target error) bool {
//	if other, ok := target.(ResourceErr); ok {
//		ignoreResource := other.Resource == ""
//		ignoreCode := other.Code == 0
//		matchResource := other.Resource == re.Resource
//		matchCode := other.Code == re.Code
//		return matchResource && matchCode ||
//			matchResource && ignoreCode ||
//			ignoreResource && matchCode
//	}
//	return false
//}

//func DoSomeThings(val1 int, val2 string) (string, error) {
//	val3, err := fmt.Println(val1)
//	if err != nil {
//		return "", fmt.Errorf("in DoSomeThings: %w", err)
//	}
//	val4, err := fmt.Println(val2)
//	if err != nil {
//		return "", fmt.Errorf("in DoSomeThings: %w", err)
//	}
//	result, err := fmt.Println(val3, val4)
//	if err != nil {
//		return "", fmt.Errorf("in DoSomeThings: %w", err)
//	}
//	//return result, nil
//}

//Wrapping Errors with defer
//func DoSomeThings(val1 int, val2 string) (_ string, err error) {
//	defer func() {
//		if err != nil {
//			err = fmt.Errorf("in DoSomeThings: %w", err)
//		}
//	}()
//	val3, err := doThing1(val1)
//	if err != nil {
//		return "", err
//	}
//	val4, err := doThing2(val2)
//	if err != nil {
//		return "", err
//	}
//	return doThing3(val3, val4)
//}

func doPanic(msg string) {
	panic(msg)
}

func div60(i int) {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println(v)
		}
	}()
	fmt.Println(60 / i)
}

func main() {
	numerator := 20
	denominator := 3
	remainder, mod, err := calcRemainderAndMod(numerator, denominator)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(remainder, mod)


	data := []byte("This is not a zip file")
	notAZipFile := bytes.NewReader(data)
	_, err = zip.NewReader(notAZipFile, int64(len(data)))
	if err == zip.ErrFormat {
		fmt.Println("Sentinel Errors-Told you so")
	}

	//err := GenerateError(true)
	//fmt.Println(err != nil)
	//err = GenerateError(false)
	//fmt.Println(err != nil)

	//err := fileChecker("not_here.txt")
	//if err != nil {
	//	fmt.Println(err)
	//	if wrappedErr := errors.Unwrap(err); wrappedErr != nil {
	//		fmt.Println(wrappedErr)
	//	}
	//}

	err = fileChecker("not_here.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("That file doesn't exist")
		}
	}

	//if errors.Is(err, ResourceErr{Resource: "Database"}) {
	//	fmt.Println("The database is broken:", err)
	//	// process the codes
	//}

	//doPanic(os.Args[0])

	for _, val := range []int{1, 2, 0, 6} {
		div60(val)
	}

}