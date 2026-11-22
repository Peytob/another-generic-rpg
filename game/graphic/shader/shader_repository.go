package shader

import "game/repository"

type Repository struct {
	idIndex repository.Index[uint32, Shader]
}

func NewShaderRepository() *Repository {
	return &Repository{
		idIndex: repository.NewIndex[uint32, Shader](),
	}
}

func (r *Repository) Put(shader Shader) bool {
	_, exists := r.idIndex.Get(shader.ID)
	if exists {
		return false
	}

	r.idIndex.Add(shader.ID, shader)

	return true
}

func (r *Repository) Delete(shader Shader) bool {
	return r.idIndex.Remove(shader.ID)
}

func (r *Repository) GetByID(id uint32) (Shader, bool) {
	return r.idIndex.Get(id)
}

func (r *Repository) Array() []Shader {
	return r.idIndex.Array()
}
