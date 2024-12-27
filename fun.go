//go:build (freebsd || linux || windows || darwin) && (amd64 || arm64)

package ffi

import (
	"reflect"
	"unsafe"
)

type Fun struct {
	Addr uintptr
	Cif  *Cif
}

func (f Fun) Call(ret any, args ...any) {
	var rV unsafe.Pointer
	if ret != nil {
		rV = reflect.ValueOf(ret).UnsafePointer()
	}

	nArgs := len(args)
	aV := make([]unsafe.Pointer, nArgs)

	for i := 0; i < nArgs; i++ {
		aV[i] = reflect.ValueOf(args[i]).UnsafePointer()
	}

	Call(f.Cif, f.Addr, rV, aV...)
}
