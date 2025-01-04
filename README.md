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
In this example we create our own library, which consists of two type definitions and one function:

```c
#include <stdbool.h>
#include <string.h>

typedef enum {
    GROCERIES,
    HOUSEHOLD,
    BEAUTY
} Category;

typedef struct {
    const char *name;
    double price;
    Category category;
} Item;

bool IsItemValid(Item item)
{
    if (!item.name || strlen(item.name) == 0)
    {
        return false;
    }

    if (item.price < 0)
    {
        return false;
    }

    if (item.category > BEAUTY)
    {
        return false;
    }

    return true;
}
```

Compile the code into a shared library:

```sh
gcc -shared -o libitem.so -fPIC item.c
```

The consuming Go code:

```golang
package main

import (
	"fmt"

	"github.com/jupiterrider/ffi"
)

type Category uint32

const (
	Groceries Category = iota
	Household
	Beauty
)

type Item struct {
	Name     *byte
	Price    float64
	Category Category
}

func main() {
	// load the library
	lib, err := ffi.Load("./libitem.so")
	if err != nil {
		panic(err)
	}

	// create a new ffi.Type which defines the fields of the Item struct
	typeItem := ffi.NewType(&ffi.TypePointer, &ffi.TypeDouble, &ffi.TypeUint32)

	// get the IsItemValid function and describe its signature
	// (for bool we use ffi.TypeUint8)
	isItemValid, err := lib.Prep("IsItemValid", &ffi.TypeUint8, &typeItem)
	if err != nil {
		panic(err)
	}

	var item Item
	// strings are null-terminated and converted into a byte pointer
	item.Name = &[]byte("Apple\x00")[0]
	item.Price = 0.22
	item.Category = Groceries

	// the return value is stored in a 64-bit integer type, because libffi
	// cannot handle smaller integer types as return value
	var result ffi.Arg

	// call the C function
	// (keep in mind that you have to pass pointers and not the values themselves)
	isItemValid.Call(&result, &item)

	if result.Bool() {
		fmt.Println("Item is valid!")
	} else {
		fmt.Println("Item is not valid!")
	}
}
```

You can find more examples inside the [examples](examples) folder of this repository.
