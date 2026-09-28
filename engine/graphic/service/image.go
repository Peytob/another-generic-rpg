package service

import (
	"engine/graphic/internal/repository"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"

	gresource "engine/graphic/resource"
)

type ImageAccessor interface {
	GetByName(name string) (gresource.Image, bool)
	All() []gresource.Image
}

type Image interface {
	ImageAccessor

	Load(name string, path string) (gresource.Image, error)
	LoadFs(name string, fs fs.FS, path string) (gresource.Image, error)

	Unload(name string) error
}

type imageService struct {
	images *repository.ImageRepository
}

func NewImage(images *repository.ImageRepository) Image {
	return imageService{
		images: images,
	}
}

func (i imageService) GetByName(name string) (gresource.Image, bool) {
	return i.images.ByName(name)
}

func (i imageService) All() []gresource.Image {
	return i.images.All()
}

func (i imageService) Load(name string, path string) (gresource.Image, error) {
	dir, file := filepath.Split(filepath.ToSlash(path))
	if dir == "" {
		dir = "."
	}

	return i.LoadFs(name, os.DirFS(dir), file)
}

func (i imageService) LoadFs(name string, fsys fs.FS, path string) (gresource.Image, error) {
	if _, exists := i.images.ByName(name); exists {
		return gresource.Image{}, fmt.Errorf("image %s already registered", name)
	}

	file, err := fsys.Open(path)
	if err != nil {
		return gresource.Image{}, fmt.Errorf("failed to open image %s: %w", path, err)
	}
	defer file.Close()

	decoded, _, err := image.Decode(file)
	if err != nil {
		return gresource.Image{}, fmt.Errorf("failed to decode image %s: %w", path, err)
	}

	imgType, data, err := rawPixels(decoded)
	if err != nil {
		return gresource.Image{}, err
	}

	img := gresource.Image{
		Name:   name,
		Type:   imgType,
		Width:  decoded.Bounds().Dx(),
		Height: decoded.Bounds().Dy(),
		Data:   data,
	}

	if !i.images.Put(img) {
		return gresource.Image{}, fmt.Errorf("image %s already registered", name)
	}

	return img, nil
}

func (i imageService) Unload(name string) error {
	if !i.images.Remove(name) {
		return fmt.Errorf("image %s not found", name)
	}

	return nil
}

// rawPixels converts decoded image into raw RGB or RGBA pixel data.
// Unsupported decoded image types result in an error.
func rawPixels(img image.Image) (gresource.ImageType, []byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	switch src := img.(type) {
	case *image.RGBA:
		return gresource.ImageTypeRGBA, copyRows(src.Pix, src.Stride, width, height, 4), nil
	case *image.Gray:
		data := make([]byte, width*height*3)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				v := src.GrayAt(x, y).Y
				offset := (y*width + x) * 3
				data[offset] = v
				data[offset+1] = v
				data[offset+2] = v
			}
		}
		return gresource.ImageTypeRGB, data, nil
	case *image.Paletted:
		hasAlpha := paletteHasAlpha(src.Palette)
		channels := 3
		if hasAlpha {
			channels = 4
		}

		data := make([]byte, width*height*channels)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
				offset := (y*width + x) * channels
				data[offset] = c.R
				data[offset+1] = c.G
				data[offset+2] = c.B
				if hasAlpha {
					data[offset+3] = c.A
				}
			}
		}

		if hasAlpha {
			return gresource.ImageTypeRGBA, data, nil
		}
		return gresource.ImageTypeRGB, data, nil
	default:
		return gresource.ImageTypeUnknown, nil, fmt.Errorf("unsupported image type %T", img)
	}
}

func copyRows(pix []byte, stride, width, height, channels int) []byte {
	if stride == width*channels {
		return pix
	}

	data := make([]byte, width*height*channels)
	for y := 0; y < height; y++ {
		copy(data[y*width*channels:(y+1)*width*channels], pix[y*stride:y*stride+width*channels])
	}

	return data
}

func paletteHasAlpha(palette color.Palette) bool {
	for _, c := range palette {
		_, _, _, a := c.RGBA()
		if a != 0xffff {
			return true
		}
	}

	return false
}
