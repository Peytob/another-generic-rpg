package textures

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"game/pkg/files"
	"game/pkg/logger"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const bytesPerPixel = 4

type Service interface {
	CreateTexture(ctx context.Context, name string, width, height int32, pixels []byte) (Texture, error)
	LoadTexture(ctx context.Context, name string) (Texture, error)
	DeleteTexture(ctx context.Context, texture Texture) error
	TextureArray() []Texture
}

type service struct {
	loader            files.Loader
	textureRepository *Repository
}

func InitializeTextures(_ context.Context, loader files.Loader, textureRepository *Repository) (Service, error) {
	return &service{
		loader:            loader,
		textureRepository: textureRepository,
	}, nil
}

func (s *service) CreateTexture(ctx context.Context, name string, width, height int32, pixels []byte) (Texture, error) {
	if width <= 0 || height <= 0 {
		return Texture{}, errors.New("texture size must be positive")
	}

	expected := int64(width) * int64(height) * int64(bytesPerPixel)
	if expected != int64(len(pixels)) {
		return Texture{}, fmt.Errorf("pixel data size does not match texture size: expected %d, got %d", expected, len(pixels))
	}

	textureID := createGlTexture(width, height, pixels)

	texture := Texture{
		ID:     textureID,
		Name:   name,
		Width:  width,
		Height: height,
	}

	if !s.textureRepository.Put(texture) {
		gl.DeleteTextures(1, &textureID)
		return Texture{}, fmt.Errorf("texture %s already registered", name)
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "created texture",
		logger.Uint32("texture_id", texture.ID),
		slog.String("texture_name", texture.Name),
		slog.Int("width", int(width)),
		slog.Int("height", int(height)))

	return texture, nil
}

func (s *service) LoadTexture(ctx context.Context, name string) (Texture, error) {
	image, err := s.loader.LoadImage(ctx, name)
	if err != nil {
		return Texture{}, fmt.Errorf("image loading: %w", err)
	}

	return s.CreateTexture(ctx, name, image.Width, image.Height, image.Pixels)
}

func (s *service) DeleteTexture(ctx context.Context, texture Texture) error {
	if !gl.IsTexture(texture.ID) {
		return errors.New("texture already deleted or id is not texture")
	}

	gl.DeleteTextures(1, &texture.ID)

	if !s.textureRepository.Delete(texture) {
		logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelWarn,
			"texture deleted from ogl but not found in repository", logger.Uint32("texture_id", texture.ID))
	}

	logger.FromCtx(ctx).LogAttrs(ctx, slog.LevelInfo, "texture deleted", logger.Uint32("texture_id", texture.ID))

	return nil
}

func (s *service) TextureArray() []Texture {
	return s.textureRepository.Array()
}

func createGlTexture(width, height int32, pixels []byte) uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return textureID
}
