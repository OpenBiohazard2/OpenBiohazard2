package main

import (
	"fmt"
	"log"
	"os"

	"github.com/samuelyuan/openbiohazard2/fileio"
)

func main() {
  if len(os.Args) < 4 {
		log.Fatal("You only entered ", len(os.Args), " arguments. Command format is invalid.")
		log.Fatal("The syntax of this command is: fileconv.exe [toolName] [inputFilename] [outputFilename]")
		log.Fatal("Tool names supported: tim2png, adt2png, sap2wav")
		log.Fatal("Example command: fileconv tim2png test.tim test.png")
  }

	toolName := os.Args[1]
	inputFilename := os.Args[2]
	outputFilename := os.Args[3]

	fmt.Println("Converting", inputFilename, "to", outputFilename)

	switch toolName {
	case "tim2png":
		timOutput := fileio.LoadTIMFile(inputFilename)
		timOutput.ConvertToPNG(outputFilename)
	case "adt2png":
		adtOutput := fileio.LoadADTFile(inputFilename)
		adtOutput.ConvertToPNG(outputFilename)
	case "sap2wav":
		sapOutput := fileio.LoadSAPFile(inputFilename)
		sapOutput.ConvertToWAV(outputFilename)
	default:
		log.Fatal("You entered an invalid tool name: ", toolName, ". Conversion failed.")
	}
}
