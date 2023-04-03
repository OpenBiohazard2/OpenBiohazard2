package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("You only entered ", len(os.Args), " arguments. Command format is invalid.")
		log.Fatal("The syntax of this command is: fileconv [fileFormat] [inputFilename]")
		log.Fatal("File formats supported: do2")
		log.Fatal("Example command: fileconv do2 door00.do2")
	}

	fileFormat := os.Args[1]
	inputFilename := os.Args[2]

	inputBase := filenameWithoutExtension(filepath.Base(inputFilename))
	outputFolder := filepath.Join(filepath.Dir(inputFilename), inputBase)
	fmt.Println("Unpacking", inputFilename, "to", outputFolder)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, 0777)
	}

	switch fileFormat {
	case "do2":
		do2Output := fileio.LoadDO2File(inputFilename)

		file, err := os.Open(inputFilename)
		if err != nil {
			log.Panicf("Failed to read file %s, error: %s", inputFilename, err)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		baseOutputFilename := filepath.Join(outputFolder, inputBase)
		do2FileFormat := do2Output.DO2FileFormat
		writeFile(baseOutputFilename+".vh", getBufferSubset(data, do2FileFormat.VHOffset, do2FileFormat.VHLength))
		writeFile(baseOutputFilename+".vb", getBufferSubset(data, do2FileFormat.VBOffset, do2FileFormat.VBLength))
		writeFile(baseOutputFilename+".md1", getBufferSubset(data, do2FileFormat.MD1Offset, do2FileFormat.MD1Length))
		writeFile(baseOutputFilename+".tim", getBufferSubset(data, do2FileFormat.TIMOffset, do2FileFormat.TIMLength))
	default:
		log.Fatal("You entered an invalid file format: ", fileFormat, ". Unpack failed.")
	}
}

func getBufferSubset(data []byte, offset int64, length int64) []byte {
	return data[offset : offset+length]
}

func writeFile(outputFilename string, buffer []byte) {
	err := ioutil.WriteFile(outputFilename, buffer, 0644)
	if err != nil {
		log.Panicf("Error creating %s, error: %s ", outputFilename, err)
	}
	fmt.Println("Written data to ", outputFilename)
}

func filenameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, path.Ext(filename))
}
