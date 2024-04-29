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
	argTypes := []*ffi.Typ{&ffi.TypDouble}
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypDouble, argTypes); ok != ffi.OK {
		panic("ffi prep failed")
	}

	// call the function
	var returnValue float64
	argValue := float64(1)
	argValues := []unsafe.Pointer{unsafe.Pointer(&argValue)}
	ffi.Call(&cif, cos, unsafe.Pointer(&returnValue), argValues)

	// prints 0.5403023058681398
	fmt.Println(returnValue)
}
