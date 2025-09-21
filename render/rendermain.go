package render

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)


type RenderDef struct {
	ShaderSystem     *shader.ShaderSystem // Grouped shader management
	ViewSystem       *ViewSystem          // Grouped camera/view components
	SceneSystem      *SceneSystem         // Grouped scene entities
	EnvironmentLight [3]float32
	VideoBuffer      *Entity2D

	// Screen image management for menu rendering
	ScreenImageManager *ScreenImageManager
	
	// OpenGL renderer instance
	Renderer *OpenGLRenderer
}

type DebugEntities struct {
	CameraSwitchDebugEntity *DebugEntity
	DebugEntities           []*DebugEntity
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
		VideoBuffer:        NewBackgroundImageEntity(),
		ScreenImageManager: NewScreenImageManager(),
		Renderer:           NewOpenGLRenderer(shaderSystem.GetUniformLocations()),
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

func (r *RenderDef) GetPerspectiveMatrix(fovDegrees float32) mgl32.Mat4 {
	return r.ViewSystem.GetPerspectiveMatrix(fovDegrees)
}
