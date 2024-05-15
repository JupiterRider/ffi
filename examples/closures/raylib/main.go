package main

import (
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

var cifVsprintf ffi.Cif
var vsprintf uintptr

func main() {
	runtime.LockOSThread()

	libc, err := purego.Dlopen("libc.so.6", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	vsprintf, err = purego.Dlsym(libc, "vsprintf")
	if err != nil {
		panic(err)
	}

	if status := ffi.PrepCif(&cifVsprintf, ffi.DefaultAbi, 3, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	raylib, err := purego.Dlopen("libraylib.so", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	// void SetTraceLogCallback(TraceLogCallback callback)
	setTraceLogCallback, err := purego.Dlsym(raylib, "SetTraceLogCallback")
	if err != nil {
		panic(err)
	}

	var cifSetTraceLogCallback ffi.Cif
	if status := ffi.PrepCif(&cifSetTraceLogCallback, ffi.DefaultAbi, 1, &ffi.TypeVoid, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// typedef void (*TraceLogCallback)(int logLevel, const char *text, va_list args);
	var cifTraceLogCallback ffi.Cif
	if status := ffi.PrepCif(&cifTraceLogCallback, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	var code unsafe.Pointer
	closure := ffi.ClosureAlloc(&code)
	// defer ffi.ClosureFree(closure)

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifTraceLogCallback, purego.NewCallback(CustomLog), nil, code); status != ffi.OK {
			panic(status)
		}
	}

	ffi.Call(&cifSetTraceLogCallback, setTraceLogCallback, nil, unsafe.Pointer(&code))

	// void TraceLog(int logLevel, const char *text, ...);
	traceLog, err := purego.Dlsym(raylib, "TraceLog")
	if err != nil {
		panic(err)
	}

	var cifTraceLog ffi.Cif
	if status := ffi.PrepCifVar(&cifTraceLog, ffi.DefaultAbi, 2, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeDouble); status != ffi.OK {
		panic(status)
	}

	logLevel := int32(3)
	text, _ := unix.BytePtrFromString("Pi is %f")
	pi := math.Pi
	ffi.Call(&cifTraceLog, traceLog, nil, unsafe.Pointer(&logLevel), unsafe.Pointer(&text), unsafe.Pointer(&pi))
}

func CustomLog(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) {
	arguments := unsafe.Slice(args, cif.NArgs)
	logLevel := *(*int32)(arguments[0])
	var charsWritten int32
	var buffer [256]byte
	b := &buffer[0]
	ffi.Call(&cifVsprintf, vsprintf, unsafe.Pointer(&charsWritten), unsafe.Pointer(&b), arguments[1], arguments[2])
	if charsWritten > 0 {
		fmt.Printf("[%d] %s\n", logLevel, unix.BytePtrToString(b))
	}
}
