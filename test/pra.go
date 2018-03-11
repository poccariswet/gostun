package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type hoge struct {
	n int
}

func main() {
	h := &hoge{n: 100}
	v := reflect.ValueOf(h)
	fv1 := v.Elem().FieldByName("n")
	i := (*int)(unsafe.Pointer(fv1.UnsafeAddr()))
	*i = 200
	fmt.Println(*i)
}
