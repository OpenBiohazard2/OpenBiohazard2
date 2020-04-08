package fileio

// .emd file - Enemy models

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

type EMDHeader struct {
	DirOffset uint32 // offset to directory
	DirCount  uint32 // number of objects in directory
}

type EMDOffsets struct {
	OffsetUnknown    uint32
	OffsetAnimation1 uint32 // .edd file
	OffsetSkeleton1  uint32 // .emr file
	OffsetAnimation2 uint32 // .edd file
	OffsetSkeleton2  uint32 // .emr file
	OffsetAnimation3 uint32 // .edd file
	OffsetSkeleton3  uint32 // .emr file
	OffsetMesh       uint32 // .md1 file
}

type EMDOutput struct {
	AnimationData1 *EDDOutput
	SkeletonData1  *EMROutput
	AnimationData2 *EDDOutput
	SkeletonData2  *EMROutput
	AnimationData3 *EDDOutput
	SkeletonData3  *EMROutput
	MeshData       *MD1Output
}

func LoadEMDFile(filename string) *EMDOutput {
	file, _ := os.Open(filename)
	defer file.Close()
	if file == nil {
		log.Fatal("EMD file doesn't exist: ", filename)
		return nil
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fileLength := fi.Size()
	fileOutput, err := LoadEMDStream(file, fileLength)
	if err != nil {
		log.Fatal("Failed to load EMD file: ", err)
		return nil
	}

	return fileOutput
}

func LoadEMDStream(r io.ReaderAt, fileLength int64) (*EMDOutput, error) {
	streamReader := io.NewSectionReader(r, int64(0), fileLength)

	emdHeader := EMDHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &emdHeader); err != nil {
		return nil, err
	}

	// Read the offset for each section
	offset := int64(emdHeader.DirOffset)
	offsetReader := io.NewSectionReader(r, offset, fileLength-offset)
	emdOffsets := EMDOffsets{}
	if err := binary.Read(offsetReader, binary.LittleEndian, &emdOffsets); err != nil {
		return nil, err
	}

	animationData1, err := loadAnimationData(r, fileLength, int64(emdOffsets.OffsetAnimation1))
	if err != nil {
		return nil, err
	}

	skeletonData1, err := loadSkeletonData(r, fileLength, int64(emdOffsets.OffsetSkeleton1), animationData1)
	if err != nil {
		return nil, err
	}

	animationData2, err := loadAnimationData(r, fileLength, int64(emdOffsets.OffsetAnimation2))
	if err != nil {
		return nil, err
	}

	skeletonData2, err := loadSkeletonData(r, fileLength, int64(emdOffsets.OffsetSkeleton2), animationData2)
	if err != nil {
		return nil, err
	}

	animationData3, err := loadAnimationData(r, fileLength, int64(emdOffsets.OffsetAnimation3))
	if err != nil {
		return nil, err
	}

	skeletonData3, err := loadSkeletonData(r, fileLength, int64(emdOffsets.OffsetSkeleton3), animationData3)
	if err != nil {
		return nil, err
	}

	meshData, err := loadMeshData(r, fileLength, int64(emdOffsets.OffsetMesh))
	if err != nil {
		return nil, err
	}

	output := &EMDOutput{
		AnimationData1: animationData1,
		SkeletonData1:  skeletonData1,
		AnimationData2: animationData2,
		SkeletonData2:  skeletonData2,
		AnimationData3: animationData3,
		SkeletonData3:  skeletonData3,
		MeshData:       meshData,
	}
	return output, nil
}
