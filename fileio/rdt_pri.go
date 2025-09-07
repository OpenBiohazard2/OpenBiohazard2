package fileio

// .pri - Camera mask sprite data

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PRIHeader struct {
	CountOffsets uint16
	CountMasks   uint16
}

type PRIRelativeOffset struct {
	MaskCount uint16 // Number of masks with which to use this structure
	Unknown   uint16
	DestX     int16 // Destination position on background image
	DestY     int16
}

type PRIMaskSquare struct {
	SrcX  uint8  // x-coordinate on source .tim image
	SrcY  uint8  // x-coordinate on source .tim image
	DestX uint8  // destination x-coordinate
	DestY uint8  // destination y-coordinate
	DestZ uint16 // destination z-coordinate (depth/z-buffer)
	Width uint16 // width of the mask tile
}

type PRIMaskRectangle struct {
	SrcX   uint8  // x-coordinate on source .tim image
	SrcY   uint8  // x-coordinate on source .tim image
	DestX  uint8  // destination x-coordinate
	DestY  uint8  // destination y-coordinate
	DestZ  uint16 //  destination z-coordinate (depth/z-buffer)
	Zero   uint16
	Width  uint16 // width of the mask tile
	Height uint16 // height of the mask tile
}

type MaskRectangle struct {
	SrcX   int
	SrcY   int
	DestX  int
	DestY  int
	Depth  int
	Zero   int
	Width  int
	Height int
}

type PRIOutput struct {
	Masks []MaskRectangle
}

func LoadRDT_PRI(r io.ReaderAt, fileLength int64) (*PRIOutput, error) {
	reader := io.NewSectionReader(r, int64(0), fileLength)

	priHeader := PRIHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &priHeader); err != nil {
		return nil, err
	}

	if priHeader.CountOffsets == 0xffff || priHeader.CountMasks == 0xffff {
		return nil, nil
	}

	relativeOffsets := make([]PRIRelativeOffset, int(priHeader.CountOffsets))
	if err := binary.Read(reader, binary.LittleEndian, &relativeOffsets); err != nil {
		return nil, err
	}

	totalCountMask := 0
	for _, offsetData := range relativeOffsets {
		totalCountMask += int(offsetData.MaskCount)
	}
	if totalCountMask > int(priHeader.CountMasks) {
		return nil, fmt.Errorf("Actual mask count exceeds expected count")
	}

	// current position in file
	currentOffset := int64(4) + (int64(len(relativeOffsets)) * int64(8))
	maskData := make([]MaskRectangle, 0)
	for numOffset := 0; numOffset < int(priHeader.CountOffsets); numOffset++ {
		for numMask := 0; numMask < int(relativeOffsets[numOffset].MaskCount); numMask++ {
			// skip 6 bytes
			reader := io.NewSectionReader(r, currentOffset+int64(6), fileLength)
			squareFlag := uint16(0)
			if err := binary.Read(reader, binary.LittleEndian, &squareFlag); err != nil {
				return nil, err
			}

			reader = io.NewSectionReader(r, currentOffset, fileLength)
			if squareFlag == 0 {
				readRect := PRIMaskRectangle{}
				if err := binary.Read(reader, binary.LittleEndian, &readRect); err != nil {
					return nil, err
				}

				// Add offset
				maskRect := MaskRectangle{
					SrcX:   int(readRect.SrcX),
					SrcY:   int(readRect.SrcY),
					DestX:  int(readRect.DestX) + int(relativeOffsets[numOffset].DestX),
					DestY:  int(readRect.DestY) + int(relativeOffsets[numOffset].DestY),
					Depth:  int(readRect.DestZ),
					Zero:   0,
					Width:  int(readRect.Width),
					Height: int(readRect.Height),
				}

				rectSizeBytes := int64(12)
				if maskRect.DestX+maskRect.Width > 320 || maskRect.DestY+maskRect.Height > 240 {
					return nil, fmt.Errorf("Mask rect is out of bounds: %+v", maskRect)
				}

				maskData = append(maskData, maskRect)
				currentOffset += rectSizeBytes
			} else {
				readSquare := PRIMaskSquare{}
				if err := binary.Read(reader, binary.LittleEndian, &readSquare); err != nil {
					return nil, err
				}
				// Add offset
				maskRect := MaskRectangle{
					SrcX:   int(readSquare.SrcX),
					SrcY:   int(readSquare.SrcY),
					DestX:  int(readSquare.DestX) + int(relativeOffsets[numOffset].DestX),
					DestY:  int(readSquare.DestY) + int(relativeOffsets[numOffset].DestY),
					Depth:  int(readSquare.DestZ),
					Zero:   0,
					Width:  int(readSquare.Width),
					Height: int(readSquare.Width),
				}

				squareSizeBytes := int64(8)
				if maskRect.DestX+maskRect.Width > 320 || maskRect.DestY+maskRect.Height > 240 {
					return nil, fmt.Errorf("Mask rect is out of bounds: %+v", maskRect)
				}

				maskData = append(maskData, maskRect)

				currentOffset += squareSizeBytes
			}
		}
	}

	output := &PRIOutput{
		Masks: maskData,
	}
	return output, nil

}
