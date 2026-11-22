package shader

import "errors"

var (
	ErrShaderCompile            = errors.New("failed to compile shader")
	ErrProgramLink              = errors.New("failed to link shader program")
	ErrVertexShaderRequired     = errors.New("vertex shader required to create shader program")
	ErrFragmentShaderRequired   = errors.New("fragment shader required to create shader program")
	ErrProgramAlreadyRegistered = errors.New("shader program already registered")
	ErrInvalidProgram           = errors.New("program already deleted or id is not program")
	ErrInvalidShader            = errors.New("shader already deleted or id is not shader")
)
