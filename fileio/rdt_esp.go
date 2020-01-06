package fileio

// .esp - Effect sprite data

import (
	"encoding/binary"
	"io"
)

type ESPHeader struct {
	Ids [8]uint8 // If value doesn't equal 0xff, there is an animation
}

type AnimBlockHeader struct {
	NumFrames  uint16
	NumSprites uint16
	Width      uint8
	Height     uint8
	Unknown    uint16
}

type AnimFrame struct {
	SpriteId uint8
	Unknown  [7]uint8
}

// x and y are top-left position of sprite in TIM image
// Offset x and y are the distance to the center of the sprite
// The total width is 2 * OffsetX and total height is 2 * OffsetY
type AnimSprite struct {
	X       uint8
	Y       uint8
	OffsetX int8
	OffsetY int8
}

func LoadRDT_ESP(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) error {
	offset := offsets.OffsetSpriteAnimations
	reader := io.NewSectionReader(r, int64(offset), fileLength-int64(offset))

	// Read the header
	espHeader := ESPHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &espHeader); err != nil {
		return err
	}

	// Read the offset to each block, which is stored separately
	blockOffsets := make([]uint32, 0)
	validSpriteCount := 0
	for i := 0; i < len(espHeader.Ids); i++ {
		if espHeader.Ids[i] == 255 {
			continue
		}
		validSpriteCount++

		reader = io.NewSectionReader(r, int64(offsets.OffsetSpriteAnimationsOffset)-int64(i*4), 4)
		blockOffset := uint32(0)
		if err := binary.Read(reader, binary.LittleEndian, &blockOffset); err != nil {
			return err
		}
		blockOffsets = append(blockOffsets, blockOffset)
	}

	for _, blockOffset := range blockOffsets {
		animationDataOffset := int64(offsets.OffsetSpriteAnimations) + int64(blockOffset)
		reader = io.NewSectionReader(r, animationDataOffset, fileLength-int64(animationDataOffset))
		animBlockHeader := AnimBlockHeader{}
		if err := binary.Read(reader, binary.LittleEndian, &animBlockHeader); err != nil {
			return err
		}

		// Read animation frames
		for j := 0; j < int(animBlockHeader.NumFrames); j++ {
			animFrame := AnimFrame{}
			if err := binary.Read(reader, binary.LittleEndian, &animFrame); err != nil {
				return err
			}
		}

		// Read animation sprites
		for j := 0; j < int(animBlockHeader.NumSprites); j++ {
			animSprite := AnimSprite{}
			if err := binary.Read(reader, binary.LittleEndian, &animSprite); err != nil {
				return err
			}
		}

		// TODO: Figure out the remaining offsets of the block
	}

	// Read Sprite TIM image
	offset = offsets.OffsetSpriteImage
	for i := 0; i < validSpriteCount; i++ {
		reader = io.NewSectionReader(r, int64(offset), fileLength-int64(offset))
		timOutput, err := LoadTIMStream(reader, fileLength-int64(offset))
		if err != nil {
			return err
		}
		offset += uint32(timOutput.NumBytes)
	}
	return nil
}
