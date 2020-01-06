package render

type SceneEntity struct {
	TextureId    uint32    // texture id in OpenGL
	VertexBuffer []float32 // 3 elements for x,y,z and 2 elements for texture u,v
}

func NewSceneEntity() *SceneEntity {
	return &SceneEntity{
		TextureId:    0xFFFFFFFF,
		VertexBuffer: []float32{},
	}
}

func (entity *SceneEntity) SetTexture(imagePixels []uint16, imageWidth int32, imageHeight int32) {
	entity.TextureId = BuildTexture(imagePixels, imageWidth, imageHeight)
}

func (entity *SceneEntity) SetMesh(vertexBuffer []float32) {
	entity.VertexBuffer = vertexBuffer
}
