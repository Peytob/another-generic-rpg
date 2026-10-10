package uniform

import "game/repository"

type BlockRepository struct {
	idIndex   repository.Index[uint32, Block]
	nameIndex repository.Index[string, Block]
}

func NewUniformBlockRepository() *BlockRepository {
	return &BlockRepository{
		idIndex:   repository.NewIndex[uint32, Block](),
		nameIndex: repository.NewIndex[string, Block](),
	}
}

func (r *BlockRepository) Put(block Block) bool {
	_, exists := r.idIndex.Get(block.ID)
	if exists {
		return false
	}

	r.idIndex.Add(block.ID, block)
	r.nameIndex.Add(block.Name, block)

	return true
}

func (r *BlockRepository) Delete(block Block) bool {
	ok := r.idIndex.Remove(block.ID)
	ok = ok || r.nameIndex.Remove(block.Name)
	return ok
}

func (r *BlockRepository) GetByID(id uint32) (Block, bool) {
	return r.idIndex.Get(id)
}

func (r *BlockRepository) GetByName(name string) (Block, bool) {
	return r.nameIndex.Get(name)
}

func (r *BlockRepository) Array() []Block {
	return r.idIndex.Array()
}
