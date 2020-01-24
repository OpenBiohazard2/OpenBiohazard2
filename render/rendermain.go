package render

import (
	"../fileio"
	"../game"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_ITEM = 5
)

type RenderDef struct {
	ProgramShader    uint32
	Camera           *Camera
	ProjectionMatrix mgl32.Mat4
	ViewMatrix       mgl32.Mat4
	SceneEntityMap   map[string]*SceneEntity
	EnvironmentLight [3]float32
	WindowWidth      int
	WindowHeight     int
}

type DebugEntities struct {
	CameraId                int
	CameraSwitches          []fileio.RVDHeader
	CollisionEntities       []fileio.CollisionEntity
	DoorTriggers            []game.ScriptDoor
	ItemTriggers            []game.ScriptItemAotSet
	CameraSwitchTransitions map[int][]int
}

type PlayerEntity struct {
	TextureId           uint32
	VertexBuffer        []float32
	PLDOutput           *fileio.PLDOutput
	Player              *game.Player
	AnimationPoseNumber int
}

type SpriteEntity struct {
	TextureIds [][]uint32
	Sprites    []game.ScriptSprite
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
	camera := NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, cameraUp, 60.0)
	renderDef := &RenderDef{
		ProgramShader:    shader.ProgramShader,
		Camera:           camera,
		ProjectionMatrix: mgl32.Perspective(mgl32.DegToRad(60.0), float32(4/3), 16, 45000),
		ViewMatrix:       mgl32.Ident4(),
		SceneEntityMap:   make(map[string]*SceneEntity),
		WindowWidth:      windowWidth,
		WindowHeight:     windowHeight,
	}
	return renderDef
}

func (r *RenderDef) RenderFrame(playerEntity PlayerEntity,
	itemEntities []SceneMD1Entity,
	debugEntities DebugEntities,
	spriteEntity SpriteEntity,
	timeElapsedSeconds float64) {
	programShader := r.ProgramShader

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate shader
	gl.UseProgram(programShader)

	viewLoc := gl.GetUniformLocation(programShader, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(programShader, gl.Str("projection\x00"))

	r.ProjectionMatrix = r.GetPerspectiveMatrix(r.Camera.CameraFov)

	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &r.ViewMatrix[0])
	gl.UniformMatrix4fv(projectionLoc, 1, false, &r.ProjectionMatrix[0])

	r.RenderSceneEntity(r.SceneEntityMap[ENTITY_BACKGROUND_ID], RENDER_TYPE_BACKGROUND)
	r.RenderSceneEntity(r.SceneEntityMap[ENTITY_CAMERA_MASK_ID], RENDER_TYPE_CAMERA_MASK)
	for _, itemEntity := range itemEntities {
		r.RenderMD1Entity(itemEntity, RENDER_TYPE_ITEM)
	}

	envLightLoc := gl.GetUniformLocation(r.ProgramShader, gl.Str("envLight\x00"))
	gl.Uniform3fv(envLightLoc, 1, &r.EnvironmentLight[0])
	RenderAnimatedEntity(programShader, playerEntity, timeElapsedSeconds)

	RenderSprites(programShader, spriteEntity.Sprites, spriteEntity.TextureIds, timeElapsedSeconds)

	// Only render for debugging
	cameraId := debugEntities.CameraId
	cameraSwitches := debugEntities.CameraSwitches
	collisionEntities := debugEntities.CollisionEntities
	cameraSwitchTransitions := debugEntities.CameraSwitchTransitions
	RenderCameraSwitches(programShader, cameraSwitches, cameraSwitchTransitions, cameraId)
	RenderCollisionEntities(programShader, collisionEntities)
	RenderSlopedSurfaces(programShader, collisionEntities)
	RenderDoorTriggers(programShader, debugEntities.DoorTriggers)
	RenderItemTriggers(programShader, debugEntities.ItemTriggers)
}

func (r *RenderDef) AddSceneEntity(entityId string, entity *SceneEntity) {
	r.SceneEntityMap[entityId] = entity
}

func (r *RenderDef) SetEnvironmentLight(light fileio.LITCameraLight) {
	lightColor := light.AmbientColor
	// Normalize color rgb values to be between 0.0 and 1.0
	red := float32(lightColor.R) / float32(255.0)
	green := float32(lightColor.G) / float32(255.0)
	blue := float32(lightColor.B) / float32(255.0)
	r.EnvironmentLight = [3]float32{red, green, blue}
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

func (r *RenderDef) RenderSceneEntity(entity *SceneEntity, renderType int32) {
	// Skip
	if entity == nil {
		return
	}

	vertexBuffer := entity.VertexBuffer
	if len(vertexBuffer) == 0 {
		return
	}

	programShader := r.ProgramShader
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, renderType)

	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV
	stride := int32(5 * floatSize)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBuffer)*floatSize, gl.Ptr(vertexBuffer), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*floatSize))
	gl.EnableVertexAttribArray(1)

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, entity.TextureId)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/5))

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
}

func (r *RenderDef) GetPerspectiveMatrix(fovDegrees float32) mgl32.Mat4 {
	ratio := float64(r.WindowWidth) / float64(r.WindowHeight)
	return mgl32.Perspective(mgl32.DegToRad(fovDegrees), float32(ratio), 16, 45000)
}
