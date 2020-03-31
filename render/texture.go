package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

func NewTextureTIM(timOutput *fileio.TIMOutput) uint32 {
	texColors := timOutput.ConvertToRenderData()
	textureId := BuildTexture(texColors, int32(timOutput.ImageWidth), int32(timOutput.ImageHeight))
	return textureId
}

func UpdateTextureADT(texId uint32, adtOutput *fileio.ADTOutput) {
	texColors := adtOutput.ConvertToRenderData()
	imageWidth := int32(320)
	imageHeight := int32(240)
	UpdateTexture(texId, texColors, imageWidth, imageHeight)
}

func BuildTexture(imagePixels []uint16, imageWidth int32, imageHeight int32) uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)

	// Image is 16 bit in A1R5G5B5 format
	gl.TexImage2D(uint32(gl.TEXTURE_2D), 0, int32(gl.RGBA), imageWidth, imageHeight,
		0, uint32(gl.RGBA), uint32(gl.UNSIGNED_SHORT_1_5_5_5_REV), gl.Ptr(imagePixels))

	// Set texture wrapping/filtering options
	gl.TexParameteri(uint32(gl.TEXTURE_2D), gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(uint32(gl.TEXTURE_2D), gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	return texId
}

func UpdateTexture(texId uint32, imagePixels []uint16, imageWidth int32, imageHeight int32) {
	gl.BindTexture(gl.TEXTURE_2D, texId)

	// Image is 16 bit in A1R5G5B5 format
	gl.TexImage2D(uint32(gl.TEXTURE_2D), 0, int32(gl.RGBA), imageWidth, imageHeight,
		0, uint32(gl.RGBA), uint32(gl.UNSIGNED_SHORT_1_5_5_5_REV), gl.Ptr(imagePixels))

	// Set texture wrapping/filtering options
	gl.TexParameteri(uint32(gl.TEXTURE_2D), gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(uint32(gl.TEXTURE_2D), gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	return
}
