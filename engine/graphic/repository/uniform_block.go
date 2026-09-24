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

func (ub *UniformBlockRepository) ByID(id uint32) (gresource.UniformBlock, bool) {
	return ub.idIndex.Get(id)
}

func (ub *UniformBlockRepository) ByName(name string) (gresource.UniformBlock, bool) {
	return ub.nameIndex.Get(name)
}

func (ub *UniformBlockRepository) Remove(id uint32) bool {
	uniformBlock, ok := ub.idIndex.Get(id)
	if !ok {
		return false
	}

	ub.idIndex.Remove(uniformBlock.ID)
	ub.nameIndex.Remove(uniformBlock.Name)

	return true
}

func (ub *UniformBlockRepository) Put(uniformBlock gresource.UniformBlock) bool {
	_, exists := ub.ByID(uniformBlock.ID)
	if exists {
		return false
	}

	ub.idIndex.Add(uniformBlock.ID, uniformBlock)
	ub.nameIndex.Add(uniformBlock.Name, uniformBlock)

	return true
}

func (ub *UniformBlockRepository) All() []gresource.UniformBlock {
	return ub.idIndex.Array()
}
