//go:build windows && (amd64 || arm64)

package main

import (
	"fmt"
	"image/color"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type Texture struct {
	ID                             uint32
	Width, Height, Mipmaps, Format int32
}

var (
	TypeTexture = ffi.Type{Type: ffi.Struct, Elements: &[]*ffi.Type{&ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, nil}[0]}
	TypeColor   = ffi.Type{Type: ffi.Struct, Elements: &[]*ffi.Type{&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, nil}[0]}
)

var (
	InitWindow        func(width, height int32, title string)
	CloseWindow       func()
	WindowShouldClose func() bool
	BeginDrawing      func()
	EndDrawing        func()
	ClearBackground   func(col color.RGBA)
	LoadTexture       func(filename string) Texture
	UnloadTexture     func(texture Texture)
	DrawTexture       func(texture Texture, posX, posY int32, col color.RGBA)
)

func init() {
	runtime.LockOSThread()

	const libname = "raylib.dll"
	raylib, err := syscall.LoadLibrary(libname)
	if err != nil {
		panic(fmt.Errorf("cannot load library %s: %w", libname, err))
	}

	// InitWindow -------------------------------
	var cifInitWindow ffi.Cif
	if status := ffi.PrepCif(&cifInitWindow, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	symInitWindow, err := syscall.GetProcAddress(raylib, "InitWindow")
	if err != nil {
		panic(err)
	}

	InitWindow = func(width, height int32, title string) {
		byteTitle := &[]byte(title + "\x00")[0] // you can also use golang.org/x/sys/windows.BytePtrFromString to create a null-terminated string
		ffi.Call(&cifInitWindow, symInitWindow, nil, unsafe.Pointer(&width), unsafe.Pointer(&height), unsafe.Pointer(&byteTitle))
	}

	// CloseWindow ------------------------------
	var cifVoidVoid ffi.Cif
	if status := ffi.PrepCif(&cifVoidVoid, ffi.DefaultAbi, 0, &ffi.TypeVoid); status != ffi.OK {
		panic(status)
	}

	symCloseWindow, err := syscall.GetProcAddress(raylib, "CloseWindow")
	if err != nil {
		panic(err)
	}

	CloseWindow = func() {
		ffi.Call(&cifVoidVoid, symCloseWindow, nil)
	}

	// WindowShouldClose ------------------------
	var cifWindowShouldClose ffi.Cif
	if status := ffi.PrepCif(&cifWindowShouldClose, ffi.DefaultAbi, 0, &ffi.TypeUint8); status != ffi.OK {
		panic(status)
	}

	symWindowShouldClose, err := syscall.GetProcAddress(raylib, "WindowShouldClose")
	if err != nil {
		panic(err)
	}

	WindowShouldClose = func() bool {
		var close ffi.Arg
		ffi.Call(&cifWindowShouldClose, symWindowShouldClose, unsafe.Pointer(&close))
		return byte(close) != 0
	}

	// BeginDrawing -----------------------------
	symBeginDrawing, err := syscall.GetProcAddress(raylib, "BeginDrawing")
	if err != nil {
		panic(err)
	}

	BeginDrawing = func() {
		ffi.Call(&cifVoidVoid, symBeginDrawing, nil)
	}

	// EndDrawing -------------------------------
	symEndDrawing, err := syscall.GetProcAddress(raylib, "EndDrawing")
	if err != nil {
		panic(err)
	}

	EndDrawing = func() {
		ffi.Call(&cifVoidVoid, symEndDrawing, nil)
	}

	// ClearBackground --------------------------
	var cifClearBackground ffi.Cif
	if status := ffi.PrepCif(&cifClearBackground, ffi.DefaultAbi, 1, &ffi.TypeVoid, &TypeColor); status != ffi.OK {
		panic(status)
	}

	symClearBackground, err := syscall.GetProcAddress(raylib, "ClearBackground")
	if err != nil {
		panic(err)
	}

	ClearBackground = func(col color.RGBA) {
		ffi.Call(&cifClearBackground, symClearBackground, nil, unsafe.Pointer(&col))
	}

	// LoadTexture ------------------------------
	var cifLoadTexture ffi.Cif
	if status := ffi.PrepCif(&cifLoadTexture, ffi.DefaultAbi, 1, &TypeTexture, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	symLoadTexture, err := syscall.GetProcAddress(raylib, "LoadTexture")
	if err != nil {
		panic(err)
	}

	LoadTexture = func(filename string) Texture {
		byteFilename := &[]byte(filename + "\x00")[0] // you can also use golang.org/x/sys/windows.BytePtrFromString to create a null-terminated string
		var texture Texture
		ffi.Call(&cifLoadTexture, symLoadTexture, unsafe.Pointer(&texture), unsafe.Pointer(&byteFilename))
		return texture
	}

	// UnloadTexture ----------------------------
	var cifUnloadTexture ffi.Cif
	if status := ffi.PrepCif(&cifUnloadTexture, ffi.DefaultAbi, 1, &ffi.TypeVoid, &TypeTexture); status != ffi.OK {
		panic(status)
	}

	symUnloadTexture, err := syscall.GetProcAddress(raylib, "UnloadTexture")
	if err != nil {
		panic(err)
	}

	UnloadTexture = func(texture Texture) {
		ffi.Call(&cifUnloadTexture, symUnloadTexture, nil, unsafe.Pointer(&texture))
	}

	// DrawTexture ------------------------------
	var cifDrawTexture ffi.Cif
	if status := ffi.PrepCif(&cifDrawTexture, ffi.DefaultAbi, 4, &ffi.TypeVoid, &TypeTexture, &ffi.TypeSint32, &ffi.TypeSint32, &TypeColor); status != ffi.OK {
		panic(status)
	}

	symDrawTexture, err := syscall.GetProcAddress(raylib, "DrawTexture")
	if err != nil {
		panic(err)
	}

	DrawTexture = func(texture Texture, posX, posY int32, col color.RGBA) {
		ffi.Call(&cifDrawTexture, symDrawTexture, nil, unsafe.Pointer(&texture), unsafe.Pointer(&posX), unsafe.Pointer(&posY), unsafe.Pointer(&col))
	}
}

func main() {
	white := color.RGBA{255, 255, 255, 255}

	const width, height = 1280, 720

	InitWindow(width, height, "raylib ffi example")
	defer CloseWindow()

	texture := LoadTexture("examples/structs/raylib/gopher-with-C-book.png")
	defer UnloadTexture(texture)

	for !WindowShouldClose() {
		BeginDrawing()
		ClearBackground(white)
		DrawTexture(texture, width/2-texture.Width/2, height/2-texture.Height/2, white)
		EndDrawing()
	}
}
