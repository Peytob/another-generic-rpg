package uniform

type Block struct {
	ID           uint32
	Name         string
	BindingPoint uint32

	Variables *Variables
}

type Variable struct {
	Name   string
	Size   int
	Offset int
}

type Variables struct {
	variables map[string]Variable
}

func NewUniformVariables() *Variables {
	return &Variables{variables: make(map[string]Variable)}
}

func (uv *Variables) Set(name string, size int, offset int) *Variables {
	uv.variables[name] = Variable{
		Name:   name,
		Size:   size,
		Offset: offset,
	}

	return uv
}

func (uv *Variables) Lookup(name string) (Variable, bool) {
	if uv == nil {
		return Variable{}, false
	}
	v, ok := uv.variables[name]
	return v, ok
}

func (uv *Variables) TotalSize() int {
	if uv == nil {
		return 0
	}

	total := 0
	for _, v := range uv.variables {
		total = max(v.Size+v.Offset, total)
	}

	return total
}
