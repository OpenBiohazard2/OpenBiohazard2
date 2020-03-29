package fileio

import (
	"io"
)

func LoadRDT_ESP(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*ESPOutput, error) {
	sectionBeginOffset := int64(offsets.OffsetSpriteAnimations)
	reader := io.NewSectionReader(r, sectionBeginOffset, fileLength-sectionBeginOffset)

	eofOffset := int64(offsets.OffsetSpriteAnimationsOffset)

	espOutput, err := LoadESPStream(reader, (eofOffset+4)-sectionBeginOffset, eofOffset-sectionBeginOffset)
	if err != nil {
		return nil, err
	}

	// Read Sprite TIM image
	timOffset := offsets.OffsetSpriteImage
	for i := 0; i < espOutput.ValidSpriteCount; i++ {
		timReader := io.NewSectionReader(r, int64(timOffset), fileLength-int64(timOffset))
		timOutput, err := LoadTIMStream(timReader, fileLength-int64(timOffset))
		if err != nil {
			return nil, err
		}
		timOffset += uint32(timOutput.NumBytes)

		espOutput.SpriteData[i].ImageData = timOutput
	}

	return espOutput, nil
}
