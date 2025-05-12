# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.1] - 2025-05-12

### Fixed
- The dependency `github.com/ebitengine/purego` had to be upgraded to version 0.8.3 due to a [bug](https://github.com/golang/go/issues/73617) in Go 1.23.9 and 1.24.3.

## [0.4.0] - 2025-03-13

### Changed
- `Fun.Call` panics now, if the number of arguments doesn't match the prepared Cif.

## [0.3.0] - 2025-01-18

### Added
- libffi's closure API has been implemented, which allows you to create C functions at runtime:
    - `func ClosureAlloc(size uintptr, code *unsafe.Pointer) *Closure`
    - `func ClosureFree(writable *Closure)`
    - `func PrepClosureLoc(closure *Closure, cif *Cif, fun uintptr, userData, codeLoc unsafe.Pointer) Status`
- The new types `Fun` and `Lib` can reduce boilerplate and eliminate platform-dependent code:
    - `func (f Fun) Call(ret any, args ...any)`
    - `func Load(name string) (l Lib, err error)`
    - `func (l Lib) Close() error`
    - `func (l Lib) Get(name string) (addr uintptr, err error)`
    - `func (l Lib) Prep(name string, ret *Type, args ...*Type) (f Fun, err error)`
    - `func (l Lib) PrepVar(name string, nFixedArgs int, ret *Type, args ...*Type) (f Fun, err error)`
- New method `func (a Arg) Bool() bool` added.
- [Changelog](https://github.com/JupiterRider/ffi/blob/main/CHANGELOG.md) file added.

### Changed
- On Linux, libffi.so.7 was loaded when libffi.so.8 could not be found. This is no longer the case.

## [0.2.2] - 2024-12-22

### Added
- Function `func NewType(elements ...*Type) Type` added.

## [0.2.1] - 2024-10-30

### Changed

- Dependency `github.com/ebitengine/purego` updated.
- The raylib example embeds the image file now.

## [0.2.0] - 2024-09-30

### Removed

- Dependency `golang.org/x/sys` removed.

## [0.1.1] - 2024-09-28

### Changed

- Dependencies `github.com/ebitengine/purego` and `golang.org/x/sys` updated.

[0.4.1]: https://github.com/JupiterRider/ffi/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/JupiterRider/ffi/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/JupiterRider/ffi/compare/v0.2.2...v0.3.0
[0.3.0]: https://github.com/JupiterRider/ffi/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/JupiterRider/ffi/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/JupiterRider/ffi/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/JupiterRider/ffi/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/JupiterRider/ffi/compare/v0.1.0...v0.1.1
