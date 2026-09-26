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
	shader, ok := sr.idIndex.Get(id)
	if !ok {
		return false
	}

	sr.idIndex.Remove(shader.ID)
	sr.nameIndex.Remove(shader.Name)

	return true
}

func (sr *ShaderRepository) Put(shader gresource.Shader) bool {
	_, exists := sr.ByID(shader.ID)
	if exists {
		return false
	}

	sr.idIndex.Add(shader.ID, shader)
	sr.nameIndex.Add(shader.Name, shader)

	return true
}

func (sr *ShaderRepository) All() []gresource.Shader {
	return sr.idIndex.Array()
}
