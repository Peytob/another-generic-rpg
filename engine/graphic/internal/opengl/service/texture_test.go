package service

import (
	"testing"

	gresource "engine/graphic/resource"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func TestImageParamsRGB(t *testing.T) {
	img := gresource.Image{
		Type:   gresource.ImageTypeRGB,
		Width:  2,
		Height: 3,
		Data:   make([]byte, 2*3*3),
	}

	format, internalFormat, err := imageParams(img)
	if err != nil {
		t.Fatalf("imageParams failed: %v", err)
	}

	if format != gl.RGB {
		t.Errorf("format = %d, want RGB", format)
	}
	if internalFormat != gl.RGB8 {
		t.Errorf("internalFormat = %d, want RGB8", internalFormat)
	}
}

func TestImageParamsRGBA(t *testing.T) {
	img := gresource.Image{
		Type:   gresource.ImageTypeRGBA,
		Width:  4,
		Height: 2,
		Data:   make([]byte, 4*2*4),
	}

	format, internalFormat, err := imageParams(img)
	if err != nil {
		t.Fatalf("imageParams failed: %v", err)
	}

	if format != gl.RGBA {
		t.Errorf("format = %d, want RGBA", format)
	}
	if internalFormat != gl.RGBA8 {
		t.Errorf("internalFormat = %d, want RGBA8", internalFormat)
	}
}

func TestImageParamsInvalidType(t *testing.T) {
	img := gresource.Image{
		Type:   gresource.ImageTypeUnknown,
		Width:  1,
		Height: 1,
		Data:   make([]byte, 1),
	}

	if _, _, err := imageParams(img); err == nil {
		t.Fatal("imageParams must fail for unknown image type")
	}
}

func TestImageParamsInvalidSize(t *testing.T) {
	img := gresource.Image{
		Type:   gresource.ImageTypeRGB,
		Width:  0,
		Height: 1,
		Data:   nil,
	}

	if _, _, err := imageParams(img); err == nil {
		t.Fatal("imageParams must fail for zero size")
	}
}

func TestImageParamsDataLengthMismatch(t *testing.T) {
	img := gresource.Image{
		Type:   gresource.ImageTypeRGB,
		Width:  2,
		Height: 2,
		Data:   make([]byte, 5),
	}

	if _, _, err := imageParams(img); err == nil {
		t.Fatal("imageParams must fail for data length mismatch")
	}
}
