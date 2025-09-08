package geometry

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewCircle(t *testing.T) {
	center := mgl32.Vec3{1, 2, 3}
	radius := float32(5.0)
	circle := NewCircle(center, radius)

	// Test basic properties
	if circle.Center != center {
		t.Errorf("Expected center %v, got %v", center, circle.Center)
	}
	if circle.Radius != radius {
		t.Errorf("Expected radius %f, got %f", radius, circle.Radius)
	}

	// Test vertex buffer length (8 triangles * 3 vertices * 3 components = 72 floats)
	expectedLength := 8 * 3 * 3
	if len(circle.VertexBuffer) != expectedLength {
		t.Errorf("Expected vertex buffer length %d, got %d", expectedLength, len(circle.VertexBuffer))
	}
}

func TestCircleVertexGeneration(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	radius := float32(1.0)
	circle := NewCircle(center, radius)

	// Test that all vertices are approximately at the correct distance from center
	tolerance := float32(0.01)
	for i := 0; i < len(circle.VertexBuffer); i += 3 {
		x := circle.VertexBuffer[i]
		z := circle.VertexBuffer[i+2]

		// Skip center vertices (every 3rd vertex starting from 0)
		if i%9 == 0 {
			continue
		}

		distance := float32(math.Sqrt(float64(x*x + z*z)))
		if math.Abs(float64(distance-radius)) > float64(tolerance) {
			t.Errorf("Vertex at index %d: expected distance ~%f, got %f", i, radius, distance)
		}
	}
}

func TestCircleAtOrigin(t *testing.T) {
	center := mgl32.Vec3{0, 0, 0}
	radius := float32(2.0)
	circle := NewCircle(center, radius)

	// Test that the first non-center vertex is at the expected position
	// First triangle: center, vertex1, vertex2
	// vertex1 should be at (radius, 0, 0) for angle 0
	vertex1X := circle.VertexBuffer[3] // First vertex after center
	vertex1Z := circle.VertexBuffer[5] // Z component of first vertex

	expectedX := radius
	expectedZ := float32(0)
	tolerance := float32(0.01)

	if math.Abs(float64(vertex1X-expectedX)) > float64(tolerance) {
		t.Errorf("Expected first vertex X ~%f, got %f", expectedX, vertex1X)
	}
	if math.Abs(float64(vertex1Z-expectedZ)) > float64(tolerance) {
		t.Errorf("Expected first vertex Z ~%f, got %f", expectedZ, vertex1Z)
	}
}

func TestCircleWithOffset(t *testing.T) {
	center := mgl32.Vec3{10, 20, 30}
	radius := float32(3.0)
	circle := NewCircle(center, radius)

	// Test that all vertices are offset by the center
	for i := 0; i < len(circle.VertexBuffer); i += 3 {
		x := circle.VertexBuffer[i]
		y := circle.VertexBuffer[i+1]
		z := circle.VertexBuffer[i+2]

		// Skip center vertices
		if i%9 == 0 {
			if x != center.X() || y != center.Y() || z != center.Z() {
				t.Errorf("Center vertex at index %d: expected %v, got (%f, %f, %f)",
					i, center, x, y, z)
			}
			continue
		}

		// Check that vertex is within radius of center
		dx := x - center.X()
		dy := y - center.Y()
		dz := z - center.Z()
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))

		tolerance := float32(0.01)
		if math.Abs(float64(distance-radius)) > float64(tolerance) {
			t.Errorf("Vertex at index %d: distance from center %f, expected ~%f",
				i, distance, radius)
		}
	}
}
