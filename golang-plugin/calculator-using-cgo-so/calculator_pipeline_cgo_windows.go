package main

/*
#include <windows.h>
#include <stdlib.h>

typedef char* (*operator_type)();
char* Operator(void* f) {
	return ((operator_type)f)();
}


typedef double (*operate_type)(double, double);
double Operate(void* f, double d1, double d2) {
	return ((operate_type)f)(d1, d2);
}
*/
import "C"
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	opFuncMap := map[string]func(float64, float64) float64{}

	populateOpFuncMap(&opFuncMap)

	fmt.Println(opFuncMap)

	for {
		fmt.Print("-> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("exit", strings.ToLower(text)) == 0 {
			break
		}
		v1, op, v2, err := parseExp(text)
		if err != nil {
			panic(err)
		}
		fmt.Printf("v1: %f, op: %s, v2: %f, out: %f\n", v1, op, v2, opFuncMap[op](v1, v2))
	}
}

func populateOpFuncMap(m *map[string]func(float64, float64) float64) {
	var soFilePaths []string
	filepath.Walk("./plugins", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if strings.HasSuffix(path, ".dll") {
			soFilePaths = append(soFilePaths, path)
		}
		return nil
	})
	for i := 0; i < len(soFilePaths); i++ {
		add(soFilePaths[i], m)
	}
}

func add(path string, m *map[string]func(float64, float64) float64) {
	pathCString := C.CString(path)
	defer C.free(unsafe.Pointer(pathCString))
	handle :=  C.LoadLibrary(pathCString)
	//defer C.free(unsafe.Pointer(handle))
	if handle == nil {
		err := fmt.Errorf("error opening %s", path)
		panic(err)
	}
	//defer func() {
	//	if r := C.dlclose(handle); r != 0 {
	//		err := fmt.Errorf("error closing %s", path)
	//		panic(err)
	//	}
	//}()

	keyPtrCString := C.CString("Operator")
	defer C.free(unsafe.Pointer(keyPtrCString))
	keyPtr := C.GetProcAddress(handle, keyPtrCString)

	if keyPtr == nil {
		err := fmt.Sprintf("No Operator for so: %s", path)
		panic(err)
	}

	println(C.GoString(C.Operator(unsafe.Pointer(keyPtr))))

	valuePtrCString := C.CString("Operate")
	defer C.free(unsafe.Pointer(valuePtrCString))
	valuePtr := C.GetProcAddress(handle, valuePtrCString)

	// if valuePtr == nil {
	// 	err := fmt.Sprintf("No Operate for so: %s", path)
	// 	panic(err)
	// }

	// key := C.Operator(unsafe.Pointer(keyPtr))
	// // defer C.free(unsafe.Pointer(key))

	(*m)[C.GoString(C.Operator(unsafe.Pointer(keyPtr)))] = func(f1 float64, f2 float64) float64 {
		return float64(C.Operate(unsafe.Pointer(valuePtr), C.double(f1), C.double(f2)))
	}
}



func parseExp(s string) (float64, string, float64, error) {
	println(s)
	expParser := regexp.MustCompile(`(?m)^([0-9]+) *([+\-\/*]) *([0-9]+) *$`)
	parts := expParser.FindStringSubmatch(s)
	fmt.Println(parts)
	if len(parts) != 4 {
		return 0, "", 0, errors.New(fmt.Sprintf("Invalid part length for expression: %s", s))
	}
	v1, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, "", 0, err
	}
	v2, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return 0, "", 0, err
	}
	return v1, parts[2], v2, nil
}
