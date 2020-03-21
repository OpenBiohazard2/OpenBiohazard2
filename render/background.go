package render

const (
	RENDER_TYPE_BACKGROUND  = 1
	ENTITY_BACKGROUND_ID    = "ENTITY_BACKGROUND"
	BACKGROUND_IMAGE_WIDTH  = 320
	BACKGROUND_IMAGE_HEIGHT = 240
)

func (backgroundImageEntity *SceneEntity) UpdateBackgroundImageEntity(renderDef *RenderDef, backgroundImageColors []uint16) {
	// The background image is a rectangle that covers the entire screen
	// It should be drawn in the back
	z := float32(0.999)
	backgroundVertexBuffer := []float32{
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
	backgroundImageEntity.SetTexture(backgroundImageColors, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
	backgroundImageEntity.SetMesh(backgroundVertexBuffer)
	return
}
