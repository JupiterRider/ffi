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

	text := &[]byte("Pi is %f\n\x00")[0] // C requires a null-terminated string
	pi := math.Pi
	var nCharsPrinted ffi.Arg
	printf.Call(&nCharsPrinted, &text, &pi)
	fmt.Printf("%d characters printed\n", int32(nCharsPrinted))
}
