package resource

import (
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// LoadTIMImages loads multiple TIM images and converts them to Image16Bit
func LoadTIMImages(filename string) []*Image16Bit {
	timOutputs, err := fileio.LoadTIMImages(filename)
	if err != nil {
		log.Fatal("Error loading TIM images: ", err)
	}

	images := make([]*Image16Bit, len(timOutputs))
	for i, timOutput := range timOutputs {
		images[i] = ConvertPixelsToImage16Bit(timOutput.PixelData)
	}

	return images
}

// LoadADTImage loads a single ADT image and converts it to Image16Bit
func LoadADTImage(filename string) *Image16Bit {
	adtOutput, err := fileio.LoadADTFile(filename)
	if err != nil {
		log.Fatal("Error loading ADT image: ", err)
	}

	return ConvertPixelsToImage16Bit(adtOutput.PixelData)
}
