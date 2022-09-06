package fileio

// .do2 file - Door file

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

type DO2FileFormat struct {
	VHOffset  int64 // .vab header
	VHLength  int64
	VBOffset  int64 // .vab data
	VBLength  int64
	MD1Offset int64 // model file
	MD1Length int64
	TIMOffset int64 // texture file
	TIMLength int64
}

type DO2Output struct {
	VABHeaderOutput *VABHeaderOutput
	MD1Output       *MD1Output
	TIMOutput       *TIMOutput
	DO2FileFormat   *DO2FileFormat
}

func LoadDO2File(filename string) *DO2Output {
	file, _ := os.Open(filename)
	defer file.Close()
	if file == nil {
		log.Fatal("DO2 file doesn't exist: ", filename)
		return nil
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fileLength := fi.Size()
	fileOutput, err := LoadDO2Stream(file, fileLength)
	if err != nil {
		log.Fatal("Failed to load DO2 file: ", err)
		return nil
	}

	return fileOutput
}

func LoadDO2Stream(r io.ReaderAt, fileLength int64) (*DO2Output, error) {
	streamReader := io.NewSectionReader(r, int64(0), fileLength)
	do2Header := make([]uint16, 8)
	if err := binary.Read(streamReader, binary.LittleEndian, &do2Header); err != nil {
		return nil, err
	}

	vabHeaderOffset := int64(16)
	vabHeaderReader := io.NewSectionReader(r, vabHeaderOffset, fileLength)
	vabHeaderOutput, err := LoadVABHeaderStream(vabHeaderReader, fileLength)
	if err != nil {
		return nil, err
	}

	vabDataOffset := vabHeaderOffset + int64(vabHeaderOutput.NumBytes) + int64(8)
	vabDataReader := io.NewSectionReader(r, vabDataOffset, fileLength)
	vabDataOutput, err := LoadVABDataStream(vabDataReader, fileLength, vabHeaderOutput)
	if err != nil {
		return nil, err
	}

	offsetAfterVab := vabDataOffset + int64(vabDataOutput.NumBytes)

	// There is a block of zeros after the vab data
	emptyReader := io.NewSectionReader(r, offsetAfterVab, fileLength)
	var intBlock uint32
	emptySectionBytes := 0
	for {
		if err := binary.Read(emptyReader, binary.LittleEndian, &intBlock); err != nil {
			return nil, err
		}

		if intBlock != 0 {
			break
		}

		emptySectionBytes += 4
	}

	// Skip this section
	unknownSectionReader := io.NewSectionReader(r, offsetAfterVab+int64(emptySectionBytes), fileLength)
	var unknownSectionSize uint32
	if err := binary.Read(unknownSectionReader, binary.LittleEndian, &unknownSectionSize); err != nil {
		return nil, err
	}

	md1Offset := offsetAfterVab + int64(emptySectionBytes) + int64(unknownSectionSize)
	md1Reader := io.NewSectionReader(r, md1Offset, fileLength-md1Offset)
	md1Output, err := LoadMD1Stream(md1Reader, fileLength-md1Offset)
	if err != nil {
		log.Fatal("Error loading md1 file in do2 file")
		return nil, err
	}

	timOffset := md1Offset + int64(md1Output.NumBytes)
	timReader := io.NewSectionReader(r, timOffset, fileLength-timOffset)
	timOutput, err := LoadTIMStream(timReader, fileLength-timOffset)
	if err != nil {
		log.Fatal("Error loading tim file in do2 file")
		return nil, err
	}

	do2FileFormat := &DO2FileFormat{
		VHOffset:  vabHeaderOffset,
		VHLength:  int64(vabHeaderOutput.NumBytes),
		VBOffset:  vabDataOffset,
		VBLength:  int64(vabDataOutput.NumBytes),
		MD1Offset: md1Offset,
		MD1Length: int64(md1Output.NumBytes),
		TIMOffset: timOffset,
		TIMLength: int64(timOutput.NumBytes),
	}

	output := &DO2Output{
		VABHeaderOutput: vabHeaderOutput,
		MD1Output:       md1Output,
		TIMOutput:       timOutput,
		DO2FileFormat:   do2FileFormat,
	}
	return output, nil
}
