package geometry

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

// NewDebugRectangle creates a debug rectangle from four corner vertices
func NewDebugRectangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) *Quad {
	return NewQuad([4]mgl32.Vec3{corner1, corner2, corner3, corner4})
}

// NewDebugTriangle creates a debug triangle from three corner vertices
func NewDebugTriangle(corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3) *Triangle {
	triBuffer := make([]float32, 0)
	vertex1 := []float32{corner1.X(), corner1.Y(), corner1.Z()}
	vertex2 := []float32{corner2.X(), corner2.Y(), corner2.Z()}
	vertex3 := []float32{corner3.X(), corner3.Y(), corner3.Z()}

	triBuffer = append(triBuffer, vertex1...)
	triBuffer = append(triBuffer, vertex2...)
	triBuffer = append(triBuffer, vertex3...)
	
	return &Triangle{
		VertexBuffer: triBuffer,
	}
}

// NewCollisionDebugEntity creates a debug entity for collision shapes
func NewCollisionDebugEntity(collisionEntities []fileio.CollisionEntity) []float32 {
	vertexBuffer := make([]float32, 0)
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 0:
			// Rectangle
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex4 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			rect := NewDebugRectangle(vertex1, vertex2, vertex3, vertex4)
			vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
		case 1:
			// Triangle \\|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			tri := NewDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri.VertexBuffer...)
		case 2:
			// Triangle |/
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			tri := NewDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri.VertexBuffer...)
		case 3:
			// Triangle /|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			tri := NewDebugTriangle(vertex1, vertex2, vertex3)
			vertexBuffer = append(vertexBuffer, tri.VertexBuffer...)
		case 6:
			// Circle
			radius := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + radius, 0, float32(entity.Z) + radius}
			circle := NewCircle(center, radius)
			vertexBuffer = append(vertexBuffer, circle.VertexBuffer...)
		case 7:
			// Ellipse, rectangle with rounded corners on the x-axis
			majorAxis := float32(entity.Width) / 2.0
			minorAxis := float32(entity.Density) / 2.0
			center := mgl32.Vec3{float32(entity.X) + majorAxis, 0, float32(entity.Z) + minorAxis}
			ellipse := NewEllipse(center, majorAxis, minorAxis, true)
			vertexBuffer = append(vertexBuffer, ellipse.VertexBuffer...)
		case 8:
			// Ellipse, rectangle with rounded corners on the z-axis
			majorAxis := float32(entity.Density) / 2.0
			minorAxis := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + minorAxis, 0, float32(entity.Z) + majorAxis}
			ellipse := NewEllipse(center, majorAxis, minorAxis, false)
			vertexBuffer = append(vertexBuffer, ellipse.VertexBuffer...)
		}
	}
	return vertexBuffer
}

// NewCameraSwitchDebugVertexBuffer creates vertex buffer for camera switch debug entities
func NewCameraSwitchDebugVertexBuffer(curCameraId int, cameraSwitches []fileio.RVDHeader, cameraSwitchTransitions map[int][]int) []float32 {
	vertexBuffer := make([]float32, 0)
	for _, regionIndex := range cameraSwitchTransitions[curCameraId] {
		cameraSwitch := cameraSwitches[regionIndex]
		corners := [4][]float32{
			{float32(cameraSwitch.X1), float32(cameraSwitch.Z1)},
			{float32(cameraSwitch.X2), float32(cameraSwitch.Z2)},
			{float32(cameraSwitch.X3), float32(cameraSwitch.Z3)},
			{float32(cameraSwitch.X4), float32(cameraSwitch.Z4)},
		}
		rect := NewQuadFourPoints(corners)
		vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
	}
	return vertexBuffer
}

// NewSlopedSurfacesDebugVertexBuffer creates vertex buffer for sloped surfaces debug entities
func NewSlopedSurfacesDebugVertexBuffer(collisionEntities []fileio.CollisionEntity) []float32 {
	vertexBuffer := make([]float32, 0)
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 11:
			// Ramp
			rect := NewSlopedRectangle(entity)
			vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
		case 12:
			// Stairs
			rect := NewSlopedRectangle(entity)
			vertexBuffer = append(vertexBuffer, rect.VertexBuffer...)
		}
	}
	return vertexBuffer
}
