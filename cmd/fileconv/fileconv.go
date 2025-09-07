package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf(`Invalid arguments. You entered %d arguments.
Usage: fileconv [toolName] [inputFilename] [outputFilename]
Tool names supported: tim2png, adt2png, sap2wav
Example: fileconv tim2png test.tim test.png`, len(os.Args))
	}

	toolName := os.Args[1]
	inputFilename := os.Args[2]
	outputFilename := os.Args[3]

	fmt.Println("Converting", inputFilename, "to", outputFilename)

	switch toolName {
	case "tim2png":
		timOutput, err := fileio.LoadTIMFile(inputFilename)
		if err != nil {
			log.Fatal("Error loading TIM file: ", err)
		}
		timOutput.ConvertToPNG(outputFilename)
	case "adt2png":
		adtOutput, err := fileio.LoadADTFile(inputFilename)
		if err != nil {
			log.Fatal("Error loading ADT file: ", err)
		}
		adtOutput.ConvertToPNG(outputFilename)
	case "sap2wav":
		sapOutput, err := fileio.LoadSAPFile(inputFilename)
		if err != nil {
			log.Fatal("Error loading SAP file: ", err)
		}
		sapOutput.ConvertToWAV(outputFilename)
	default:
		log.Fatalf("Invalid tool name: %s. Conversion failed.", toolName)
	}
}
