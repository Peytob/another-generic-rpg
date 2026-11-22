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
	if _, exists := r.idIndex.Get(program.ID); exists {
		return false
	}

	if _, exists := r.nameIndex.Get(program.Name); exists {
		return false
	}

	r.idIndex.Add(program.ID, program)
	r.nameIndex.Add(program.Name, program)

	return true
}

func (r *ProgramRepository) Delete(program Program) bool {
	okID := r.idIndex.Remove(program.ID)
	okName := r.nameIndex.Remove(program.Name)

	return okID || okName
}

func (r *ProgramRepository) GetByID(id uint32) (Program, bool) {
	return r.idIndex.Get(id)
}

func (r *ProgramRepository) GetByName(name string) (Program, bool) {
	return r.nameIndex.Get(name)
}

func (r *ProgramRepository) Array() []Program {
	return r.idIndex.Array()
}
