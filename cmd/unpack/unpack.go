package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: unpack <fileFormat> <inputFilename>")
		fmt.Println("")
		fmt.Println("Supported file formats:")
		fmt.Println("  do2    - Door files")
		fmt.Println("  pld    - Player model files")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  unpack do2 door00.do2")
		fmt.Println("  unpack pld leon.pld")
		fmt.Println("  unpack do2 data/Pl0/Door/door01.do2")
		fmt.Println("  unpack pld data/Pl0/Pld/leon.pld")
		fmt.Println("")
		os.Exit(1)
	}

	fileFormat := os.Args[1]
	inputFilename := os.Args[2]

	// Validate input file exists
	if _, err := os.Stat(inputFilename); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist\n", inputFilename)
		os.Exit(1)
	}

	inputBase := filenameWithoutExtension(filepath.Base(inputFilename))
	outputFolder := filepath.Join(".", inputBase)
	
	fmt.Printf("Unpacking: %s\n", inputFilename)
	fmt.Printf("Output to: %s\n", outputFolder)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err := os.Mkdir(outputFolder, 0777); err != nil {
			fmt.Printf("Error: Failed to create output directory '%s': %v\n", outputFolder, err)
			os.Exit(1)
		}
		fmt.Printf("Created output directory: %s\n", outputFolder)
	}

	switch fileFormat {
	case "do2":
		fmt.Println("Processing DO2 file...")
		fmt.Println("Loading DO2 file structure...")
		do2Output := fileio.LoadDO2File(inputFilename)

		file, err := os.Open(inputFilename)
		if err != nil {
			fmt.Printf("Error: Failed to read file '%s': %v\n", inputFilename, err)
			os.Exit(1)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error: Failed to read file data: %v\n", err)
			os.Exit(1)
		}

		baseOutputFilename := filepath.Join(outputFolder, inputBase)
		do2FileFormat := do2Output.DO2FileFormat
		
		fmt.Printf("Found DO2 structure - VH: %d+%d, VB: %d+%d, MD1: %d+%d, TIM: %d+%d\n",
			do2FileFormat.VHOffset, do2FileFormat.VHLength,
			do2FileFormat.VBOffset, do2FileFormat.VBLength,
			do2FileFormat.MD1Offset, do2FileFormat.MD1Length,
			do2FileFormat.TIMOffset, do2FileFormat.TIMLength)
		
		fmt.Println("Extracting components:")
		writeFile(baseOutputFilename+".vh", getBufferSubset(data, do2FileFormat.VHOffset, do2FileFormat.VHLength), "VH")
		writeFile(baseOutputFilename+".vb", getBufferSubset(data, do2FileFormat.VBOffset, do2FileFormat.VBLength), "VB")
		writeFile(baseOutputFilename+".md1", getBufferSubset(data, do2FileFormat.MD1Offset, do2FileFormat.MD1Length), "Model Data")
		writeFile(baseOutputFilename+".tim", getBufferSubset(data, do2FileFormat.TIMOffset, do2FileFormat.TIMLength), "Texture Image")
		
		fmt.Printf("\nSuccessfully unpacked %s to %s\n", inputFilename, outputFolder)
		
	case "pld":
		fmt.Println("Processing PLD file...")
		
		file, err := os.Open(inputFilename)
		if err != nil {
			fmt.Printf("Error: Failed to read file '%s': %v\n", inputFilename, err)
			os.Exit(1)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error: Failed to read file data: %v\n", err)
			os.Exit(1)
		}

		// Parse PLD header to get offsets
		fmt.Println("Parsing PLD header...")
		pldOffsets, err := extractPLDOffsets(data)
		if err != nil {
			fmt.Printf("Error: Failed to parse PLD offsets: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Found offsets - Animation: %d, Skeleton: %d, Mesh: %d, Texture: %d\n",
			pldOffsets.OffsetAnimation, pldOffsets.OffsetSkeleton, pldOffsets.OffsetMesh, pldOffsets.OffsetTexture)

		baseOutputFilename := filepath.Join(outputFolder, inputBase)
		
		fmt.Println("Extracting components:")
		
		// Extract animation data (.edd) - we need to determine length
		if pldOffsets.OffsetAnimation > 0 {
			// For now, extract a reasonable chunk - this could be improved
			animationLength := int64(1024 * 1024) // 1MB default
			if pldOffsets.OffsetSkeleton > pldOffsets.OffsetAnimation {
				animationLength = int64(pldOffsets.OffsetSkeleton - pldOffsets.OffsetAnimation)
			}
			fmt.Printf("  Extracting animation data from offset %d, length %d\n", pldOffsets.OffsetAnimation, animationLength)
			writeFile(baseOutputFilename+".edd", getBufferSubset(data, int64(pldOffsets.OffsetAnimation), animationLength), "Animation Data")
		} else {
			fmt.Println("  Skipping animation data (no offset found)")
		}
		
		// Extract skeleton data (.emr)
		if pldOffsets.OffsetSkeleton > 0 {
			skeletonLength := int64(1024 * 1024) // 1MB default
			if pldOffsets.OffsetMesh > pldOffsets.OffsetSkeleton {
				skeletonLength = int64(pldOffsets.OffsetMesh - pldOffsets.OffsetSkeleton)
			}
			fmt.Printf("  Extracting skeleton data from offset %d, length %d\n", pldOffsets.OffsetSkeleton, skeletonLength)
			writeFile(baseOutputFilename+".emr", getBufferSubset(data, int64(pldOffsets.OffsetSkeleton), skeletonLength), "Skeleton Data")
		} else {
			fmt.Println("  Skipping skeleton data (no offset found)")
		}
		
		// Extract mesh data (.md1)
		if pldOffsets.OffsetMesh > 0 {
			meshLength := int64(1024 * 1024) // 1MB default
			if pldOffsets.OffsetTexture > pldOffsets.OffsetMesh {
				meshLength = int64(pldOffsets.OffsetTexture - pldOffsets.OffsetMesh)
			}
			fmt.Printf("  Extracting mesh data from offset %d, length %d\n", pldOffsets.OffsetMesh, meshLength)
			writeFile(baseOutputFilename+".md1", getBufferSubset(data, int64(pldOffsets.OffsetMesh), meshLength), "Model Data")
		} else {
			fmt.Println("  Skipping mesh data (no offset found)")
		}
		
		// Extract texture data (.tim) - extract to end of file
		if pldOffsets.OffsetTexture > 0 {
			textureLength := int64(len(data)) - int64(pldOffsets.OffsetTexture)
			fmt.Printf("  Extracting texture data from offset %d, length %d\n", pldOffsets.OffsetTexture, textureLength)
			writeFile(baseOutputFilename+".tim", getBufferSubset(data, int64(pldOffsets.OffsetTexture), textureLength), "Texture Image")
		} else {
			fmt.Println("  Skipping texture data (no offset found)")
		}
		
		fmt.Printf("\nSuccessfully unpacked %s to %s\n", inputFilename, outputFolder)
		
	default:
		fmt.Printf("Error: Unsupported file format '%s'\n", fileFormat)
		fmt.Println("Supported formats: do2, pld")
		os.Exit(1)
	}
}

func getBufferSubset(data []byte, offset int64, length int64) []byte {
	return data[offset : offset+length]
}

func writeFile(outputFilename string, buffer []byte, componentName string) {
	err := os.WriteFile(outputFilename, buffer, 0644)
	if err != nil {
		fmt.Printf("Error: Failed to write %s to '%s': %v\n", componentName, outputFilename, err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ %s (%d bytes) -> %s\n", componentName, len(buffer), filepath.Base(outputFilename))
}

func filenameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, path.Ext(filename))
}

// PLD structures for parsing offsets
type PLDHeader struct {
	DirOffset uint32
	DirCount  uint32
}

type PLDOffsets struct {
	OffsetAnimation uint32
	OffsetSkeleton  uint32
	OffsetMesh      uint32
	OffsetTexture   uint32
}

func extractPLDOffsets(data []byte) (*PLDOffsets, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("file too small to contain PLD header")
	}
	
	// Read PLD header
	dirOffset := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	dirCount := uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16 | uint32(data[7])<<24
	
	// Log header info for debugging
	fmt.Printf("  PLD Header - Directory offset: %d, Count: %d\n", dirOffset, dirCount)
	
	if int64(dirOffset) >= int64(len(data)) {
		return nil, fmt.Errorf("directory offset %d exceeds file size %d", dirOffset, len(data))
	}
	
	// Read offsets from directory
	offsetStart := int(dirOffset)
	if offsetStart+16 > len(data) {
		return nil, fmt.Errorf("directory section too small")
	}
	
	offsets := &PLDOffsets{}
	offsets.OffsetAnimation = uint32(data[offsetStart]) | uint32(data[offsetStart+1])<<8 | uint32(data[offsetStart+2])<<16 | uint32(data[offsetStart+3])<<24
	offsets.OffsetSkeleton = uint32(data[offsetStart+4]) | uint32(data[offsetStart+5])<<8 | uint32(data[offsetStart+6])<<16 | uint32(data[offsetStart+7])<<24
	offsets.OffsetMesh = uint32(data[offsetStart+8]) | uint32(data[offsetStart+9])<<8 | uint32(data[offsetStart+10])<<16 | uint32(data[offsetStart+11])<<24
	offsets.OffsetTexture = uint32(data[offsetStart+12]) | uint32(data[offsetStart+13])<<8 | uint32(data[offsetStart+14])<<16 | uint32(data[offsetStart+15])<<24
	
	return offsets, nil
}
