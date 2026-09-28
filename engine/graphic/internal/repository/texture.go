package repository

import (
	gresource "engine/graphic/resource"
	"engine/repository"
)

type TextureRepository struct {
	idIndex repository.Index[uint32, gresource.Texture]
}

func NewTextureRepository() *TextureRepository {
	return &TextureRepository{
		idIndex: repository.NewIndex[uint32, gresource.Texture](),
	}
}

func (sr *TextureRepository) ByID(id uint32) (gresource.Texture, bool) {
	return sr.idIndex.Get(id)
}

func (sr *TextureRepository) Remove(id uint32) bool {
	return sr.idIndex.Remove(id)
}

func (sr *TextureRepository) Put(texture gresource.Texture) bool {
	return sr.idIndex.Add(texture.Id, texture)
}

func (sr *TextureRepository) All() []gresource.Texture {
	return sr.idIndex.Array()
}
