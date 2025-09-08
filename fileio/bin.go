package fileio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type ImageFile struct {
	Offset uint32
	Length uint32
}

type BinOutput struct {
	ImagesIndex []ImageFile
	FileLength  int64
}

type RoomImageOutput struct {
	BackgroundImage *ADTOutput
	ImageMask       *TIMOutput
}

func LoadBINFile(inputFilename string) (*BinOutput, error) {
	binFile, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open BIN file %s: %w", inputFilename, err)
	}
	defer binFile.Close()

	fi, err := binFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat BIN file %s: %w", inputFilename, err)
	}
	archiveLength := fi.Size()

	imagesIndex, err := LoadBIN(binFile, archiveLength)
	if err != nil {
		return nil, fmt.Errorf("failed to load BIN data from %s: %w", inputFilename, err)
	}

	return &BinOutput{
		ImagesIndex: imagesIndex,
		FileLength:  archiveLength,
	}, nil
}

func LoadBIN(r io.ReaderAt, archiveLength int64) ([]ImageFile, error) {
	streamReader := NewStreamReader(io.NewSectionReader(r, int64(0), archiveLength))

	firstOffset, err := streamReader.ReadUint32()
	if err != nil {
		return []ImageFile{}, fmt.Errorf("failed to read first offset: %w", err)
	}

	numImages := firstOffset / 4
	imagesIndex := make([]ImageFile, 0)

	// Read offsets to each image
	imageFile := ImageFile{
		Offset: firstOffset,
		Length: 0,
	}
	imagesIndex = append(imagesIndex, imageFile)
	for i := 1; i < int(numImages); i++ {
		offset, err := streamReader.ReadUint32()
		if err != nil {
			return []ImageFile{}, fmt.Errorf("failed to read offset %d: %w", i, err)
		}

		// Zero offset is invalid
		if offset == 0 {
			continue
		}
		imageFile := ImageFile{
			Offset: offset,
			Length: 0,
		}
		imagesIndex = append(imagesIndex, imageFile)
	}

	// calculate length of image data
	for i := 0; i < len(imagesIndex)-1; i++ {
		imagesIndex[i].Length = imagesIndex[i+1].Offset - imagesIndex[i].Offset
	}
	imagesIndex[len(imagesIndex)-1].Length = uint32(int(archiveLength) - int(imagesIndex[len(imagesIndex)-1].Offset))
	return imagesIndex, nil
}

func LoadTIMImages(inputFilename string) ([]*TIMOutput, error) {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("Failed to load TIM images. BIN file doesn't exist", inputFilename)
		return nil, nil
	}

	fi, err := binFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	archiveLength := fi.Size()

	images := make([]*TIMOutput, 0)
	totalBytesRead := 0

	for totalBytesRead < int(archiveLength) {
		timReader := io.NewSectionReader(binFile, int64(totalBytesRead), archiveLength)
		timOutput, err := LoadTIMStream(timReader, archiveLength)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, timOutput)
		totalBytesRead += timOutput.NumBytes
	}

	return images, nil
}

func ExtractItemImage(inputFilename string, binOutput *BinOutput, imageId int) (*RoomImageOutput, error) {
	binFile, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open BIN file %s: %w", inputFilename, err)
	}
	defer binFile.Close()

	binReader := io.NewSectionReader(binFile, int64(0), binOutput.FileLength)

	imageBlock := binOutput.ImagesIndex[imageId]
	if imageBlock.Length == 0 {
		fmt.Println("Warning: Image has no data")
		return nil, nil
	}

	blockLength := int(imageBlock.Length)
	if blockLength > int(binOutput.FileLength) {
		blockLength = int(binOutput.FileLength)
	}

	adtReader := io.NewSectionReader(binReader, int64(imageBlock.Offset), int64(blockLength))
	adtOutput, err := LoadADTStream(adtReader)
	if err != nil {
		return nil, fmt.Errorf("failed to load ADT stream: %w", err)
	}
	timReader := bytes.NewReader(adtOutput.RawData)
	timOutput, err := LoadTIMStream(timReader, int64(len(adtOutput.RawData)))
	if err != nil {
		return nil, fmt.Errorf("invalid BIN data: %w", err)
	}

	return &RoomImageOutput{
		BackgroundImage: nil,
		ImageMask:       timOutput,
	}, nil
}

// Room image is stored as an ADT file
func ExtractRoomBackground(inputFilename string, binOutput *BinOutput, roomId int) (*RoomImageOutput, error) {
	binFile, err := os.Open(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open BIN file %s: %w", inputFilename, err)
	}
	defer binFile.Close()

	binReader := io.NewSectionReader(binFile, int64(0), binOutput.FileLength)

	imageBlock := binOutput.ImagesIndex[roomId]
	if imageBlock.Length == 0 {
		fmt.Println("Warning: Room has no data")
		return nil, nil
	}

	// The first part is the background image, which is an .adt file
	adtReader := io.NewSectionReader(binReader, int64(imageBlock.Offset), int64(imageBlock.Length))
	adtOutput, err := LoadADTStream(adtReader)
	if err != nil {
		return nil, fmt.Errorf("failed to load ADT stream: %w", err)
	}

	// The next part is an image mask, which is a .tim file
	beginOffset := (320 * 256 * 2)

	// The background image doesn't contain an image mask
	if len(adtOutput.RawData) == beginOffset {
		return &RoomImageOutput{
			BackgroundImage: adtOutput,
			ImageMask:       nil,
		}, nil
	}
	timReader := bytes.NewReader(adtOutput.RawData[beginOffset:])
	timOutput, err := LoadTIMStream(timReader, int64(len(adtOutput.RawData))-int64(beginOffset))
	if err != nil {
		return nil, fmt.Errorf("failed to load TIM stream: %w", err)
	}

	return &RoomImageOutput{
		BackgroundImage: adtOutput,
		ImageMask:       timOutput,
	}, nil
}
