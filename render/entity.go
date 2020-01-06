package render

import (
	"../fileio"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	RENDER_TYPE_ENTITY = 3
	FRAME_TIME         = 30 // time in milliseconds
	VERTEX_LEN         = 8
)

var (
	totalTime   = float64(0)
	frameIndex  = 0 // points to frame number
	frameNumber = 0 // corresponds to rotation
	curPose     = -1
)

// Offset in vertex buffer
type ComponentOffsets struct {
	StartIndex int
	EndIndex   int
}

func RenderEntity(programShader uint32, playerEntity PlayerEntity, timeElapsedSeconds float64) {
	texId := playerEntity.TextureId
	pldOutput := playerEntity.PLDOutput
	entityVertexBuffer := playerEntity.VertexBuffer

	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_ENTITY)

	updateAnimationFrame(playerEntity, timeElapsedSeconds)

	// The root of the skeleton is component 0
	transforms := make([]mgl32.Mat4, len(pldOutput.MeshData.Components))
	buildComponentTransforms(pldOutput, 0, -1, transforms)

	// Build vertex and texture data
	componentOffsets := calculateComponentOffsets(pldOutput)
	floatSize := 4

	// 3 floats for vertex, 2 floats for texture UV, 3 float for normals
	stride := int32(VERTEX_LEN * floatSize)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(entityVertexBuffer)*floatSize, gl.Ptr(entityVertexBuffer), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(3*floatSize))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(5*floatSize))
	gl.EnableVertexAttribArray(2)

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texId)

	// Draw triangles
	for i := 0; i < len(componentOffsets); i++ {
		// Set offset to translate each component relative to model origin
		boneOffsetLoc := gl.GetUniformLocation(programShader, gl.Str("boneOffset\x00"))
		boneOffset := transforms[i]
		gl.UniformMatrix4fv(boneOffsetLoc, 1, false, &boneOffset[0])

		startIndex := componentOffsets[i].StartIndex
		endIndex := componentOffsets[i].EndIndex

		// Render model component
		vertOffset := int32(startIndex / VERTEX_LEN)
		numVertices := int32((endIndex - startIndex) / VERTEX_LEN)
		gl.DrawArrays(gl.TRIANGLES, vertOffset, numVertices)
	}

	// Cleanup
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
}

func updateAnimationFrame(playerEntity PlayerEntity, timeElapsedSeconds float64) {
	pldOutput := playerEntity.PLDOutput
	poseNumber := playerEntity.AnimationPoseNumber

	if curPose != -1 {
		totalTime += timeElapsedSeconds * 1000
	} else {
		totalTime = 0
	}

	if curPose != poseNumber {
		frameIndex = 0
		if poseNumber != -1 {
			frameData := pldOutput.AnimationData.AnimationIndexFrames[poseNumber]
			frameNumber = frameData[frameIndex].FrameId
		}
		curPose = poseNumber
	}

	if totalTime >= FRAME_TIME && curPose != -1 {
		totalTime = 0
		frameIndex++
		if poseNumber != -1 {
			frameData := pldOutput.AnimationData.AnimationIndexFrames[poseNumber]
			if frameIndex >= len(frameData) {
				frameIndex = 0
			}
			frameNumber = frameData[frameIndex].FrameId
		}
	}
}

func buildComponentTransforms(pldOutput *fileio.PLDOutput, curId int, parentId int, transforms []mgl32.Mat4) {
	transformMatrix := mgl32.Ident4()
	if parentId != -1 {
		transformMatrix = transforms[parentId]
	}

	offsetFromParent := pldOutput.SkeletonData.RelativePositionData[curId]

	// Translate from parent offset
	translate := mgl32.Translate3D(float32(offsetFromParent.X), float32(offsetFromParent.Y), float32(offsetFromParent.Z))
	transformMatrix = transformMatrix.Mul4(translate)

	// Rotate if there is an animation pose
	if curPose != -1 {
		quat := mgl32.QuatIdent()
		frameRotation := pldOutput.SkeletonData.FrameData[frameNumber].RotationAngles[curId]
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.X(), mgl32.Vec3{1.0, 0.0, 0.0}))
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.Y(), mgl32.Vec3{0.0, 1.0, 0.0}))
		quat = quat.Mul(mgl32.QuatRotate(frameRotation.Z(), mgl32.Vec3{0.0, 0.0, 1.0}))
		transformMatrix = transformMatrix.Mul4(quat.Mat4())
	}

	transforms[curId] = transformMatrix

	for i := 0; i < len(pldOutput.SkeletonData.ArmatureChildren[curId]); i++ {
		newParent := curId
		newChild := int(pldOutput.SkeletonData.ArmatureChildren[curId][i])
		buildComponentTransforms(pldOutput, newChild, newParent, transforms)
	}
}

func calculateComponentOffsets(pldOutput *fileio.PLDOutput) []ComponentOffsets {
	componentOffsets := make([]ComponentOffsets, len(pldOutput.MeshData.Components))
	startIndex := 0
	endIndex := 0
	for i, entityModel := range pldOutput.MeshData.Components {
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

func BuildEntityComponentVertices(pldOutput *fileio.PLDOutput) []float32 {
	vertexBuffer := make([]float32, 0)
	textureData := pldOutput.TextureData

	for _, entityModel := range pldOutput.MeshData.Components {
		// Triangles
		for j := 0; j < len(entityModel.TriangleIndices); j++ {
			triangleIndex := entityModel.TriangleIndices[j]
			textureInfo := entityModel.TriangleTextures[j]

			vertex0 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex0])
			uv0 := buildTextureUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page, textureData)
			normal0 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal0])

			vertex1 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex1])
			uv1 := buildTextureUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page, textureData)
			normal1 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal1])

			vertex2 := buildModelVertex(entityModel.TriangleVertices[triangleIndex.IndexVertex2])
			uv2 := buildTextureUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page, textureData)
			normal2 := buildModelNormal(entityModel.TriangleNormals[triangleIndex.IndexNormal2])

			// v0, v1, v2
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex1...)
			vertexBuffer = append(vertexBuffer, uv1...)
			vertexBuffer = append(vertexBuffer, normal1...)

			vertexBuffer = append(vertexBuffer, vertex2...)
			vertexBuffer = append(vertexBuffer, uv2...)
			vertexBuffer = append(vertexBuffer, normal2...)
		}

		// Quads
		for j := 0; j < len(entityModel.QuadIndices); j++ {
			quadIndex := entityModel.QuadIndices[j]
			textureInfo := entityModel.QuadTextures[j]

			vertex0 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex0])
			uv0 := buildTextureUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page, textureData)
			normal0 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal0])

			vertex1 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex1])
			uv1 := buildTextureUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page, textureData)
			normal1 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal1])

			vertex2 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex2])
			uv2 := buildTextureUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page, textureData)
			normal2 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal2])

			vertex3 := buildModelVertex(entityModel.QuadVertices[quadIndex.IndexVertex3])
			uv3 := buildTextureUV(float32(textureInfo.U3), float32(textureInfo.V3), textureInfo.Page, textureData)
			normal3 := buildModelNormal(entityModel.QuadNormals[quadIndex.IndexNormal3])

			// v0, v1, v3
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex1...)
			vertexBuffer = append(vertexBuffer, uv1...)
			vertexBuffer = append(vertexBuffer, normal1...)

			vertexBuffer = append(vertexBuffer, vertex3...)
			vertexBuffer = append(vertexBuffer, uv3...)
			vertexBuffer = append(vertexBuffer, normal3...)

			// v0, v2, v3
			vertexBuffer = append(vertexBuffer, vertex0...)
			vertexBuffer = append(vertexBuffer, uv0...)
			vertexBuffer = append(vertexBuffer, normal0...)

			vertexBuffer = append(vertexBuffer, vertex2...)
			vertexBuffer = append(vertexBuffer, uv2...)
			vertexBuffer = append(vertexBuffer, normal2...)

			vertexBuffer = append(vertexBuffer, vertex3...)
			vertexBuffer = append(vertexBuffer, uv3...)
			vertexBuffer = append(vertexBuffer, normal3...)
		}
	}
	return vertexBuffer
}

func buildModelVertex(vertex fileio.MD1Vertex) []float32 {
	return []float32{float32(vertex.X), float32(vertex.Y), float32(vertex.Z)}
}

func buildTextureUV(u float32, v float32, texturePage uint16, textureData *fileio.TIMOutput) []float32 {
	textureOffsetUnit := float32(textureData.ImageWidth) / float32(textureData.NumPalettes)
	textureCoordOffset := textureOffsetUnit * float32(texturePage&3)

	newU := (float32(u) + textureCoordOffset) / float32(textureData.ImageWidth)
	newV := float32(v) / float32(textureData.ImageHeight)
	return []float32{newU, newV}
}

func buildModelNormal(normal fileio.MD1Vertex) []float32 {
	return []float32{float32(normal.X), float32(normal.Y), float32(normal.Z)}
}
