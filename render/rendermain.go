package render

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
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
	ProgramShader    uint32
	Camera           *Camera
	ProjectionMatrix mgl32.Mat4
	ViewMatrix       mgl32.Mat4
	EnvironmentLight [3]float32
	WindowWidth      int
	WindowHeight     int
	VideoBuffer      *Surface2D

	SpriteGroupEntity     *SpriteGroupEntity
	BackgroundImageEntity *SceneEntity
	CameraMaskEntity      *SceneEntity
	ItemGroupEntity       *ItemGroupEntity

	// Cached uniform locations for performance
	UniformLocations UniformLocations
}

type UniformLocations struct {
	// Main rendering uniforms
	GameState  int32
	View       int32
	Projection int32
	EnvLight   int32

	// Entity rendering uniforms
	RenderType int32
	Model      int32
	Diffuse    int32

	// Debug rendering uniforms
	DebugColor int32
	BoneOffset int32
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

	shader := NewShader("render/openbiohazard2.vert", "render/openbiohazard2.frag")

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	cameraUp := mgl32.Vec3{0, -1, 0}
	camera := NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, cameraUp, DEFAULT_FOV_DEGREES)
	renderDef := &RenderDef{
		ProgramShader:    shader.ProgramShader,
		Camera:           camera,
		ProjectionMatrix: mgl32.Perspective(mgl32.DegToRad(DEFAULT_FOV_DEGREES), float32(ASPECT_RATIO), NEAR_PLANE, FAR_PLANE),
		ViewMatrix:       mgl32.Ident4(),
		WindowWidth:      windowWidth,
		WindowHeight:     windowHeight,
		VideoBuffer:      NewSurface2D(),

		BackgroundImageEntity: NewBackgroundImageEntity(),
		CameraMaskEntity:      NewSceneEntity(),
		ItemGroupEntity:       NewItemGroupEntity(),
	}

	// Cache all uniform locations for performance
	renderDef.cacheUniformLocations()

	return renderDef
}

// cacheUniformLocations caches all uniform locations to avoid expensive gl.GetUniformLocation calls every frame
func (r *RenderDef) cacheUniformLocations() {
	programShader := r.ProgramShader

	// Main rendering uniforms
	r.UniformLocations.GameState = gl.GetUniformLocation(programShader, gl.Str("gameState\x00"))
	r.UniformLocations.View = gl.GetUniformLocation(programShader, gl.Str("view\x00"))
	r.UniformLocations.Projection = gl.GetUniformLocation(programShader, gl.Str("projection\x00"))
	r.UniformLocations.EnvLight = gl.GetUniformLocation(programShader, gl.Str("envLight\x00"))

	// Entity rendering uniforms
	r.UniformLocations.RenderType = gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	r.UniformLocations.Model = gl.GetUniformLocation(programShader, gl.Str("model\x00"))
	r.UniformLocations.Diffuse = gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))

	// Debug rendering uniforms
	r.UniformLocations.DebugColor = gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
	r.UniformLocations.BoneOffset = gl.GetUniformLocation(programShader, gl.Str("boneOffset\x00"))
}

func (r *RenderDef) RenderFrame(playerEntity PlayerEntity,
	debugEntities DebugEntities,
	timeElapsedSeconds float64) {

	programShader := r.ProgramShader

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate shader
	gl.UseProgram(programShader)

	// Use cached uniform locations for better performance
	gl.Uniform1i(r.UniformLocations.GameState, RENDER_GAME_STATE_MAIN)

	r.ProjectionMatrix = r.GetPerspectiveMatrix(r.Camera.CameraFov)

	// Pass the matrices to the shader using cached locations
	gl.UniformMatrix4fv(r.UniformLocations.View, 1, false, &r.ViewMatrix[0])
	gl.UniformMatrix4fv(r.UniformLocations.Projection, 1, false, &r.ProjectionMatrix[0])

	r.RenderBackground()
	for _, itemEntity := range r.ItemGroupEntity.ModelObjectData {
		r.RenderStaticEntity(*itemEntity, RENDER_TYPE_ITEM)
	}

	// Use cached uniform location for environment light
	gl.Uniform3fv(r.UniformLocations.EnvLight, 1, &r.EnvironmentLight[0])
	RenderAnimatedEntity(r, playerEntity, timeElapsedSeconds)

	// RenderSprites(r, r.SpriteGroupEntity, timeElapsedSeconds)

	// Only render for debugging
	RenderCameraSwitches(r, debugEntities.CameraSwitchDebugEntity)
	RenderDebugEntities(r, debugEntities.DebugEntities)
}

func (r *RenderDef) RenderBackground() {
	r.RenderSceneEntity(r.BackgroundImageEntity, RENDER_GAME_STATE_BACKGROUND_SOLID)
	r.RenderSceneEntity(r.CameraMaskEntity, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)
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
	ratio := float64(r.WindowWidth) / float64(r.WindowHeight)
	return mgl32.Perspective(mgl32.DegToRad(fovDegrees), float32(ratio), NEAR_PLANE, FAR_PLANE)
}
