package main

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

func main() {
	var closure *ffi.Closure
	var bound_puts unsafe.Pointer
	fmt.Printf("%v\n", bound_puts)
	c := ffi.ClosureAlloc(uint64(unsafe.Sizeof(closure)), &bound_puts)
	closure = *(**ffi.Closure)(unsafe.Pointer(&c))
	fmt.Printf("%v\n", closure)

	ffi.ClosureFree(uintptr(unsafe.Pointer(closure)))
	fmt.Printf("%v\n", closure)
	closure.Tramp[0] = 255
	fmt.Printf("%v\n", closure)
}
