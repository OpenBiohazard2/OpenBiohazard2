package render

import (
	"image"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/geometry"
)

const (
	CAMERA_MASK_WIDTH  = 256
	CAMERA_MASK_HEIGHT = 256

	BACKGROUND_IMAGE_WIDTH  = 320
	BACKGROUND_IMAGE_HEIGHT = 240
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

func (cameraMaskImageEntity *SceneEntity) UpdateCameraImageMaskEntity(
	renderDef *RenderDef,
	roomOutput *fileio.RoomImageOutput,
	cameraMasks []fileio.MaskRectangle) {
	if roomOutput.ImageMask == nil {
		// Clear previous mask
		cameraMaskColors := make([]uint16, BACKGROUND_IMAGE_WIDTH*BACKGROUND_IMAGE_HEIGHT)
		cameraMaskImageEntity.SetTexture(cameraMaskColors, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
		return
	}

	cameraMaskImage := BuildCameraMask(roomOutput, cameraMasks)
	cameraMaskDepthBuffer := BuildCameraMaskDepthBuffer(roomOutput, cameraMasks, renderDef)

	cameraMaskImageEntity.SetTexture(cameraMaskImage.GetPixelsForRendering(), BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
	cameraMaskImageEntity.SetMesh(cameraMaskDepthBuffer)
	return
}

func BuildCameraMask(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle) *Image16Bit {
	// Combine original background image with mask data
	backgroundImage := ConvertPixelsToImage16Bit(roomImageOutput.BackgroundImage.PixelData)
	imageMask := ConvertPixelsToImage16Bit(roomImageOutput.ImageMask.PixelData)

	for _, cameraMask := range cameraMasks {
		sourceRect := image.Rect(cameraMask.SrcX, cameraMask.SrcY, 
			cameraMask.SrcX+cameraMask.Width, cameraMask.SrcY+cameraMask.Height)
		backgroundImage.ApplyMask(image.Point{cameraMask.DestX, cameraMask.DestY}, imageMask, sourceRect) 
	}

	return backgroundImage
}

func BuildCameraMaskDepthBuffer(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle, renderDef *RenderDef) []float32 {
	maskBuffer := make([]float32, 0)
	for _, cameraMask := range cameraMasks {
		// x, y are rectangle coordinates from the background image
		destX := float32(cameraMask.DestX)
		destY := float32(cameraMask.DestY)
		// z is normalized depth
		depth := normalizeMaskDepth(float32(cameraMask.Depth), renderDef)
		maskWidth := float32(cameraMask.Width)
		maskHeight := float32(cameraMask.Height)

		// Create a rectangle for each mask
		corners := [4][]float32{
			{destX, destY},
			{destX + maskWidth, destY},
			{destX + maskWidth, destY + maskHeight},
			{destX, destY + maskHeight},
		}

		vertices := [4][]float32{}
		uvs := [4][]float32{}
		for i, corner := range corners {
			x := corner[0]
			y := corner[1]
			vertices[i] = []float32{convertToScreenX(x), convertToScreenY(y), depth}
			uvs[i] = []float32{convertToTextureU(x), convertToTextureV(y)}
		}
		bufferRect := geometry.NewTexturedRectangle(vertices, uvs)
		maskBuffer = append(maskBuffer, bufferRect.VertexBuffer...)
	}
	return maskBuffer
}
