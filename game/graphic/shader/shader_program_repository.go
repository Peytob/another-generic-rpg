package shader

import "game/repository"

type ProgramRepository struct {
	idIndex   repository.Index[uint32, Program]
	nameIndex repository.Index[string, Program]
}

func NewProgramRepository() *ProgramRepository {
	return &ProgramRepository{
		idIndex:   repository.NewIndex[uint32, Program](),
		nameIndex: repository.NewIndex[string, Program](),
	}
}

func (r *ProgramRepository) Put(program Program) bool {
	_, exists := r.idIndex.Get(program.ID)
	if exists {
		return false
	}

	r.idIndex.Add(program.ID, program)
	r.nameIndex.Add(program.Name, program)

	return true
}

func (r *ProgramRepository) Delete(program Program) bool {
	return r.idIndex.Remove(program.ID)
}

func (r *ProgramRepository) GetById(id uint32) (Program, bool) {
	return r.idIndex.Get(id)
}

func (r *ProgramRepository) GetByName(id uint32) (Program, bool) {
	return r.idIndex.Get(id)
}

func (r *ProgramRepository) Array() []Program {
	return r.idIndex.Array()
}
