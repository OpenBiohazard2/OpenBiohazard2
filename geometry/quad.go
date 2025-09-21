package geometry

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

type Quad struct {
	Vertices     [4]mgl32.Vec3
	VertexBuffer []float32
}

func NewQuad(corners [4]mgl32.Vec3) *Quad {
	vertices := make([][]float32, 4)
	for i := 0; i < 4; i++ {
		vertices[i] = []float32{corners[i].X(), corners[i].Y(), corners[i].Z()}
	}

	vertexBuffer := make([]float32, 0)
	// v0, v1, v2
	tri1Indices := [3]int{0, 1, 2}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
	}
	// v0, v3, v2
	tri2Indices := [3]int{0, 3, 2}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
	}
	return &Quad{
		Vertices:     corners,
		VertexBuffer: vertexBuffer,
	}
}

func NewQuadFourPoints(xzPairs [4][]float32) *Quad {
	corners := [4]mgl32.Vec3{}
	for i := 0; i < 4; i++ {
		corners[i] = mgl32.Vec3{xzPairs[i][0], 0, xzPairs[i][1]}
	}
	return NewQuad(corners)
}

func NewRectangle(x float32, z float32, width float32, depth float32) *Quad {
	corners := [4]mgl32.Vec3{
		{x, 0, z},
		{x, 0, z + depth},
		{x + width, 0, z + depth},
		{x + width, 0, z},
	}
	return NewQuad(corners)
}

func NewTexturedRectangle(vertices [4][]float32, uvs [4][]float32) *Quad {
	vertexBuffer := make([]float32, 0)

	// v0, v1, v2
	tri1Indices := [3]int{0, 1, 2}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
	}
	// v0, v3, v2
	tri2Indices := [3]int{0, 3, 2}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
	}

	return &Quad{
		VertexBuffer: vertexBuffer,
	}
}

func NewQuadMD1(vertices [4][]float32, uvs [4][]float32, normals [4][]float32) *Quad {
	vertexBuffer := make([]float32, 0)

	// MD1 vertex order is different (v0, v1, v3, v2)
	// Vertex order of other quads is (v0, v1, v2, v3)
	// v0, v1, v3
	tri1Indices := [3]int{0, 1, 3}
	for _, index := range tri1Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
		vertexBuffer = append(vertexBuffer, normals[index]...)
	}
	// v0, v2, v3
	tri2Indices := [3]int{0, 2, 3}
	for _, index := range tri2Indices {
		vertexBuffer = append(vertexBuffer, vertices[index]...)
		vertexBuffer = append(vertexBuffer, uvs[index]...)
		vertexBuffer = append(vertexBuffer, normals[index]...)
	}

	return &Quad{
		VertexBuffer: vertexBuffer,
	}
}

// NewFullScreenQuad creates a full-screen quad for 2D rendering
func NewFullScreenQuad(z float32) *Quad {
	vertices := [4][]float32{
		{-1.0, 1.0, z},
		{-1.0, -1.0, z},
		{1.0, -1.0, z},
		{1.0, 1.0, z},
	}
	uvs := [4][]float32{
		{0.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
	}
	return NewTexturedRectangle(vertices, uvs)
}

// NewCameraMaskQuad creates a textured rectangle for camera mask rendering
// It takes a MaskRectangle object and generates vertices and UVs using coordinate utilities
func NewCameraMaskQuad(cameraMask fileio.MaskRectangle, depth float32) *Quad {
	// Extract fields from the camera mask object
	destX := float32(cameraMask.DestX)
	destY := float32(cameraMask.DestY)
	maskWidth := float32(cameraMask.Width)
	maskHeight := float32(cameraMask.Height)
	
	// Create corners in screen space
	corners := [4][]float32{
		{destX, destY},
		{destX + maskWidth, destY},
		{destX + maskWidth, destY + maskHeight},
		{destX, destY + maskHeight},
	}

	vertices := [4][]float32{}
	uvs := [4][]float32{}
	for i, corner := range corners {
		x := corner[0]
		y := corner[1]
		vertices[i] = []float32{ConvertToScreenX(x), ConvertToScreenY(y), depth}
		uvs[i] = []float32{ConvertToTextureU(x), ConvertToTextureV(y)}
	}
	return NewTexturedRectangle(vertices, uvs)
}

// NewBillboardSprite creates a camera-aligned billboard sprite
func NewBillboardSprite(spriteCenter mgl32.Vec3, spriteWidth float32, viewMatrix mgl32.Mat4) *Quad {
	// Generate billboard sprite vertices (unit square)
	squareVertices := [4]mgl32.Vec3{
		{0, 1, 0},
		{1, 1, 0},
		{1, 0, 0},
		{0, 0, 0},
	}

	// Extract camera orientation from view matrix
	cameraRight := mgl32.Vec3{viewMatrix.At(0, 0), viewMatrix.At(1, 0), viewMatrix.At(2, 0)}
	cameraUp := mgl32.Vec3{viewMatrix.At(0, 1), viewMatrix.At(1, 1), viewMatrix.At(2, 1)}

	// Calculate world-space positions for camera-aligned billboard
	renderVertices := [4][]float32{}
	for i := 0; i < 4; i++ {
		x := squareVertices[i].X()
		y := squareVertices[i].Y()
		worldspacePosition := spriteCenter.Add(cameraRight.Mul(x * spriteWidth)).Add(cameraUp.Mul(y * spriteWidth))
		renderVertices[i] = []float32{worldspacePosition.X(), worldspacePosition.Y(), worldspacePosition.Z()}
	}

	// Standard UV coordinates for full texture
	uvs := [4][]float32{
		{0.0, 0.0},
		{1.0, 0.0},
		{1.0, 1.0},
		{0.0, 1.0},
	}

	return NewTexturedRectangle(renderVertices, uvs)
}