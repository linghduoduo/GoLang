package main

import (
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"math/bits"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Foo struct {
	A int
	B string
}

var x Foo

func DoSomething(f Foo) {
	fmt.Println(f.A, f.B)
}

func changeInt(i *int) {
	*i = 20
}

func changeIntReflect(i *int) {
	iv := reflect.ValueOf(i)
	iv.Elem().SetInt(20)
}

func hasNoValue(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func,
		reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

// Unmarshal maps all of the rows of data in a slice of slice of strings
// into a slice of structs.
// The first row is assumed to be the header with the column names.
func Marshal(v interface{}) ([][]string, error) {
	sliceVal := reflect.ValueOf(v)
	if sliceVal.Kind() != reflect.Slice {
		return nil, errors.New("must be a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("must be a slice of structs")
	}
	var out [][]string
	header := marshalHeader(structType)
	out = append(out, header)
	for i := 0; i < sliceVal.Len(); i++ {
		row, err := marshalOne(sliceVal.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

// Marshal maps all of the structs in a slice of structs to a slice of slice
// of strings.
// The first row written is the header with the column names.
func marshalHeader(vt reflect.Type) []string {
	var row []string
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		if curTag, ok := field.Tag.Lookup("csv"); ok {
			row = append(row, curTag)
		}
	}
	return row
}

func marshalOne(vv reflect.Value) ([]string, error) {
	var row []string
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		fieldVal := vv.Field(i)
		if _, ok := vt.Field(i).Tag.Lookup("csv"); !ok {
			continue
		}
		switch fieldVal.Kind() {
		case reflect.Int:
			row = append(row, strconv.FormatInt(fieldVal.Int(), 10))
		case reflect.String:
			row = append(row, fieldVal.String())
		case reflect.Bool:
			row = append(row, strconv.FormatBool(fieldVal.Bool()))
		default:
			return nil, fmt.Errorf("cannot handle field of kind %v",
				fieldVal.Kind())
		}
	}
	return row, nil
}

func Unmarshal(data [][]string, v interface{}) error {
	sliceValPtr := reflect.ValueOf(v)
	if sliceValPtr.Kind() != reflect.Ptr {
		return errors.New("must be a pointer to a slice of structs")
	}
	sliceVal := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return errors.New("must be a pointer to a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return errors.New("must be a pointer to a slice of structs")
	}

	// assume the first row is a header
	header := data[0]
	namePos := make(map[string]int, len(header))
	for k, v := range header {
		namePos[v] = k
	}

	for _, row := range data[1:] {
		newVal := reflect.New(structType).Elem()
		err := unmarshalOne(row, namePos, newVal)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, newVal))
	}
	return nil
}

func unmarshalOne(row []string, namePos map[string]int, vv reflect.Value) error {
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		typeField := vt.Field(i)
		pos, ok := namePos[typeField.Tag.Get("csv")]
		if !ok {
			continue
		}
		val := row[pos]
		field := vv.Field(i)
		switch field.Kind() {
		case reflect.Int:
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		case reflect.String:
			field.SetString(val)
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			field.SetBool(b)
		default:
			return fmt.Errorf("cannot handle field of kind %v",
				field.Kind())
		}
	}
	return nil
}

type MyData struct {
	Name   string `csv:"name"`
	Age    int    `csv:"age"`
	HasPet bool   `csv:"has_pet"`
}

func MakeTimedFunction(f interface{}) interface{} {
	ft := reflect.TypeOf(f)
	fv := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(ft, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := fv.Call(in)
		end := time.Now()
		fmt.Println(end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func timeMe(a int) int {
	time.Sleep(time.Duration(a) * time.Second)
	result := a * 2
	return result
}

func Filter(slice interface{}, filter interface{}) interface{} {
	sv := reflect.ValueOf(slice)
	fv := reflect.ValueOf(filter)

	sliceLen := sv.Len()
	out := reflect.MakeSlice(sv.Type(), 0, sliceLen)
	for i := 0; i < sliceLen; i++ {
		curVal := sv.Index(i)
		values := fv.Call([]reflect.Value{curVal})
		if values[0].Bool() {
			out = reflect.Append(out, curVal)
		}
	}
	return out.Interface()
}

func main() {
	//TYPES AND KINDS
	var x int
	xt := reflect.TypeOf(x)
	fmt.Println(xt.Name())     // returns int

	f := Foo{}
	ft := reflect.TypeOf(f)
	fmt.Println(ft.Name())     // returns Foo
	xpt := reflect.TypeOf(&x)
	fmt.Println(xpt.Name())    // returns an empty string

	fmt.Println(xpt.Name())        // returns an empty string
	fmt.Println(xpt.Kind())        // returns reflect.Ptr
	fmt.Println(xpt.Elem().Name()) // returns "int"
	fmt.Println(xpt.Elem().Kind()) // returns reflect.Int

	type Foo struct {
		A int    `myTag:"value"`
		B string `myTag:"value2"`
	}

	var f2 Foo
	ft2 := reflect.TypeOf(f2)
	for i := 0; i < ft2.NumField(); i++ {
		curField := ft2.Field(i)
		fmt.Println(curField.Name, curField.Type.Name(),
			curField.Tag.Get("myTag"))
	}

	s := []string{"a", "b", "c"}
	sv := reflect.ValueOf(s)        // sv is of type reflect.Value
	s3 := sv.Interface().([]string) // s2 is of type []string
	fmt.Println(s3)

	i := 10
	iv := reflect.ValueOf(&i)
	ivv :=  iv.Elem()
	ivv.SetInt(20)
	fmt.Println(i) // prints 20

	changeInt(&i)
	fmt.Println(i)

	changeIntReflect(&i)
	fmt.Println(i)

	//var stringType = reflect.TypeOf((*string)(nil)).Elem()
	//var stringSliceType = reflect.TypeOf([]string(nil))
	//ssv := reflect.MakeSlice(stringSliceType, 0, 10)
	//sv2 := reflect.New(stringType).Elem()
	//sv.SetString("hello")
	//ssv = reflect.Append(ssv, sv2)
	//ss := ssv.Interface().([]string)
	//fmt.Println(ss) // prints [hello]

	//Use Reflection to Write a Data Marshaler

	data := `name,age,has_pet
Jon,"100",true
"Fred ""The Hammer"" Smith",42,false
Martha,37,"true"
`
	r := csv.NewReader(strings.NewReader(data))
	allData, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	var entries []MyData
	Unmarshal(allData, &entries)
	fmt.Println(entries)

	//now to turn entries into output
	out, err := Marshal(entries)
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	w := csv.NewWriter(sb)
	w.WriteAll(out)
	fmt.Println(sb)

	//Build Functions with Reflection to Automate Repetitive Tasks
	timed:= MakeTimedFunction(timeMe).(func(int) int)
	fmt.Println(timed(2))

	//Only Use Reflection If It’s Worthwhile
	names := []string{"Andrew", "Bob", "Clara", "Hortense"}
	longNames := Filter(names, func(s string) bool {
		return len(s) > 3
	}).([]string)
	fmt.Println(longNames)

	ages := []int{20, 50, 13}
	adults := Filter(ages, func(age int) bool {
		return age >= 18
	}).([]int)
	fmt.Println(adults)

	//unsafe Strings and Slices
	s2 := "hello"
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s2))
	fmt.Println(sHdr.Len) // prints 5

	for i := 0; i < sHdr.Len; i++ {
		bp := *(*byte)(unsafe.Pointer(sHdr.Data + uintptr(i)))
		fmt.Print(string(bp))
	}
	fmt.Println()
	runtime.KeepAlive(s)

	//s4 := []int{10, 20, 30}
	//sHdr := (*reflect.SliceHeader)(unsafe.Pointer(&s4))
	//fmt.Println(sHdr.Len) // prints 3
	//fmt.Println(sHdr.Cap) // prints 3
	//
	//intByteSize := unsafe.Sizeof(s[0])
	//fmt.Println(intByteSize)
	//for i := 0; i < sHdr.Len; i++ {
	//	intVal := *(*int)(unsafe.Pointer(sHdr.Data + intByteSize*uintptr(i)))
	//	fmt.Println(intVal)
	//}
	//runtime.KeepAlive(s)




}

type Data struct {
	Value  uint32   // 4 bytes
	Label  [10]byte // 10 bytes
	Active bool     // 1 byte
	// Go padded this with 1 byte to make it align
}


func DataFromBytes(b [16]byte) Data {
	d := Data{}
	d.Value = binary.BigEndian.Uint32(b[:4])
	copy(d.Label[:], b[4:14])
	d.Active = b[14] != 0
	return d
}

func DataFromBytesUnsafe(b [16]byte) Data {
	data := *(*Data)(unsafe.Pointer(&b))
	if isLE {
		data.Value = bits.ReverseBytes32(data.Value)
	}
	return data
}

var isLE bool

func init() {
	var x uint16 = 0xFF00
	xb := *(*[2]byte)(unsafe.Pointer(&x))
	isLE = (xb[0] == 0x00)
}


func BytesFromData(d Data) [16]byte {
	out := [16]byte{}
	binary.BigEndian.PutUint32(out[:4], d.Value)
	copy(out[4:14], d.Label[:])
	if d.Active {
		out[14] = 1
	}
	return out
}


func BytesFromDataUnsafe(d Data) [16]byte {
	if isLE {
		d.Value = bits.ReverseBytes32(d.Value)
	}
	b := *(*[16]byte)(unsafe.Pointer(&d))
	return b
}
