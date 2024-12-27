//go:build windows && (amd64 || arm64)

package ffi

import (
	"fmt"
	"syscall"
)

func Load(name string) (l Lib, err error) {
	handle, err := syscall.LoadLibrary(name)
	if err != nil {
		err = fmt.Errorf("%s: error loading library: %w", name, err)
	}
	l.Addr = uintptr(handle)
	return
}

func (l Lib) Get(name string) (addr uintptr, err error) {
	return syscall.GetProcAddress(syscall.Handle(l.Addr), name)
}

func (l Lib) Close() error {
	if l.Addr == 0 {
		return nil
	}

	return syscall.FreeLibrary(syscall.Handle(l.Addr))
}
