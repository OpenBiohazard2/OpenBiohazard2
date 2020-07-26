package fileio

import (
	"bytes"
	"encoding/binary"
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

func LoadBINFile(inputFilename string) *BinOutput {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("Load BIN file failed. BIN file doesn't exist: ", inputFilename)
		return nil
	}

	fi, err := binFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	archiveLength := fi.Size()

	imagesIndex, err := LoadBIN(binFile, archiveLength)
	if err != nil {
		log.Fatal(err)
	}

	return &BinOutput{
		ImagesIndex: imagesIndex,
		FileLength:  archiveLength,
	}
}

func LoadBIN(r io.ReaderAt, archiveLength int64) ([]ImageFile, error) {
	reader := io.NewSectionReader(r, int64(0), archiveLength)
	firstOffset := uint32(0)
	if err := binary.Read(reader, binary.LittleEndian, &firstOffset); err != nil {
		return []ImageFile{}, err
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
		offset := uint32(0)
		if err := binary.Read(reader, binary.LittleEndian, &offset); err != nil {
			return []ImageFile{}, err
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

func ExtractItemImage(inputFilename string, binOutput *BinOutput, imageId int) *RoomImageOutput {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("Item BIN file doesn't exist: ", inputFilename)
		return nil
	}
	binReader := io.NewSectionReader(binFile, int64(0), binOutput.FileLength)

	imageBlock := binOutput.ImagesIndex[imageId]
	if imageBlock.Length == 0 {
		fmt.Println("Warning: Image has no data")
		return nil
	}

	blockLength := int(imageBlock.Length)
	if blockLength > int(binOutput.FileLength) {
		blockLength = int(binOutput.FileLength)
	}

	adtReader := io.NewSectionReader(binReader, int64(imageBlock.Offset), int64(blockLength))
	adtOutput := LoadADTStream(adtReader)
	timReader := bytes.NewReader(adtOutput.RawData)
	timOutput, err := LoadTIMStream(timReader, int64(len(adtOutput.RawData)))
	if err != nil {
		log.Fatal("Invalid BIN data. ", err)
	}

	return &RoomImageOutput{
		BackgroundImage: nil,
		ImageMask:       timOutput,
	}
}

// Room image is stored as an ADT file
func ExtractRoomBackground(inputFilename string, binOutput *BinOutput, roomId int) *RoomImageOutput {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("Room BIN file doesn't exist: ", inputFilename)
		return nil
	}
	binReader := io.NewSectionReader(binFile, int64(0), binOutput.FileLength)

	imageBlock := binOutput.ImagesIndex[roomId]
	if imageBlock.Length == 0 {
		fmt.Println("Warning: Room has no data")
		return nil
	}

	// The first part is the background image, which is an .adt file
	adtReader := io.NewSectionReader(binReader, int64(imageBlock.Offset), int64(imageBlock.Length))
	adtOutput := LoadADTStream(adtReader)

	// The next part is an image mask, which is a .tim file
	beginOffset := (320 * 256 * 2)

	// The background image doesn't contain an image mask
	if len(adtOutput.RawData) == beginOffset {
		return &RoomImageOutput{
			BackgroundImage: adtOutput,
			ImageMask:       nil,
		}
	}
	timReader := bytes.NewReader(adtOutput.RawData[beginOffset:])
	timOutput, err := LoadTIMStream(timReader, int64(len(adtOutput.RawData))-int64(beginOffset))
	if err != nil {
		log.Fatal(err)
	}

	return &RoomImageOutput{
		BackgroundImage: adtOutput,
		ImageMask:       timOutput,
	}
}
