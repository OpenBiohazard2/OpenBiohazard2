package render

import (
	"../fileio"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	RENDER_GAME_STATE_MAIN_MENU = 2
	ENTITY_MAIN_MENU_ID         = "MAIN_MENU_IMAGE"
)

func (renderDef *RenderDef) GenerateMainMenuImageEntity(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput) {
	newImageColors := NewSurface2D()
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildMainMenuText(menuBackgroundTextOutput, newImageColors, 0)

	imageEntity := NewSceneEntity()
	imageEntity.SetTexture(newImageColors, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT)
	imageEntity.SetMesh(buildSurface2DVertexBuffer())
	renderDef.AddSceneEntity(ENTITY_MAIN_MENU_ID, imageEntity)
}

func buildMainMenuBackground(backgroundImageOutput *fileio.ADTOutput, newImageColors []uint16) {
	copyPixelsTransparent(backgroundImageOutput.PixelData, 0, 0, 320, 240, newImageColors, 0, 0)
}

func buildMainMenuText(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, mainMenuOption int) {
	buildTitleText(menuBackgroundTextOutput, newImageColors)
	buildMainMenuOptions(menuBackgroundTextOutput, newImageColors, mainMenuOption)
}

func buildTitleText(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16) {
	copyPixelsTransparent(menuBackgroundTextOutput[1].PixelData, 0, 0, 128, 81, newImageColors, 18, 30)
	copyPixelsTransparent(menuBackgroundTextOutput[1].PixelData, 0, 81, 128, 47, newImageColors, 146, 31)
	copyPixelsTransparent(menuBackgroundTextOutput[2].PixelData, 0, 0, 128, 34, newImageColors, 146, 78)
	copyPixelsTransparent(menuBackgroundTextOutput[2].PixelData, 0, 34, 46, 82, newImageColors, 274, 31)
}

func buildMainMenuOptions(menuBackgroundTextOutput []*fileio.TIMOutput, newImageColors []uint16, mainMenuOption int) {
	selectedOption := 1.0
	otherOption := 0.2

	optionsBrightness := [3]float64{otherOption, otherOption, otherOption}
	optionsBrightness[mainMenuOption] = selectedOption

	// Load Game
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 70, 29, 106, 13, newImageColors, 114, 134, optionsBrightness[0])

	// New Game
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 54, 17, 147, 12, newImageColors, 95, 154, optionsBrightness[1])

	// Option
	copyPixelsBrightness(menuBackgroundTextOutput[0].PixelData, 88, 43, 74, 14, newImageColors, 130, 174, optionsBrightness[2])
}

func (renderDef *RenderDef) UpdateMainMenu(
	menuBackgroundImageOutput *fileio.ADTOutput,
	menuBackgroundTextOutput []*fileio.TIMOutput,
	mainMenuOption int) {
	newImageColors := NewSurface2D()
	buildMainMenuBackground(menuBackgroundImageOutput, newImageColors)
	buildMainMenuText(menuBackgroundTextOutput, newImageColors, mainMenuOption)
	renderDef.SceneEntityMap[ENTITY_MAIN_MENU_ID].SetTexture(newImageColors, IMAGE_SURFACE_WIDTH, IMAGE_SURFACE_HEIGHT)
}
func (renderDef *RenderDef) RenderMainMenu() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	programShader := renderDef.ProgramShader

	// Activate shader
	gl.UseProgram(programShader)

	renderGameStateUniform := gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	gl.Uniform1i(renderGameStateUniform, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)

	renderDef.RenderSceneEntity(renderDef.SceneEntityMap[ENTITY_MAIN_MENU_ID], RENDER_TYPE_BACKGROUND)
}
