package fileio

// .rdt - Room data

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
	OffsetItems                  uint32
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

type RDTItemOffsets struct {
	OffsetTexture uint32 // .tim file
	OffsetModel   uint32 // .md1 file
}

type RDTOutput struct {
	Header           RDTHeader
	RIDOutput        *RIDOutput // camera positions
	CameraSwitchData *RVDOutput
	LightData        *LITOutput
	CollisionData    *SCAOutput
	InitScriptData   *SCDOutput
	RoomScriptData   *SCDOutput
	SpriteOutput     *ESPOutput
	ItemTextureData  []*TIMOutput
	ItemModelData    []*MD1Output
}

func LoadRDTFile(filename string) (*RDTOutput, error) {
	rdtFile, _ := os.Open(filename)
	defer rdtFile.Close()

	if rdtFile == nil {
		log.Fatal("RDT file doesn't exist. Filename:", filename)
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

	// Read item models and textures
	itemTextureData := make([]*TIMOutput, rdtHeader.NumModels)
	itemModelData := make([]*MD1Output, rdtHeader.NumModels)
	if rdtHeader.NumModels > 0 {
		// Get the offsets
		offset := int64(offsets.OffsetItems)
		tempReader := io.NewSectionReader(r, offset, fileLength-offset)
		modelItemData := make([]RDTItemOffsets, rdtHeader.NumModels)
		if err := binary.Read(tempReader, binary.LittleEndian, &modelItemData); err != nil {
			log.Fatal("Error reading item model data ", err)
		}

		// Read item texture
		modelTextureLength := fileLength - int64(offsets.OffsetModelImage)
		for i := 0; i < int(rdtHeader.NumModels); i++ {
			timReader := io.NewSectionReader(r, int64(modelItemData[i].OffsetTexture), modelTextureLength)
			timOutput, err := LoadTIMStream(timReader, modelTextureLength)
			if err != nil {
				log.Fatal("Error reading item texture: ", err)
			}
			itemTextureData[i] = timOutput
		}

		// Read item model
		for i := 0; i < int(rdtHeader.NumModels); i++ {
			offset = int64(modelItemData[i].OffsetModel)
			// Invalid offset
			if offset == 0 {
				continue
			}
			modelLength := fileLength - offset
			timReader := io.NewSectionReader(r, offset, modelLength)
			md1Output, err := LoadMD1Stream(timReader, modelLength)
			if err != nil {
				log.Fatal("Error reading item model: ", err)
			}
			itemModelData[i] = md1Output
		}
	}

	offset := int64(offsets.OffsetLang1)
	if offset > 0 {
		lang1MsgReader := io.NewSectionReader(r, offset, fileLength-offset)
		LoadRDT_MSGStream(lang1MsgReader, fileLength)
	}

	// Script data
	// Run once when the level loads
	offset = int64(offsets.OffsetInitScript)
	initSCDReader := io.NewSectionReader(r, offset, fileLength-offset)
	initSCDOutput, err := LoadRDT_SCDStream(initSCDReader, fileLength)
	if err != nil {
		return nil, err
	}

	// Run during the game
	offset = int64(offsets.OffsetExecuteScript)
	roomSCDReader := io.NewSectionReader(r, offset, fileLength-offset)
	roomSCDOutput, err := LoadRDT_SCDStream(roomSCDReader, fileLength)
	if err != nil {
		return nil, err
	}

	// Sprite animations
	espOutput, err := LoadRDT_ESP(r, fileLength, rdtHeader, offsets)
	if err != nil {
		return nil, err
	}

	output := &RDTOutput{
		Header:           rdtHeader,
		RIDOutput:        ridOutput,
		CameraSwitchData: rvdOutput,
		LightData:        litOutput,
		CollisionData:    scaOutput,
		InitScriptData:   initSCDOutput,
		RoomScriptData:   roomSCDOutput,
		SpriteOutput:     espOutput,
		ItemTextureData:  itemTextureData,
		ItemModelData:    itemModelData,
	}
	return output, nil
}
