package uniform

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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

func newUniformVariables() *Variables {
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

func setUniformVariableMat4(uniformBlock Block, variableName string, mat mgl32.Mat4) error {
	variable, found := uniformBlock.Variables.Lookup(variableName)
	if !found {
		return ErrVariableNotFound
	}

	gl.BindBuffer(gl.UNIFORM_BUFFER, uniformBlock.ID)
	gl.BufferSubData(gl.UNIFORM_BUFFER, variable.Offset, variable.Size, gl.Ptr(&mat[0]))

	return nil
}
