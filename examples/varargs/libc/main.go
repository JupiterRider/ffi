//go:build (freebsd || linux || darwin || windows) && (amd64 || arm64)

package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/jupiterrider/ffi"
)

func main() {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "libc.so.6"
	case "freebsd":
		filename = "libc.so.7"
	case "windows":
		filename = "msvcrt.dll"
	case "darwin":
		filename = "libc.dylib"
	}

	libc, err := ffi.Load(filename)
	if err != nil {
		panic(err)
	}

	// printf is a variadic function, with one fixed argument
	printf, err := libc.PrepVar("printf", 1, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeDouble)
	if err != nil {
		panic(err)
	}

	fflush, err := libc.Prep("fflush", &ffi.TypeSint32, &ffi.TypePointer)
	if err != nil {
		panic(err)
	}

	// C requires a null-terminated string
	text := &[]byte("Pi is %f\n\x00")[0]
	pi := math.Pi
	var nCharsPrinted ffi.Arg
	printf.Call(&nCharsPrinted, &text, &pi)

	// we call fflush with NULL as argument to flush all open streams,
	// because printf is buffered
	var ok ffi.Arg
	var stream uintptr
	fflush.Call(&ok, &stream)

	fmt.Printf("%d characters printed\n", int32(nCharsPrinted))
}
