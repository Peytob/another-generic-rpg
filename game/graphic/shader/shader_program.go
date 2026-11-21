package shader

type Program struct {
	ID   uint32
	Name string
}

type ProgramBuilder struct {
	stages map[Stage]Shader
}

func NewProgramBuilder() *ProgramBuilder {
	return &ProgramBuilder{stages: make(map[Stage]Shader)}
}

func (sb *ProgramBuilder) Set(shader Shader) *ProgramBuilder {
	sb.stages[shader.Stage] = shader
	return sb
}

func (sb *ProgramBuilder) Stage(stage Stage) (Shader, bool) {
	s, ok := sb.stages[stage]
	return s, ok
}
