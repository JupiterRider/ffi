//go:build windows && (amd64 || arm64)

package ffi

import (
	"fmt"
	"syscall"
)

func init() {
	const filename = "libffi-8.dll"

	libffi, err := syscall.LoadLibrary(filename)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", filename, err))
	}

	prepCif, err = syscall.GetProcAddress(libffi, "ffi_prep_cif")
	if err != nil {
		panic(err)
	}

	prepCifVar, err = syscall.GetProcAddress(libffi, "ffi_prep_cif_var")
	if err != nil {
		panic(err)
	}

	call, err = syscall.GetProcAddress(libffi, "ffi_call")
	if err != nil {
		panic(err)
	}
}
