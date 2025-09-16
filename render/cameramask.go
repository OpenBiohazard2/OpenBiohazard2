package render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
)

const (
	CAMERA_MASK_WIDTH  = 256
	CAMERA_MASK_HEIGHT = 256
)

func (cameraMaskImageEntity *SceneEntity) UpdateCameraImageMaskEntity(
	viewSystem *ViewSystem,
	roomOutput *fileio.RoomImageOutput,
	cameraMasks []fileio.MaskRectangle) {
	if roomOutput.ImageMask == nil {
		// Clear previous mask
		cameraMaskColors := make([]uint16, BACKGROUND_IMAGE_WIDTH*BACKGROUND_IMAGE_HEIGHT)
		cameraMaskImageEntity.SetTexture(cameraMaskColors, BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
		return
	}

	cameraMaskImage := BuildCameraMask(roomOutput, cameraMasks)
	cameraMaskDepthBuffer := BuildCameraMaskDepthBuffer(roomOutput, cameraMasks, viewSystem)

	cameraMaskImageEntity.SetTexture(cameraMaskImage.GetPixelsForRendering(), BACKGROUND_IMAGE_WIDTH, BACKGROUND_IMAGE_HEIGHT)
	cameraMaskImageEntity.SetMesh(cameraMaskDepthBuffer)
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

func BuildCameraMaskDepthBuffer(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle, viewSystem *ViewSystem) []float32 {
	maskBuffer := make([]float32, 0)
	for _, cameraMask := range cameraMasks {
		// x, y are rectangle coordinates from the background image
		destX := float32(cameraMask.DestX)
		destY := float32(cameraMask.DestY)
		// z is normalized depth
		depth := viewSystem.Camera.NormalizeMaskDepth(float32(cameraMask.Depth), viewSystem.ProjectionMatrix, viewSystem.ViewMatrix)
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
			vertices[i] = []float32{ConvertToScreenX(x), ConvertToScreenY(y), depth}
			uvs[i] = []float32{ConvertToTextureU(x), ConvertToTextureV(y)}
		}
		bufferRect := geometry.NewTexturedRectangle(vertices, uvs)
		maskBuffer = append(maskBuffer, bufferRect.VertexBuffer...)
	}
	return maskBuffer
}
