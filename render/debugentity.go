package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/geometry"
	"github.com/samuelyuan/openbiohazard2/world"
)

const (
	RENDER_TYPE_DEBUG = -1
)

type DebugEntity struct {
	Color              [4]float32
	VertexBuffer       []float32
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func RenderCameraSwitches(programShader uint32, cameraSwitchDebugEntity *DebugEntity) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

	RenderDebugEntities(programShader, []*DebugEntity{cameraSwitchDebugEntity})
}

func RenderDebugEntities(programShader uint32, debugEntities []*DebugEntity) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

	floatSize := 4

	for _, debugEntity := range debugEntities {
		entityVertexBuffer := debugEntity.VertexBuffer
		if len(entityVertexBuffer) == 0 {
			continue
		}

		// 3 floats for vertex
		stride := int32(3 * floatSize)

		vao := debugEntity.VertexArrayObject
		gl.BindVertexArray(vao)

		vbo := debugEntity.VertexBufferObject
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(entityVertexBuffer)*floatSize, gl.Ptr(entityVertexBuffer), gl.STATIC_DRAW)

		// Position attribute
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)

		diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
		gl.Uniform1i(diffuseUniform, 0)

		debugColorLoc := gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
		color := debugEntity.Color
		gl.Uniform4f(debugColorLoc, color[0], color[1], color[2], color[3])

		// Draw triangles
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(entityVertexBuffer)/3))

		// Cleanup
		gl.DisableVertexAttribArray(0)
	}
}

func BuildAllDebugEntities(gameWorld *world.GameWorld) []*DebugEntity {
	debugEntities := make([]*DebugEntity, 0)
	debugEntities = append(debugEntities, NewDoorTriggerDebugEntity(gameWorld.AotManager.Doors))
	debugEntities = append(debugEntities, NewCollisionDebugEntity(gameWorld.GameRoom.CollisionEntities))
	debugEntities = append(debugEntities, NewSlopedSurfacesDebugEntity(gameWorld.GameRoom.CollisionEntities))
	debugEntities = append(debugEntities, NewItemTriggerDebugEntity(gameWorld.AotManager.Items))
	debugEntities = append(debugEntities, NewAotTriggerDebugEntity(gameWorld.AotManager.AotTriggers))
	return debugEntities
}

func NewCollisionDebugEntity(collisionEntities []fileio.CollisionEntity) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 0:
			// Rectangle
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex4 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			rect := buildDebugRectangle(vertex1, vertex2, vertex3, vertex4)
			vertexBuffer = append(vertexBuffer, rect...)
		case 1:
			// Triangle \\|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			tri := buildDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri...)
		case 2:
			// Triangle |/
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			tri := buildDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri...)
		case 3:
			// Triangle /|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			tri := buildDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri...)
		case 6:
			// Circle
			radius := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + radius, 0, float32(entity.Z) + radius}
			circle := geometry.NewCircle(center, radius)
			vertexBuffer = append(vertexBuffer, circle.VertexBuffer...)
		case 7:
			// Ellipse, rectangle with rounded corners on the x-axis
			majorAxis := float32(entity.Width) / 2.0
			minorAxis := float32(entity.Density) / 2.0
			center := mgl32.Vec3{float32(entity.X) + majorAxis, 0, float32(entity.Z) + minorAxis}
			ellipse := geometry.NewEllipse(center, majorAxis, minorAxis, true)
			vertexBuffer = append(vertexBuffer, ellipse.VertexBuffer...)
		case 8:
			// Ellipse, rectangle with rounded corners on the z-axis
			majorAxis := float32(entity.Density) / 2.0
			minorAxis := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + minorAxis, 0, float32(entity.Z) + majorAxis}
			ellipse := geometry.NewEllipse(center, majorAxis, minorAxis, false)
			vertexBuffer = append(vertexBuffer, ellipse.VertexBuffer...)
		}
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{1.0, 0.0, 0.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewCameraSwitchDebugEntity(curCameraId int,
	cameraSwitches []fileio.RVDHeader,
	cameraSwitchTransitions map[int][]int) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, regionIndex := range cameraSwitchTransitions[curCameraId] {
		cameraSwitch := cameraSwitches[regionIndex]
		corners := [4][]float32{
			[]float32{float32(cameraSwitch.X1), float32(cameraSwitch.Z1)},
			[]float32{float32(cameraSwitch.X2), float32(cameraSwitch.Z2)},
			[]float32{float32(cameraSwitch.X3), float32(cameraSwitch.Z3)},
			[]float32{float32(cameraSwitch.X4), float32(cameraSwitch.Z4)},
		}
		rect := geometry.NewQuadFourPoints(corners)
		vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{0.0, 1.0, 0.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewDoorTriggerDebugEntity(doors []world.AotDoor) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range doors {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{0.0, 0.0, 1.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewItemTriggerDebugEntity(items []world.AotItem) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range items {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{0.0, 1.0, 1.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewAotTriggerDebugEntity(aotTriggers []world.AotObject) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range aotTriggers {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{0.0, 1.0, 1.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewSlopedSurfacesDebugEntity(collisionEntities []fileio.CollisionEntity) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 11:
			// Ramp
			rect := geometry.NewSlopedRectangle(entity)
			vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
		case 12:
			// Stairs
			rect := geometry.NewSlopedRectangle(entity)
			vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
		}
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              [4]float32{1.0, 0.0, 1.0, 0.3},
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func buildDebugRectangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) []float32 {
	quad := geometry.NewQuad([4]mgl32.Vec3{corner1, corner2, corner3, corner4})
	return quad.VertexBuffer
}

func buildDebugTriangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3) []float32 {
	triBuffer := make([]float32, 0)
	vertex1 := []float32{corner1.X(), corner1.Y(), corner1.Z()}
	vertex2 := []float32{corner2.X(), corner2.Y(), corner2.Z()}
	vertex3 := []float32{corner3.X(), corner3.Y(), corner3.Z()}

	triBuffer = append(triBuffer, vertex1...)
	triBuffer = append(triBuffer, vertex2...)
	triBuffer = append(triBuffer, vertex3...)
	return triBuffer
}
