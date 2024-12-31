package main

import (
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
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
	closure := ffi.ClosureAlloc(&code)
	defer ffi.ClosureFree(closure)

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 1, &ffi.TypeFloat, &ffi.TypeFloat); status != ffi.OK {
		panic(status)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, purego.NewCallback(F), nil, code); status != ffi.OK {
			panic(status)
		}
	}

	var ret float32
	invoke.Call(&ret, &code)
	fmt.Println(ret)
	fmt.Println(float32(math.Pi * 2))
}

func F(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	arguments := unsafe.Slice(args, cif.NArgs)
	*(*float32)(ret) = *(*float32)(arguments[0]) * 2
	return 0
}
