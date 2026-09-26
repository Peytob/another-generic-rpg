package repository

import (
	gresource "engine/graphic/resource"
	"engine/repository"
)

type UniformBlockRepository struct {
	idIndex   repository.Index[uint32, gresource.UniformBlock]
	nameIndex repository.Index[string, gresource.UniformBlock]
}

func NewUniformBlockRepository() *UniformBlockRepository {
	return &UniformBlockRepository{
		idIndex:   repository.NewIndex[uint32, gresource.UniformBlock](),
		nameIndex: repository.NewIndex[string, gresource.UniformBlock](),
	}
}

func (ubr *UniformBlockRepository) ByID(id uint32) (gresource.UniformBlock, bool) {
	return ubr.idIndex.Get(id)
}

func (ubr *UniformBlockRepository) ByName(name string) (gresource.UniformBlock, bool) {
	return ubr.nameIndex.Get(name)
}

func (ubr *UniformBlockRepository) Remove(id uint32) bool {
	uniformBlock, ok := ubr.idIndex.Get(id)
	if !ok {
		return false
	}

	ubr.idIndex.Remove(uniformBlock.ID)
	ubr.nameIndex.Remove(uniformBlock.Name)

	return true
}

func (ubr *UniformBlockRepository) Put(uniformBlock gresource.UniformBlock) bool {
	_, exists := ubr.ByID(uniformBlock.ID)
	if exists {
		return false
	}

	ubr.idIndex.Add(uniformBlock.ID, uniformBlock)
	ubr.nameIndex.Add(uniformBlock.Name, uniformBlock)

	return true
}

func (ubr *UniformBlockRepository) All() []gresource.UniformBlock {
	return ubr.idIndex.Array()
}
