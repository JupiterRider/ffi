//go:build (freebsd || linux) && (amd64 || arm64)

package ffi

import (
	"reflect"
	"unsafe"

	"github.com/ebitengine/purego"
)

var prepCif, prepCifVar, call uintptr

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

type Abi uint32

type Status uint32

const (
	OK Status = iota
	BadTypedef
	BadAbi
	BadArgType
)

func (s Status) String() string {
	status := map[Status]string{OK: "OK", BadTypedef: "bad type definition", BadAbi: "bad ABI", BadArgType: "bad argument type"}
	return status[s]
}

// These constants are used for the Type field of [Type].
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

// Type is used to describe the structure of a data type.
//
// Example:
//
//	typedef struct Point {
//		int x;
//		int y;
//	} Point;
//
//	typePoint := ffi.Type{Type: ffi.Struct, Elements: &[]*ffi.Type{&ffi.TypeSint32, &ffi.TypeSint32, nil}[0]}
//
// Primitive data types are already defined (e.g. [TypeDouble] for float64).
type Type struct {
	Size      uint64 // Initialize to 0 (automatically set by libffi as needed).
	Alignment uint16 // Initialize to 0 (automatically set by libffi as needed).
	Type      uint16 // Use ffi.Struct for struct types.
	Elements  **Type // Pointer to the first element of a nil-terminated slice.
}

// Cif stands for "Call InterFace". It describes the signature of a function.
//
// Use [PrepCif] to initialize it.
type Cif struct {
	Abi      uint32
	NArgs    uint32
	ArgTypes **Type
	RType    *Type
	Bytes    uint32
	Flags    uint32
}

// PrepCif initializes cif.
//   - abi is the ABI to use. Normally [DefaultAbi] is what you want.
//   - nArgs is the number of arguments. Use 0 if the function has none.
//   - rType is the return type. Use [TypeVoid] if the function has none.
//   - aTypes are the arguments. Leave empty or provide nil if the function has none.
//
// The returned status code will be [OK], if everything worked properly.
//
// Example:
//
//	double cos(double x);
//
//	var cif ffi.Cif
//	status := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeDouble, &ffi.TypeDouble)
//	if status != ffi.OK {
//		panic(status)
//	}
func PrepCif(cif *Cif, abi Abi, nArgs uint32, rType *Type, aTypes ...*Type) Status {
	ret, _, _ := purego.SyscallN(prepCif, uintptr(unsafe.Pointer(cif)), uintptr(abi), uintptr(nArgs), uintptr(unsafe.Pointer(rType)), uintptr(reflect.ValueOf(aTypes).UnsafePointer()))
	return Status(ret)
}

// PrepCifVar initializes cif for a call to a variadic function.
//
// In general its operation is the same as for [PrepCif] except that:
//   - nFixedArgs is the number of fixed arguments, prior to any variadic arguments. It must be greater than zero.
//   - nTotalArgs is the total number of arguments, including variadic and fixed arguments. aTypes must have this many elements.
//
// This function will return [BadArgType] if any of the variable argument types is [TypeFloat].
// Same goes for integer types smaller than 4 bytes. See [issue 608].
//
// Note that, different cif's must be prepped for calls to the same function when different numbers of arguments are passed.
//
// Also note that a call to this function with nFixedArgs = nTotalArgs is NOT equivalent to a call to [PrepCif].
//
// Example:
//
//	int printf(const char *restrict format, ...);
//
//	var cif ffi.Cif
//	status := ffi.PrepCifVar(&cif, ffi.DefaultAbi, 1, 2, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeDouble)
//	if status != ffi.OK {
//		panic(status)
//	}
//
//	text, _ := unix.BytePtrFromString("Pi is %f\n")
//	pi := math.Pi
//	var nCharsPrinted int32
//	ffi.Call(&cif, printf, unsafe.Pointer(&nCharsPrinted), unsafe.Pointer(&text), unsafe.Pointer(&pi))
//
// [issue 608]: https://github.com/libffi/libffi/issues/608
func PrepCifVar(cif *Cif, abi Abi, nFixedArgs, nTotalArgs uint32, rType *Type, aTypes ...*Type) Status {
	const intSize = 4

	// This check has been rebuild according to the original: https://github.com/libffi/libffi/blob/v3.4.6/src/prep_cif.c#L244
	//
	// Without rebuild, the type check wouldn't work for float,
	// because libffi compares the pointer to ffi_type_float instead of value equality.
	for i := nFixedArgs; i < nTotalArgs; i++ {
		argType := *aTypes[i]
		if argType == TypeFloat || ((argType.Type != Struct && argType.Type != Complex) && argType.Size < intSize) {
			return BadArgType
		}
	}

	ret, _, _ := purego.SyscallN(prepCifVar, uintptr(unsafe.Pointer(cif)), uintptr(abi), uintptr(nFixedArgs), uintptr(nTotalArgs), uintptr(unsafe.Pointer(rType)), uintptr(reflect.ValueOf(aTypes).UnsafePointer()))
	return Status(ret)
}

// Call calls the function fn according to the description given in cif. cif must have already been prepared using [PrepCif].
//   - fn is the address of the desired function. Use [purego.Dlsym] to get one.
//   - rValue is a pointer to a variable that will hold the result of the function call. Provide nil if the function has no return value.
//   - aValues are pointers to the argument values. Leave empty or provide nil if the function takes none.
//
// Example:
//
//	double cos(double x);
//
//	cosine, x := 0.0, 1.0
//	ffi.Call(&cif, cos, unsafe.Pointer(&cosine), unsafe.Pointer(&x))
func Call(cif *Cif, fn uintptr, rValue unsafe.Pointer, aValues ...unsafe.Pointer) {
	purego.SyscallN(call, uintptr(unsafe.Pointer(cif)), fn, uintptr(rValue), uintptr(reflect.ValueOf(aValues).UnsafePointer()))
}
