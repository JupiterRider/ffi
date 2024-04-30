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
	OK Status = 0
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

type Typ struct {
	Size      uint64
	Alignment uint16
	Typ       uint16
	Elements  **Typ
}

type Cif struct {
	Abi      uint32
	Nargs    uint32
	ArgTypes **Typ
	RTyp     *Typ
	Bytes    uint32
	Flags    uint32
}

var PrepCif func(cif *Cif, abi Abi, nargs uint32, rtyp *Typ, atypes []*Typ) Status

var Call func(cif *Cif, fn uintptr, rvalue unsafe.Pointer, avalue []unsafe.Pointer)
