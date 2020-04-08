package fileio

// .do2 file - Door file

import (
	"io"
	"log"
	"os"
)

type DO2Output struct {
	VABHeaderOutput *VABHeaderOutput
}

func LoadDO2File(filename string) *DO2Output {
	file, _ := os.Open(filename)
	defer file.Close()
	if file == nil {
		log.Fatal("DO2 file doesn't exist: ", filename)
		return nil
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fileLength := fi.Size()
	fileOutput, err := LoadDO2Stream(file, fileLength)
	if err != nil {
		log.Fatal("Failed to load DO2 file: ", err)
		return nil
	}

	return fileOutput
}

func LoadDO2Stream(r io.ReaderAt, fileLength int64) (*DO2Output, error) {
	vabHeaderReader := io.NewSectionReader(r, int64(16), fileLength)
	vabHeaderOutput, err := LoadVABHeaderStream(vabHeaderReader, fileLength)
	if err != nil {
		return nil, err
	}

	output := &DO2Output{
		VABHeaderOutput: vabHeaderOutput,
	}
	return output, nil
}
