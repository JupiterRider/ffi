//go:build (freebsd || linux || windows || darwin) && (amd64 || arm64)

package ffi

import "fmt"

type Lib struct {
	Addr uintptr
}

func (l Lib) Prep(name string, ret *Type, args ...*Type) (f Fun, err error) {
	if f.Addr, err = l.Get(name); err != nil {
		return
	}

	f.Cif = new(Cif)
	if status := PrepCif(f.Cif, DefaultAbi, uint32(len(args)), ret, args...); status != OK {
		return f, fmt.Errorf("%s: error preparing function: %s", name, status)
	}

	return
}
