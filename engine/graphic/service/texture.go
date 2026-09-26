package service

import gresource "engine/graphic/resource"

type Texture interface {
	LoadTexture(image gresource.Image) (gresource.Texture, error)
	GetTexture(id uint32) (gresource.Texture, bool)
	DeleteTexture(id uint32) error
}
