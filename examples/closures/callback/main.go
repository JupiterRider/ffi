package main

import (
	"fmt"
	"math"
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

	lib, err := ffi.Load(filename)
	if err != nil {
		panic(err)
	}

	invoke, err := lib.Prep("Invoke", &ffi.TypeFloat, &ffi.TypePointer)
	if err != nil {
		panic(err)
	}

	var code unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &code)
	defer ffi.ClosureFree(closure)

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 1, &ffi.TypeFloat, &ffi.TypeFloat); status != ffi.OK {
		panic(status)
	}

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		arguments := unsafe.Slice(args, cif.NArgs)
		*(*float32)(ret) = *(*float32)(arguments[0]) * *(*float32)(userData)
		return 0
	})

	multiplier := float32(2)
	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, unsafe.Pointer(&multiplier), code); status != ffi.OK {
			panic(status)
		}
	}

	var ret float32
	invoke.Call(&ret, &code)
	fmt.Println(ret)
	fmt.Println(float32(math.Pi * multiplier))
}
