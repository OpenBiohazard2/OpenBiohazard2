package fileio

// .md1 - 3D model file

import (
	"encoding/binary"
	"io"
)

type MD1Header struct {
	SectionLengthBytes uint32
	Unknown            uint32
	NumObj             uint32
}

type MD1ObjectHeader struct {
	TrianglesHeader MD1TrianglesHeader
	QuadsHeader     MD1QuadsHeader
}

type MD1TrianglesHeader struct {
	VertexOffset        uint32
	VertexCount         uint32
	NormalOffset        uint32
	NormalCount         uint32
	TriangleIndexOffset uint32
	TriangleIndexCount  uint32
	TextureOffset       uint32
}

type MD1QuadsHeader struct {
	VertexOffset    uint32
	VertexCount     uint32
	NormalOffset    uint32
	NormalCount     uint32
	QuadIndexOffset uint32
	QuadIndexCount  uint32
	TextureOffset   uint32
}

type MD1Object struct {
	TriangleVertices []MD1Vertex
	TriangleNormals  []MD1Vertex
	TriangleIndices  []MD1TriangleIndex
	TriangleTextures []MD1TriangleTexture
	QuadVertices     []MD1Vertex
	QuadNormals      []MD1Vertex
	QuadIndices      []MD1QuadIndex
	QuadTextures     []MD1QuadTexture
}

type MD1Vertex struct {
	X    int16
	Y    int16
	Z    int16
	Zero uint16
}

type MD1TriangleIndex struct {
	IndexNormal0 uint16
	IndexVertex0 uint16
	IndexNormal1 uint16
	IndexVertex1 uint16
	IndexNormal2 uint16
	IndexVertex2 uint16
}

type MD1TriangleTexture struct {
	U0     uint8 // UV coordinates for vertex 0
	V0     uint8
	ClutId uint16 // Texture clut id, bits 0-5
	U1     uint8  // UV coordinates for vertex 1
	V1     uint8
	Page   uint16 // Texture page
	U2     uint8  // UV coordinates for vertex 2
	V2     uint8
	Zero   uint16
}

type MD1QuadIndex struct {
	IndexNormal0 uint16
	IndexVertex0 uint16
	IndexNormal1 uint16
	IndexVertex1 uint16
	IndexNormal2 uint16
	IndexVertex2 uint16
	IndexNormal3 uint16
	IndexVertex3 uint16
}

type MD1QuadTexture struct {
	U0     uint8 // UV coordinates for vertex 0
	V0     uint8
	ClutId uint16 // Texture clut id, bits 0-5
	U1     uint8  // UV coordinates for vertex 1
	V1     uint8
	Page   uint16 // Texture page
	U2     uint8  // UV coordinates for vertex 2
	V2     uint8
	Zero1  uint16
	U3     uint8 // UV coordinates for vertex 3
	V3     uint8
	Zero2  uint16
}

type MD1Output struct {
	Components []MD1Object
}

func LoadMD1Stream(r io.ReaderAt, fileLength int64) (*MD1Output, error) {
	fileReader := io.NewSectionReader(r, int64(0), fileLength)

	// Read header
	md1Header := MD1Header{}
	if err := binary.Read(fileReader, binary.LittleEndian, &md1Header); err != nil {
		return nil, err
	}

	// Read header offsets
	modelObjectHeaders := make([]MD1ObjectHeader, int(md1Header.NumObj)/2)
	if err := binary.Read(fileReader, binary.LittleEndian, &modelObjectHeaders); err != nil {
		return nil, err
	}

	// Offsets are after model header, which is 12 bytes
	beginOffset := int64(12)
	objects := make([]MD1Object, len(modelObjectHeaders))
	for i := 0; i < len(modelObjectHeaders); i++ {
		modelObjectHeader := modelObjectHeaders[i]
		// Triangle data
		offset := beginOffset + int64(modelObjectHeader.TrianglesHeader.VertexOffset)
		reader := io.NewSectionReader(r, offset, fileLength-offset)
		triangleVertices := make([]MD1Vertex, modelObjectHeader.TrianglesHeader.VertexCount)
		if err := binary.Read(reader, binary.LittleEndian, &triangleVertices); err != nil {
			return nil, err
		}

		offset = beginOffset + int64(modelObjectHeader.TrianglesHeader.NormalOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		triangleNormals := make([]MD1Vertex, modelObjectHeader.TrianglesHeader.NormalCount)
		if err := binary.Read(reader, binary.LittleEndian, &triangleNormals); err != nil {
			return nil, err
		}

		offset = beginOffset + int64(modelObjectHeader.TrianglesHeader.TriangleIndexOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		triangleIndices := make([]MD1TriangleIndex, modelObjectHeader.TrianglesHeader.TriangleIndexCount)
		if err := binary.Read(reader, binary.LittleEndian, &triangleIndices); err != nil {
			return nil, err
		}

		offset = beginOffset + int64(modelObjectHeader.TrianglesHeader.TextureOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		triangleTextures := make([]MD1TriangleTexture, modelObjectHeader.TrianglesHeader.TriangleIndexCount)
		if err := binary.Read(reader, binary.LittleEndian, &triangleTextures); err != nil {
			return nil, err
		}

		// Quad data
		offset = beginOffset + int64(modelObjectHeader.QuadsHeader.VertexOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		quadVertices := make([]MD1Vertex, modelObjectHeader.QuadsHeader.VertexCount)
		if err := binary.Read(reader, binary.LittleEndian, &quadVertices); err != nil {
			return nil, err
		}

		offset = beginOffset + int64(modelObjectHeader.QuadsHeader.NormalOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		quadNormals := make([]MD1Vertex, modelObjectHeader.QuadsHeader.NormalCount)
		if err := binary.Read(reader, binary.LittleEndian, &quadNormals); err != nil {
			return nil, err
		}

		// A quad has 2 triangles
		offset = beginOffset + int64(modelObjectHeader.QuadsHeader.QuadIndexOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		quadIndices := make([]MD1QuadIndex, modelObjectHeader.QuadsHeader.QuadIndexCount)
		if err := binary.Read(reader, binary.LittleEndian, &quadIndices); err != nil {
			return nil, err
		}

		offset = beginOffset + int64(modelObjectHeader.QuadsHeader.TextureOffset)
		reader = io.NewSectionReader(r, offset, fileLength-offset)
		quadTextures := make([]MD1QuadTexture, modelObjectHeader.QuadsHeader.QuadIndexCount)
		if err := binary.Read(reader, binary.LittleEndian, &quadTextures); err != nil {
			return nil, err
		}

		objects[i] = MD1Object{
			TriangleVertices: triangleVertices,
			TriangleNormals:  triangleNormals,
			TriangleIndices:  triangleIndices,
			TriangleTextures: triangleTextures,
			QuadVertices:     quadVertices,
			QuadNormals:      quadNormals,
			QuadIndices:      quadIndices,
			QuadTextures:     quadTextures,
		}
	}

	md1Output := &MD1Output{
		Components: objects,
	}
	return md1Output, nil
}
