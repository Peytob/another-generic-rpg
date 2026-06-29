package resource

type UniformBlock struct {
	Id           uint32
	Name         string
	BindingPoint uint32

	Variables UniformVariables
}

type UniformVariable struct {
	Name   string
	Size   int
	Offset int
}

type UniformVariables map[string]UniformVariable

func NewUniformVariables() UniformVariables {
	return make(map[string]UniformVariable)
}

func (uv UniformVariables) Set(name string, size int, offset int) UniformVariables {
	uv[name] = UniformVariable{
		Name:   name,
		Size:   size,
		Offset: offset,
	}

	return uv
}

func (uv UniformVariables) Get(name string) (UniformVariable, bool) {
	v, ok := uv[name]
	return v, ok
}
