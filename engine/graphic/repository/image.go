package repository

import (
	gresource "engine/graphic/resource"
	"engine/repository"
)

type ImageRepository struct {
	nameIndex repository.Index[string, gresource.Image]
}

func NewImageRepository() *ImageRepository {
	return &ImageRepository{
		nameIndex: repository.NewIndex[string, gresource.Image](),
	}
}

func (ir *ImageRepository) ByName(name string) (gresource.Image, bool) {
	return ir.nameIndex.Get(name)
}

func (ir *ImageRepository) Remove(name string) bool {
	return ir.nameIndex.Remove(name)
}

func (ir *ImageRepository) Put(image gresource.Image) bool {
	return ir.nameIndex.Add(image.Name, image)
}

func (ir *ImageRepository) All() []gresource.Image {
	return ir.nameIndex.Array()
}

func (ir *ImageRepository) Clear() {
	ir.nameIndex.Clear()
}
