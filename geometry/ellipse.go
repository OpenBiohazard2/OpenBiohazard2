package geometry

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Ellipse struct {
	Center       mgl32.Vec3
	MajorAxis    float32
	MinorAxis    float32
	XAxisMajor   bool
	VertexBuffer []float32
}

func NewEllipse(centerVertex mgl32.Vec3, majorAxis float32, minorAxis float32, xAxisMajor bool) *Ellipse {
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
	return &Ellipse{
		Center:       centerVertex,
		MajorAxis:    majorAxis,
		MinorAxis:    minorAxis,
		XAxisMajor:   xAxisMajor,
		VertexBuffer: ellipseBuffer,
	}
}
