//go:build darwin && (amd64 || arm64)

package ffi

import (
	"github.com/ebitengine/purego"
)

func init() {
	const filename = "libffi.8.dylib"

	libffi, err := purego.Dlopen(filename, purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	prepCif, err = purego.Dlsym(libffi, "ffi_prep_cif")
	if err != nil {
		panic(err)
	}

	prepCifVar, err = purego.Dlsym(libffi, "ffi_prep_cif_var")
	if err != nil {
		panic(err)
	}

	call, err = purego.Dlsym(libffi, "ffi_call")
	if err != nil {
		panic(err)
	}
}
