package render

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_GAME_STATE_MAIN                   = 0
	RENDER_GAME_STATE_BACKGROUND_SOLID       = 1
	RENDER_GAME_STATE_BACKGROUND_TRANSPARENT = 2
	RENDER_TYPE_ITEM                         = 5

	// Camera Constants
	DEFAULT_FOV_DEGREES = 60.0
	NEAR_PLANE          = 16.0
	FAR_PLANE           = 45000.0
	ASPECT_RATIO        = 4.0 / 3.0
)

type RenderDef struct {
	ShaderSystem     *shader.ShaderSystem // Grouped shader management
	ViewSystem       *ViewSystem          // Grouped camera/view components
	SceneSystem      *SceneSystem         // Grouped scene entities
	EnvironmentLight [3]float32
	VideoBuffer      *Surface2D

	// Screen image management for menu rendering
	ScreenImageManager *ScreenImageManager
}

type DebugEntities struct {
	CameraSwitchDebugEntity *DebugEntity
	DebugEntities           []*DebugEntity
}

type RenderRoom struct {
	CameraMaskData  [][]fileio.MaskRectangle
	LightData       []fileio.LITCameraLight
	ItemTextureData []*fileio.TIMOutput
	ItemModelData   []*fileio.MD1Output
	SpriteData      []fileio.SpriteData
}

func InitRenderer(windowWidth int, windowHeight int) *RenderDef {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Initialize shader system
	shaderSystem := shader.NewShaderSystem()
	err := shaderSystem.Initialize("shader/openbiohazard2.vert", "shader/openbiohazard2.frag")
	if err != nil {
		panic(err)
	}

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	renderDef := &RenderDef{
		ShaderSystem:       shaderSystem,
		ViewSystem:         NewViewSystem(windowWidth, windowHeight),
		SceneSystem:        NewSceneSystem(),
		VideoBuffer:        NewSurface2D(),
		ScreenImageManager: NewScreenImageManager(),
	}

	return renderDef
}

func (r *RenderDef) RenderFrame(playerEntity PlayerEntity,
	debugEntities DebugEntities,
	timeElapsedSeconds float64) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate shader
	r.ShaderSystem.Use()

	// Use cached uniform locations for better performance
	r.ShaderSystem.SetGameState(RENDER_GAME_STATE_MAIN)

	r.ViewSystem.UpdateMatrices()

	// Pass the matrices to the shader using cached locations
	viewMatrix := r.ViewSystem.GetViewMatrix()
	projectionMatrix := r.ViewSystem.GetProjectionMatrix()
	r.ShaderSystem.SetViewMatrix(viewMatrix)
	r.ShaderSystem.SetProjectionMatrix(projectionMatrix)

	r.SceneSystem.RenderBackground(r)
	r.SceneSystem.RenderItems(r)

	// Use cached uniform location for environment light
	r.ShaderSystem.SetEnvironmentLight(r.EnvironmentLight)
	RenderAnimatedEntity(r, playerEntity, timeElapsedSeconds)

	// RenderSprites(r, r.SpriteGroupEntity, timeElapsedSeconds)

	// Only render for debugging
	RenderCameraSwitches(r, debugEntities.CameraSwitchDebugEntity)
	RenderDebugEntities(r, debugEntities.DebugEntities)
}

// UpdateCameraMask updates the camera mask entity
func (r *RenderDef) UpdateCameraMask(roomOutput *fileio.RoomImageOutput, masks []fileio.MaskRectangle) {
	r.SceneSystem.UpdateCameraMask(r, roomOutput, masks)
}

func NewRenderRoom(rdtOutput *fileio.RDTOutput) RenderRoom {
	return RenderRoom{
		CameraMaskData:  rdtOutput.RIDOutput.CameraMasks,
		LightData:       rdtOutput.LightData.Lights,
		ItemTextureData: rdtOutput.ItemTextureData,
		ItemModelData:   rdtOutput.ItemModelData,
		SpriteData:      rdtOutput.SpriteOutput.SpriteData,
	}
}

func BuildEnvironmentLight(light fileio.LITCameraLight) [3]float32 {
	lightColor := light.AmbientColor
	// Normalize color rgb values to be between 0.0 and 1.0
	red := float32(lightColor.R) / 255.0
	green := float32(lightColor.G) / 255.0
	blue := float32(lightColor.B) / 255.0
	return [3]float32{red, green, blue}
}

func (r *RenderDef) GetPerspectiveMatrix(fovDegrees float32) mgl32.Mat4 {
	return r.ViewSystem.GetPerspectiveMatrix(fovDegrees)
}
