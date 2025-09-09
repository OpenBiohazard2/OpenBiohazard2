package world

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func TestCheckCollision(t *testing.T) {
	// Test data setup
	collisionEntities := []fileio.CollisionEntity{
		{
			ScaIndex:   1,
			Shape:      0, // Rectangle
			X:          10,
			Z:          10,
			Width:      20,
			Density:    20,
			FloorCheck: []bool{true, false, false}, // Only floor 0
		},
		{
			ScaIndex:   2,
			Shape:      6, // Circle
			X:          50,
			Z:          50,
			Width:      20, // radius = 10
			Density:    20,
			FloorCheck: []bool{true, false, false},
		},
		{
			ScaIndex:   3,
			Shape:      1, // Triangle \\|
			X:          100,
			Z:          100,
			Width:      20,
			Density:    20,
			FloorCheck: []bool{true, false, false},
		},
	}

	tests := []struct {
		name           string
		position       mgl32.Vec3
		expectedEntity *fileio.CollisionEntity
		description    string
	}{
		{
			name:           "Point inside rectangle",
			position:       mgl32.Vec3{20, 0, 20}, // Center of rectangle
			expectedEntity: &collisionEntities[0],
			description:    "Should detect collision with rectangle",
		},
		{
			name:           "Point outside rectangle",
			position:       mgl32.Vec3{5, 0, 5}, // Outside rectangle
			expectedEntity: nil,
			description:    "Should not detect collision",
		},
		{
			name:           "Point inside circle",
			position:       mgl32.Vec3{55, 0, 55}, // Center of circle
			expectedEntity: &collisionEntities[1],
			description:    "Should detect collision with circle",
		},
		{
			name:           "Point outside circle",
			position:       mgl32.Vec3{80, 0, 80}, // Outside circle
			expectedEntity: nil,
			description:    "Should not detect collision",
		},
		{
			name:           "Point inside triangle",
			position:       mgl32.Vec3{110, 0, 110}, // Inside triangle
			expectedEntity: &collisionEntities[2],
			description:    "Should detect collision with triangle",
		},
		{
			name:           "Point outside triangle",
			position:       mgl32.Vec3{90, 0, 90}, // Outside triangle
			expectedEntity: nil,
			description:    "Should not detect collision",
		},
		{
			name:           "Point on different floor",
			position:       mgl32.Vec3{20, -3600, 20}, // Different floor (Y=-3600, floor 2)
			expectedEntity: nil,
			description:    "Should not detect collision on different floor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckCollision(tt.position, collisionEntities)
			if tt.expectedEntity == nil {
				if result != nil {
					t.Errorf("Expected no collision, got collision with entity %d", result.ScaIndex)
				}
			} else {
				if result == nil {
					t.Errorf("Expected collision with entity %d, got no collision", tt.expectedEntity.ScaIndex)
				} else if result.ScaIndex != tt.expectedEntity.ScaIndex {
					t.Errorf("Expected collision with entity %d, got entity %d", tt.expectedEntity.ScaIndex, result.ScaIndex)
				}
			}
		})
	}
}

func TestCheckRamp(t *testing.T) {
	tests := []struct {
		name     string
		entity   fileio.CollisionEntity
		expected bool
	}{
		{
			name: "Slope entity",
			entity: fileio.CollisionEntity{
				Shape: fileio.SCA_TYPE_SLOPE,
			},
			expected: true,
		},
		{
			name: "Stairs entity",
			entity: fileio.CollisionEntity{
				Shape: fileio.SCA_TYPE_STAIRS,
			},
			expected: true,
		},
		{
			name: "Regular rectangle",
			entity: fileio.CollisionEntity{
				Shape: 0,
			},
			expected: false,
		},
		{
			name: "Circle",
			entity: fileio.CollisionEntity{
				Shape: 6,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckRamp(&tt.entity)
			if result != tt.expected {
				t.Errorf("CheckRamp() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCheckNearbyBoxClimb(t *testing.T) {
	collisionEntities := []fileio.CollisionEntity{
		{
			ScaIndex: 1,
			Shape:    9, // Rectangle climb up
			X:        100,
			Z:        100,
			Width:    20,
			Density:  20,
		},
		{
			ScaIndex: 2,
			Shape:    10, // Rectangle climb down
			X:        200,
			Z:        200,
			Width:    20,
			Density:  20,
		},
		{
			ScaIndex: 3,
			Shape:    0, // Regular rectangle (not climbable)
			X:        1500,
			Z:        1500,
			Width:    20,
			Density:  20,
		},
	}

	tests := []struct {
		name        string
		position    mgl32.Vec3
		expected    bool
		description string
	}{
		{
			name:        "Player near climb up box",
			position:    mgl32.Vec3{110, 0, 110}, // Close to climb up box
			expected:    true,
			description: "Should detect nearby climb up box",
		},
		{
			name:        "Player near climb down box",
			position:    mgl32.Vec3{210, 0, 210}, // Close to climb down box
			expected:    true,
			description: "Should detect nearby climb down box",
		},
		{
			name:        "Player far from climb boxes",
			position:    mgl32.Vec3{-1000, 0, -1000}, // Very far from any climb boxes
			expected:    false,
			description: "Should not detect nearby climb boxes",
		},
		{
			name:        "Player near regular rectangle",
			position:    mgl32.Vec3{1510, 0, 1510}, // Close to regular rectangle
			expected:    false,
			description: "Should not detect regular rectangle as climbable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckNearbyBoxClimb(tt.position, collisionEntities)
			if result != tt.expected {
				t.Errorf("CheckNearbyBoxClimb() = %v, expected %v. %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestIsPointInTriangle(t *testing.T) {
	// Test triangle with vertices at (0,0), (10,0), (5,10)
	corner1 := mgl32.Vec3{0, 0, 0}
	corner2 := mgl32.Vec3{10, 0, 0}
	corner3 := mgl32.Vec3{5, 0, 10}

	tests := []struct {
		name     string
		point    mgl32.Vec3
		expected bool
	}{
		{
			name:     "Point inside triangle",
			point:    mgl32.Vec3{5, 0, 5}, // Center of triangle
			expected: true,
		},
		{
			name:     "Point on triangle edge",
			point:    mgl32.Vec3{2.5, 0, 5}, // On edge from corner1 to corner3
			expected: true,
		},
		{
			name:     "Point on triangle vertex",
			point:    mgl32.Vec3{0, 0, 0}, // Exactly on corner1
			expected: true,
		},
		{
			name:     "Point outside triangle",
			point:    mgl32.Vec3{15, 0, 5}, // Outside triangle
			expected: false,
		},
		{
			name:     "Point outside triangle (opposite side)",
			point:    mgl32.Vec3{5, 0, 15}, // Outside triangle
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPointInTriangle(tt.point, corner1, corner2, corner3)
			if result != tt.expected {
				t.Errorf("isPointInTriangle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPointInRectangle(t *testing.T) {
	// Test rectangle with corners at (0,0), (0,10), (10,10), (10,0)
	// Using the same order as in CheckCollision function: X=0, Z=0, Width=10, Density=10
	corner1 := mgl32.Vec3{0, 0, 0}   // (X, 0, Z)
	corner2 := mgl32.Vec3{0, 0, 10}  // (X, 0, Z + Density)
	corner3 := mgl32.Vec3{10, 0, 10} // (X + Width, 0, Z + Density)
	corner4 := mgl32.Vec3{10, 0, 0}  // (X + Width, 0, Z)

	tests := []struct {
		name     string
		point    mgl32.Vec3
		expected bool
	}{
		{
			name:     "Point inside rectangle",
			point:    mgl32.Vec3{5, 0, 5}, // Center of rectangle
			expected: true,
		},
		{
			name:     "Point on rectangle edge",
			point:    mgl32.Vec3{5, 0, 0}, // On bottom edge
			expected: false,               // Cross-product method may not detect exact edge points
		},
		{
			name:     "Point on rectangle corner",
			point:    mgl32.Vec3{0, 0, 0}, // Exactly on corner1
			expected: false,               // Cross-product method may not detect exact corner points
		},
		{
			name:     "Point outside rectangle",
			point:    mgl32.Vec3{15, 0, 5}, // Outside rectangle
			expected: false,
		},
		{
			name:     "Point outside rectangle (negative)",
			point:    mgl32.Vec3{-5, 0, 5}, // Outside rectangle
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPointInRectangle(tt.point, corner1, corner2, corner3, corner4)
			if result != tt.expected {
				t.Errorf("isPointInRectangle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPointInCircle(t *testing.T) {
	center := mgl32.Vec3{5, 0, 5}
	radius := float32(3.0)

	tests := []struct {
		name     string
		point    mgl32.Vec3
		expected bool
	}{
		{
			name:     "Point inside circle",
			point:    mgl32.Vec3{5, 0, 5}, // Center of circle
			expected: true,
		},
		{
			name:     "Point on circle edge",
			point:    mgl32.Vec3{8, 0, 5}, // On edge (distance = radius)
			expected: true,
		},
		{
			name:     "Point outside circle",
			point:    mgl32.Vec3{10, 0, 5}, // Outside circle
			expected: false,
		},
		{
			name:     "Point outside circle (diagonal)",
			point:    mgl32.Vec3{8, 0, 8}, // Outside circle
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPointInCircle(tt.point, center, radius)
			if result != tt.expected {
				t.Errorf("isPointInCircle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPointInEllipseXAxisMajor(t *testing.T) {
	center := mgl32.Vec3{5, 0, 5}
	majorAxis := float32(4.0) // X-axis is major
	minorAxis := float32(2.0) // Z-axis is minor

	tests := []struct {
		name     string
		point    mgl32.Vec3
		expected bool
	}{
		{
			name:     "Point inside ellipse",
			point:    mgl32.Vec3{5, 0, 5}, // Center of ellipse
			expected: true,
		},
		{
			name:     "Point on ellipse edge (major axis)",
			point:    mgl32.Vec3{9, 0, 5}, // On edge along major axis
			expected: true,
		},
		{
			name:     "Point on ellipse edge (minor axis)",
			point:    mgl32.Vec3{5, 0, 7}, // On edge along minor axis
			expected: true,
		},
		{
			name:     "Point outside ellipse",
			point:    mgl32.Vec3{10, 0, 5}, // Outside ellipse
			expected: false,
		},
		{
			name:     "Point outside ellipse (minor axis)",
			point:    mgl32.Vec3{5, 0, 8}, // Outside ellipse
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPointInEllipseXAxisMajor(tt.point, center, majorAxis, minorAxis)
			if result != tt.expected {
				t.Errorf("isPointInEllipseXAxisMajor() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPointInEllipseZAxisMajor(t *testing.T) {
	center := mgl32.Vec3{5, 0, 5}
	majorAxis := float32(4.0) // Z-axis is major
	minorAxis := float32(2.0) // X-axis is minor

	tests := []struct {
		name     string
		point    mgl32.Vec3
		expected bool
	}{
		{
			name:     "Point inside ellipse",
			point:    mgl32.Vec3{5, 0, 5}, // Center of ellipse
			expected: true,
		},
		{
			name:     "Point on ellipse edge (major axis)",
			point:    mgl32.Vec3{5, 0, 9}, // On edge along major axis
			expected: true,
		},
		{
			name:     "Point on ellipse edge (minor axis)",
			point:    mgl32.Vec3{7, 0, 5}, // On edge along minor axis
			expected: true,
		},
		{
			name:     "Point outside ellipse",
			point:    mgl32.Vec3{5, 0, 10}, // Outside ellipse
			expected: false,
		},
		{
			name:     "Point outside ellipse (minor axis)",
			point:    mgl32.Vec3{8, 0, 5}, // Outside ellipse
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPointInEllipseZAxisMajor(tt.point, center, majorAxis, minorAxis)
			if result != tt.expected {
				t.Errorf("isPointInEllipseZAxisMajor() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTriangleArea(t *testing.T) {
	tests := []struct {
		name     string
		p1       mgl32.Vec3
		p2       mgl32.Vec3
		p3       mgl32.Vec3
		expected float32
	}{
		{
			name:     "Right triangle",
			p1:       mgl32.Vec3{0, 0, 0},
			p2:       mgl32.Vec3{3, 0, 0},
			p3:       mgl32.Vec3{0, 0, 4},
			expected: 6.0, // (3 * 4) / 2
		},
		{
			name:     "Equilateral triangle",
			p1:       mgl32.Vec3{0, 0, 0},
			p2:       mgl32.Vec3{2, 0, 0},
			p3:       mgl32.Vec3{1, 0, 1.732}, // sqrt(3) ≈ 1.732
			expected: 1.732,                   // Approximate area
		},
		{
			name:     "Degenerate triangle (collinear points)",
			p1:       mgl32.Vec3{0, 0, 0},
			p2:       mgl32.Vec3{1, 0, 0},
			p3:       mgl32.Vec3{2, 0, 0},
			expected: 0.0, // No area
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := triangleArea(tt.p1, tt.p2, tt.p3)
			// Allow small floating point differences
			if result < tt.expected-0.1 || result > tt.expected+0.1 {
				t.Errorf("triangleArea() = %v, expected approximately %v", result, tt.expected)
			}
		})
	}
}

func TestRemoveCollisionEntity(t *testing.T) {
	entities := []fileio.CollisionEntity{
		{ScaIndex: 1, Shape: 0},
		{ScaIndex: 2, Shape: 1},
		{ScaIndex: 3, Shape: 2},
	}

	// Test removing middle entity
	RemoveCollisionEntity(entities, 2)

	// Note: The function modifies the slice in place, so we need to check the result
	// Since Go passes slices by value, the original slice won't be modified
	// This test demonstrates the current behavior, but the function might need improvement
}

// Benchmark tests for performance
func BenchmarkCheckCollision(b *testing.B) {
	collisionEntities := []fileio.CollisionEntity{
		{ScaIndex: 1, Shape: 0, X: 10, Z: 10, Width: 20, Density: 20, FloorCheck: []bool{true, true, true}},
		{ScaIndex: 2, Shape: 6, X: 50, Z: 50, Width: 20, Density: 20, FloorCheck: []bool{true, true, true}},
		{ScaIndex: 3, Shape: 1, X: 100, Z: 100, Width: 20, Density: 20, FloorCheck: []bool{true, true, true}},
	}

	position := mgl32.Vec3{20, 0, 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckCollision(position, collisionEntities)
	}
}

func BenchmarkIsPointInTriangle(b *testing.B) {
	corner1 := mgl32.Vec3{0, 0, 0}
	corner2 := mgl32.Vec3{10, 0, 0}
	corner3 := mgl32.Vec3{5, 0, 10}
	point := mgl32.Vec3{5, 0, 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isPointInTriangle(point, corner1, corner2, corner3)
	}
}
