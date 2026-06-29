package resource

type ShaderProgram uint32

func (sp ShaderProgram) Id() uint32 {
	return uint32(sp)
}

type ShaderProgramBuilder struct {
	Vertex   Shader
	Fragment Shader
}
