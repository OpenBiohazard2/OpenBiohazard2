package geometry

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewEllipseXAxisMajor(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	majorAxis := float32(4.0)
	minorAxis := float32(2.0)
	xAxisMajor := true

	ellipse := NewEllipse(center, majorAxis, minorAxis, xAxisMajor)

	// Test basic properties
	if ellipse.Center != center {
		t.Errorf("Expected center %v, got %v", center, ellipse.Center)
	}
	if ellipse.MajorAxis != majorAxis {
		t.Errorf("Expected major axis %f, got %f", majorAxis, ellipse.MajorAxis)
	}
	if ellipse.MinorAxis != minorAxis {
		t.Errorf("Expected minor axis %f, got %f", minorAxis, ellipse.MinorAxis)
	}
	if ellipse.XAxisMajor != xAxisMajor {
		t.Errorf("Expected XAxisMajor %v, got %v", xAxisMajor, ellipse.XAxisMajor)
	}

	// Test vertex buffer length (8 triangles * 3 vertices * 3 components = 72 floats)
	expectedLength := 8 * 3 * 3
	if len(ellipse.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(ellipse.VertexBuffer))
	}
}

func TestNewEllipseZAxisMajor(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	majorAxis := float32(3.0)
	minorAxis := float32(1.5)
	xAxisMajor := false

	ellipse := NewEllipse(center, majorAxis, minorAxis, xAxisMajor)

	// Test basic properties
	if ellipse.XAxisMajor != xAxisMajor {
		t.Errorf("Expected XAxisMajor %v, got %v", xAxisMajor, ellipse.XAxisMajor)
	}
}

func TestEllipseVertexGenerationXMajor(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	majorAxis := float32(2.0)
	minorAxis := float32(1.0)
	xAxisMajor := true

	ellipse := NewEllipse(center, majorAxis, minorAxis, xAxisMajor)

	// Test that vertices follow ellipse equation: (x/a)² + (z/b)² = 1
	// where a = majorAxis, b = minorAxis for X-major ellipse
	tolerance := float32(0.01)
	for i := 0; i < len(ellipse.VertexBuffer); i += 3 {
		x := ellipse.VertexBuffer[i]
		z := ellipse.VertexBuffer[i+2]

		// Skip center vertices (every 3rd vertex starting from 0)
		if i%9 == 0 {
			continue
		}

		// Check ellipse equation
		ellipseValue := (x*x)/(majorAxis*majorAxis) + (z*z)/(minorAxis*minorAxis)
		if math.Abs(float64(ellipseValue-1.0)) > float64(tolerance) {
			t.Errorf("Vertex at index %d: ellipse equation value %f, expected ~1.0", i, ellipseValue)
		}
	}
}

func TestEllipseVertexGenerationZMajor(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	majorAxis := float32(3.0)
	minorAxis := float32(1.5)
	xAxisMajor := false

	ellipse := NewEllipse(center, majorAxis, minorAxis, xAxisMajor)

	// Test that vertices follow ellipse equation: (x/b)² + (z/a)² = 1
	// where a = majorAxis, b = minorAxis for Z-major ellipse
	tolerance := float32(0.01)
	for i := 0; i < len(ellipse.VertexBuffer); i += 3 {
		x := ellipse.VertexBuffer[i]
		z := ellipse.VertexBuffer[i+2]

		// Skip center vertices
		if i%9 == 0 {
			continue
		}

		// Check ellipse equation (swapped for Z-major)
		ellipseValue := (x*x)/(minorAxis*minorAxis) + (z*z)/(majorAxis*majorAxis)
		if math.Abs(float64(ellipseValue-1.0)) > float64(tolerance) {
			t.Errorf("Vertex at index %d: ellipse equation value %f, expected ~1.0", i, ellipseValue)
		}
	}
}

func TestEllipseWithOffset(t *testing.T) {
	center := mgl32.Vec3{5, 10, 15}
	majorAxis := float32(2.0)
	minorAxis := float32(1.0)
	xAxisMajor := true

	ellipse := NewEllipse(center, majorAxis, minorAxis, xAxisMajor)

	// Test that all vertices are offset by the center
	for i := 0; i < len(ellipse.VertexBuffer); i += 3 {
		x := ellipse.VertexBuffer[i]
		y := ellipse.VertexBuffer[i+1]
		z := ellipse.VertexBuffer[i+2]

		// Skip center vertices
		if i%9 == 0 {
			if x != center.X() || y != center.Y() || z != center.Z() {
				t.Errorf("Center vertex at index %d: expected %v, got (%f, %f, %f)",
					i, center, x, y, z)
			}
			continue
		}

		// Check that vertex is within ellipse bounds relative to center
		dx := x - center.X()
		dz := z - center.Z()
		ellipseValue := (dx*dx)/(majorAxis*majorAxis) + (dz*dz)/(minorAxis*minorAxis)

		tolerance := float32(0.01)
		if math.Abs(float64(ellipseValue-1.0)) > float64(tolerance) {
			t.Errorf("Vertex at index %d: ellipse equation value %f, expected ~1.0", i, ellipseValue)
		}
	}
}
