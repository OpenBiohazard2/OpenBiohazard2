package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type RDTHeader struct {
	NumSprites uint8
	NumCameras uint8
	NumModels  uint8
	NumItems   uint8
	NumDoors   uint8
	NumRooms   uint8
	NumReverb  uint8 // related to sound
	SpriteMax  uint8 // max number of .pri sprites used by one of the room's cameras
}

type RDTOffsets struct {
	OffsetRoomSound              uint32 // offset to room .snd sound table data
	OffsetRoomVABHeader          uint32 // .vh file
	OffsetRoomVABData            uint32 // .vb file
	OffsetEnemyVABHeader         uint32 // .vh file
	OffsetEnemyVABData           uint32 // .vb file
	OffsetOTA                    uint32
	OffsetCollisionData          uint32 // .sca file
	OffsetCameraPosition         uint32 // .rid file
	OffsetCameraSwitches         uint32 // .rvd file
	OffsetLights                 uint32 // .lit file
	OffsetModels                 uint32 // .md1 file
	OffsetFloors                 uint32 // .flr file
	OffsetBlocks                 uint32 // .blk file
	OffsetLang1                  uint32 // .msg file
	OffsetLang2                  uint32 // .msg file
	OffsetScrollTexture          uint32 // .tim file
	OffsetInitScript             uint32 // .scd file
	OffsetExecuteScript          uint32 // .scd file
	OffsetSpriteAnimations       uint32 // .esp file
	OffsetSpriteAnimationsOffset uint32 // .esp file
	OffsetSpriteImage            uint32 // .tim file
	OffsetModelImage             uint32 // .tim file
	OffsetRBJ                    uint32 // .rbj file
}

type RDTOutput struct {
	Header           RDTHeader
	RIDOutput        *RIDOutput // camera positions
	CameraSwitchData *RVDOutput
	LightData        *LITOutput
	CollisionData    *SCAOutput
}

func LoadRDTFile(filename string) (*RDTOutput, error) {
	rdtFile, _ := os.Open(filename)
	defer rdtFile.Close()

	if rdtFile == nil {
		log.Fatal("RDT file doesn't exist")
		return nil, fmt.Errorf("RDT file doesn't exist")
	}

	fi, err := rdtFile.Stat()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fileLength := fi.Size()
	return LoadRDT(rdtFile, fileLength)
}

func LoadRDT(r io.ReaderAt, fileLength int64) (*RDTOutput, error) {
	reader := io.NewSectionReader(r, int64(0), fileLength)

	rdtHeader := RDTHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &rdtHeader); err != nil {
		return nil, err
	}

	offsets := RDTOffsets{}
	if err := binary.Read(reader, binary.LittleEndian, &offsets); err != nil {
		return nil, err
	}

	// Camera position data
	ridOutput, err := LoadRDT_RID(r, fileLength, rdtHeader, offsets)
	if err != nil {
		return nil, err
	}

	// Camera switch data
	rvdOutput, err := LoadRDT_RVD(r, fileLength, rdtHeader, offsets)
	if err != nil {
		return nil, err
	}

	// Collision data
	scaOutput, err := LoadRDT_SCA(r, fileLength, rdtHeader, offsets)
	if err != nil {
		return nil, err
	}

	// Light data
	litOutput, err := LoadRDT_LIT(r, fileLength, rdtHeader, offsets)
	if err != nil {
		return nil, err
	}

	// Sprite animations
	LoadRDT_ESP(r, fileLength, rdtHeader, offsets)

	output := &RDTOutput{
		Header:           rdtHeader,
		RIDOutput:        ridOutput,
		CameraSwitchData: rvdOutput,
		LightData:        litOutput,
		CollisionData:    scaOutput,
	}
	return output, nil
}
