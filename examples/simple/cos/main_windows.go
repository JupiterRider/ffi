//go:build windows && (amd64 || arm64)

package main

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/windows"
)

func main() {
	// open the shared library
	const libname = "ntdll.dll"
	ntdll, err := windows.LoadLibrary(libname)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	}

	// get the function's address
	cos, err := windows.GetProcAddress(ntdll, "cos")
	if err != nil {
		panic(err)
	}

	// describe the function's signature
	var cif ffi.Cif
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeDouble, &ffi.TypeDouble); status != ffi.OK {
		panic(status)
	}

	// call the function
	cosine, x := 0.0, 1.0
	ffi.Call(&cif, cos, unsafe.Pointer(&cosine), unsafe.Pointer(&x))

	// prints 0.5403023058681398
	fmt.Println(cosine)
}
