//go:build (freebsd || linux) && (amd64 || arm64)

package ffi

import (
	"github.com/ebitengine/purego"
)

func init() {
	filename := "libffi.so.8"
Load:
	libffi, err := purego.Dlopen(filename, purego.RTLD_LAZY)
	if err != nil {
		if err.Error() == "libffi.so.8: cannot open shared object file: No such file or directory" {
			filename = "libffi.so.7"
			goto Load
		}
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
