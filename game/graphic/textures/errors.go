package textures

import "errors"

var (
	ErrTextureSize              = errors.New("texture size must be positive")
	ErrPixelDataSize            = errors.New("pixel data size does not match texture size")
	ErrTextureAlreadyRegistered = errors.New("texture already registered")
	ErrInvalidTexture           = errors.New("texture already deleted or id is not texture")
)
