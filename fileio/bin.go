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
		log.Fatal("File doesn't exist")
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

	imagesIndex := make([]ImageFile, numImages)

	// Offsets
	imagesIndex[0].Offset = firstOffset
	for i := 1; i < int(numImages); i++ {
		Offset := uint32(0)
		if err := binary.Read(reader, binary.LittleEndian, &Offset); err != nil {
			return []ImageFile{}, err
		}
		imagesIndex[i].Offset = Offset
	}

	// calculate length of image data
	for i := 0; i < int(numImages)-1; i++ {
		imagesIndex[i].Length = imagesIndex[i+1].Offset - imagesIndex[i].Offset
	}
	imagesIndex[numImages-1].Length = uint32(int(archiveLength) - int(imagesIndex[numImages-1].Offset))
	return imagesIndex, nil
}

func LoadTIMImages(inputFilename string) ([]*TIMOutput, error) {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("File doesn't exist")
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

// Room image is stored as an ADT file
func ExtractRoomBackground(inputFilename string, binOutput *BinOutput, roomId int) *RoomImageOutput {
	binFile, _ := os.Open(inputFilename)
	defer binFile.Close()

	if binFile == nil {
		log.Fatal("File doesn't exist")
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
