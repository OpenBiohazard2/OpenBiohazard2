package resource

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// LoadTIMImages loads multiple TIM images and converts them to Image16Bit
func LoadTIMImages(filename string) ([]*Image16Bit, error) {
	timOutputs, err := fileio.LoadTIMImages(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading TIM images from %s: %w", filename, err)
	}

	images := make([]*Image16Bit, len(timOutputs))
	for i, timOutput := range timOutputs {
		images[i] = ConvertPixelsToImage16Bit(timOutput.PixelData)
	}

	return images, nil
}

// LoadADTImage loads a single ADT image and converts it to Image16Bit
func LoadADTImage(filename string) (*Image16Bit, error) {
	adtOutput, err := fileio.LoadADTFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading ADT image from %s: %w", filename, err)
	}

	return ConvertPixelsToImage16Bit(adtOutput.PixelData), nil
}
