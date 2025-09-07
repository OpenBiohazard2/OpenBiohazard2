package fileio

// .esp - Effect sprite data

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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
	SpriteData       []SpriteData
	ValidSpriteCount int
}

type SpriteData struct {
	Id             int
	FrameData      []AnimFrame
	FramePositions []AnimSprite
	AnimMovements  []AnimMovement
	ImageData      *TIMOutput
}

func LoadESPFile(filename string) (*ESPOutput, error) {
	espFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open ESP file %s: %w", filename, err)
	}
	defer espFile.Close()

	fi, err := espFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat ESP file %s: %w", filename, err)
	}
	fileLength := fi.Size()
	espOutput, err := LoadESPStream(espFile, fileLength, fileLength-4)
	if err != nil {
		return nil, fmt.Errorf("failed to load ESP file %s: %w", filename, err)
	}

	return espOutput, nil
}

func LoadESPStream(r io.ReaderAt, fileLength int64, eofOffset int64) (*ESPOutput, error) {
	streamReader := io.NewSectionReader(r, int64(0), fileLength)

	// Read the header
	espHeader := ESPHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &espHeader); err != nil {
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

		offsetReader := io.NewSectionReader(r, eofOffset-int64(i*4), 4)
		blockOffset := uint32(0)
		if err := binary.Read(offsetReader, binary.LittleEndian, &blockOffset); err != nil {
			return nil, err
		}
		blockOffsets = append(blockOffsets, blockOffset)
	}

	spriteData := make([]SpriteData, len(blockOffsets))
	for i, blockOffset := range blockOffsets {
		animationDataOffset := int64(blockOffset)
		animationDataReader := io.NewSectionReader(r, animationDataOffset, fileLength-int64(animationDataOffset))
		animBlockHeader := AnimBlockHeader{}
		if err := binary.Read(animationDataReader, binary.LittleEndian, &animBlockHeader); err != nil {
			return nil, err
		}

		// Read animation frames
		frames := make([]AnimFrame, int(animBlockHeader.NumFrames))
		for j := 0; j < int(animBlockHeader.NumFrames); j++ {
			animFrame := AnimFrame{}
			if err := binary.Read(animationDataReader, binary.LittleEndian, &animFrame); err != nil {
				return nil, err
			}
			frames[j] = animFrame
		}

		// Read animation sprites
		positions := make([]AnimSprite, int(animBlockHeader.NumSprites))
		for j := 0; j < int(animBlockHeader.NumSprites); j++ {
			animSprite := AnimSprite{}
			if err := binary.Read(animationDataReader, binary.LittleEndian, &animSprite); err != nil {
				return nil, err
			}
			positions[j] = animSprite
		}

		// Read animation movement
		movementDataOffsets := make([]uint16, 8)
		if err := binary.Read(animationDataReader, binary.LittleEndian, &movementDataOffsets); err != nil {
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
			ImageData:      nil, // needs to be loaded separately
		}
	}

	output := &ESPOutput{
		SpriteData:       spriteData,
		ValidSpriteCount: validSpriteCount,
	}
	return output, nil
}
