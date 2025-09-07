package fileio

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strconv"
	"unsafe"
)

const (
	NODE_LEFT          = 0
	NODE_RIGHT         = 1
	TOTAL_IMAGE_WIDTH  = 320
	TOTAL_IMAGE_HEIGHT = 240
)

type UnpackArray8 struct {
	start  int64
	length int64
}

type Node struct {
	ChildNodes [2]int64
}

type UnpackArray struct {
	DataCount uint64
	Ptr8      []UnpackArray8
	Tree      []Node
}

type ADTOutput struct {
	PixelData [][]uint16
	RawData   []uint8
}

func LoadADTFile(inputFilename string) (*ADTOutput, error) {
	imgFile, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open ADT file %s: %w", inputFilename, err)
	}
	defer imgFile.Close()

	return LoadADTStream(imgFile)
}

func LoadADTStream(adtReader io.ReaderAt) (*ADTOutput, error) {
	imgArr, rawData, err := unpackADT(adtReader)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ADT: %w", err)
	}

	if len(imgArr) < 320*256 {
		fmt.Println("Warning: the ADT file doesn't contain a 320x240 image")
		return &ADTOutput{
			RawData: rawData,
		}, nil
	}
	pixelData := restoreImage(imgArr)
	return &ADTOutput{
		PixelData: pixelData,
		RawData:   rawData,
	}, nil
}

func newUnpackArray(arrayLength int) UnpackArray {
	array := UnpackArray{
		DataCount: uint64(arrayLength),
		Ptr8:      make([]UnpackArray8, arrayLength),
		Tree:      make([]Node, arrayLength*2),
	}
	for i := 0; i < len(array.Ptr8); i++ {
		array.Ptr8[i].start = -1
		array.Ptr8[i].length = -1
	}
	for i := 0; i < len(array.Tree); i++ {
		array.Tree[i].ChildNodes = [2]int64{-1, -1}
	}
	return array
}

func readBitFieldArray(bitReader *BitReader, array *UnpackArray, curIndex int) int {
	// Descend down the tree
	for {
		if bitReader.UnsafeReadBit() == 1 {
			curIndex = int(array.Tree[curIndex].ChildNodes[NODE_RIGHT])
		} else {
			curIndex = int(array.Tree[curIndex].ChildNodes[NODE_LEFT])
		}
		if curIndex < int(array.DataCount) {
			break
		}
	}

	return curIndex
}

func readBinaryNumber(bitReader *BitReader) int {
	// Read a list of zero bits terminated by a one bit
	numZeroBits := int(0)
	for bitReader.UnsafeReadBit() == 0 {
		numZeroBits++
	}

	// Read in a binary number with 'numZeroBits'
	// Convert to decimal
	binaryNumber := int(1)
	for i := 0; i < numZeroBits; i++ {
		binaryNumber = bitReader.UnsafeReadBit() + (binaryNumber << 1)
	}

	return binaryNumber
}

// Determine the start of each byte array
func initArrayStart(array *UnpackArray) {
	// Frequency of byte array lengths
	freqArray := [17]uint16{}
	for i := 0; i < int(array.DataCount); i++ {
		numValues := array.Ptr8[i].length
		if numValues <= 16 {
			freqArray[numValues]++
		}
	}

	var tmp [18]uint16
	for i := 0; i < 16; i++ {
		tmp[i+2] = (tmp[i+1] + freqArray[i+1]) << 1
	}

	for i := 0; i < 18; i++ {
		for j := 0; j < int(array.DataCount); j++ {
			if int(array.Ptr8[j].length) == i {
				array.Ptr8[j].start = int64(tmp[i] & 0xFFFF)
				tmp[i]++
			}
		}
	}
}

// Build binary tree
func initArrayTree(array *UnpackArray) {
	curLength := array.DataCount
	curArrayIndex := curLength + 1

	array.Tree[curLength].ChildNodes[NODE_LEFT] = -1
	array.Tree[curLength].ChildNodes[NODE_RIGHT] = -1
	array.Tree[curArrayIndex].ChildNodes[NODE_LEFT] = -1
	array.Tree[curArrayIndex].ChildNodes[NODE_RIGHT] = -1

	for i := 0; i < int(array.DataCount); i++ {
		curPtr8Start := array.Ptr8[i].start
		curPtr8Length := array.Ptr8[i].length

		curLength := array.DataCount

		for j := 0; j < int(curPtr8Length); j++ {
			curMask := 1 << uint(int(curPtr8Length)-j-1)
			var arrayOffset int

			if (curMask & int(curPtr8Start)) != 0 {
				arrayOffset = NODE_RIGHT
			} else {
				arrayOffset = NODE_LEFT
			}

			if j+1 == int(curPtr8Length) {
				array.Tree[curLength].ChildNodes[arrayOffset] = int64(i)
				break
			}

			// node at 'curLength' has an empty child
			if array.Tree[curLength].ChildNodes[arrayOffset] == -1 {
				// node at 'curLength' is parent of node at 'curArrayIndex'
				array.Tree[curLength].ChildNodes[arrayOffset] = int64(curArrayIndex)
				array.Tree[curArrayIndex].ChildNodes[NODE_RIGHT] = -1
				array.Tree[curArrayIndex].ChildNodes[NODE_LEFT] = -1
				curLength = curArrayIndex
				curArrayIndex++
			} else {
				// traverse down to its child node
				curLength = uint64(array.Tree[curLength].ChildNodes[arrayOffset])
			}
		}
	}
}

func initUnpackBlock(bitReader *BitReader) (UnpackArray, UnpackArray, error) {
	// Array1
	array1 := newUnpackArray(16)

	prevValue := 0
	for i := 0; i < 16; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			return UnpackArray{}, UnpackArray{}, err
		}

		if bit == 1 {
			prevValue ^= readBinaryNumber(bitReader)
		}
		array1.Ptr8[i].length = int64(prevValue)
	}

	initArrayStart(&array1)
	initArrayTree(&array1)

	// Array 2
	array2 := newUnpackArray(512)
	array2Tmp := make([]int, 512)

	curBit, err := bitReader.ReadBit()
	if err != nil {
		return UnpackArray{}, UnpackArray{}, err
	}

	j := 0
	for j < int(array2.DataCount) {
		curBitField := readBinaryNumber(bitReader)
		if curBit == 1 {
			for i := 0; i < curBitField; i++ {
				array2Tmp[j+i] = readBitFieldArray(bitReader, &array1, int(array1.DataCount))
			}
			j += curBitField
			curBit = 0
		} else if curBit == 0 {
			for i := 0; i < curBitField; i++ {
				array2Tmp[j+i] = 0
			}
			j += curBitField
			curBit = 1
		}
	}

	j = 0
	for i := 0; i < int(array2.DataCount); i++ {
		j = j ^ array2Tmp[i]
		array2.Ptr8[i].length = int64(j)
	}

	initArrayStart(&array2)
	initArrayTree(&array2)

	// Array 3
	array3 := newUnpackArray(16)

	prevValue = 0
	for i := 0; i < 16; i++ {
		bit, err := bitReader.ReadBit()
		if err != nil {
			return UnpackArray{}, UnpackArray{}, err
		}

		if bit == 1 {
			prevValue ^= readBinaryNumber(bitReader)
		}
		array3.Ptr8[i].length = int64(prevValue)
	}

	initArrayStart(&array3)
	initArrayTree(&array3)

	return array2, array3, nil
}

// Load image file
// Return an array of 16bit colors
func unpackADT(r io.ReaderAt) ([]uint16, []uint8, error) {
	maxFileSize := 320 * 256 * 2
	// Skip the first uint32
	reader := io.NewSectionReader(r, int64(unsafe.Sizeof(uint32(0))), int64(unsafe.Sizeof(uint16(0)))*int64(maxFileSize))
	tmp16kOffset := 0
	imageByteData := make([]uint8, 0)
	tmp16k := make([]uint8, 16384)

	bitReader := NewBitReader(reader)

	for {
		blockLen := binary.LittleEndian.Uint16([]byte{bitReader.UnsafeReadByte(), bitReader.UnsafeReadByte()})

		if blockLen == 0 {
			break
		}

		array2, array3, err := initUnpackBlock(bitReader)
		if err != nil {
			return []uint16{}, []uint8{}, err
		}

		for i := 0; i < int(blockLen); i++ {
			curBitField := readBitFieldArray(bitReader, &array2, int(array2.DataCount))

			// Check if the bit field can fit within a byte
			if curBitField < 256 {
				imageByteData = append(imageByteData, uint8(curBitField))
				tmp16k[tmp16kOffset] = uint8(curBitField)
				tmp16kOffset = (tmp16kOffset + 1) % len(tmp16k)
			} else {
				numValues := curBitField - 0xfd
				curBitField = readBitFieldArray(bitReader, &array3, int(array3.DataCount))
				if curBitField != 0 {
					numBits := curBitField - 1
					curBitField = int(bitReader.UnsafeReadNumBits(numBits) & 0xffff)
					curBitField += 1 << uint(numBits)
				}

				// copy from start offset
				startOffset := (tmp16kOffset - curBitField - 1) & 0x3fff
				for j := 0; j < numValues; j++ {
					tmp16k[tmp16kOffset] = tmp16k[startOffset]
					imageByteData = append(imageByteData, tmp16k[tmp16kOffset])
					startOffset = (startOffset + 1) % len(tmp16k)
					tmp16kOffset = (tmp16kOffset + 1) % len(tmp16k)
				}
			}
		}
	}

	// Convert 8 bit array to 16 bit array, since colors are 16 bit
	image16BitData := make([]uint16, 1+len(imageByteData)/2)
	for i := 0; i < len(imageByteData); i += 2 {
		// Combine 2 bytes to get a 16 bit number
		image16BitData[i/2] = binary.LittleEndian.Uint16(imageByteData[i : i+2])
	}
	return image16BitData, imageByteData, nil
}

func restoreImage(colorArr []uint16) [][]uint16 {
	// Create new image and save each pixel
	wholeImageData := make([][]uint16, TOTAL_IMAGE_HEIGHT)
	for y := 0; y < len(wholeImageData); y++ {
		wholeImageData[y] = make([]uint16, TOTAL_IMAGE_WIDTH)
	}

	// The first part is a 256x240 image on the left side
	for y := 0; y < TOTAL_IMAGE_HEIGHT; y++ {
		for x := 0; x < 256; x++ {
			arrayPosition := (256 * y) + x
			wholeImageData[y][x] = colorArr[arrayPosition]
		}
	}

	// The second part is a 64x128 image on the top right
	offsetY := 256
	for y := 0; y < 128; y += 2 {
		for offsetX := 0; offsetX < 64; offsetX++ {
			arrayPosition := offsetX + (256 * offsetY)
			wholeImageData[y][256+offsetX] = colorArr[arrayPosition]
		}

		for offsetX := 0; offsetX < 64; offsetX++ {
			arrayPosition := (128 + offsetX) + (256 * offsetY)
			wholeImageData[y+1][256+offsetX] = colorArr[arrayPosition]
		}

		offsetY++
	}

	// The third part is a 64x112 image on the bottom right
	offsetY = 256
	for y := 128; y < TOTAL_IMAGE_HEIGHT; y += 2 {
		for offsetX := 0; offsetX < 64; offsetX++ {
			arrayPosition := (64 + offsetX) + (256 * offsetY)
			wholeImageData[y][256+offsetX] = colorArr[arrayPosition]
		}

		for offsetX := 0; offsetX < 64; offsetX++ {
			arrayPosition := (192 + offsetX) + (256 * offsetY)
			wholeImageData[y+1][256+offsetX] = colorArr[arrayPosition]
		}

		offsetY++
	}

	return wholeImageData
}

func (adtOutput *ADTOutput) ConvertToRenderData() []uint16 {
	pixelData2D := adtOutput.PixelData
	pixelData1D := make([]uint16, len(pixelData2D)*len(pixelData2D[0]))

	for y := 0; y < len(pixelData2D); y++ {
		for x := 0; x < len(pixelData2D[y]); x++ {
			index := (y * len(pixelData2D[y])) + x
			pixelData1D[index] = pixelData2D[y][x]
		}
	}
	return pixelData1D
}

func (adtOutput *ADTOutput) ConvertToPNG(outputFilename string) {
	pixelData := adtOutput.PixelData

	imageOutputData := image.NewRGBA(image.Rect(0, 0, TOTAL_IMAGE_WIDTH, TOTAL_IMAGE_HEIGHT))
	for y := 0; y < len(pixelData); y++ {
		for x := 0; x < len(pixelData[y]); x++ {
			colorBits := fmt.Sprintf("%016b", pixelData[y][x])
			// color is in A1B5G5R5 format
			a, _ := strconv.ParseInt(string(colorBits[0]), 2, 1)
			a = 255
			b, _ := strconv.ParseInt(string(colorBits[1:6]), 2, 5)
			g, _ := strconv.ParseInt(string(colorBits[6:11]), 2, 5)
			r, _ := strconv.ParseInt(string(colorBits[11:16]), 2, 5)

			b *= 8
			g *= 8
			r *= 8

			imageOutputData.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	imageOutputFile, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}
	defer imageOutputFile.Close()
	png.Encode(imageOutputFile, imageOutputData)

	fmt.Println("Written image data to " + outputFilename)
}
