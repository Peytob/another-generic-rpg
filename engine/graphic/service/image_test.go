package service

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"engine/graphic/repository"
	gresource "engine/graphic/resource"
)

func encodePng(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}

	return buf.Bytes()
}

func encodeJpeg(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to encode jpeg: %v", err)
	}

	return buf.Bytes()
}

func newGrayPng(t *testing.T, width, height int) []byte {
	t.Helper()

	img := image.NewGray(image.Rect(0, 0, width, height))
	img.SetGray(0, 0, color.Gray{Y: 255})

	return encodePng(t, img)
}

func newImageService() Image {
	return NewImage(repository.NewImageRepository())
}

func TestImageLoadFsPngGray(t *testing.T) {
	src := image.NewGray(image.Rect(0, 0, 2, 2))
	src.SetGray(0, 0, color.Gray{Y: 255})

	fsys := fstest.MapFS{"img.png": &fstest.MapFile{Data: encodePng(t, src)}}

	svc := newImageService()
	img, err := svc.LoadFs("hero", fsys, "img.png")
	if err != nil {
		t.Fatalf("LoadFs failed: %v", err)
	}

	if img.Name != "hero" {
		t.Errorf("Name = %q, want %q", img.Name, "hero")
	}
	if img.Type != gresource.ImageTypeRGB {
		t.Errorf("Type = %v, want ImageTypeRGB", img.Type)
	}
	if img.Width != 2 || img.Height != 2 {
		t.Errorf("size = %dx%d, want 2x2", img.Width, img.Height)
	}
	if len(img.Data) != 2*2*3 {
		t.Fatalf("len(Data) = %d, want %d", len(img.Data), 2*2*3)
	}
	if got := img.Data[0:3]; !bytes.Equal(got, []byte{255, 255, 255}) {
		t.Errorf("pixel(0,0) = %v, want [255 255 255]", got)
	}

	stored, ok := svc.GetByName("hero")
	if !ok || stored.Name != "hero" {
		t.Errorf("GetByName(hero) = (%v, %v), want stored image", stored, ok)
	}
	if all := svc.All(); len(all) != 1 {
		t.Errorf("All() len = %d, want 1", len(all))
	}
}

func TestImageLoadFsPngPaletted(t *testing.T) {
	src := image.NewPaletted(
		image.Rect(0, 0, 2, 2),
		color.Palette{color.RGBA{R: 255, A: 255}, color.RGBA{B: 128, A: 128}},
	)
	src.SetColorIndex(0, 0, 0)
	src.SetColorIndex(1, 1, 1)

	fsys := fstest.MapFS{"img.png": &fstest.MapFile{Data: encodePng(t, src)}}

	svc := newImageService()
	img, err := svc.LoadFs("pal", fsys, "img.png")
	if err != nil {
		t.Fatalf("LoadFs failed: %v", err)
	}

	if img.Type != gresource.ImageTypeRGBA {
		t.Errorf("Type = %v, want ImageTypeRGBA", img.Type)
	}
	if len(img.Data) != 2*2*4 {
		t.Fatalf("len(Data) = %d, want %d", len(img.Data), 2*2*4)
	}
	if got := img.Data[0:4]; !bytes.Equal(got, []byte{255, 0, 0, 255}) {
		t.Errorf("pixel(0,0) = %v, want [255 0 0 255]", got)
	}
	if got := img.Data[12:16]; !bytes.Equal(got, []byte{0, 0, 255, 128}) {
		t.Errorf("pixel(1,1) = %v, want [0 0 255 128]", got)
	}
}

func TestImageLoadFsUnsupportedNrgba(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	fsys := fstest.MapFS{"img.png": &fstest.MapFile{Data: encodePng(t, src)}}

	svc := newImageService()
	if _, err := svc.LoadFs("x", fsys, "img.png"); err == nil {
		t.Fatal("LoadFs of NRGBA image must fail")
	}
}

func TestImageLoadFsUnsupportedJpeg(t *testing.T) {
	src := image.NewYCbCr(image.Rect(0, 0, 4, 4), image.YCbCrSubsampleRatio444)
	fsys := fstest.MapFS{"img.jpg": &fstest.MapFile{Data: encodeJpeg(t, src)}}

	svc := newImageService()
	if _, err := svc.LoadFs("x", fsys, "img.jpg"); err == nil {
		t.Fatal("LoadFs of YCbCr (JPEG) image must fail")
	}

	if _, ok := svc.GetByName("x"); ok {
		t.Error("failed image must not be registered")
	}
}

func TestImageLoadFsChecksNameBeforeLoad(t *testing.T) {
	svc := newImageService()
	if _, err := svc.LoadFs("dup", fstest.MapFS{"img.png": &fstest.MapFile{Data: newGrayPng(t, 1, 1)}}, "img.png"); err != nil {
		t.Fatalf("first LoadFs failed: %v", err)
	}

	_, err := svc.LoadFs("dup", fstest.MapFS{}, "missing.png")
	if err == nil {
		t.Fatal("second LoadFs with same name must fail")
	}
	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("error = %q, want name conflict detected before load", err)
	}
}

func TestImageLoadFromDisk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "img.png")
	if err := os.WriteFile(path, newGrayPng(t, 3, 1), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	svc := newImageService()
	img, err := svc.Load("disk", path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if img.Width != 3 || img.Height != 1 {
		t.Errorf("size = %dx%d, want 3x1", img.Width, img.Height)
	}
	if _, ok := svc.GetByName("disk"); !ok {
		t.Error("image not registered in repository")
	}
}

func TestImageUnload(t *testing.T) {
	fsys := fstest.MapFS{"img.png": &fstest.MapFile{Data: newGrayPng(t, 1, 1)}}

	svc := newImageService()
	if _, err := svc.LoadFs("a", fsys, "img.png"); err != nil {
		t.Fatalf("LoadFs failed: %v", err)
	}
	if _, err := svc.LoadFs("b", fsys, "img.png"); err != nil {
		t.Fatalf("LoadFs failed: %v", err)
	}

	if err := svc.Unload("a"); err != nil {
		t.Fatalf("Unload(a) failed: %v", err)
	}

	if _, ok := svc.GetByName("a"); ok {
		t.Error("image a still present after Unload")
	}
	if err := svc.Unload("a"); err == nil {
		t.Error("Unload of missing image must fail")
	}
	if _, ok := svc.GetByName("b"); !ok {
		t.Error("image b must stay after unloading image a")
	}
	if all := svc.All(); len(all) != 1 {
		t.Errorf("All() len = %d, want 1", len(all))
	}
}
