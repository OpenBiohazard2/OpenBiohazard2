package fileio

// .sap - Sound file

import (
	"fmt"
	"io/ioutil"
	"log"
)

type SAPOutput struct {
	AudioData []byte
}

func LoadSAPFile(filename string) *SAPOutput {
	// Skip first 8 bytes
	// The rest is a .wav file
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading " + filename)
	}

	return &SAPOutput{
		AudioData: buffer[8:],
	}
}

func (sapOutput *SAPOutput) ConvertToWAV(outputFilename string) {
	err := ioutil.WriteFile(outputFilename, sapOutput.AudioData, 0644)
	if err != nil {
		log.Fatal("Error writing .sap file to .wav file: ", err)
	}

	fmt.Println("Written audio data to " + outputFilename)
}
