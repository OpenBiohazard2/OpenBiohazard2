package render

import (
	"image"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
)

const (
	CAMERA_MASK_WIDTH  = 256
	CAMERA_MASK_HEIGHT = 256
)

func (cameraMaskImageEntity *Entity2D) UpdateCameraImageMaskEntity(
	viewSystem *ViewSystem,
	roomOutput *fileio.RoomImageOutput,
	cameraMasks []fileio.MaskRectangle) {
	if roomOutput.ImageMask == nil {
		// Clear previous mask
		cameraMaskColors := make([]uint16, geometry.BACKGROUND_IMAGE_WIDTH*geometry.BACKGROUND_IMAGE_HEIGHT)
		cameraMaskImageEntity.SetTexture(cameraMaskColors, geometry.BACKGROUND_IMAGE_WIDTH, geometry.BACKGROUND_IMAGE_HEIGHT)
		return
	}

	cameraMaskImage := BuildCameraMask(roomOutput, cameraMasks)
	cameraMaskDepthBuffer := BuildCameraMaskDepthBuffer(roomOutput, cameraMasks, viewSystem)

	cameraMaskImageEntity.SetTexture(cameraMaskImage.GetPixelsForRendering(), geometry.BACKGROUND_IMAGE_WIDTH, geometry.BACKGROUND_IMAGE_HEIGHT)
	cameraMaskImageEntity.SetMesh(cameraMaskDepthBuffer)
}

func BuildCameraMask(roomImageOutput *fileio.RoomImageOutput, cameraMasks []fileio.MaskRectangle) *resource.Image16Bit {
	// Combine original background image with mask data
	backgroundImage := resource.ConvertPixelsToImage16Bit(roomImageOutput.BackgroundImage.PixelData)
	imageMask := resource.ConvertPixelsToImage16Bit(roomImageOutput.ImageMask.PixelData)

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
		// z is normalized depth
		depth := viewSystem.Camera.NormalizeMaskDepth(float32(cameraMask.Depth), viewSystem.ProjectionMatrix, viewSystem.ViewMatrix)

		// Create a rectangle for each mask using geometry package
		bufferRect := geometry.NewCameraMaskQuad(cameraMask, depth)
		maskBuffer = append(maskBuffer, bufferRect.VertexBuffer...)
	}
	return maskBuffer
}
