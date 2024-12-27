//go:build (freebsd || linux || darwin) && (amd64 || arm64)

package ffi

import (
	"github.com/ebitengine/purego"
)

func Load(name string) (l Lib, err error) {
	l.Addr, err = purego.Dlopen(name, purego.RTLD_LAZY)
	return
}

func (l Lib) Get(name string) (addr uintptr, err error) {
	return purego.Dlsym(l.Addr, name)
}

func (l Lib) Close() error {
	if l.Addr == 0 {
		return nil
	}

	return purego.Dlclose(l.Addr)
}
