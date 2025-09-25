package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// JSON-compatible versions of the data structures for debugging

type PLDOutputJSON struct {
	AnimationData *EDDOutputJSON `json:"animation_data"`
	SkeletonData  *EMROutputJSON `json:"skeleton_data"`
	MeshData      *MD1OutputJSON `json:"mesh_data"`
	TextureData   *TIMOutputJSON `json:"texture_data"`
}

type EDDOutputJSON struct {
	AnimationIndexFrames [][]fileio.EDDTableElement `json:"animation_index_frames"`
	NumFrames            int                        `json:"num_frames"`
}

type EMROutputJSON struct {
	RelativePositionData []fileio.EMRRelativePosition `json:"relative_position_data"`
	ArmatureData         []fileio.EMRArmature         `json:"armature_data"`
	ArmatureChildren     [][]int                      `json:"armature_children"`
	MeshIdList           []int                        `json:"mesh_id_list"`
	FrameData            []AnimationFrameJSON         `json:"frame_data"`
}

type AnimationFrameJSON struct {
	FrameHeader    fileio.EMRFrame `json:"frame_header"`
	RotationAngles []Vec3JSON      `json:"rotation_angles"`
}

type Vec3JSON struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type MD1OutputJSON struct {
	Components []MD1ObjectJSON `json:"components"`
	NumBytes   int64           `json:"num_bytes"`
}

type MD1ObjectJSON struct {
	TriangleVertices []fileio.MD1Vertex           `json:"triangle_vertices"`
	TriangleNormals  []fileio.MD1Vertex           `json:"triangle_normals"`
	TriangleIndices  []fileio.MD1TriangleIndex    `json:"triangle_indices"`
	TriangleTextures []fileio.MD1TriangleTexture  `json:"triangle_textures"`
	QuadVertices     []fileio.MD1Vertex           `json:"quad_vertices"`
	QuadNormals      []fileio.MD1Vertex           `json:"quad_normals"`
	QuadIndices      []fileio.MD1QuadIndex        `json:"quad_indices"`
	QuadTextures     []fileio.MD1QuadTexture      `json:"quad_textures"`
}

type TIMOutputJSON struct {
	PixelData   [][]uint16 `json:"pixel_data"`
	ImageWidth  int        `json:"image_width"`
	ImageHeight int        `json:"image_height"`
	NumPalettes int        `json:"num_palettes"`
	NumBytes    int        `json:"num_bytes"`
}

// EMD JSON structures
type EMDOutputJSON struct {
	AnimationData1 *EDDOutputJSON `json:"animation_data_1"`
	SkeletonData1  *EMROutputJSON `json:"skeleton_data_1"`
	AnimationData2 *EDDOutputJSON `json:"animation_data_2"`
	SkeletonData2  *EMROutputJSON `json:"skeleton_data_2"`
	AnimationData3 *EDDOutputJSON `json:"animation_data_3"`
	SkeletonData3  *EMROutputJSON `json:"skeleton_data_3"`
	MeshData       *MD1OutputJSON `json:"mesh_data"`
}

// Convert PLDOutput to JSON-compatible format
func convertPLDToJSON(pld *fileio.PLDOutput) *PLDOutputJSON {
	jsonOutput := &PLDOutputJSON{}

	// Convert animation data
	if pld.AnimationData != nil {
		jsonOutput.AnimationData = &EDDOutputJSON{
			AnimationIndexFrames: pld.AnimationData.AnimationIndexFrames,
			NumFrames:            pld.AnimationData.NumFrames,
		}
	}

	// Convert skeleton data
	if pld.SkeletonData != nil {
		jsonOutput.SkeletonData = convertEMRToJSON(pld.SkeletonData)
	}

	// Convert mesh data
	if pld.MeshData != nil {
		jsonOutput.MeshData = convertMD1ToJSON(pld.MeshData)
	}

	// Convert texture data
	if pld.TextureData != nil {
		jsonOutput.TextureData = &TIMOutputJSON{
			PixelData:   pld.TextureData.PixelData,
			ImageWidth:  pld.TextureData.ImageWidth,
			ImageHeight: pld.TextureData.ImageHeight,
			NumPalettes: pld.TextureData.NumPalettes,
			NumBytes:    pld.TextureData.NumBytes,
		}
	}

	return jsonOutput
}

// Convert EMDOutput to JSON-compatible format
func convertEMDToJSON(emd *fileio.EMDOutput) *EMDOutputJSON {
	jsonOutput := &EMDOutputJSON{}

	// Convert animation data 1
	if emd.AnimationData1 != nil {
		jsonOutput.AnimationData1 = &EDDOutputJSON{
			AnimationIndexFrames: emd.AnimationData1.AnimationIndexFrames,
			NumFrames:            emd.AnimationData1.NumFrames,
		}
	}

	// Convert skeleton data 1
	if emd.SkeletonData1 != nil {
		jsonOutput.SkeletonData1 = convertEMRToJSON(emd.SkeletonData1)
	}

	// Convert animation data 2
	if emd.AnimationData2 != nil {
		jsonOutput.AnimationData2 = &EDDOutputJSON{
			AnimationIndexFrames: emd.AnimationData2.AnimationIndexFrames,
			NumFrames:            emd.AnimationData2.NumFrames,
		}
	}

	// Convert skeleton data 2
	if emd.SkeletonData2 != nil {
		jsonOutput.SkeletonData2 = convertEMRToJSON(emd.SkeletonData2)
	}

	// Convert animation data 3
	if emd.AnimationData3 != nil {
		jsonOutput.AnimationData3 = &EDDOutputJSON{
			AnimationIndexFrames: emd.AnimationData3.AnimationIndexFrames,
			NumFrames:            emd.AnimationData3.NumFrames,
		}
	}

	// Convert skeleton data 3
	if emd.SkeletonData3 != nil {
		jsonOutput.SkeletonData3 = convertEMRToJSON(emd.SkeletonData3)
	}

	// Convert mesh data
	if emd.MeshData != nil {
		jsonOutput.MeshData = convertMD1ToJSON(emd.MeshData)
	}

	return jsonOutput
}

// Helper function to convert EMR data to JSON
func convertEMRToJSON(emr *fileio.EMROutput) *EMROutputJSON {
	frameDataJSON := make([]AnimationFrameJSON, len(emr.FrameData))
	for i, frame := range emr.FrameData {
		rotationAnglesJSON := make([]Vec3JSON, len(frame.RotationAngles))
		for j, angle := range frame.RotationAngles {
			rotationAnglesJSON[j] = Vec3JSON{
				X: angle.X(),
				Y: angle.Y(),
				Z: angle.Z(),
			}
		}
		frameDataJSON[i] = AnimationFrameJSON{
			FrameHeader:    frame.FrameHeader,
			RotationAngles: rotationAnglesJSON,
		}
	}

	// Convert byte arrays to int arrays for better readability
	armatureChildrenInt := make([][]int, len(emr.ArmatureChildren))
	for i, children := range emr.ArmatureChildren {
		armatureChildrenInt[i] = make([]int, len(children))
		for j, child := range children {
			armatureChildrenInt[i][j] = int(child)
		}
	}

	meshIdListInt := make([]int, len(emr.MeshIdList))
	for i, id := range emr.MeshIdList {
		meshIdListInt[i] = int(id)
	}

	return &EMROutputJSON{
		RelativePositionData: emr.RelativePositionData,
		ArmatureData:         emr.ArmatureData,
		ArmatureChildren:     armatureChildrenInt,
		MeshIdList:           meshIdListInt,
		FrameData:            frameDataJSON,
	}
}

// Helper function to convert MD1 data to JSON
func convertMD1ToJSON(md1 *fileio.MD1Output) *MD1OutputJSON {
	componentsJSON := make([]MD1ObjectJSON, len(md1.Components))
	for i, component := range md1.Components {
		componentsJSON[i] = MD1ObjectJSON{
			TriangleVertices: component.TriangleVertices,
			TriangleNormals:  component.TriangleNormals,
			TriangleIndices:  component.TriangleIndices,
			TriangleTextures: component.TriangleTextures,
			QuadVertices:     component.QuadVertices,
			QuadNormals:      component.QuadNormals,
			QuadIndices:      component.QuadIndices,
			QuadTextures:     component.QuadTextures,
		}
	}

	return &MD1OutputJSON{
		Components: componentsJSON,
		NumBytes:   md1.NumBytes,
	}
}

func main() {
	var inputFile string
	var outputFile string
	var prettyPrint bool
	var format string

	flag.StringVar(&inputFile, "input", "", "Input model file path")
	flag.StringVar(&outputFile, "output", "", "Output JSON file path (optional, defaults to stdout)")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty print JSON output")
	flag.StringVar(&format, "format", "", "File format (pld, emd) - auto-detected from extension if not specified")
	flag.Parse()

	if inputFile == "" {
		fmt.Println("Usage: modeldumper -input <model_file> [-output <json_file>] [-pretty=true] [-format=<format>]")
		fmt.Println("Supported formats: pld, emd")
		fmt.Println("Example: modeldumper -input data/PL0/PLD/LEON.PLD -output leon_data.json")
		fmt.Println("Example: modeldumper -input data/PL0/EMD/LEON.EMD -output leon_data.json")
		os.Exit(1)
	}

	// Auto-detect format from file extension if not specified
	if format == "" {
		ext := strings.ToLower(filepath.Ext(inputFile))
		switch ext {
		case ".pld":
			format = "pld"
		case ".emd":
			format = "emd"
		default:
			log.Fatalf("Unknown file format: %s. Please specify format with -format flag", ext)
		}
	}

	// Load model file based on format
	fmt.Printf("Loading %s file: %s\n", strings.ToUpper(format), inputFile)
	
	var jsonData []byte
	var err error

	switch format {
	case "pld":
		pldData, err := fileio.LoadPLDFile(inputFile)
		if err != nil {
			log.Fatalf("Failed to load PLD file: %v", err)
		}
		
		jsonOutput := convertPLDToJSON(pldData)
		if prettyPrint {
			jsonData, err = json.MarshalIndent(jsonOutput, "", "  ")
		} else {
			jsonData, err = json.Marshal(jsonOutput)
		}
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		fmt.Printf("Successfully processed PLD file with %d components\n", len(pldData.MeshData.Components))

	case "emd":
		emdData := fileio.LoadEMDFile(inputFile)
		if emdData == nil {
			log.Fatalf("Failed to load EMD file")
		}
		
		jsonOutput := convertEMDToJSON(emdData)
		if prettyPrint {
			jsonData, err = json.MarshalIndent(jsonOutput, "", "  ")
		} else {
			jsonData, err = json.Marshal(jsonOutput)
		}
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		componentCount := 0
		if emdData.MeshData != nil {
			componentCount = len(emdData.MeshData.Components)
		}
		fmt.Printf("Successfully processed EMD file with %d components\n", componentCount)

	default:
		log.Fatalf("Unsupported format: %s", format)
	}

	// Output JSON
	if outputFile == "" {
		// Output to stdout
		fmt.Println(string(jsonData))
	} else {
		// Output to file
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Model data dumped to: %s\n", outputFile)
	}
}
