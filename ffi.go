//go:build linux && amd64
// +build linux,amd64

package ffi

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/ebitengine/purego"
)

var prepCif, call uintptr

func init() {
	handle, err := purego.Dlopen("libffi.so", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	prepCif, err = purego.Dlsym(handle, "ffi_prep_cif")
	if err != nil {
		panic(err)
	}

	call, err = purego.Dlsym(handle, "ffi_call")
	if err != nil {
		panic(err)
	}
}

type Abi uint32

const (
	DefaultAbi Abi = 2
)

type Status uint32

const (
	OK Status = iota
	BadTypedef
	BadAbi
	BadArgType
)

const (
	Void = iota
	Int
	Float
	Double
	Longdouble
	Uint8
	Sint8
	Uint16
	Sint16
	Uint32
	Sint32
	Uint64
	Sint64
	Struct
	Pointer
	Complex
)

type Type struct {
	Size      uint64
	Alignment uint16
	Type      uint16
	Elements  **Type
}

type Cif struct {
	Abi      uint32
	NArgs    uint32
	ArgTypes **Type
	RType    *Type
	Bytes    uint32
	Flags    uint32
}

func PrepCif(cif *Cif, abi Abi, nArgs uint32, rType *Type, aTypes []*Type) Status {
	ret, _, err := purego.SyscallN(prepCif, uintptr(unsafe.Pointer(cif)), uintptr(abi), uintptr(nArgs), uintptr(unsafe.Pointer(rType)), uintptr(reflect.ValueOf(aTypes).UnsafePointer()))
	if err != 0 {
		panic(fmt.Sprintf("syscall failed with error code %d", err))
	}
	return Status(ret)
}

func Call(cif *Cif, fn uintptr, rValue unsafe.Pointer, aValues []unsafe.Pointer) {
	_, _, err := purego.SyscallN(call, uintptr(unsafe.Pointer(cif)), fn, uintptr(rValue), uintptr(reflect.ValueOf(aValues).UnsafePointer()))
	if err != 0 {
		panic(fmt.Sprintf("syscall failed with error code %d", err))
	}
}
