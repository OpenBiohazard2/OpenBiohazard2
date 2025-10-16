package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}

	toolName := os.Args[1]
	inputFilename := os.Args[2]
	outputFilename := os.Args[3]

	// Parse optional flags (skeleton is default for better user experience)
	useSkeleton := true
	if len(os.Args) > 4 {
		for _, arg := range os.Args[4:] {
			switch arg {
			case "--skeleton", "-s":
				useSkeleton = true
			case "--raw", "-r":
				useSkeleton = false
			}
		}
	}

	// Check if input file exists
	if _, err := os.Stat(inputFilename); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist\n", inputFilename)
		os.Exit(1)
	}

	fmt.Printf("Converting %s to %s using %s...\n", inputFilename, outputFilename, toolName)
	if useSkeleton {
		fmt.Println("Using skeleton data for full character model")
	} else {
		fmt.Println("Using raw MD1 data without skeleton")
	}

	switch toolName {
	case "tim2png":
		convertTIMToPNG(inputFilename, outputFilename)
	case "adt2png":
		convertADTToPNG(inputFilename, outputFilename)
	case "sap2wav":
		convertSAPToWAV(inputFilename, outputFilename)
	case "pld2obj":
		convertPLDToOBJ(inputFilename, outputFilename, useSkeleton)
	case "emd2obj":
		convertEMDToOBJ(inputFilename, outputFilename, useSkeleton)
	default:
		fmt.Printf("Error: Invalid tool name '%s'\n", toolName)
		fmt.Println("Supported tools: tim2png, adt2png, sap2wav, pld2obj, emd2obj")
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("File Converter - Convert Biohazard 2 files to standard formats")
	fmt.Println("")
	fmt.Println("Usage: fileconv [toolName] [inputFilename] [outputFilename] [flags]")
	fmt.Println("")
	fmt.Println("Supported tools:")
	fmt.Println("  tim2png  - Convert TIM texture to PNG")
	fmt.Println("  adt2png  - Convert ADT image to PNG") 
	fmt.Println("  sap2wav  - Convert SAP audio to WAV")
	fmt.Println("  pld2obj  - Convert PLD mesh to OBJ")
	fmt.Println("  emd2obj  - Convert EMD mesh to OBJ")
	fmt.Println("")
	fmt.Println("OBJ Export Flags:")
	fmt.Println("  --skeleton, -s  - Use skeleton data for full character model (default)")
	fmt.Println("  --raw, -r       - Export raw MD1 data without skeleton")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  fileconv tim2png data/Pl0/Emd0/EM000.TIM em000.png")
	fmt.Println("  fileconv adt2png data/Pl0/Emd0/EM000.ADT em000.png")
	fmt.Println("  fileconv sap2wav data/Pl0/Voice/STAGE0/0000.SAP voice.wav")
	fmt.Println("  fileconv pld2obj data/PL0/PLD/PL00.PLD leon.obj")
	fmt.Println("  fileconv emd2obj data/PL0/EMD0/EM000.EMD enemy.obj --raw")
	fmt.Println("")
	fmt.Printf("Error: You provided %d arguments, but 4 are required\n", len(os.Args))
}

func convertTIMToPNG(inputFilename, outputFilename string) {
	fmt.Println("Loading TIM file...")
	timOutput, err := fileio.LoadTIMFile(inputFilename)
	if err != nil {
		fmt.Printf("Error loading TIM file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("TIM file loaded: %dx%d pixels, %d palettes\n", 
		timOutput.ImageWidth, timOutput.ImageHeight, timOutput.NumPalettes)
	
	fmt.Println("Converting to PNG...")
	if err := timOutput.ConvertToPNG(outputFilename); err != nil {
		fmt.Printf("Error converting to PNG: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully converted to %s\n", outputFilename)
}

func convertADTToPNG(inputFilename, outputFilename string) {
	fmt.Println("Loading ADT file...")
	adtOutput, err := fileio.LoadADTFile(inputFilename)
	if err != nil {
		fmt.Printf("Error loading ADT file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ADT file loaded: %d pixel data entries\n", len(adtOutput.PixelData))
	
	fmt.Println("Converting to PNG...")
	adtOutput.ConvertToPNG(outputFilename)
	fmt.Printf("Successfully converted to %s\n", outputFilename)
}

func convertSAPToWAV(inputFilename, outputFilename string) {
	fmt.Println("Loading SAP file...")
	sapOutput, err := fileio.LoadSAPFile(inputFilename)
	if err != nil {
		fmt.Printf("Error loading SAP file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SAP file loaded: %d bytes of audio data\n", len(sapOutput.AudioData))
	
	fmt.Println("Converting to WAV...")
	sapOutput.ConvertToWAV(outputFilename)
	fmt.Printf("Successfully converted to %s\n", outputFilename)
}

func convertPLDToOBJ(inputFilename, outputFilename string, useSkeleton bool) {
	fmt.Println("Loading PLD file...")
	pld, err := fileio.LoadPLDFile(inputFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if pld.MeshData == nil {
		fmt.Println("Error: no mesh data")
		os.Exit(1)
	}

	fmt.Println("Converting to OBJ...")
	
	// Export single texture file
	texturePNG := outputFilename[:len(outputFilename)-4] + ".png"
	if pld.TextureData != nil {
		fmt.Println("Exporting texture...")
		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(texturePNG), 0755); err != nil {
			fmt.Printf("Warning: Failed to create output directory: %v\n", err)
		} else if err := pld.TextureData.ConvertToPNG(texturePNG); err != nil {
			fmt.Printf("Warning: Failed to export texture %s: %v\n", texturePNG, err)
		}
	}
	
	// Build meshes with single texture reference
	textureBase := filepath.Base(texturePNG)
	var meshes []mesh
	var materials map[MatKey]Material
	
	if useSkeleton && pld.SkeletonData != nil {
		fmt.Printf("Using skeleton data for full character model (%d skeleton components)...\n", len(pld.SkeletonData.RelativePositionData))
		// Precompute all skeleton transforms
		skeletonTransforms := make([]mgl32.Mat4, len(pld.SkeletonData.RelativePositionData))
		buildComponentTransformsRecursive(pld.SkeletonData, 0, -1, skeletonTransforms)
		meshes, materials = buildMeshesFromMD1(pld.MeshData, pld.TextureData, textureBase, pld.SkeletonData, skeletonTransforms)
	} else if useSkeleton {
		fmt.Println("Warning: Skeleton requested but no skeleton data found, using raw MD1 data...")
		meshes, materials = buildMeshesFromMD1(pld.MeshData, pld.TextureData, textureBase, nil, nil)
	} else {
		fmt.Println("Using raw MD1 data without skeleton...")
		meshes, materials = buildMeshesFromMD1(pld.MeshData, pld.TextureData, textureBase, nil, nil)
	}
	
	// Write MTL file
	mtlPath := outputFilename[:len(outputFilename)-4] + ".mtl"
	if err := writeMTL(mtlPath, materials); err != nil {
		fmt.Printf("Error writing MTL: %v\n", err)
		os.Exit(1)
	}
	
	// Write OBJ file
	mtlBase := filepath.Base(mtlPath)
	if err := writeOBJWithMTL(outputFilename, mtlBase, meshes); err != nil {
		fmt.Printf("Error writing OBJ: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully converted to %s with %d materials\n", outputFilename, len(materials))
}

func convertEMDToOBJ(inputFilename, outputFilename string, useSkeleton bool) {
	fmt.Println("Loading EMD file...")
	emd := fileio.LoadEMDFile(inputFilename)
	if emd == nil {
		fmt.Println("Error: failed to load EMD")
		os.Exit(1)
	}
	if emd.MeshData == nil {
		fmt.Println("Error: no mesh data")
		os.Exit(1)
	}

	fmt.Println("Converting to OBJ...")
	
	// Export single texture file
	texturePNG := outputFilename[:len(outputFilename)-4] + ".png"
	
	// For EMD files, we need to load the corresponding TIM file
	// Assume TIM file is in the same directory with same base name
	timPath := inputFilename[:len(inputFilename)-4] + ".TIM"
	if _, err := os.Stat(timPath); err == nil {
		fmt.Println("Loading TIM texture...")
		timOutput, err := fileio.LoadTIMFile(timPath)
		if err == nil {
			fmt.Println("Exporting texture...")
			// Ensure output directory exists
			if err := os.MkdirAll(filepath.Dir(texturePNG), 0755); err != nil {
				fmt.Printf("Warning: Failed to create output directory: %v\n", err)
			} else if err := timOutput.ConvertToPNG(texturePNG); err != nil {
				fmt.Printf("Warning: Failed to export texture %s: %v\n", texturePNG, err)
			}
		} else {
			fmt.Printf("Warning: Failed to load TIM file %s: %v\n", timPath, err)
		}
	} else {
		fmt.Printf("Warning: TIM file not found: %s\n", timPath)
	}
	
	// Build meshes with single texture reference
	textureBase := filepath.Base(texturePNG)
	var timData *fileio.TIMOutput
	if _, err := os.Stat(timPath); err == nil {
		timData, _ = fileio.LoadTIMFile(timPath)
	}
	
	var meshes []mesh
	var materials map[MatKey]Material
	
	if useSkeleton && emd.SkeletonData1 != nil {
		fmt.Printf("Using first skeleton data for full character model (%d skeleton components)...\n", len(emd.SkeletonData1.RelativePositionData))
		// Precompute all skeleton transforms
		skeletonTransforms := make([]mgl32.Mat4, len(emd.SkeletonData1.RelativePositionData))
		buildComponentTransformsRecursive(emd.SkeletonData1, 0, -1, skeletonTransforms)
		meshes, materials = buildMeshesFromMD1(emd.MeshData, timData, textureBase, emd.SkeletonData1, skeletonTransforms)
	} else if useSkeleton {
		fmt.Println("Warning: Skeleton requested but no skeleton data found, using raw MD1 data...")
		meshes, materials = buildMeshesFromMD1(emd.MeshData, timData, textureBase, nil, nil)
	} else {
		fmt.Println("Using raw MD1 data without skeleton...")
		meshes, materials = buildMeshesFromMD1(emd.MeshData, timData, textureBase, nil, nil)
	}
	
	// Write MTL file
	mtlPath := outputFilename[:len(outputFilename)-4] + ".mtl"
	if err := writeMTL(mtlPath, materials); err != nil {
		fmt.Printf("Error writing MTL: %v\n", err)
		os.Exit(1)
	}
	
	// Write OBJ file
	mtlBase := filepath.Base(mtlPath)
	if err := writeOBJWithMTL(outputFilename, mtlBase, meshes); err != nil {
		fmt.Printf("Error writing OBJ: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully converted to %s with %d materials\n", outputFilename, len(materials))
}
