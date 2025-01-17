package ffi_test

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

func ExamplePrepClosureLoc() {
	var sin unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &sin)
	defer ffi.ClosureFree(closure)

	fun := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		arg := unsafe.Slice(args, cif.NArgs)[0]
		sine := math.Sin(*(*float64)(arg))
		*(*float64)(ret) = sine
		return 0
	})

	var cif ffi.Cif
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeDouble, &ffi.TypeDouble); status != ffi.OK {
		panic(status)
	}

	if status := ffi.PrepClosureLoc(closure, &cif, fun, nil, sin); status != ffi.OK {
		panic(status)
	}

	var sine float64
	var x float64 = 1
	ffi.Call(&cif, uintptr(sin), unsafe.Pointer(&sine), unsafe.Pointer(&x))
	fmt.Println(sine)
	// Output: 0.8414709848078965
}
