package main

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

func main() {
	// open the shared library
	libm, err := purego.Dlopen("libm.so.6", purego.RTLD_LAZY)
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
