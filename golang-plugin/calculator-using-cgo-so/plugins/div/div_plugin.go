package main

import "C"

//export Operator
func Operator() *C.char {
	return C.CString("/")
}

//export Operate
func Operate(v1 C.double, v2 C.double) C.double {
	return v1 / v2
}

func main() {}
