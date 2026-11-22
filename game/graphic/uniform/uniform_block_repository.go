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
	if _, exists := r.idIndex.Get(block.ID); exists {
		return false
	}

	if _, exists := r.nameIndex.Get(block.Name); exists {
		return false
	}

	r.idIndex.Add(block.ID, block)
	r.nameIndex.Add(block.Name, block)

	return true
}

func (r *BlockRepository) Delete(block Block) bool {
	okID := r.idIndex.Remove(block.ID)
	okName := r.nameIndex.Remove(block.Name)

	return okID || okName
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
