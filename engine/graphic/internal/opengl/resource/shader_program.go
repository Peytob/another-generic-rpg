package resource

type ShaderProgram uint32

func (sp ShaderProgram) ID() uint32 {
	return uint32(sp)
}
