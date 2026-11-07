package textures

import "game/repository"

type Repository struct {
	idIndex   repository.Index[uint32, Texture]
	nameIndex repository.Index[string, Texture]
}

func NewTextureRepository() *Repository {
	return &Repository{
		idIndex:   repository.NewIndex[uint32, Texture](),
		nameIndex: repository.NewIndex[string, Texture](),
	}
}

func (r *Repository) Put(texture Texture) bool {
	_, exists := r.idIndex.Get(texture.ID)
	if exists {
		return false
	}

	r.idIndex.Add(texture.ID, texture)
	r.nameIndex.Add(texture.Name, texture)

	return true
}

func (r *Repository) Delete(texture Texture) bool {
	ok := r.idIndex.Remove(texture.ID)
	ok = ok || r.nameIndex.Remove(texture.Name)

	return ok
}

func (r *Repository) GetByID(id uint32) (Texture, bool) {
	return r.idIndex.Get(id)
}

func (r *Repository) GetByName(name string) (Texture, bool) {
	return r.nameIndex.Get(name)
}

func (r *Repository) Array() []Texture {
	return r.idIndex.Array()
}
