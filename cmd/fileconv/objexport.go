package main

import (
	"fmt"
	"math"
	"os"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

// Basic vector types for 3D geometry
type vec2 struct{ u, v float32 }
type vec3 struct{ x, y, z float32 }

// Face indexing - p=position, t=texture, n=normal indices
type faceIndex struct{ p, t, n int } // -1 if not used
type face struct{ a, b, c, d faceIndex } // d=-1 for triangles, valid for quads

// Mesh data structure
type mesh struct {
	name      string
	positions []vec3
	uvs       []vec2
	normals   []vec3
	faces     []face
	faceMats  []string // material name for each face
}

// Material system for OBJ/MTL export
type MatKey struct{ TPage, Palette int }
type Material struct {
	Name string
	PNG  string // relative path to texture file
}

// Convert MD1 vertex to vec3 (raw coordinates, no scaling)
func buildModelVertex(vertex fileio.MD1Vertex) vec3 {
	return vec3{
		x: float32(vertex.X),
		y: float32(vertex.Y),
		z: float32(vertex.Z),
	}
}

// Convert texture coordinates with page offset and V-flip for OBJ format
func buildTextureUV(u, v float32, texturePage uint16, textureData *fileio.TIMOutput) vec2 {
	if textureData == nil {
		// Fallback: simple normalization with V flip
		return vec2{u: u / 255.0, v: 1.0 - (v / 255.0)}
	}
	// Normalize UV coordinates with texture page offset and V flip
	textureOffsetUnit := float32(textureData.ImageWidth) / float32(textureData.NumPalettes)
	textureCoordOffset := textureOffsetUnit * float32(texturePage&3)
	newU := (u + textureCoordOffset) / float32(textureData.ImageWidth)
	newV := 1.0 - (v / float32(textureData.ImageHeight))
	return vec2{u: newU, v: newV}
}

// Convert and normalize normal vector to unit length
func buildModelNormal(normal fileio.MD1Vertex) vec3 {
	x, y, z := float32(normal.X), float32(normal.Y), float32(normal.Z)
	magnitude := float32(math.Sqrt(float64(x*x + y*y + z*z)))
	if magnitude > 0 {
		return vec3{x: x / magnitude, y: y / magnitude, z: z / magnitude}
	}
	return vec3{x: 0, y: 0, z: 1} // Default up vector if magnitude is 0
}

// Find maximum vertex and normal indices in a component
func findMaxIndices(entityModel *fileio.MD1Object) (maxVertexIdx, maxNormalIdx int) {
	// Check triangles
	for _, triangleIndex := range entityModel.TriangleIndices {
		if int(triangleIndex.IndexVertex0) > maxVertexIdx { maxVertexIdx = int(triangleIndex.IndexVertex0) }
		if int(triangleIndex.IndexVertex1) > maxVertexIdx { maxVertexIdx = int(triangleIndex.IndexVertex1) }
		if int(triangleIndex.IndexVertex2) > maxVertexIdx { maxVertexIdx = int(triangleIndex.IndexVertex2) }
		if int(triangleIndex.IndexNormal0) > maxNormalIdx { maxNormalIdx = int(triangleIndex.IndexNormal0) }
		if int(triangleIndex.IndexNormal1) > maxNormalIdx { maxNormalIdx = int(triangleIndex.IndexNormal1) }
		if int(triangleIndex.IndexNormal2) > maxNormalIdx { maxNormalIdx = int(triangleIndex.IndexNormal2) }
	}
	
	// Check quads
	for _, quadIndex := range entityModel.QuadIndices {
		if int(quadIndex.IndexVertex0) > maxVertexIdx { maxVertexIdx = int(quadIndex.IndexVertex0) }
		if int(quadIndex.IndexVertex1) > maxVertexIdx { maxVertexIdx = int(quadIndex.IndexVertex1) }
		if int(quadIndex.IndexVertex2) > maxVertexIdx { maxVertexIdx = int(quadIndex.IndexVertex2) }
		if int(quadIndex.IndexVertex3) > maxVertexIdx { maxVertexIdx = int(quadIndex.IndexVertex3) }
		if int(quadIndex.IndexNormal0) > maxNormalIdx { maxNormalIdx = int(quadIndex.IndexNormal0) }
		if int(quadIndex.IndexNormal1) > maxNormalIdx { maxNormalIdx = int(quadIndex.IndexNormal1) }
		if int(quadIndex.IndexNormal2) > maxNormalIdx { maxNormalIdx = int(quadIndex.IndexNormal2) }
		if int(quadIndex.IndexNormal3) > maxNormalIdx { maxNormalIdx = int(quadIndex.IndexNormal3) }
	}
	return maxVertexIdx, maxNormalIdx
}

// Build vertex and normal arrays from component data
func buildVertexArrays(entityModel *fileio.MD1Object, maxVertexIdx, maxNormalIdx int) ([]vec3, []vec3) {
	pos := make([]vec3, maxVertexIdx+1)
	nrms := make([]vec3, maxNormalIdx+1)
	
	// Fill vertex positions from triangle vertices
	for i, vertex := range entityModel.TriangleVertices {
		if i < len(pos) {
			pos[i] = buildModelVertex(vertex)
		}
	}
	
	// Fill vertex positions from quad vertices (may overwrite triangle vertices)
	for i, vertex := range entityModel.QuadVertices {
		if i < len(pos) {
			pos[i] = buildModelVertex(vertex)
		}
	}
	
	// Fill normals from triangle normals
	for i, normal := range entityModel.TriangleNormals {
		if i < len(nrms) {
			nrms[i] = buildModelNormal(normal)
		}
	}
	
	// Fill normals from quad normals (may overwrite triangle normals)
	for i, normal := range entityModel.QuadNormals {
		if i < len(nrms) {
			nrms[i] = buildModelNormal(normal)
		}
	}
	
	return pos, nrms
}

// Apply skeleton transforms to vertices and normals
func applySkeletonTransforms(pos, nrms []vec3, skeleton *fileio.EMROutput, skeletonTransforms []mgl32.Mat4, ci int) {
	if skeleton != nil && ci < len(skeleton.RelativePositionData) {
		transform := skeletonTransforms[ci]
		
		// Transform vertices
		for i := range pos {
			vertex := mgl32.Vec3{pos[i].x, pos[i].y, pos[i].z}
			transformed := transform.Mul4x1(vertex.Vec4(1.0))
			// Apply Y-flip for OBJ coordinate system
			pos[i] = vec3{x: transformed.X(), y: -transformed.Y(), z: transformed.Z()}
		}
		
		// Transform normals (rotation only)
		for i := range nrms {
			normal := mgl32.Vec3{nrms[i].x, nrms[i].y, nrms[i].z}
			rotation := mgl32.Mat3{
				transform[0], transform[1], transform[2],
				transform[4], transform[5], transform[6],
				transform[8], transform[9], transform[10],
			}
			transformed := rotation.Mul3x1(normal)
			transformed = transformed.Normalize()
			// Flip normals to match vertex Y-flip
			nrms[i] = vec3{x: transformed.X(), y: -transformed.Y(), z: transformed.Z()}
		}
	} else if skeleton != nil {
		fmt.Printf("  Warning: Component %d has no skeleton data (skeleton has %d components)\n", ci, len(skeleton.RelativePositionData))
	} else {
		// No skeleton - apply Y-flip for OBJ coordinate system
		for i := range pos {
			pos[i] = vec3{x: pos[i].x, y: -pos[i].y, z: pos[i].z}
		}
		for i := range nrms {
			nrms[i] = vec3{x: nrms[i].x, y: -nrms[i].y, z: nrms[i].z}
		}
	}
}

// Process triangles and add them to the mesh
func processTriangles(m *mesh, entityModel *fileio.MD1Object, tex *fileio.TIMOutput, materials map[MatKey]Material, texturePNG string, addUV func(float32, float32, uint16) int) {
	for j := 0; j < len(entityModel.TriangleIndices); j++ {
		triangleIndex := entityModel.TriangleIndices[j]
		textureInfo := entityModel.TriangleTextures[j]

		// Add UV coordinates
		uv0Idx := addUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page)
		uv1Idx := addUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page)
		uv2Idx := addUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page)

		// Create triangle face
		m.faces = append(m.faces, face{
			a: faceIndex{p: int(triangleIndex.IndexVertex0), t: uv0Idx, n: int(triangleIndex.IndexNormal0)},
			b: faceIndex{p: int(triangleIndex.IndexVertex1), t: uv1Idx, n: int(triangleIndex.IndexNormal1)},
			c: faceIndex{p: int(triangleIndex.IndexVertex2), t: uv2Idx, n: int(triangleIndex.IndexNormal2)},
			d: faceIndex{p: -1, t: -1, n: -1}, // -1 indicates triangle
		})

		// Add material
		mk := MatKey{TPage: int(textureInfo.Page), Palette: int(textureInfo.ClutId)}
		matName := fmt.Sprintf("mat_tp%02d_clut%02d", mk.TPage, mk.Palette)
		m.faceMats = append(m.faceMats, matName)
		if _, ok := materials[mk]; !ok {
			materials[mk] = Material{Name: matName, PNG: texturePNG}
		}
	}
}

// Process quads and add them to the mesh
func processQuads(m *mesh, entityModel *fileio.MD1Object, tex *fileio.TIMOutput, materials map[MatKey]Material, texturePNG string, addUV func(float32, float32, uint16) int) {
	for j := 0; j < len(entityModel.QuadIndices); j++ {
		quadIndex := entityModel.QuadIndices[j]
		textureInfo := entityModel.QuadTextures[j]

		// Add UV coordinates
		uv0Idx := addUV(float32(textureInfo.U0), float32(textureInfo.V0), textureInfo.Page)
		uv1Idx := addUV(float32(textureInfo.U1), float32(textureInfo.V1), textureInfo.Page)
		uv2Idx := addUV(float32(textureInfo.U2), float32(textureInfo.V2), textureInfo.Page)
		uv3Idx := addUV(float32(textureInfo.U3), float32(textureInfo.V3), textureInfo.Page)

		// Create quad face
		m.faces = append(m.faces, face{
			a: faceIndex{p: int(quadIndex.IndexVertex2), t: uv2Idx, n: int(quadIndex.IndexNormal2)}, 
			b: faceIndex{p: int(quadIndex.IndexVertex3), t: uv3Idx, n: int(quadIndex.IndexNormal3)}, 
			c: faceIndex{p: int(quadIndex.IndexVertex1), t: uv1Idx, n: int(quadIndex.IndexNormal1)}, 
			d: faceIndex{p: int(quadIndex.IndexVertex0), t: uv0Idx, n: int(quadIndex.IndexNormal0)},
		})

		// Add material for quad
		mk := MatKey{TPage: int(textureInfo.Page), Palette: int(textureInfo.ClutId)}
		matName := fmt.Sprintf("mat_tp%02d_clut%02d", mk.TPage, mk.Palette)
		m.faceMats = append(m.faceMats, matName)
		if _, ok := materials[mk]; !ok {
			materials[mk] = Material{Name: matName, PNG: texturePNG}
		}
	}
}

// Log export statistics
func logExportStatistics(md1 *fileio.MD1Output, out []mesh, materials map[MatKey]Material, skeleton *fileio.EMROutput) {
	totalTriangles := 0
	totalQuads := 0
	totalVertices := 0
	totalUVs := 0
	totalNormals := 0
	for _, m := range out {
		totalVertices += len(m.positions)
		totalUVs += len(m.uvs)
		totalNormals += len(m.normals)
	}
	for _, c := range md1.Components {
		totalTriangles += len(c.TriangleIndices)
		totalQuads += len(c.QuadIndices)
	}
	totalPolygons := totalTriangles + totalQuads

	if skeleton != nil {
		fmt.Printf("\n=== OBJ Export Statistics (with skeleton) ===\n")
	} else {
		fmt.Printf("\n=== OBJ Export Statistics ===\n")
	}
	fmt.Printf("Objects: %d\n", len(md1.Components))
	fmt.Printf("Triangles: %d\n", totalTriangles)
	fmt.Printf("Quads: %d\n", totalQuads)
	fmt.Printf("Total Faces: %d\n", totalPolygons)
	fmt.Printf("Vertexes: %d\n", totalVertices)
	fmt.Printf("Normals: %d\n", totalNormals)
	fmt.Printf("Materials: %d\n", len(materials))
	fmt.Printf("=============================\n\n")
}

// Convert MD1 model data to mesh format for OBJ export
func buildMeshesFromMD1(md1 *fileio.MD1Output, tex *fileio.TIMOutput, texturePNG string, skeleton *fileio.EMROutput, skeletonTransforms []mgl32.Mat4) ([]mesh, map[MatKey]Material) {
	out := make([]mesh, 0, len(md1.Components))
	materials := make(map[MatKey]Material)

	for ci, entityModel := range md1.Components {
		m := mesh{name: fmt.Sprintf("component_%d", ci)}

		// Find max indices and build vertex arrays
		maxVertexIdx, maxNormalIdx := findMaxIndices(&entityModel)
		pos, nrms := buildVertexArrays(&entityModel, maxVertexIdx, maxNormalIdx)

		// Apply skeleton transforms
		applySkeletonTransforms(pos, nrms, skeleton, skeletonTransforms, ci)

		// Build UV coordinates
		var uvs []vec2
		addUV := func(u, v float32, page uint16) int {
			idx := len(uvs)
			uvs = append(uvs, buildTextureUV(u, v, page, tex))
			return idx
		}

		// Process triangles and quads
		processTriangles(&m, &entityModel, tex, materials, texturePNG, addUV)
		processQuads(&m, &entityModel, tex, materials, texturePNG, addUV)

		m.positions = pos
		m.uvs = uvs
		m.normals = nrms

		// Log component statistics
		fmt.Printf("Component %d: %d triangles, %d quads, %d vertices, %d UVs, %d normals\n", 
			ci, len(entityModel.TriangleIndices), len(entityModel.QuadIndices), 
			len(pos), len(uvs), len(nrms))

		out = append(out, m)
	}

	// Log total statistics
	logExportStatistics(md1, out, materials, skeleton)

	return out, materials
}

// Write MTL material file
func writeMTL(path string, mats map[MatKey]Material) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, m := range mats {
		fmt.Fprintf(f, "newmtl %s\n", m.Name)
		fmt.Fprintln(f, "Ka 1 1 1")
		fmt.Fprintln(f, "Kd 1 1 1")
		fmt.Fprintln(f, "Ks 0 0 0")
		fmt.Fprintln(f, "d 1")
		fmt.Fprintln(f, "illum 2")
		fmt.Fprintf(f, "map_Kd %s\n\n", m.PNG)
	}
	return nil
}

// Write OBJ file with optional material support
func writeOBJ(outPath string, meshes []mesh, mtlBase string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	fmt.Fprintln(f, "# generated by fileconv")
	useMaterials := mtlBase != ""
	if useMaterials {
		fmt.Fprintf(f, "mtllib %s\n", mtlBase)
	}
	fmt.Fprintln(f, "")
	vOfs, vtOfs, vnOfs := 0, 0, 0
	for meshIdx, m := range meshes {
		// Write object name
		fmt.Fprintf(f, "o %s\n", m.name)

		// Write vertices
		for _, p := range m.positions {
			fmt.Fprintf(f, "v %.6f %.6f %.6f\n", p.x, p.y, p.z)
		}
		fmt.Fprintf(f, "# %d vertices\n\n", len(m.positions))

		// Write texture coordinates
		for _, t := range m.uvs {
			fmt.Fprintf(f, "vt %.6f %.6f\n", t.u, t.v)
		}
		fmt.Fprintf(f, "# %d texture coordinates\n\n", len(m.uvs))

		// Write normals
		for _, n := range m.normals {
			fmt.Fprintf(f, "vn %.6f %.6f %.6f\n", n.x, n.y, n.z)
		}
		fmt.Fprintf(f, "# %d vertex normals\n\n", len(m.normals))

		// Object group and smoothing
		fmt.Fprintf(f, "g Object_%d\n", meshIdx)
		fmt.Fprintln(f, "s 1")

		// Face formatter with accumulated offsets
		fstr := func(x faceIndex) string {
			v := x.p + 1 + vOfs  // 1-based indexing
			vt := x.t + 1 + vtOfs
			vn := x.n + 1 + vnOfs
			return fmt.Sprintf("%d/%d/%d", v, vt, vn)
		}

		// Write faces with optional material switching
		currentMat := ""
		for faceIdx, fc := range m.faces {
			// Material switching (only for material format)
			if useMaterials && faceIdx < len(m.faceMats) {
				matName := m.faceMats[faceIdx]
				if matName != currentMat {
					fmt.Fprintf(f, "usemtl %s\n", matName)
					currentMat = matName
				}
			}
			
			// Write face - handle triangles and quads
			if fc.d.p == -1 {
				// Triangle - normal order for both formats
				fmt.Fprintf(f, "f %s %s %s\n", fstr(fc.a), fstr(fc.b), fstr(fc.c))
			} else {
				// Quad - write all 4 vertices
				fmt.Fprintf(f, "f %s %s %s %s\n", fstr(fc.a), fstr(fc.b), fstr(fc.c), fstr(fc.d))
			}
		}
		fmt.Fprintf(f, "# %d faces\n\n", len(m.faces))

		vOfs += len(m.positions)
		vtOfs += len(m.uvs)
		vnOfs += len(m.normals)
	}
	return nil
}

// Write OBJ file with material support
func writeOBJWithMTL(outOBJ, mtlBase string, meshes []mesh) error {
	return writeOBJ(outOBJ, meshes, mtlBase)
}

// Build skeleton transforms recursively
func buildComponentTransformsRecursive(skeleton *fileio.EMROutput, curId int, parentId int, transforms []mgl32.Mat4) {
	// Start with identity matrix
	transformMatrix := mgl32.Ident4()
	
	// Apply parent transform if not root
	if parentId != -1 && parentId < len(transforms) {
		transformMatrix = transforms[parentId]
	}

	// Apply component translation
	if curId < len(skeleton.RelativePositionData) {
		offsetFromParent := skeleton.RelativePositionData[curId]
		translate := mgl32.Translate3D(float32(offsetFromParent.X), float32(offsetFromParent.Y), float32(offsetFromParent.Z))
		transformMatrix = transformMatrix.Mul4(translate)
		
		// Add rotation (identity for now, can be extended for animation)
		quat := mgl32.QuatIdent()
		transformMatrix = transformMatrix.Mul4(quat.Mat4())
	}

	// Store transform for this component
	if curId < len(transforms) {
		transforms[curId] = transformMatrix
	}

	// Process children recursively
	if curId < len(skeleton.ArmatureChildren) {
		for i := 0; i < len(skeleton.ArmatureChildren[curId]); i++ {
			childId := int(skeleton.ArmatureChildren[curId][i])
			buildComponentTransformsRecursive(skeleton, childId, curId, transforms)
		}
	}
}
