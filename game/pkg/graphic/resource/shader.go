package resource

// ShaderStage Graphic pipeline stage
type ShaderStage uint32

const (
	// UnknownStage placeholder for unknown stage shaders
	UnknownStage = ShaderStage(iota)

	// VertexStage vertex stage shader
	// VERTEX_SHADER in OGL backend
	VertexStage

	// TesselationControlStage vertex stage shader
	// GL_TESS_CONTROL_SHADER in OGL backend
	TesselationControlStage

	// TesselationEvaluationStage vertex stage shader
	// GL_TESS_EVALUATION_SHADER in OGL backend
	TesselationEvaluationStage

	// GeometryStage geometry stage shader
	// GEOMETRY_SHADER in OGL backend
	GeometryStage

	// FragmentStage fragment stage shader
	// GL_FRAGMENT_SHADER in OGL backend
	FragmentStage
)

// LoadedShaderStage temporary type used while shaders loading. Contains information about loaded shader stage
type LoadedShaderStage struct {
	ID    uint32
	Stage ShaderStage
}

// Shader facade for compiled shader programs (or pipelines in other terms), not only separated stage code
type Shader struct {
	ID   uint32
	Name string
}

type ShaderBuilder struct {
	stages map[ShaderStage]LoadedShaderStage
}

func NewShaderBuilder() *ShaderBuilder {
	return &ShaderBuilder{stages: make(map[ShaderStage]LoadedShaderStage, 6)}
}

func (sb *ShaderBuilder) Set(stage LoadedShaderStage) *ShaderBuilder {
	sb.stages[stage.Stage] = stage
	return sb
}

func (sb *ShaderBuilder) Stage(stage ShaderStage) LoadedShaderStage {
	return sb.stages[stage]
}
