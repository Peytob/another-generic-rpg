package resource

type UniformBlock struct {
	ID           uint32
	Name         string
	BindingPoint uint32

	Variables *UniformVariables
}

type UniformVariable struct {
	Name   string
	Size   int
	Offset int
}

type UniformVariables struct {
	variables map[string]UniformVariable
}

func NewUniformVariables() *UniformVariables {
	return &UniformVariables{variables: make(map[string]UniformVariable)}
}

func (uv *UniformVariables) Set(name string, size int, offset int) *UniformVariables {
	uv.variables[name] = UniformVariable{
		Name:   name,
		Size:   size,
		Offset: offset,
	}

	return uv
}

func (uv *UniformVariables) Lookup(name string) (UniformVariable, bool) {
	if uv == nil {
		return UniformVariable{}, false
	}
	v, ok := uv.variables[name]
	return v, ok
}

func (uv *UniformVariables) TotalSize() int {
	if uv == nil {
		return 0
	}

	total := 0
	for _, v := range uv.variables {
		total = max(v.Size+v.Offset, total)
	}

	return total
}
