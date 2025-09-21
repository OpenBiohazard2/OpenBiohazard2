package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// RenderRoom contains all the data needed to render a room
type RenderRoom struct {
	CameraMaskData  [][]fileio.MaskRectangle
	LightData       []fileio.LITCameraLight
	ItemTextureData []*fileio.TIMOutput
	ItemModelData   []*fileio.MD1Output
	SpriteData      []fileio.SpriteData
}

// NewRenderRoom creates a new RenderRoom from RDT output data
func NewRenderRoom(rdtOutput *fileio.RDTOutput) RenderRoom {
	return RenderRoom{
		CameraMaskData:  rdtOutput.RIDOutput.CameraMasks,
		LightData:       rdtOutput.LightData.Lights,
		ItemTextureData: rdtOutput.ItemTextureData,
		ItemModelData:   rdtOutput.ItemModelData,
		SpriteData:      rdtOutput.SpriteOutput.SpriteData,
	}
}

// BuildEnvironmentLight converts a LIT camera light to normalized RGB values
func BuildEnvironmentLight(light fileio.LITCameraLight) [3]float32 {
	lightColor := light.AmbientColor
	// Normalize color rgb values to be between 0.0 and 1.0
	red := float32(lightColor.R) / 255.0
	green := float32(lightColor.G) / 255.0
	blue := float32(lightColor.B) / 255.0
	return [3]float32{red, green, blue}
}
