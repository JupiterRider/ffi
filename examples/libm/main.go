package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

func main() {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "libm.so.6"
	case "freebsd":
		filename = "libm.so.5"
	}

	// open the shared library
	libm, err := purego.Dlopen(filename, purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	// get the function's address
	cos, err := purego.Dlsym(libm, "cos")
	if err != nil {
		panic(err)
	}

	// describe the function's signature
	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeDouble, &ffi.TypeDouble); ok != ffi.OK {
		panic("ffi prep failed")
	}

	// call the function
	returnValue, argValue := 0.0, 1.0
	ffi.Call(&cif, cos, unsafe.Pointer(&returnValue), unsafe.Pointer(&argValue))

	// prints 0.5403023058681398
	fmt.Println(returnValue)
}
