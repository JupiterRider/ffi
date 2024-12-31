//go:build (freebsd || linux || windows || darwin) && (amd64 || arm64)

package ffi

import (
	"runtime"
)

func init() {
	var filename string

	switch runtime.GOOS {
	case "freebsd", "linux":
		filename = "libffi.so.8"
	case "windows":
		filename = "libffi-8.dll"
	case "darwin":
		filename = "libffi.8.dylib"
	}

	libffi, err := Load(filename)
	if err != nil {
		panic(err)
	}

	prepCif, err = libffi.Get("ffi_prep_cif")
	if err != nil {
		panic(err)
	}

	prepCifVar, err = libffi.Get("ffi_prep_cif_var")
	if err != nil {
		panic(err)
	}

	call, err = libffi.Get("ffi_call")
	if err != nil {
		panic(err)
	}

	closureAlloc, err = libffi.Get("ffi_closure_alloc")
	if err != nil {
		panic(err)
	}

	closureFree, err = libffi.Get("ffi_closure_free")
	if err != nil {
		panic(err)
	}

	prepClosureLoc, err = libffi.Get("ffi_prep_closure_loc")
	if err != nil {
		panic(err)
	}
}
