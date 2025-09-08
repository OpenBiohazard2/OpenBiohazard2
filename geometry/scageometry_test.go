package geometry

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func TestNewSlopedRectangleType0(t *testing.T) {
	// Test slope type 0: starts from x-axis, slopes up
	entity := fileio.CollisionEntity{
		X:           10,
		Z:           20,
		Width:       5,
		Density:     3,
		SlopeHeight: 2,
		SlopeType:   0,
	}

	quad := NewSlopedRectangle(entity)

	if quad == nil {
		t.Fatal("Expected quad, got nil")
	}

	// Check that vertices are in correct order for type 0
	expectedVertices := [4]mgl32.Vec3{
		{10, 0, 20}, // bottom-left
		{10, 0, 23}, // bottom-right (Z + Density)
		{15, 2, 23}, // top-right (X + Width, SlopeHeight, Z + Density)
		{15, 2, 20}, // top-left (X + Width, SlopeHeight, Z)
	}

	for i, expected := range expectedVertices {
		if quad.Vertices[i] != expected {
			t.Errorf("Type 0 vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}
}

func TestNewSlopedRectangleType1(t *testing.T) {
	// Test slope type 1: starts from x-axis, slopes down
	entity := fileio.CollisionEntity{
		X:           5,
		Z:           15,
		Width:       8,
		Density:     4,
		SlopeHeight: 3,
		SlopeType:   1,
	}

	quad := NewSlopedRectangle(entity)

	if quad == nil {
		t.Fatal("Expected quad, got nil")
	}

	// Check that vertices are in correct order for type 1
	expectedVertices := [4]mgl32.Vec3{
		{5, 3, 15},  // top-left (X, SlopeHeight, Z)
		{5, 3, 19},  // top-right (X, SlopeHeight, Z + Density)
		{13, 0, 19}, // bottom-right (X + Width, 0, Z + Density)
		{13, 0, 15}, // bottom-left (X + Width, 0, Z)
	}

	for i, expected := range expectedVertices {
		if quad.Vertices[i] != expected {
			t.Errorf("Type 1 vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}
}

func TestNewSlopedRectangleType2(t *testing.T) {
	// Test slope type 2: starts from z-axis, slopes up
	entity := fileio.CollisionEntity{
		X:           0,
		Z:           0,
		Width:       10,
		Density:     5,
		SlopeHeight: 4,
		SlopeType:   2,
	}

	quad := NewSlopedRectangle(entity)

	if quad == nil {
		t.Fatal("Expected quad, got nil")
	}

	// Check that vertices are in correct order for type 2
	expectedVertices := [4]mgl32.Vec3{
		{0, 0, 0},  // bottom-left (X, 0, Z)
		{0, 4, 5},  // top-left (X, SlopeHeight, Z + Density)
		{10, 4, 5}, // top-right (X + Width, SlopeHeight, Z + Density)
		{10, 0, 0}, // bottom-right (X + Width, 0, Z)
	}

	for i, expected := range expectedVertices {
		if quad.Vertices[i] != expected {
			t.Errorf("Type 2 vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}
}

func TestNewSlopedRectangleType3(t *testing.T) {
	// Test slope type 3: starts from z-axis, slopes down
	entity := fileio.CollisionEntity{
		X:           2,
		Z:           3,
		Width:       6,
		Density:     2,
		SlopeHeight: 1,
		SlopeType:   3,
	}

	quad := NewSlopedRectangle(entity)

	if quad == nil {
		t.Fatal("Expected quad, got nil")
	}

	// Check that vertices are in correct order for type 3
	expectedVertices := [4]mgl32.Vec3{
		{2, 1, 3}, // top-left (X, SlopeHeight, Z)
		{2, 0, 5}, // bottom-left (X, 0, Z + Density)
		{8, 0, 5}, // bottom-right (X + Width, 0, Z + Density)
		{8, 1, 3}, // top-right (X + Width, SlopeHeight, Z)
	}

	for i, expected := range expectedVertices {
		if quad.Vertices[i] != expected {
			t.Errorf("Type 3 vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}
}

func TestNewSlopedRectangleInvalidType(t *testing.T) {
	// Test invalid slope type
	entity := fileio.CollisionEntity{
		X:           0,
		Z:           0,
		Width:       1,
		Density:     1,
		SlopeHeight: 1,
		SlopeType:   99, // Invalid type
	}

	quad := NewSlopedRectangle(entity)

	if quad != nil {
		t.Error("Expected nil for invalid slope type, got quad")
	}
}

func TestNewSlopedRectangleEdgeCases(t *testing.T) {
	// Test with zero values
	entity := fileio.CollisionEntity{
		X:           0,
		Z:           0,
		Width:       0,
		Density:     0,
		SlopeHeight: 0,
		SlopeType:   0,
	}

	quad := NewSlopedRectangle(entity)

	if quad == nil {
		t.Fatal("Expected quad, got nil")
	}

	// All vertices should be at origin
	expectedVertices := [4]mgl32.Vec3{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	for i, expected := range expectedVertices {
		if quad.Vertices[i] != expected {
			t.Errorf("Zero values vertex %d: expected %v, got %v", i, expected, quad.Vertices[i])
		}
	}

	// Test with negative values
	entityNegative := fileio.CollisionEntity{
		X:           -5,
		Z:           -3,
		Width:       2,
		Density:     1,
		SlopeHeight: -2,
		SlopeType:   0,
	}

	quadNegative := NewSlopedRectangle(entityNegative)

	if quadNegative == nil {
		t.Fatal("Expected quad, got nil")
	}

	// Check that negative values are handled correctly
	expectedNegative := [4]mgl32.Vec3{
		{-5, 0, -3},
		{-5, 0, -2},
		{-3, -2, -2},
		{-3, -2, -3},
	}

	for i, expected := range expectedNegative {
		if quadNegative.Vertices[i] != expected {
			t.Errorf("Negative values vertex %d: expected %v, got %v", i, expected, quadNegative.Vertices[i])
		}
	}
}
