package engine

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type Image struct {
	format         uint32
	internalFormat int32
	width          int32
	height         int32
	pixelType      uint32
	data           unsafe.Pointer
}

func MakeImage(path string) (Image, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return Image{}, err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return Image{}, err
	}

	// exctract rgba values
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Image{}, fmt.Errorf("Image not power of 2")
	}

	return Image{
		uint32(gl.RGBA),
		int32(gl.RGBA),
		//int32(gl.SRGB_ALPHA),
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		uint32(gl.UNSIGNED_BYTE),
		gl.Ptr(rgba.Pix),
	}, nil
}

type RawImageData struct {
	data   []uint8
	width  int32
	height int32
}

func MakeRawImageData(path string) (RawImageData, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return RawImageData{}, err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return RawImageData{}, err
	}

	// exctract rgba values
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return RawImageData{}, fmt.Errorf("Image not power of 2")
	}

	return RawImageData{
		rgba.Pix,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
	}, nil
}

func (data *RawImageData) GetWidth() int32 {
	return data.width
}
func (data *RawImageData) GetHeight() int32 {
	return data.height
}
func (data *RawImageData) GetR(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx]
}
func (data *RawImageData) GetG(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+1]
}
func (data *RawImageData) GetB(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+2]
}
func (data *RawImageData) GetA(x, y int32) uint8 {
	idx := data.getIdx(x, y)
	return data.data[idx+3]
}
func (data *RawImageData) GetRGB(x, y int32) (uint8, uint8, uint8) {
	idx := data.getIdx(x, y)
	return data.data[idx], data.data[idx+1], data.data[idx+2]
}
func (data *RawImageData) GetRGBA(x, y int32) (uint8, uint8, uint8, uint8) {
	idx := data.getIdx(x, y)
	return data.data[idx], data.data[idx+1], data.data[idx+2], data.data[idx+3]
}

func (data *RawImageData) getIdx(x, y int32) int32 {
	return (y*data.width + x) * 4
}
