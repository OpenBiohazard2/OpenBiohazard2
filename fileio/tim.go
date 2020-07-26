package fileio

// .tim - Playstation 1 Texture format

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	TIM_BPP_4  = 8
	TIM_BPP_8  = 9
	TIM_BPP_16 = 2
	TIM_BPP_24 = 3
)

type TIMHeader struct {
	Magic     uint32 // always 16
	BPP       uint32 // (8 = 4 bit), (9 = 8 bit), (2 = 16 bit), (3 = 24 bit)
	Offset    uint32
	OriginX   uint16
	OriginY   uint16
	NumColors uint16
	NumCluts  uint16
}

type TIMImageHeader struct {
	Size    uint32 // total image data size in bytes (image header + image data)
	OriginX uint16 // image x origin
	OriginY uint16 // image y origin
	Width   uint16 // *4 for 4 bit, *2 for 8 bit
	Height  uint16 // image height
}

type TIMOutput struct {
	PixelData   [][]uint16
	ImageWidth  int
	ImageHeight int
	NumPalettes int
	NumBytes    int
}

func LoadTIMFile(filename string) *TIMOutput {
	timFile, _ := os.Open(filename)
	defer timFile.Close()
	if timFile == nil {
		log.Fatal("TIM file doesn't exist: ", filename)
		return nil
	}
	fi, err := timFile.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fileLength := fi.Size()
	timOutput, err := LoadTIMStream(timFile, fileLength)
	if err != nil {
		log.Fatal("Failed to load TIM file: ", err)
		return nil
	}
	return timOutput
}

func LoadTIMStream(r io.ReaderAt, fileLength int64) (*TIMOutput, error) {
	reader := io.NewSectionReader(r, int64(0), fileLength)

	timHeader := TIMHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &timHeader); err != nil {
		return nil, err
	}

	if timHeader.Magic != 16 {
		log.Fatal("TIM header is invalid: ", timHeader)
	}

	// Read TIM cluts
	if timHeader.BPP == TIM_BPP_4 {
		return read4BPP(reader, timHeader)
	} else if timHeader.BPP == TIM_BPP_8 {
		return read8BPP(reader, timHeader)
	} else {
		log.Fatal(fmt.Errorf("BPP %v is not supported.\n", timHeader.BPP))
	}

	return nil, nil
}

func read4BPP(reader *io.SectionReader, timHeader TIMHeader) (*TIMOutput, error) {
	numColors := timHeader.NumColors
	palettes := make([][]uint16, int(timHeader.NumCluts))

	if numColors != 16 {
		log.Fatal("4BPP TIM image does not have 16 colors in palette. It has ", numColors, "colors.")
	}

	for i := 0; i < int(timHeader.NumCluts); i++ {
		palettes[i] = make([]uint16, numColors)
		if err := binary.Read(reader, binary.LittleEndian, &palettes[i]); err != nil {
			return nil, err
		}
	}
	timImageHeader := TIMImageHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &timImageHeader); err != nil {
		return nil, err
	}

	totalImageWidth := int(timImageHeader.Width) * 4
	totalImageHeight := int(timImageHeader.Height)

	imageDataLength := totalImageWidth * totalImageHeight
	imageData := make([]uint8, imageDataLength/2)
	if err := binary.Read(reader, binary.LittleEndian, &imageData); err != nil {
		return nil, err
	}

	pixelData2D := make([][]uint16, totalImageHeight)
	for i := 0; i < totalImageHeight; i++ {
		pixelData2D[i] = make([]uint16, totalImageWidth)
	}

	for i := 0; i < imageDataLength; i += 2 {
		index := imageData[i/2]
		colorPalette := palettes[(i%totalImageWidth)/(totalImageWidth/len(palettes))]

		color1Index := index & 0x0F
		color2Index := (index & 0xF0) >> 4

		// color 1
		x := i % totalImageWidth
		y := i / totalImageWidth
		pixelData2D[y][x] = colorPalette[color1Index]

		// color 2
		x = (i + 1) % totalImageWidth
		y = (i + 1) / totalImageWidth
		pixelData2D[y][x] = colorPalette[color2Index]
	}

	headerBytes := 32
	paletteBytes := int(timHeader.NumColors) * int(timHeader.NumCluts) * 2
	imageBytes := (totalImageWidth / 2) * totalImageHeight
	numBytes := imageBytes + paletteBytes + headerBytes

	timOutput := &TIMOutput{
		PixelData:   pixelData2D,
		ImageWidth:  totalImageWidth,
		ImageHeight: totalImageHeight,
		NumPalettes: int(timHeader.NumCluts),
		NumBytes:    numBytes,
	}
	return timOutput, nil
}

func read8BPP(reader *io.SectionReader, timHeader TIMHeader) (*TIMOutput, error) {
	numColors := timHeader.NumColors
	palettes := make([][]uint16, int(timHeader.NumCluts))

	if numColors != 256 {
		log.Fatal("8BPP TIM image does not have 256 colors in palette. It has ", numColors, "colors.")
	}

	for i := 0; i < int(timHeader.NumCluts); i++ {
		palettes[i] = make([]uint16, numColors)
		if err := binary.Read(reader, binary.LittleEndian, &palettes[i]); err != nil {
			return nil, err
		}
	}
	timImageHeader := TIMImageHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &timImageHeader); err != nil {
		return nil, err
	}

	totalImageWidth := int(timImageHeader.Width) * 2
	totalImageHeight := int(timImageHeader.Height)

	imageDataLength := totalImageWidth * totalImageHeight
	imageData := make([]uint8, imageDataLength)
	if err := binary.Read(reader, binary.LittleEndian, &imageData); err != nil {
		return nil, err
	}

	pixelData2D := make([][]uint16, totalImageHeight)
	for i := 0; i < totalImageHeight; i++ {
		pixelData2D[i] = make([]uint16, totalImageWidth)
	}

	for i := 0; i < imageDataLength; i++ {
		index := imageData[i]
		paletteIndex := int(math.Floor(float64(i%totalImageWidth) * float64(len(palettes)) / float64(totalImageWidth)))
		colorPalette := palettes[paletteIndex]
		x := i % totalImageWidth
		y := i / totalImageWidth
		pixelData2D[y][x] = colorPalette[index]
	}

	headerBytes := 32
	paletteBytes := int(timHeader.NumColors) * int(timHeader.NumCluts) * 2
	imageBytes := totalImageWidth * totalImageHeight
	numBytes := imageBytes + paletteBytes + headerBytes

	timOutput := &TIMOutput{
		PixelData:   pixelData2D,
		ImageWidth:  totalImageWidth,
		ImageHeight: totalImageHeight,
		NumPalettes: int(timHeader.NumCluts),
		NumBytes:    numBytes,
	}
	return timOutput, nil
}

func (timOutput *TIMOutput) ConvertToRenderData() []uint16 {
	pixelData2D := timOutput.PixelData
	pixelData1D := make([]uint16, len(pixelData2D)*len(pixelData2D[0]))

	for y := 0; y < len(pixelData2D); y++ {
		for x := 0; x < len(pixelData2D[y]); x++ {
			index := (y * len(pixelData2D[y])) + x
			pixelData1D[index] = pixelData2D[y][x]
		}
	}
	return pixelData1D
}

func (timOutput *TIMOutput) ConvertToPNG(outputFilename string) error {
	pixelData2D := timOutput.PixelData
	totalImageWidth := timOutput.ImageWidth
	totalImageHeight := timOutput.ImageHeight
	imageOutputData := image.NewRGBA(image.Rect(0, 0, totalImageWidth, totalImageHeight))
	for y := 0; y < totalImageHeight; y++ {
		for x := 0; x < totalImageWidth; x++ {
			colorBits := fmt.Sprintf("%016b", pixelData2D[y][x])
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
	return nil
}
