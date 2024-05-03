# ffi
[![Go Reference](https://pkg.go.dev/badge/github.com/jupiterrider/ffi.svg)](https://pkg.go.dev/github.com/jupiterrider/ffi)

A libffi wrapper for purego.

## Purpose
You can use [purego](https://github.com/ebitengine/purego) to call C code without cgo. ffi provides extra functionality (e.g. passing and returning structs by value).

## Requirements
### OS
Any 64-bit Linux distribution should work.

### Software
[libffi](https://github.com/libffi/libffi) is preinstalled on most distributions, because it also is a dependency of Python and Ruby. If not, you can install it explicitly:

#### Debian 12, Ubuntu 22.04, Ubuntu 24.04
```bash
sudo apt install libffi8
```

#### Arch Linux
```bash
sudo pacman -S libffi
```

## Examples
In this example we use the puts function inside the standard C library to print "Hello World!" to the console:

```c
extern int puts (const char *__s);
```

```golang
package main

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

func main() {
	// open the C library
	libm, err := purego.Dlopen("libc.so.6", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	// get the address of puts
	puts, err := purego.Dlsym(libm, "puts")
	if err != nil {
		panic(err)
	}

	// describe the function's signature
	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeSint32, &ffi.TypePointer); ok != ffi.OK {
		panic("ffi prep failed")
	}

	// convert the go string into a pointer
	text, _ := unix.BytePtrFromString("Hello World!")

	// call the puts function
	ffi.Call(&cif, puts, nil, unsafe.Pointer(&text))
}
```

You can find more examples inside the [examples](examples) folder of this repository.
