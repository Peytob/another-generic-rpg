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

func (sr *ShaderRepository) ByID(id uint32) (gresource.Shader, bool) {
	return sr.idIndex.Get(id)
}

func (sr *ShaderRepository) ByName(name string) (gresource.Shader, bool) {
	return sr.nameIndex.Get(name)
}

func (sr *ShaderRepository) Remove(id uint32) bool {
	uniformBlock, ok := sr.idIndex.Get(id)
	if !ok {
		return false
	}

	sr.idIndex.Remove(uniformBlock.ID)
	sr.nameIndex.Remove(uniformBlock.Name)

	return true
}

func (sr *ShaderRepository) Put(uniformBlock gresource.Shader) bool {
	_, exists := sr.ByID(uniformBlock.ID)
	if exists {
		return false
	}

	sr.idIndex.Add(uniformBlock.ID, uniformBlock)
	sr.nameIndex.Add(uniformBlock.Name, uniformBlock)

	return true
}

func (sr *ShaderRepository) All() []gresource.Shader {
	return sr.idIndex.Array()
}
