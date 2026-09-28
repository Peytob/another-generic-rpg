package repository

type Repositories struct {
	Shader  *ShaderRepository
	Uniform *UniformBlockRepository
	Image   *ImageRepository
	Texture *TextureRepository
}
