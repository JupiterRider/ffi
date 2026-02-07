package main

import (
	"fmt"
	"runtime"

	"github.com/jupiterrider/ffi"
)

func main() {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "libm.so.6"
	case "freebsd":
		filename = "libm.so.5"
	case "darwin":
		filename = "libm.dylib"
	case "windows":
		filename = "ntdll.dll"
	}

	// open the shared library
	libm, err := ffi.Load(filename)
	if err != nil {
		panic(err)
	}

	// get the function's address and describe its signature
	cos, err := libm.Prep("cos", &ffi.TypeDouble, &ffi.TypeDouble)
	if err != nil {
		panic(err)
	}

	// call the function
	cosine, x := 0.0, 1.0
	cos.Call(&cosine, &x)

	// prints 0.5403023058681398
	fmt.Println(cosine)
}
