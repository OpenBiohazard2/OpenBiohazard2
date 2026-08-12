package fileio

// .sap - Sound file

import (
	"fmt"
	"os"
)

type SAPOutput struct {
	AudioData []byte
}

func LoadSAPFile(filename string) (*SAPOutput, error) {
	// Skip first 8 bytes
	// The rest is a .wav file
	buffer, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read SAP file %s: %w", filename, err)
	}

	return &SAPOutput{
		AudioData: buffer[8:],
	}, nil
}

func (sapOutput *SAPOutput) ConvertToWAV(outputFilename string) error {
	err := os.WriteFile(outputFilename, sapOutput.AudioData, 0644)
	if err != nil {
		return fmt.Errorf("error writing .sap file to .wav file %s: %w", outputFilename, err)
	}

	fmt.Println("Written audio data to " + outputFilename)
	return nil
}
