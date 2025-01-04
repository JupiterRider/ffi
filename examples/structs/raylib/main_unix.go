//go:build (freebsd || linux) && (amd64 || arm64)

package main

import (
	"image/color"
	"runtime"
	"unsafe"

	_ "embed"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

type Texture struct {
	ID                             uint32
	Width, Height, Mipmaps, Format int32
}

type Image struct {
	Data                           unsafe.Pointer
	Width, Height, Mipmaps, Format int32
}

//go:embed gopher-with-C-book.png
var gopher []byte

var (
	TypeImage   = ffi.NewType(&ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	TypeTexture = ffi.NewType(&ffi.TypeUint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32)
	TypeColor   = ffi.NewType(&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
)

var (
	InitWindow           func(width, height int32, title string)
	CloseWindow          func()
	WindowShouldClose    func() bool
	BeginDrawing         func()
	EndDrawing           func()
	ClearBackground      func(col color.RGBA)
	LoadImageFromMemory  func(fileType string, fileData []byte) Image
	LoadTextureFromImage func(img Image) Texture
	UnloadImage          func(img Image)
	UnloadTexture        func(texture Texture)
	DrawTexture          func(texture Texture, posX, posY int32, col color.RGBA)
)

func init() {
	runtime.LockOSThread()

	raylib, err := purego.Dlopen("libraylib.so", purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	// InitWindow -------------------------------
	var cifInitWindow ffi.Cif
	if status := ffi.PrepCif(&cifInitWindow, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	symInitWindow, err := purego.Dlsym(raylib, "InitWindow")
	if err != nil {
		panic(err)
	}

	InitWindow = func(width, height int32, title string) {
		byteTitle := &[]byte(title + "\x00")[0] // you can also use golang.org/x/sys/unix.BytePtrFromString to create a null-terminated string
		ffi.Call(&cifInitWindow, symInitWindow, nil, unsafe.Pointer(&width), unsafe.Pointer(&height), unsafe.Pointer(&byteTitle))
	}

	// CloseWindow ------------------------------
	var cifVoidVoid ffi.Cif
	if status := ffi.PrepCif(&cifVoidVoid, ffi.DefaultAbi, 0, &ffi.TypeVoid); status != ffi.OK {
		panic(status)
	}

	symCloseWindow, err := purego.Dlsym(raylib, "CloseWindow")
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

	symWindowShouldClose, err := purego.Dlsym(raylib, "WindowShouldClose")
	if err != nil {
		panic(err)
	}

	WindowShouldClose = func() bool {
		var close ffi.Arg
		ffi.Call(&cifWindowShouldClose, symWindowShouldClose, unsafe.Pointer(&close))
		return close.Bool()
	}

	// BeginDrawing -----------------------------
	symBeginDrawing, err := purego.Dlsym(raylib, "BeginDrawing")
	if err != nil {
		panic(err)
	}

	BeginDrawing = func() {
		ffi.Call(&cifVoidVoid, symBeginDrawing, nil)
	}

	// EndDrawing -------------------------------
	symEndDrawing, err := purego.Dlsym(raylib, "EndDrawing")
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

	symClearBackground, err := purego.Dlsym(raylib, "ClearBackground")
	if err != nil {
		panic(err)
	}

	ClearBackground = func(col color.RGBA) {
		ffi.Call(&cifClearBackground, symClearBackground, nil, unsafe.Pointer(&col))
	}

	// LoadImageFromMemory ----------------------
	var cifLoadImageFromMemory ffi.Cif
	if status := ffi.PrepCif(&cifLoadImageFromMemory, ffi.DefaultAbi, 3, &TypeImage, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); status != ffi.OK {
		panic(status)
	}

	symLoadImageFromMemory, err := purego.Dlsym(raylib, "LoadImageFromMemory")
	if err != nil {
		panic(err)
	}

	LoadImageFromMemory = func(fileType string, fileData []byte) Image {
		byteFileType := &[]byte(fileType + "\x00")[0]
		ptrFileData := &fileData[0]
		dataSize := int32(len(fileData))
		var img Image
		ffi.Call(&cifLoadImageFromMemory, symLoadImageFromMemory, unsafe.Pointer(&img), unsafe.Pointer(&byteFileType), unsafe.Pointer(&ptrFileData), unsafe.Pointer(&dataSize))
		return img
	}

	// LoadTextureFromImage ---------------------
	var cifLoadTextureFromImage ffi.Cif
	if status := ffi.PrepCif(&cifLoadTextureFromImage, ffi.DefaultAbi, 1, &TypeTexture, &TypeImage); status != ffi.OK {
		panic(status)
	}

	symLoadTextureFromImage, err := purego.Dlsym(raylib, "LoadTextureFromImage")
	if err != nil {
		panic(err)
	}

	LoadTextureFromImage = func(img Image) Texture {
		var texture Texture
		ffi.Call(&cifLoadTextureFromImage, symLoadTextureFromImage, unsafe.Pointer(&texture), unsafe.Pointer(&img))
		return texture
	}

	// UnloadImage ----------------------------
	var cifUnloadImage ffi.Cif
	if status := ffi.PrepCif(&cifUnloadImage, ffi.DefaultAbi, 1, &ffi.TypeVoid, &TypeImage); status != ffi.OK {
		panic(status)
	}

	symUnloadImage, err := purego.Dlsym(raylib, "UnloadImage")
	if err != nil {
		panic(err)
	}

	UnloadImage = func(img Image) {
		ffi.Call(&cifUnloadImage, symUnloadImage, nil, unsafe.Pointer(&img))
	}

	// UnloadTexture ----------------------------
	var cifUnloadTexture ffi.Cif
	if status := ffi.PrepCif(&cifUnloadTexture, ffi.DefaultAbi, 1, &ffi.TypeVoid, &TypeTexture); status != ffi.OK {
		panic(status)
	}

	symUnloadTexture, err := purego.Dlsym(raylib, "UnloadTexture")
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

	symDrawTexture, err := purego.Dlsym(raylib, "DrawTexture")
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

	img := LoadImageFromMemory(".png", gopher)
	texture := LoadTextureFromImage(img)
	defer UnloadTexture(texture)
	UnloadImage(img)

	for !WindowShouldClose() {
		BeginDrawing()
		ClearBackground(white)
		DrawTexture(texture, width/2-texture.Width/2, height/2-texture.Height/2, white)
		EndDrawing()
	}
}
