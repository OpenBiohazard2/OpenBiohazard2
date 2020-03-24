package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

func NewItemEntities(items []fileio.ScriptInstrItemAotSet,
	itemTextureData []*fileio.TIMOutput,
	itemModelData []*fileio.MD1Output) []SceneMD1Entity {
	itemEntities := make([]SceneMD1Entity, 0)
	for _, item := range items {
		// skip rendering
		if item.Md1ModelId == 255 {
			continue
		}

		itemTextureData := itemTextureData[item.Md1ModelId]
		itemMeshData := itemModelData[item.Md1ModelId]
		itemTexColors := itemTextureData.ConvertToRenderData()
		itemTextureId := BuildTexture(itemTexColors, int32(itemTextureData.ImageWidth), int32(itemTextureData.ImageHeight))
		itemEntityVertexBuffer := BuildEntityComponentVertices(itemMeshData, itemTextureData)

		// position in the center of the trigger region
		modelPosition := mgl32.Vec3{float32(item.X) + float32(item.Width)/2.0, 0.0, float32(item.Z) + float32(item.Depth)/2.0}

		var vao uint32
		gl.GenVertexArrays(1, &vao)

		var vbo uint32
		gl.GenBuffers(1, &vbo)

		sceneEntity := SceneMD1Entity{
			TextureId:          itemTextureId,
			VertexBuffer:       itemEntityVertexBuffer,
			ModelPosition:      modelPosition,
			VertexArrayObject:  vao,
			VertexBufferObject: vbo,
		}
		itemEntities = append(itemEntities, sceneEntity)
	}

	return itemEntities
}
