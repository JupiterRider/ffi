package main

import (
	"image/color"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
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

	raylib, err := purego.Dlopen("libraylib.so", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	// InitWindow -------------------------------
	var cifInitWindow ffi.Cif
	if ok := ffi.PrepCif(&cifInitWindow, ffi.DefaultAbi, 3, &ffi.TypeVoid, []*ffi.Type{&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, nil}); ok != ffi.OK {
		panic("prep failed")
	}

	symInitWindow, err := purego.Dlsym(raylib, "InitWindow")
	if err != nil {
		panic(err)
	}

	InitWindow = func(width, height int32, title string) {
		byteTitle := []byte(title)
		ffi.Call(&cifInitWindow, symInitWindow, nil, []unsafe.Pointer{unsafe.Pointer(&width), unsafe.Pointer(&height), unsafe.Pointer(&byteTitle)})
	}

	// CloseWindow ------------------------------
	var cifVoidVoid ffi.Cif
	if ok := ffi.PrepCif(&cifVoidVoid, ffi.DefaultAbi, 0, &ffi.TypeVoid, nil); ok != ffi.OK {
		panic("prep failed")
	}

	symCloseWindow, err := purego.Dlsym(raylib, "CloseWindow")
	if err != nil {
		panic(err)
	}

	CloseWindow = func() {
		ffi.Call(&cifVoidVoid, symCloseWindow, nil, nil)
	}

	// WindowShouldClose ------------------------
	var cifWindowShouldClose ffi.Cif
	if ok := ffi.PrepCif(&cifWindowShouldClose, ffi.DefaultAbi, 0, &ffi.TypeUint32, nil); ok != ffi.OK {
		panic("prep failed")
	}

	symWindowShouldClose, err := purego.Dlsym(raylib, "WindowShouldClose")
	if err != nil {
		panic(err)
	}

	WindowShouldClose = func() bool {
		close := uint32(0)
		ffi.Call(&cifWindowShouldClose, symWindowShouldClose, unsafe.Pointer(&close), nil)
		return close != 0
	}

	// BeginDrawing -----------------------------
	symBeginDrawing, err := purego.Dlsym(raylib, "BeginDrawing")
	if err != nil {
		panic(err)
	}

	BeginDrawing = func() {
		ffi.Call(&cifVoidVoid, symBeginDrawing, nil, nil)
	}

	// EndDrawing -------------------------------
	symEndDrawing, err := purego.Dlsym(raylib, "EndDrawing")
	if err != nil {
		panic(err)
	}

	EndDrawing = func() {
		ffi.Call(&cifVoidVoid, symEndDrawing, nil, nil)
	}

	// ClearBackground --------------------------
	var cifClearBackground ffi.Cif
	if ok := ffi.PrepCif(&cifClearBackground, ffi.DefaultAbi, 1, &ffi.TypeVoid, []*ffi.Type{&TypeColor, nil}); ok != ffi.OK {
		panic("prep failed")
	}

	symClearBackground, err := purego.Dlsym(raylib, "ClearBackground")
	if err != nil {
		panic(err)
	}

	ClearBackground = func(col color.RGBA) {
		ffi.Call(&cifClearBackground, symClearBackground, nil, []unsafe.Pointer{unsafe.Pointer(&col)})
	}

	// LoadTexture ------------------------------
	var cifLoadTexture ffi.Cif
	if ok := ffi.PrepCif(&cifLoadTexture, ffi.DefaultAbi, 1, &TypeTexture, []*ffi.Type{&ffi.TypePointer, nil}); ok != ffi.OK {
		panic("prep failed")
	}

	symLoadTexture, err := purego.Dlsym(raylib, "LoadTexture")
	if err != nil {
		panic(err)
	}

	LoadTexture = func(filename string) Texture {
		byteFilename := []byte(filename)
		var texture Texture
		ffi.Call(&cifLoadTexture, symLoadTexture, unsafe.Pointer(&texture), []unsafe.Pointer{unsafe.Pointer(&byteFilename)})
		return texture
	}

	// UnloadTexture ----------------------------
	var cifUnloadTexture ffi.Cif
	if ok := ffi.PrepCif(&cifUnloadTexture, ffi.DefaultAbi, 1, &ffi.TypeVoid, []*ffi.Type{&TypeTexture, nil}); ok != ffi.OK {
		panic("prep failed")
	}

	symUnloadTexture, err := purego.Dlsym(raylib, "UnloadTexture")
	if err != nil {
		panic(err)
	}

	UnloadTexture = func(texture Texture) {
		ffi.Call(&cifUnloadTexture, symUnloadTexture, nil, []unsafe.Pointer{unsafe.Pointer(&texture)})
	}

	// DrawTexture ------------------------------
	var cifDrawTexture ffi.Cif
	if ok := ffi.PrepCif(&cifDrawTexture, ffi.DefaultAbi, 4, &ffi.TypeVoid, []*ffi.Type{&TypeTexture, &ffi.TypeSint32, &ffi.TypeSint32, &TypeColor, nil}); ok != ffi.OK {
		panic("prep failed")
	}

	symDrawTexture, err := purego.Dlsym(raylib, "DrawTexture")
	if err != nil {
		panic(err)
	}

	DrawTexture = func(texture Texture, posX, posY int32, col color.RGBA) {
		args := []unsafe.Pointer{unsafe.Pointer(&texture), unsafe.Pointer(&posX), unsafe.Pointer(&posY), unsafe.Pointer(&col)}
		ffi.Call(&cifDrawTexture, symDrawTexture, nil, args)
	}
}

func main() {
	white := color.RGBA{255, 255, 255, 255}

	const width, height = 1280, 720

	InitWindow(width, height, "raylib ffi example")
	defer CloseWindow()

	texture := LoadTexture("examples/raylib/gopher-with-C-book.png")
	defer UnloadTexture(texture)

	for !WindowShouldClose() {
		BeginDrawing()
		ClearBackground(white)
		DrawTexture(texture, width/2-texture.Width/2, height/2-texture.Height/2, white)
		EndDrawing()
	}
}
