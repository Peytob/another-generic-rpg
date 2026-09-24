package repository

import (
	gresource "engine/graphic/resource"
	"engine/repository"
)

type ShaderRepository struct {
	idIndex   repository.Index[uint32, gresource.Shader]
	nameIndex repository.Index[string, gresource.Shader]
}

func NewShaderRepository() *ShaderRepository {
	return &ShaderRepository{
		idIndex:   repository.NewIndex[uint32, gresource.Shader](),
		nameIndex: repository.NewIndex[string, gresource.Shader](),
	}
}

func (ub *ShaderRepository) ByID(id uint32) (gresource.Shader, bool) {
	return ub.idIndex.Get(id)
}

func (ub *ShaderRepository) ByName(name string) (gresource.Shader, bool) {
	return ub.nameIndex.Get(name)
}

func (ub *ShaderRepository) Remove(id uint32) bool {
	uniformBlock, ok := ub.idIndex.Get(id)
	if !ok {
		return false
	}

	ub.idIndex.Remove(uniformBlock.ID)
	ub.nameIndex.Remove(uniformBlock.Name)

	return true
}

func (ub *ShaderRepository) Put(uniformBlock gresource.Shader) bool {
	_, exists := ub.ByID(uniformBlock.ID)
	if exists {
		return false
	}

	ub.idIndex.Add(uniformBlock.ID, uniformBlock)
	ub.nameIndex.Add(uniformBlock.Name, uniformBlock)

	return true
}

func (ub *ShaderRepository) All() []gresource.Shader {
	return ub.idIndex.Array()
}
