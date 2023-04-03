package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
)

type ItemGroupEntity struct {
	ItemTextureData []*fileio.TIMOutput
	ItemModelData   []*fileio.MD1Output
	ModelObjectData []*SceneMD1Entity
}

func NewItemGroupEntity() *ItemGroupEntity {
	modelObjectData := make([]*SceneMD1Entity, 32)
	for i := 0; i < len(modelObjectData); i++ {
		var vao uint32
		gl.GenVertexArrays(1, &vao)

		var vbo uint32
		gl.GenBuffers(1, &vbo)

		modelObjectData[i] = &SceneMD1Entity{
			VertexArrayObject:  vao,
			VertexBufferObject: vbo,
			TextureId:          0,
			VertexBuffer:       []float32{},
			ModelPosition:      mgl32.Vec3{},
			RotationAngle:      0,
		}
	}

	return &ItemGroupEntity{
		ItemTextureData: make([]*fileio.TIMOutput, 0),
		ItemModelData:   make([]*fileio.MD1Output, 0),
		ModelObjectData: modelObjectData,
	}
}

func (renderDef *RenderDef) SetItemEntity(instruction fileio.ScriptInstrObjModelSet) {
	modelIndex := int(int(instruction.ObjectIndex))
	position := mgl32.Vec3{float32(instruction.Position[0]), float32(instruction.Position[1]), float32(instruction.Position[2])}
	rotationAngle := (float32(instruction.Direction[1]) / 4096.0) * 360.0

	// skip rendering
	if modelIndex == 255 {
		return
	}

	itemTextureData := renderDef.ItemGroupEntity.ItemTextureData[modelIndex]
	itemMeshData := renderDef.ItemGroupEntity.ItemModelData[modelIndex]

	itemTextureId := NewTextureTIM(itemTextureData)
	itemEntityVertexBuffer := geometry.NewMD1Geometry(itemMeshData, itemTextureData)

	// Update model object
	itemEntity := renderDef.ItemGroupEntity.ModelObjectData[modelIndex]
	itemEntity.TextureId = itemTextureId
	itemEntity.VertexBuffer = itemEntityVertexBuffer
	itemEntity.ModelPosition = position
	itemEntity.RotationAngle = rotationAngle
	renderDef.ItemGroupEntity.ModelObjectData[modelIndex] = itemEntity
}
