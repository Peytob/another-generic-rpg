package shader

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Stage uint32

const (
	Unknown       = Stage(0)
	VertexStage   = Stage(gl.VERTEX_SHADER)
	GeometryStage = Stage(gl.GEOMETRY_SHADER)
	FragmentStage = Stage(gl.FRAGMENT_SHADER)
)

type Shader struct {
	ID    uint32
	Stage Stage
}
