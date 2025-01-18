package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

func main() {
	var filename string
	switch runtime.GOOS {
	case "linux", "freebsd":
		filename = "./libcallback.so"
	case "windows":
		filename = "callback.dll"
	case "darwin":
		filename = "./libcallback.dylib"
	}

	// load the C library
	lib, err := ffi.Load(filename)
	if err != nil {
		panic(err)
	}

	// get the C function and describe its signature
	invoke, err := lib.Prep("Invoke", &ffi.TypeVoid, &ffi.TypePointer)
	if err != nil {
		panic(err)
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)

	// describe the closure's signature
	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 0, &ffi.TypeVoid); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		fmt.Println("Hello, World!")
		return 0
	})

	// prepare the closure
	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
	}

	// prints "Hello, World!"
	invoke.Call(nil, &callback)

	// the closure can be freed now
	ffi.ClosureFree(closure)
}
