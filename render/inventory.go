package render

import (
	"../fileio"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	ENTITY_INVENTORY_ID = "INVENTORY_IMAGE"
)

func (renderDef *RenderDef) GenerateInventoryImageEntity(
	inventoryImages []*fileio.TIMOutput,
	inventoryItemImages []*fileio.TIMOutput) {
	newImageColors := NewSurface2D()
	buildBackground(inventoryImages, newImageColors)
	buildItems(inventoryItemImages, newImageColors)

	imageEntity := NewSceneEntity()
	imageEntity.SetTexture(newImageColors, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT)
	imageEntity.SetMesh(buildSurface2DVertexBuffer())
	renderDef.AddSceneEntity(ENTITY_INVENTORY_ID, imageEntity)
}

func buildItems(inventoryItemImages []*fileio.TIMOutput, newImageColors []uint16) {
	// Item in top right corner
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 35)

	// Empty inventory slots
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 73)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 73)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 103)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 103)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 133)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 133)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 225, 163)
	copyPixels(inventoryItemImages[0].PixelData, 0, 0, 40, 30, newImageColors, 265, 163)

	// Equipped item
	copyPixels(inventoryItemImages[2].PixelData, 40, 90, 80, 30, newImageColors, 172, 35)
}

func buildBackground(inventoryImages []*fileio.TIMOutput, newImageColors []uint16) {
	// The inventory image is split up into many small components
	// Combine them manually back into a single image
	// source image is 256x256
	// dest image is 320x240
	backgroundColor := [3]int{5, 5, 31}
	fillPixels(newImageColors, 0, 0, 320, 240, backgroundColor[0], backgroundColor[1], backgroundColor[2])

	// Player
	copyPixels(inventoryImages[0].PixelData, 106, 152, 4, 60, newImageColors, 7, 16)  // left
	copyPixels(inventoryImages[0].PixelData, 0, 140, 39, 4, newImageColors, 11, 16)   // top
	copyPixels(inventoryImages[0].PixelData, 109, 152, 4, 60, newImageColors, 49, 16) // right
	copyPixels(inventoryImages[0].PixelData, 0, 140, 39, 4, newImageColors, 11, 72)   // bottom
	copyPixels(inventoryImages[1].PixelData, 1, 74, 37, 7, newImageColors, 12, 21)    // player name
	copyPixels(inventoryImages[1].PixelData, 0, 85, 38, 42, newImageColors, 11, 31)   // player image
	copyPixels(inventoryImages[0].PixelData, 56, 164, 38, 1, newImageColors, 11, 30)  // line between name and image

	// Pipes to the left of player image
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 17)
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 33)
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 0, 49)

	// Pipes to the right of player image
	copyPixels(inventoryImages[0].PixelData, 56, 186, 7, 7, newImageColors, 53, 32)
	copyPixels(inventoryImages[0].PixelData, 56, 186, 7, 7, newImageColors, 53, 60)

	// Health bar
	copyPixels(inventoryImages[0].PixelData, 0, 92, 99, 47, newImageColors, 60, 29)
	for i := 0; i < 8; i++ {
		fillPixels(newImageColors, 129-i, 68+i, 30+i, 1, backgroundColor[0], backgroundColor[1], backgroundColor[2])
	}

	// Equipped item
	copyPixels(inventoryImages[0].PixelData, 50, 211, 11, 39, newImageColors, 161, 29) // left
	copyPixels(inventoryImages[0].PixelData, 0, 158, 80, 6, newImageColors, 172, 29)   // top
	copyPixels(inventoryImages[0].PixelData, 91, 164, 5, 39, newImageColors, 252, 29)  // right
	copyPixels(inventoryImages[0].PixelData, 0, 155, 80, 3, newImageColors, 172, 65)   // bottom

	// Extra item
	copyPixels(inventoryImages[0].PixelData, 0, 211, 50, 41, newImageColors, 260, 29)

	// File
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 111, 16)

	// Map
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 160, 16)

	// Item
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 209, 16)

	// Exit
	copyPixels(inventoryImages[0].PixelData, 3, 164, 49, 12, newImageColors, 258, 16)

	// Item slots
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, 220, 73) // left
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 3, newImageColors, 220, 70)   // top
	copyPixels(inventoryImages[0].PixelData, 114, 92, 5, 120, newImageColors, 305, 73) // right
	copyPixels(inventoryImages[0].PixelData, 0, 140, 90, 4, newImageColors, 220, 193)  // bottom

	// Description
	descriptionColor := [3]int{6, 13, 23}
	fillPixels(newImageColors, 13, 174, 201, 49, descriptionColor[0], descriptionColor[1], descriptionColor[2])
	copyPixels(inventoryImages[0].PixelData, 106, 163, 5, 49, newImageColors, 8, 174)   // left
	copyPixels(inventoryImages[0].PixelData, 0, 147, 83, 4, newImageColors, 8, 170)     // top left
	copyPixels(inventoryImages[0].PixelData, 0, 80, 128, 4, newImageColors, 91, 170)    // top right
	copyPixels(inventoryImages[0].PixelData, 106, 163, 5, 49, newImageColors, 214, 174) // right
	copyPixels(inventoryImages[0].PixelData, 0, 147, 83, 4, newImageColors, 8, 223)     // bottom left
	copyPixels(inventoryImages[0].PixelData, 0, 80, 128, 4, newImageColors, 91, 223)    // bottom right

	// Pipes to the right of description
	copyPixels(inventoryImages[0].PixelData, 107, 242, 7, 14, newImageColors, 219, 212)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 35, 7, newImageColors, 226, 215)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 35, 7, newImageColors, 261, 215)
	copyPixels(inventoryImages[0].PixelData, 56, 178, 24, 7, newImageColors, 296, 215)
}

func (renderDef *RenderDef) RenderInventory() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	renderGameStateUniform := gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	gl.Uniform1i(renderGameStateUniform, RENDER_GAME_STATE_BACKGROUND_SOLID)

	renderDef.RenderSceneEntity(renderDef.SceneEntityMap[ENTITY_INVENTORY_ID], RENDER_TYPE_BACKGROUND)
}
