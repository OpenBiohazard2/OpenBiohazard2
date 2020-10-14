package fileio

// .do2 file - Door file

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

type DO2Output struct {
	VABHeaderOutput *VABHeaderOutput
	MD1Output       *MD1Output
	TIMOutput       *TIMOutput
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

	vabHeaderReader := io.NewSectionReader(r, int64(16), fileLength)
	vabHeaderOutput, err := LoadVABHeaderStream(vabHeaderReader, fileLength)
	if err != nil {
		return nil, err
	}

	vabDataOffset := int64(16) + int64(vabHeaderOutput.NumBytes) + int64(8)
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

	output := &DO2Output{
		VABHeaderOutput: vabHeaderOutput,
		MD1Output:       md1Output,
		TIMOutput:       timOutput,
	}
	return output, nil
}
