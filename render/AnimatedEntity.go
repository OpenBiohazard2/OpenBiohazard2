package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_ENTITY = 3
	VERTEX_LEN         = 8
)

type PlayerEntity struct {
	TextureId           uint32
	VertexBuffer        []float32
	PLDOutput           *fileio.PLDOutput
	Player              *game.Player
	AnimationPoseNumber int
	VertexArrayObject   uint32
	VertexBufferObject  uint32

	// Pre-allocated arrays to avoid allocations every frame
	Transforms       []mgl32.Mat4
	ComponentOffsets []ComponentOffsets
	LastPoseNumber   int  // Track when pose changes
	BufferUploaded   bool // Track if buffer has been uploaded to GPU

	Animation *Animation
}

// Offset in vertex buffer
type ComponentOffsets struct {
	StartIndex int
	EndIndex   int
}

func NewPlayerEntity(pldOutput *fileio.PLDOutput) *PlayerEntity {
	// Generate buffers
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	textureId := NewTextureTIM(pldOutput.TextureData)
	vertexBuffer := geometry.NewMD1Geometry(pldOutput.MeshData, pldOutput.TextureData)

	// Pre-allocate arrays based on mesh data
	transforms := make([]mgl32.Mat4, len(pldOutput.MeshData.Components))
	componentOffsets := calculateComponentOffsets(pldOutput.MeshData)

	return &PlayerEntity{
		TextureId:           textureId,
		VertexBuffer:        vertexBuffer,
		PLDOutput:           pldOutput,
		Player:              nil,
		AnimationPoseNumber: -1,
		VertexArrayObject:   vao,
		VertexBufferObject:  vbo,
		Transforms:          transforms,
		ComponentOffsets:    componentOffsets,
		LastPoseNumber:      -1,
		BufferUploaded:      false,
		Animation:           NewAnimation(),
	}
}

func (playerEntity *PlayerEntity) UpdatePlayerEntity(player *game.Player, animationPoseNumber int) {
	playerEntity.Player = player
	playerEntity.AnimationPoseNumber = animationPoseNumber
}

func RenderAnimatedEntity(r *RenderDef, playerEntity PlayerEntity, timeElapsedSeconds float64) {
	// Early return if no player
	if playerEntity.Player == nil {
		return
	}

	// Update animation and transforms
	playerEntity.updateAnimation(timeElapsedSeconds)
	playerEntity.updateTransforms()

	// Set up rendering state
	playerEntity.setupRendering(r)

	// Render all components
	playerEntity.renderComponents(r)

	// Clean up
	playerEntity.cleanup()
}

// updateAnimation handles animation frame updates
func (pe *PlayerEntity) updateAnimation(timeElapsedSeconds float64) {
	pe.Animation.UpdateAnimationFrame(pe.AnimationPoseNumber, pe.PLDOutput.AnimationData, timeElapsedSeconds)
}

// updateTransforms recalculates bone transforms when needed
func (pe *PlayerEntity) updateTransforms() {
	needsUpdate := pe.LastPoseNumber != pe.AnimationPoseNumber || !pe.BufferUploaded
	if needsUpdate {
		buildComponentTransforms(pe.PLDOutput.SkeletonData, 0, -1, pe.Transforms, pe.Animation)
		pe.LastPoseNumber = pe.AnimationPoseNumber
	}
}

// setupRendering configures OpenGL state for rendering
func (pe *PlayerEntity) setupRendering(r *RenderDef) {
	// Set uniforms using ShaderSystem methods
	r.ShaderSystem.SetRenderType(RENDER_TYPE_ENTITY)
	modelMatrix := pe.Player.GetModelMatrix()
	r.ShaderSystem.SetModelMatrix(modelMatrix)
	r.ShaderSystem.SetDiffuse(0)

	// Bind buffers
	gl.BindVertexArray(pe.VertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, pe.VertexBufferObject)

	// Upload buffer data if needed
	if !pe.BufferUploaded {
		const floatSize = FLOAT_SIZE_BYTES
		gl.BufferData(gl.ARRAY_BUFFER, len(pe.VertexBuffer)*floatSize, gl.Ptr(pe.VertexBuffer), gl.STATIC_DRAW)
		pe.BufferUploaded = true
	}

	// Set up vertex attributes
	pe.setupVertexAttributes()

	// Bind texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, pe.TextureId)
}

// setupVertexAttributes configures vertex attribute pointers
func (pe *PlayerEntity) setupVertexAttributes() {
	const (
		floatSize = FLOAT_SIZE_BYTES
		stride    = int32(VERTEX_LEN * floatSize)
	)

	// Position (3 floats)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture UV (2 floats)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*floatSize))
	gl.EnableVertexAttribArray(1)

	// Normal (3 floats)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(5*floatSize))
	gl.EnableVertexAttribArray(2)
}

// renderComponents draws all mesh components
func (pe *PlayerEntity) renderComponents(r *RenderDef) {
	for i, offset := range pe.ComponentOffsets {
		// Set bone transform using ShaderSystem method
		r.ShaderSystem.SetBoneOffset(pe.Transforms[i])

		// Calculate vertex range
		vertOffset := int32(offset.StartIndex / VERTEX_LEN)
		numVertices := int32((offset.EndIndex - offset.StartIndex) / VERTEX_LEN)

		// Draw component
		gl.DrawArrays(gl.TRIANGLES, vertOffset, numVertices)
	}
}

// cleanup disables vertex attributes
func (pe *PlayerEntity) cleanup() {
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
}

func buildComponentTransforms(skeletonData *fileio.EMROutput, curId int, parentId int, transforms []mgl32.Mat4, animation *Animation) {
	transformMatrix := mgl32.Ident4()
	if parentId != -1 {
		transformMatrix = transforms[parentId]
	}

	offsetFromParent := skeletonData.RelativePositionData[curId]

	// Translate from parent offset
	translate := mgl32.Translate3D(float32(offsetFromParent.X), float32(offsetFromParent.Y), float32(offsetFromParent.Z))
	transformMatrix = transformMatrix.Mul4(translate)

	// Rotate if there is an animation pose
	if animation.CurPose != -1 {
		quat := mgl32.QuatIdent()
		frameRotation := skeletonData.FrameData[animation.FrameNumber].RotationAngles[curId]
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.X(), mgl32.Vec3{1.0, 0.0, 0.0}))
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.Y(), mgl32.Vec3{0.0, 1.0, 0.0}))
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.Z(), mgl32.Vec3{0.0, 0.0, 1.0}))
		transformMatrix = transformMatrix.Mul4(quat.Mat4())
	}

	transforms[curId] = transformMatrix

	for i := 0; i < len(skeletonData.ArmatureChildren[curId]); i++ {
		newParent := curId
		newChild := int(skeletonData.ArmatureChildren[curId][i])
		buildComponentTransforms(skeletonData, newChild, newParent, transforms, animation)
	}
}

func calculateComponentOffsets(meshData *fileio.MD1Output) []ComponentOffsets {
	componentOffsets := make([]ComponentOffsets, len(meshData.Components))
	startIndex := 0
	endIndex := 0
	for i, entityModel := range meshData.Components {
		startIndex = endIndex
		triangleBufferCount := len(entityModel.TriangleIndices) * 3 * VERTEX_LEN
		quadBufferCount := len(entityModel.QuadIndices) * 3 * 2 * VERTEX_LEN
		endIndex = startIndex + (triangleBufferCount + quadBufferCount)

		componentOffsets[i] = ComponentOffsets{
			StartIndex: startIndex,
			EndIndex:   endIndex,
		}
	}
	return componentOffsets
}
