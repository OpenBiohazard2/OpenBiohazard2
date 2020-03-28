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
	SpriteId   uint8
	Count      uint8
	Time       uint8
	SquareSide uint8 // length and width are the same
	X          int16
	Y          int16
}

// x and y are top-left position of sprite in TIM sprite sheet
// Offset is used to shift sprite frame to align it with other frames with different dimensions
type AnimSprite struct {
	ImageX  uint8
	ImageY  uint8
	OffsetX int8
	OffsetY int8
}

type AnimMovement struct {
	FunctionId0  uint8
	FunctionId1  uint8
	Unknown0     [2]uint8
	TranslateX   uint16
	TranslateY   uint16
	Acceleration [3]uint8
	Unknown1     uint8
	Speed        [3]int16
	Unknown2     [3]uint16
}

type ESPOutput struct {
	SpriteData []SpriteData
}

type SpriteData struct {
	Id             int
	FrameData      []AnimFrame
	FramePositions []AnimSprite
	AnimMovements  []AnimMovement
	ImageData      *TIMOutput
}

func LoadRDT_ESP(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*ESPOutput, error) {
	offset := offsets.OffsetSpriteAnimations
	reader := io.NewSectionReader(r, int64(offset), fileLength-int64(offset))

	// Read the header
	espHeader := ESPHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &espHeader); err != nil {
		return nil, err
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
			return nil, err
		}
		blockOffsets = append(blockOffsets, blockOffset)
	}

	spriteData := make([]SpriteData, len(blockOffsets))
	for i, blockOffset := range blockOffsets {
		animationDataOffset := int64(offsets.OffsetSpriteAnimations) + int64(blockOffset)
		reader = io.NewSectionReader(r, animationDataOffset, fileLength-int64(animationDataOffset))
		animBlockHeader := AnimBlockHeader{}
		if err := binary.Read(reader, binary.LittleEndian, &animBlockHeader); err != nil {
			return nil, err
		}

		// Read animation frames
		frames := make([]AnimFrame, int(animBlockHeader.NumFrames))
		for j := 0; j < int(animBlockHeader.NumFrames); j++ {
			animFrame := AnimFrame{}
			if err := binary.Read(reader, binary.LittleEndian, &animFrame); err != nil {
				return nil, err
			}
			frames[j] = animFrame
		}

		// Read animation sprites
		positions := make([]AnimSprite, int(animBlockHeader.NumSprites))
		for j := 0; j < int(animBlockHeader.NumSprites); j++ {
			animSprite := AnimSprite{}
			if err := binary.Read(reader, binary.LittleEndian, &animSprite); err != nil {
				return nil, err
			}
			positions[j] = animSprite
		}

		// Read animation movement
		movementDataOffsets := make([]uint16, 8)
		if err := binary.Read(reader, binary.LittleEndian, &movementDataOffsets); err != nil {
			return nil, err
		}

		animBlockHeaderSize := 8
		frameBlockSize := 8 * int(animBlockHeader.NumFrames)
		spriteBlockSize := 4 * int(animBlockHeader.NumSprites)
		movementReaderOffset := animationDataOffset + int64(animBlockHeaderSize+frameBlockSize+spriteBlockSize)

		animMovements := make([]AnimMovement, 0)
		for j := 0; j < 8; j++ {
			if movementDataOffsets[j] == 0 {
				break
			}
			newOffset := movementReaderOffset + int64(movementDataOffsets[j]*4)
			movementReader := io.NewSectionReader(r, newOffset, fileLength-newOffset)

			movementCounts := make([]uint16, 2)
			if err := binary.Read(movementReader, binary.LittleEndian, &movementCounts); err != nil {
				return nil, err
			}

			movementData := AnimMovement{}
			if err := binary.Read(movementReader, binary.LittleEndian, &movementData); err != nil {
				return nil, err
			}
			animMovements = append(animMovements, movementData)

			// TODO: Figure out the remaining offsets of the block
		}

		spriteData[i] = SpriteData{
			Id:             int(espHeader.Ids[i]),
			FrameData:      frames,
			FramePositions: positions,
			AnimMovements:  animMovements,
		}
	}

	// Read Sprite TIM image
	offset = offsets.OffsetSpriteImage
	for i := 0; i < validSpriteCount; i++ {
		reader = io.NewSectionReader(r, int64(offset), fileLength-int64(offset))
		timOutput, err := LoadTIMStream(reader, fileLength-int64(offset))
		if err != nil {
			return nil, err
		}
		offset += uint32(timOutput.NumBytes)

		spriteData[i].ImageData = timOutput
	}

	output := &ESPOutput{
		SpriteData: spriteData,
	}
	return output, nil
}
