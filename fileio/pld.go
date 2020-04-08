package fileio

// .pld - Player models

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

type PLDHeader struct {
	DirOffset uint32 // offset to directory
	DirCount  uint32 // number of objects in directory
}

type PLDOffsets struct {
	OffsetAnimation uint32 // .edd file
	OffsetSkeleton  uint32 // .emr file
	OffsetMesh      uint32 // .md1 file
	OffsetTexture   uint32 // .tim file
}

type PLDOutput struct {
	AnimationData *EDDOutput
	SkeletonData  *EMROutput
	MeshData      *MD1Output
	TextureData   *TIMOutput
}

func LoadPLDFile(filename string) (*PLDOutput, error) {
	file, _ := os.Open(filename)
	defer file.Close()
	if file == nil {
		log.Fatal("PLD file doesn't exist:", filename)
		return nil, nil
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileLength := fi.Size()
	return LoadPLDStream(file, fileLength)
}

func LoadPLDStream(r io.ReaderAt, fileLength int64) (*PLDOutput, error) {
	reader := io.NewSectionReader(r, int64(0), fileLength)

	pldHeader := PLDHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &pldHeader); err != nil {
		return nil, err
	}

	// Read the offset for each section
	offset := int64(pldHeader.DirOffset)
	reader = io.NewSectionReader(r, offset, fileLength-offset)
	pldOffsets := PLDOffsets{}
	if err := binary.Read(reader, binary.LittleEndian, &pldOffsets); err != nil {
		return nil, err
	}

	animationData, err := loadAnimationData(r, fileLength, int64(pldOffsets.OffsetAnimation))
	if err != nil {
		return nil, err
	}

	skeletonData, err := loadSkeletonData(r, fileLength, int64(pldOffsets.OffsetSkeleton), animationData)
	if err != nil {
		return nil, err
	}

	meshData, err := loadMeshData(r, fileLength, int64(pldOffsets.OffsetMesh))
	if err != nil {
		return nil, err
	}

	timOutput, err := loadTexture(r, fileLength, int64(pldOffsets.OffsetTexture))
	if err != nil {
		return nil, err
	}

	pldOutput := &PLDOutput{
		AnimationData: animationData,
		SkeletonData:  skeletonData,
		MeshData:      meshData,
		TextureData:   timOutput,
	}
	return pldOutput, nil
}

func loadAnimationData(fileReader io.ReaderAt, fileLength int64, offset int64) (*EDDOutput, error) {
	eddReader := io.NewSectionReader(fileReader, offset, fileLength-offset)
	return LoadEDDStream(eddReader, fileLength-offset)
}

func loadSkeletonData(fileReader io.ReaderAt, fileLength int64, offset int64, animationData *EDDOutput) (*EMROutput, error) {
	emrReader := io.NewSectionReader(fileReader, offset, fileLength-offset)
	return LoadEMRStream(emrReader, fileLength-offset, animationData)
}

func loadMeshData(fileReader io.ReaderAt, fileLength int64, offset int64) (*MD1Output, error) {
	md1Reader := io.NewSectionReader(fileReader, offset, fileLength-offset)
	return LoadMD1Stream(md1Reader, fileLength-offset)
}

func loadTexture(fileReader io.ReaderAt, fileLength int64, offset int64) (*TIMOutput, error) {
	TIMReader := io.NewSectionReader(fileReader, offset, fileLength-int64(offset))
	return LoadTIMStream(TIMReader, fileLength-int64(offset))
}
