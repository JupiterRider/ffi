//go:build (freebsd || linux) && (amd64 || arm64)

package main

import (
	"math"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

func main() {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "libc.so.6"
	case "freebsd":
		filename = "libc.so.7"
	}

	libc, err := purego.Dlopen(filename, purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	printf, err := purego.Dlsym(libc, "printf")
	if err != nil {
		panic(err)
	}

	var cif ffi.Cif
	status := ffi.PrepCifVar(&cif, ffi.DefaultAbi, 1, 2, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeDouble)
	if status != ffi.OK {
		panic(status)
	}

	text, _ := unix.BytePtrFromString("Pi is %f\n")
	pi := math.Pi
	var nCharsPrinted int32
	ffi.Call(&cif, printf, unsafe.Pointer(&nCharsPrinted), unsafe.Pointer(&text), unsafe.Pointer(&pi))
}
