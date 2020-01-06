package render

import (
	"../fileio"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	RENDER_TYPE_DEBUG = 4
)

func RenderCameraSwitches(programShader uint32, cameraSwitches []fileio.RVDHeader, curCameraId int) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

	// Build vertex data
	entityVertexBuffer := buildCameraSwitchVertices(cameraSwitches, curCameraId)

	floatSize := 4

	// 3 floats for vertex
	stride := int32(3 * floatSize)

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

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	// Green rectangle
	debugColorLoc := gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
	gl.Uniform4f(debugColorLoc, 0.0, 1.0, 0.0, 0.3)

	// Draw triangles
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(entityVertexBuffer)/3))

	// Cleanup
	gl.DisableVertexAttribArray(0)
}

func buildCameraSwitchVertices(cameraSwitches []fileio.RVDHeader, curCameraId int) []float32 {
	vertexBuffer := make([]float32, 0)

	for _, cameraSwitch := range cameraSwitches {
		// Don't display regions for other cameras
		if int(cameraSwitch.Cam0) != curCameraId {
			continue
		}

		// Don't display current region
		if int(cameraSwitch.Cam1) == 0 {
			continue
		}

		vertex1 := mgl32.Vec3{float32(cameraSwitch.X1), 0, float32(cameraSwitch.Z1)}
		vertex2 := mgl32.Vec3{float32(cameraSwitch.X2), 0, float32(cameraSwitch.Z2)}
		vertex3 := mgl32.Vec3{float32(cameraSwitch.X3), 0, float32(cameraSwitch.Z3)}
		vertex4 := mgl32.Vec3{float32(cameraSwitch.X4), 0, float32(cameraSwitch.Z4)}
		rect := buildDebugRectangle(vertex1, vertex2, vertex3, vertex4)
		vertexBuffer = append(vertexBuffer, rect...)
	}
	return vertexBuffer
}

func RenderCollisionEntities(programShader uint32, collisionEntities []fileio.CollisionEntity) {
	renderTypeUniform := gl.GetUniformLocation(programShader, gl.Str("renderType\x00"))
	gl.Uniform1i(renderTypeUniform, RENDER_TYPE_DEBUG)

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
			circle := buildDebugCircle(center, radius)
			vertexBuffer = append(vertexBuffer, circle...)
		case 7:
			// Ellipse, rectangle with rounded corners on the x-axis
			majorAxis := float32(entity.Width) / 2.0
			minorAxis := float32(entity.Density) / 2.0
			center := mgl32.Vec3{float32(entity.X) + majorAxis, 0, float32(entity.Z) + minorAxis}
			ellipse := buildDebugEllipse(center, majorAxis, minorAxis, true)
			vertexBuffer = append(vertexBuffer, ellipse...)
		case 8:
			// Ellipse, rectangle with rounded corners on the z-axis
			majorAxis := float32(entity.Density) / 2.0
			minorAxis := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + minorAxis, 0, float32(entity.Z) + majorAxis}
			ellipse := buildDebugEllipse(center, majorAxis, minorAxis, false)
			vertexBuffer = append(vertexBuffer, ellipse...)
		}
	}

	floatSize := 4

	// 3 floats for vertex
	stride := int32(3 * floatSize)

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

	diffuseUniform := gl.GetUniformLocation(programShader, gl.Str("diffuse\x00"))
	gl.Uniform1i(diffuseUniform, 0)

	// Collision entities should be drawn in red
	debugColorLoc := gl.GetUniformLocation(programShader, gl.Str("debugColor\x00"))
	gl.Uniform4f(debugColorLoc, 1.0, 0.0, 0.0, 0.3)

	// Draw triangles
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexBuffer)/3))

	// Cleanup
	gl.DisableVertexAttribArray(0)
}

func buildDebugRectangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) []float32 {
	rectBuffer := make([]float32, 0)
	vertex1 := []float32{corner1.X(), corner1.Y(), corner1.Z()}
	vertex2 := []float32{corner2.X(), corner2.Y(), corner2.Z()}
	vertex3 := []float32{corner3.X(), corner3.Y(), corner3.Z()}
	vertex4 := []float32{corner4.X(), corner4.Y(), corner4.Z()}

	rectBuffer = append(rectBuffer, vertex1...)
	rectBuffer = append(rectBuffer, vertex2...)
	rectBuffer = append(rectBuffer, vertex3...)

	rectBuffer = append(rectBuffer, vertex1...)
	rectBuffer = append(rectBuffer, vertex4...)
	rectBuffer = append(rectBuffer, vertex3...)
	return rectBuffer
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

func buildDebugCircle(centerVertex mgl32.Vec3, radius float32) []float32 {
	circleBuffer := make([]float32, 0)
	center := []float32{centerVertex.X(), centerVertex.Y(), centerVertex.Z()}

	// Approximate circle using a polygon
	numVertices := 8
	for i := 0; i < numVertices; i++ {
		angle1 := float64(mgl32.DegToRad(float32(i) * 360.0 / float32(numVertices)))
		angle2 := float64(mgl32.DegToRad(float32(i+1) * 360.0 / float32(numVertices)))
		deltaX1 := radius * float32(math.Cos(angle1))
		deltaZ1 := radius * float32(math.Sin(angle1))
		deltaX2 := radius * float32(math.Cos(angle2))
		deltaZ2 := radius * float32(math.Sin(angle2))

		vertex1 := []float32{centerVertex.X() + deltaX1, centerVertex.Y(), centerVertex.Z() + deltaZ1}
		vertex2 := []float32{centerVertex.X() + deltaX2, centerVertex.Y(), centerVertex.Z() + deltaZ2}
		circleBuffer = append(circleBuffer, center...)
		circleBuffer = append(circleBuffer, vertex1...)
		circleBuffer = append(circleBuffer, vertex2...)
	}

	return circleBuffer
}

func buildDebugEllipse(centerVertex mgl32.Vec3, majorAxis float32, minorAxis float32, xAxisMajor bool) []float32 {
	ellipseBuffer := make([]float32, 0)
	center := []float32{centerVertex.X(), centerVertex.Y(), centerVertex.Z()}

	// Approximate ellipse using a polygon
	numVertices := 8
	for i := 0; i < numVertices; i++ {
		angle1 := float64(mgl32.DegToRad(float32(i) * 360.0 / float32(numVertices)))
		angle2 := float64(mgl32.DegToRad(float32(i+1) * 360.0 / float32(numVertices)))
		// Check if x-axis is major axis
		var deltaX1, deltaZ1, deltaX2, deltaZ2 float32
		if xAxisMajor {
			deltaX1 = majorAxis * float32(math.Cos(angle1))
			deltaZ1 = minorAxis * float32(math.Sin(angle1))
			deltaX2 = majorAxis * float32(math.Cos(angle2))
			deltaZ2 = minorAxis * float32(math.Sin(angle2))
		} else {
			deltaX1 = minorAxis * float32(math.Cos(angle1))
			deltaZ1 = majorAxis * float32(math.Sin(angle1))
			deltaX2 = minorAxis * float32(math.Cos(angle2))
			deltaZ2 = majorAxis * float32(math.Sin(angle2))
		}

		vertex1 := []float32{centerVertex.X() + deltaX1, centerVertex.Y(), centerVertex.Z() + deltaZ1}
		vertex2 := []float32{centerVertex.X() + deltaX2, centerVertex.Y(), centerVertex.Z() + deltaZ2}
		ellipseBuffer = append(ellipseBuffer, center...)
		ellipseBuffer = append(ellipseBuffer, vertex1...)
		ellipseBuffer = append(ellipseBuffer, vertex2...)
	}
	return ellipseBuffer
}
