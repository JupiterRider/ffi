//go:build linux && amd64
// +build linux,amd64

package ffi

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

func init() {
	handle, err := purego.Dlopen("libffi.so", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&PrepCif, handle, "ffi_prep_cif")
	purego.RegisterLibFunc(&Call, handle, "ffi_call")
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

var PrepCif func(cif *Cif, abi Abi, nArgs uint32, rType *Type, aTypes []*Type) Status

var Call func(cif *Cif, fn uintptr, rValue unsafe.Pointer, aValues []unsafe.Pointer)
