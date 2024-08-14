# ffi
[![Go Reference](https://pkg.go.dev/badge/github.com/jupiterrider/ffi.svg)](https://pkg.go.dev/github.com/jupiterrider/ffi)

A purego binding for libffi.

## Purpose
You can use [purego](https://github.com/ebitengine/purego) to call C code without cgo. ffi provides extra functionality (e.g. passing and returning structs by value).

## Requirements
### OS/Architecture
- darwin/amd64
- darwin/arm64
- freebsd/amd64
- freebsd/arm64
- linux/amd64
- linux/arm64
- windows/amd64
- windows/arm64

### Software
[libffi](https://github.com/libffi/libffi) is preinstalled on most distributions, because it also is a dependency of Python and Ruby. If not, you can install it explicitly:

#### Arch Linux
```sh
sudo pacman -S libffi
```

#### Debian 12, Ubuntu 22.04, Ubuntu 24.04
```sh
sudo apt install libffi8
```

#### Debian 11, Ubuntu 20.04
```sh
sudo apt install libffi7
```

#### FreeBSD
```sh
pkg install libffi
```
Note: Use this `-gcflags="github.com/ebitengine/purego/internal/fakecgo=-std"` build flag when cross compiling or having CGO_ENABLED set to 0 (FreeBSD only).

#### Windows
You need a `libffi-8.dll` next to the executable/root folder of your project or inside C:\Windows\System32. If you don't want to build libffi from source, you can find this dll for example inside the [Windows embeddable package](https://www.python.org/downloads/windows/) of Python.

#### macOS
You can use [Homebrew](https://brew.sh/) to install libffi:
```sh
brew install libffi
```
Note: If dlopen can't find the libffi.8.dylib file, you can try setting this environment variable:
```sh
export DYLD_FALLBACK_LIBRARY_PATH=$DYLD_FALLBACK_LIBRARY_PATH:/opt/homebrew/opt/libffi/lib
```

## Examples
In this example we use the puts function inside the standard C library to print "Hello World!" to the console:

```c
int puts(const char *s);
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
	if status := ffi.PrepCif(&cif, ffi.DefaultAbi, 1, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// convert the go string into a pointer
	text, _ := unix.BytePtrFromString("Hello World!")

	// call the puts function
	var result ffi.Arg
	ffi.Call(&cif, puts, unsafe.Pointer(&result), unsafe.Pointer(&text))
}
```

You can find more examples inside the [examples](examples) folder of this repository.
