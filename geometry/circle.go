package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Circle struct {
	Center       mgl32.Vec3
	Radius       float32
	VertexBuffer []float32
}

func NewCircle(centerVertex mgl32.Vec3, radius float32) *Circle {
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

	return &Circle{
		Center:       centerVertex,
		Radius:       radius,
		VertexBuffer: circleBuffer,
	}
}
