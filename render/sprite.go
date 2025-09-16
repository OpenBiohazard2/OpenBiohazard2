package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_SPRITE = 4
	SPRITE_FRAME_TIME  = 0.5 // in seconds
)

var (
	totalSpriteRuntime = float64(0)
	curSpriteFrame     = 0
)

type SpriteGroupEntity struct {
	SpriteTextureIndexMap map[int]int
	TextureIdPool         [][]uint32
	VertexBuffer          []float32
	VertexArrayObject     uint32
	VertexBufferObject    uint32
}

func NewSpriteGroupEntity(spriteData []fileio.SpriteData) *SpriteGroupEntity {
	spriteTextureIds := make([][]uint32, 0)
	for i := 0; i < len(spriteData); i++ {
		spriteFrames := BuildSpriteTexture(spriteData[i])
		spriteTextureIds = append(spriteTextureIds, spriteFrames)
	}

	spriteTextureIndexMap := make(map[int]int)
	for i := 0; i < len(spriteData); i++ {
		spriteTextureIndexMap[spriteData[i].Id] = i
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &SpriteGroupEntity{
		SpriteTextureIndexMap: spriteTextureIndexMap,
		TextureIdPool:         spriteTextureIds,
		VertexBuffer:          make([]float32, 0),
		VertexArrayObject:     vao,
		VertexBufferObject:    vbo,
	}
}

// Each sprite id has its own texture
// Build a texture for each frame
func BuildSpriteTexture(spriteData fileio.SpriteData) []uint32 {
	allFrameTextures := make([]uint32, 0)

	for _, frameData := range spriteData.FrameData {
		spriteId := frameData.SpriteId
		framePosition := spriteData.FramePositions[spriteId]

		frameHeight := int(frameData.SquareSide)
		frameWidth := int(frameData.SquareSide)

		if frameHeight == 0 || frameWidth == 0 {
			continue
		}

		startX := int(framePosition.ImageX)
		startY := int(framePosition.ImageY)

		frameImageColors := make([]uint16, 0)
		for y := startY; y < startY+frameHeight; y++ {
			for x := startX; x < startX+frameWidth; x++ {
				curColor := spriteData.ImageData.PixelData[y][x]

				// Determine if pixel should be transparent
				// Set black to be transparent color
				newTextureColor := curColor
				if curColor > 0 {
					// Set alpha bit to 1
					newTextureColor = uint16(curColor) | (1 << 15)
				}
				frameImageColors = append(frameImageColors, newTextureColor)
			}
		}
		textureId := BuildTexture(frameImageColors, int32(frameWidth), int32(frameHeight))
		allFrameTextures = append(allFrameTextures, textureId)
	}

	return allFrameTextures
}

func (renderDef *RenderDef) AddSprite(sprite fileio.ScriptInstrSceEsprOn) {
	spriteWidth := float32(1024 * 2)

	// Generate billboard sprite
	spriteCenter := mgl32.Vec3{float32(sprite.X), float32(sprite.Y), float32(sprite.Z)}
	squareVertices := [4]mgl32.Vec3{
		{0, 1, 0},
		{1, 1, 0},
		{1, 0, 0},
		{0, 0, 0},
	}

	viewMatrix := renderDef.ViewSystem.Camera.BuildViewMatrix()
	cameraRight := mgl32.Vec3{viewMatrix.At(0, 0), viewMatrix.At(1, 0), viewMatrix.At(2, 0)}
	cameraUp := mgl32.Vec3{viewMatrix.At(0, 1), viewMatrix.At(1, 1), viewMatrix.At(2, 1)}

	renderVertices := [4][]float32{}
	for i := 0; i < 4; i++ {
		x := squareVertices[i].X()
		y := squareVertices[i].Y()
		worldspacePosition := spriteCenter.Add(cameraRight.Mul(x * spriteWidth)).Add(cameraUp.Mul(y * spriteWidth))
		renderVertices[i] = []float32{worldspacePosition.X(), worldspacePosition.Y(), worldspacePosition.Z()}
	}
	uvs := [4][]float32{
		{0.0, 0.0},
		{1.0, 0.0},
		{1.0, 1.0},
		{0.0, 1.0},
	}
	rect := geometry.NewTexturedRectangle(renderVertices, uvs)
	renderDef.SceneSystem.SpriteGroupEntity.VertexBuffer = append(renderDef.SceneSystem.SpriteGroupEntity.VertexBuffer, rect.VertexBuffer...)
}

func RenderSprites(r *RenderDef, spriteGroupEntity *SpriteGroupEntity, timeElapsedSeconds float64) {
	vertexBuffer := spriteGroupEntity.VertexBuffer
	if len(vertexBuffer) == 0 {
		return
	}
	textureIds := spriteGroupEntity.TextureIdPool
	if len(textureIds) == 0 {
		return
	}

	// TODO: Calculate index based on id
	spriteIndex := 0

	// Check when to move on to the next frame
	totalSpriteRuntime += timeElapsedSeconds
	if totalSpriteRuntime > SPRITE_FRAME_TIME {
		totalSpriteRuntime = 0
		curSpriteFrame++
		if curSpriteFrame >= len(textureIds[spriteIndex]) {
			curSpriteFrame = 0
		}
	}

	// Create renderer
	renderer := NewOpenGLRenderer(r.ShaderSystem.GetUniformLocations())

	// Create render config for 2D sprite (position + texture)
	config := renderer.Create2DEntityConfig(
		spriteGroupEntity.VertexArrayObject,
		spriteGroupEntity.VertexBufferObject,
		vertexBuffer,
		textureIds[spriteIndex][curSpriteFrame],
		nil, // No model matrix for sprites
		RENDER_TYPE_SPRITE,
	)

	// Render the sprite
	renderer.RenderEntity(config)
}
