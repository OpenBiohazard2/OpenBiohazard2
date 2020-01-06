package render

import (
	"../fileio"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_CAMERA_MASK = 2
	ENTITY_CAMERA_MASK_ID   = "ENTITY_CAMERA_MASK"
	CAMERA_MASK_WIDTH       = 256
	CAMERA_MASK_HEIGHT      = 256
)

// Normalize the z coordinate to be between 0 and 1
// 0 is closer to the camera, 1 is farther from the camera
func normalizeMaskDepth(depth float32, renderDef *RenderDef) float32 {
	cameraDir := renderDef.Camera.GetDirection().Normalize()
	cameraFrom := renderDef.Camera.CameraFrom
	transformMatrix := renderDef.ProjectionMatrix.Mul4(renderDef.ViewMatrix)

	// Actual distance from camera is 32 * depth
	projectedPosition := cameraFrom.Add(cameraDir.Mul(depth * float32(32.0)))

	// Get its z coordinate on the screen
	renderPosition := transformMatrix.Mul4x1(mgl32.Vec4{projectedPosition.X(), projectedPosition.Y(), projectedPosition.Z(), 1})
	renderPosition = renderPosition.Mul(1 / renderPosition.W())
	return renderPosition.Z()
}

// Normalize xy coordinates between -1 and 1
func convertToScreenX(x float32) float32 {
	return 2.0*(x/float32(BACKGROUND_IMAGE_WIDTH)) - 1.0
}

func convertToScreenY(y float32) float32 {
	return -1.0 * (2.0*(y/float32(BACKGROUND_IMAGE_HEIGHT)) - 1.0)
}

// Normalize uv coordinates between 0 and 1
func convertToTextureU(u float32) float32 {
	return u / float32(BACKGROUND_IMAGE_WIDTH)
}

func convertToTextureV(v float32) float32 {
	return v / float32(BACKGROUND_IMAGE_HEIGHT)
}

func GenerateCameraImageMaskEntity(
	renderDef *RenderDef,
	roomOutput *fileio.RoomImageOutput,
	cameraMasks []fileio.MaskRectangle) {
	if roomOutput.ImageMask == nil {
		renderDef.AddSceneEntity(ENTITY_CAMERA_MASK_ID, nil)
		return
	}

	cameraMaskColors := BuildCameraMaskPixels(roomOutput, cameraMasks)
	cameraMaskDepthBuffer := BuildCameraMaskDepthBuffer(roomOutput, cameraMasks, renderDef)

	cameraMaskImageEntity := NewSceneEntity()
	cameraMaskImageEntity.SetTexture(cameraMaskColors, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
	cameraMaskImageEntity.SetMesh(cameraMaskDepthBuffer)
	renderDef.AddSceneEntity(ENTITY_CAMERA_MASK_ID, cameraMaskImageEntity)
}

func BuildCameraMaskPixels(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle) []uint16 {
	// Check if background mask exists
	if roomImageOutput.ImageMask == nil {
		return []uint16{}
	}

	// Combine original background image with mask data
	backgroundImageColors := roomImageOutput.BackgroundImage.PixelData
	imageMaskColors := roomImageOutput.ImageMask.PixelData

	textureBuffer := make([]uint16, BACKGROUND_IMAGE_WIDTH*BACKGROUND_IMAGE_HEIGHT)
	for _, cameraMask := range cameraMasks {
		for offsetY := 0; offsetY < cameraMask.Height; offsetY++ {
			for offsetX := 0; offsetX < cameraMask.Width; offsetX++ {
				backgroundColor := backgroundImageColors[cameraMask.DestY+offsetY][cameraMask.DestX+offsetX]

				// Determine if pixel should be transparent
				maskColor := imageMaskColors[cameraMask.SrcY+offsetY][cameraMask.SrcX+offsetX]
				var alpha int
				if maskColor > 0 {
					alpha = 1
				} else {
					alpha = 0
				}

				newTextureColor := int(backgroundColor) | int(alpha<<15)
				textureBuffer[((cameraMask.DestY+offsetY)*BACKGROUND_IMAGE_WIDTH)+(cameraMask.DestX+offsetX)] = uint16(newTextureColor)
			}
		}
	}

	return textureBuffer
}

func BuildCameraMaskDepthBuffer(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle, renderDef *RenderDef) []float32 {
	// Check if background mask exists
	if roomImageOutput.ImageMask == nil {
		return []float32{}
	}

	maskBuffer := make([]float32, 0)
	for _, cameraMask := range cameraMasks {
		destX := float32(cameraMask.DestX)
		destY := float32(cameraMask.DestY)
		depth := normalizeMaskDepth(float32(cameraMask.Depth), renderDef)
		maskWidth := float32(cameraMask.Width)
		maskHeight := float32(cameraMask.Height)

		// Create a rectangle for each mask
		// x, y are rectangle coordinates from the background image
		// z is normalized depth
		// u, v are texture coordinates

		// corner 1
		v1 := make([]float32, 0)
		rectX := convertToScreenX(destX)
		rectY := convertToScreenY(destY)
		textureU := convertToTextureU(destX)
		textureV := convertToTextureV(destY)
		v1 = append(v1, rectX, rectY, depth, textureU, textureV)

		// corner 2
		v2 := make([]float32, 0)
		rectX = convertToScreenX(destX + maskWidth)
		rectY = convertToScreenY(destY)
		textureU = convertToTextureU(destX + maskWidth)
		textureV = convertToTextureV(destY)
		v2 = append(v2, rectX, rectY, depth, textureU, textureV)

		// corner 3
		v3 := make([]float32, 0)
		rectX = convertToScreenX(destX + maskWidth)
		rectY = convertToScreenY(destY + maskHeight)
		textureU = convertToTextureU(destX + maskWidth)
		textureV = convertToTextureV(destY + maskHeight)
		v3 = append(v3, rectX, rectY, depth, textureU, textureV)

		// corner 4
		v4 := make([]float32, 0)
		rectX = convertToScreenX(destX)
		rectY = convertToScreenY(destY + maskHeight)
		textureU = convertToTextureU(destX)
		textureV = convertToTextureV(destY + maskHeight)
		v4 = append(v4, rectX, rectY, depth, textureU, textureV)

		// Triangle 1
		maskBuffer = append(maskBuffer, v1...)
		maskBuffer = append(maskBuffer, v2...)
		maskBuffer = append(maskBuffer, v3...)

		// Triangle 2
		maskBuffer = append(maskBuffer, v1...)
		maskBuffer = append(maskBuffer, v4...)
		maskBuffer = append(maskBuffer, v3...)
	}
	return maskBuffer
}
