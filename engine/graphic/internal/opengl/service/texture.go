package service

import (
	"engine/graphic/repository"
	gresource "engine/graphic/resource"
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	textureRepository *repository.TextureRepository
}

func NewTexture(textureRepository *repository.TextureRepository) *Texture {
	return &Texture{
		textureRepository: textureRepository,
	}
}

func (t *Texture) LoadTexture(image gresource.Image) (gresource.Texture, error) {
	format, internalFormat, err := imageParams(image)
	if err != nil {
		return gresource.Texture{}, err
	}

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		int32(internalFormat),
		int32(image.Width),
		int32(image.Height),
		0,
		format,
		gl.UNSIGNED_BYTE,
		gl.Ptr(image.Data),
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	texture := gresource.Texture{Id: tex}

	if !t.textureRepository.Put(texture) {
		gl.DeleteTextures(1, &tex)
		return gresource.Texture{}, fmt.Errorf("texture %d already registered", tex)
	}

	return texture, nil
}

func (t *Texture) GetTexture(id uint32) (gresource.Texture, bool) {
	return t.textureRepository.ByID(id)
}

func (t *Texture) DeleteTexture(id uint32) error {
	if !t.textureRepository.Remove(id) {
		return fmt.Errorf("texture %d not found", id)
	}

	gl.DeleteTextures(1, &id)

	return nil
}

// imageParams validates image resource and returns matching OpenGL pixel format
// and internal format.
func imageParams(image gresource.Image) (format, internalFormat uint32, err error) {
	if image.Width <= 0 || image.Height <= 0 {
		return 0, 0, fmt.Errorf("invalid image size %dx%d", image.Width, image.Height)
	}

	var channels int
	switch image.Type {
	case gresource.ImageTypeRGB:
		format, internalFormat, channels = gl.RGB, gl.RGB8, 3
	case gresource.ImageTypeRGBA:
		format, internalFormat, channels = gl.RGBA, gl.RGBA8, 4
	default:
		return 0, 0, fmt.Errorf("unsupported image type %v", image.Type)
	}

	expected := image.Width * image.Height * channels
	if len(image.Data) != expected {
		return 0, 0, fmt.Errorf("image data length %d does not match size %dx%d", len(image.Data), image.Width, image.Height)
	}

	return format, internalFormat, nil
}
