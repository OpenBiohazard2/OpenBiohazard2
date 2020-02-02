package fileio

// .msg - Message data

import (
	"encoding/binary"
	"io"
	"log"
)

type MSGOutput struct {
}

var (
	convertText = [6]string{
		" .___()_____0123",
		"456789:_,\"!?_ABC",
		"DEFGHIJKLMNOPQRS",
		"TUVWXYZ[/]'-_abc",
		"defghijklmnopqrs",
		"tuvwxyz_________",
	}
)

func LoadRDT_MSGStream(fileReader io.ReaderAt, fileLength int64) (*MSGOutput, error) {
	streamReader := io.NewSectionReader(fileReader, int64(0), fileLength)

	offsets := make([]uint16, 0)
	firstOffset := uint16(0)
	if err := binary.Read(streamReader, binary.LittleEndian, &firstOffset); err != nil {
		return nil, err
	}

	offsets = append(offsets, firstOffset)
	for i := 2; i < int(firstOffset); i += 2 {
		nextOffset := uint16(0)
		if err := binary.Read(streamReader, binary.LittleEndian, &nextOffset); err != nil {
			return nil, err
		}
		offsets = append(offsets, nextOffset)
	}

	for i := 0; i < len(offsets)-1; i++ {
		if offsets[i] >= offsets[i+1] {
			log.Fatal("MSG offsets are not sorted")
		}

		textData := make([]uint8, offsets[i+1]-offsets[i])
		if err := binary.Read(streamReader, binary.LittleEndian, &textData); err != nil {
			return nil, err
		}
		convertBytesToText(textData)
	}

	// Read last message
	textData := make([]uint8, 0)
	for i := 0; i < 200; i++ {
		nextChar := uint8(0)
		if err := binary.Read(streamReader, binary.LittleEndian, &nextChar); err != nil {
			return nil, err
		}

		textData = append(textData, nextChar)
		if nextChar == 0xFE {
			break
		}
	}
	convertBytesToText(textData)

	return nil, nil
}

func convertBytesToText(byteData []uint8) []string {
	textData := make([]string, len(byteData))

	for i, number := range byteData {
		if number >= 96 {
			if number == 0xF3 {
				textData[i] = string("?")
			}
			if number == 0xFC {
				textData[i] = string("\n")
			}

			continue
		}

		row := number / 16
		column := number % 16
		textData[i] = string(convertText[row][column])
	}

	return textData
}
