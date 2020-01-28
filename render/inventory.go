package render

import (
	"../fileio"
	"github.com/go-gl/gl/v4.1-core/gl"
	"math"
)

const (
	RENDER_GAME_STATE_INVENTORY = 1
	ENTITY_INVENTORY_ID         = "INVENTORY_IMAGE"
	// Original dimensions are 320x240
	INVENTORY_IMAGE_WIDTH  = 320
	INVENTORY_IMAGE_HEIGHT = 240
)

func (renderDef *RenderDef) GenerateInventoryImageEntity(inventoryImages []*fileio.TIMOutput) {
	z := float32(0.999)
	vertexBuffer := []float32{
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
		// (-1, -1, z)
		-1.0, -1.0, z, 0.0, 1.0,
		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,

		// (1, -1, z)
		1.0, -1.0, z, 1.0, 1.0,
		// (1, 1, z)
		1.0, 1.0, z, 1.0, 0.0,
		// (-1, 1, z)
		-1.0, 1.0, z, 0.0, 0.0,
	}

	// Add entity to scene
	imageEntity := NewSceneEntity()
	newImageColors := make([]uint16, INVENTORY_IMAGE_WIDTH*INVENTORY_IMAGE_HEIGHT)

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

	imageEntity.SetTexture(newImageColors, INVENTORY_IMAGE_WIDTH, INVENTORY_IMAGE_HEIGHT)
	imageEntity.SetMesh(vertexBuffer)
	renderDef.AddSceneEntity(ENTITY_INVENTORY_ID, imageEntity)
}

func copyPixels(pixelData2D [][]uint16, startX int, startY int, width int, height int,
	newImageColors []uint16, destX int, destY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImageColors[((destY+y)*INVENTORY_IMAGE_WIDTH)+(destX+x)] = pixelData2D[startY+y][startX+x]
		}
	}
}

func fillPixels(newImageColors []uint16, destX int, destY int, width int, height int,
	r int, g int, b int) {
	// Convert color to A1R5G5B5 format
	newR := uint16(math.Round(float64(r) / 8.0))
	newG := uint16(math.Round(float64(g) / 8.0))
	newB := uint16(math.Round(float64(b) / 8.0))
	color := (1 << 15) | (newB << 10) | (newG << 5) | newR

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImageColors[((destY+y)*INVENTORY_IMAGE_WIDTH)+(destX+x)] = color
		}
	}
}

func (renderDef *RenderDef) RenderInventory() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	renderGameStateUniform := gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	gl.Uniform1i(renderGameStateUniform, RENDER_GAME_STATE_INVENTORY)

	renderDef.RenderSceneEntity(renderDef.SceneEntityMap[ENTITY_INVENTORY_ID], RENDER_TYPE_BACKGROUND)
}
