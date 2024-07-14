//go:build windows && (amd64 || arm64)

package ffi

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func init() {
	const filename = "libffi-8.dll"

	libffi, err := windows.LoadLibrary(filename)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", filename, err))
	}

	prepCif, err = windows.GetProcAddress(libffi, "ffi_prep_cif")
	if err != nil {
		panic(err)
	}

	prepCifVar, err = windows.GetProcAddress(libffi, "ffi_prep_cif_var")
	if err != nil {
		panic(err)
	}

	call, err = windows.GetProcAddress(libffi, "ffi_call")
	if err != nil {
		panic(err)
	}
}
