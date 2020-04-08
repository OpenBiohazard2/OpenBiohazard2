package fileio

// .vab - Playstation 1 audio format

import (
	"encoding/binary"
	"io"
	"log"
)

type VABHeader struct {
	Magic          [4]byte // "pBAV"
	Version        uint32  // format version
	Id             uint32  // vab id
	Filesize       uint32  // waveform size in bytes
	Reserved0      uint16  // system reserved value
	ProgramCount   uint16  // total number of programs used
	ToneCount      uint16  // total number of tones used
	WaveformCount  uint16  // total number of waveforms (.vag files) used
	MasterVolume   uint8   // master volume used
	MasterPan      uint8   // master pan used
	BankAttribute1 uint8   // user defined attribute of bank 1
	BankAttribute2 uint8   // user defined attribute of bank 2
	Reserved1      uint32  // system reserved value
}

type VABProgram struct {
	Tones     uint8 // number of effective tones which compose the program
	Volume    uint8
	Priority  uint8
	Mode      uint8
	Pan       uint8
	Reserved0 uint8
	Attribute int16
	Reserved1 [2]uint32
}

type VABTone struct {
	Priority        uint8 // tone priority (0 - 127); used for controlling allocation when more voices than can be keyed on are requested
	Mode            uint8 // tone mode (0 = normal; 4 = reverb applied)
	Volume          uint8 // tone volume
	Pan             uint8 // tone pan
	Center          uint8 // center note (0~127)
	Shift           uint8 // pitch correction (0~127,cent units)
	NoteMin         uint8 // minimum note limit (0~127)
	NoteMax         uint8 // maximum note limit (0~127, provided min < max)
	VibratoWidth    uint8 // vibrato width (1/128 rate, 0~127)
	VibratoTime     uint8 // 1 cycle time of vibrato (tick units)
	PortamentoWidth uint8 // portamento width (1/128 rate, 0~127)
	PortamentoTime  uint8 // portamento holding time (tick units)
	PitchBendMin    uint8 // pitch bend (-0~127, 127 = 1 octave)
	PitchBendMax    uint8 // pitch bend (+0~127, 127 = 1 octave)
	Reserved1       uint8
	Reserved2       uint8
	Adsr1           uint16
	Adsr2           uint16
	Program         int16 // parent program
	Vag             int16 // waveform (VAG) used
	Reserved3       [4]int16
}

type VABHeaderOutput struct {
	VABHeader  VABHeader
	AudioSizes []uint16
	NumBytes   int
}

type VABDataOutput struct {
	RawADPCMData [][]uint8
}

func LoadVABHeaderStream(r io.ReaderAt, fileLength int64) (*VABHeaderOutput, error) {
	vabHeaderReader := io.NewSectionReader(r, int64(0), fileLength)

	vabHeader := VABHeader{}
	if err := binary.Read(vabHeaderReader, binary.LittleEndian, &vabHeader); err != nil {
		return nil, err
	}

	if string(vabHeader.Magic[:]) != "pBAV" {
		log.Fatal("VAB header is invalid: ", vabHeader.Magic)
	}
	if vabHeader.ProgramCount > 128 {
		log.Fatal("Too many programs: ", vabHeader.ProgramCount)
	}

	programData := make([]VABProgram, 128)
	if err := binary.Read(vabHeaderReader, binary.LittleEndian, &programData); err != nil {
		return nil, err
	}

	for i := 0; i < int(vabHeader.ProgramCount); i++ {
		tones := make([]VABTone, 16)
		if err := binary.Read(vabHeaderReader, binary.LittleEndian, &tones); err != nil {
			return nil, err
		}
	}

	audioSizes := make([]uint16, vabHeader.WaveformCount+1)
	if err := binary.Read(vabHeaderReader, binary.LittleEndian, &audioSizes); err != nil {
		return nil, err
	}

	totalProgramSize := 128 * 32
	totalToneSize := int(vabHeader.ProgramCount) * 16 * 32
	totalWaveformSize := int(vabHeader.WaveformCount+1) * 2
	totalVabHeaderSize := totalProgramSize + totalToneSize + totalWaveformSize
	vabHeaderOutput := &VABHeaderOutput{
		VABHeader:  vabHeader,
		AudioSizes: audioSizes,
		NumBytes:   totalVabHeaderSize,
	}
	return vabHeaderOutput, nil
}

func LoadVABDataStream(r io.ReaderAt, fileLength int64, vabHeaderOutput *VABHeaderOutput) (*VABDataOutput, error) {
	vabDataReader := io.NewSectionReader(r, int64(0), fileLength)

	rawADPCMData := make([][]uint8, 0)
	for i := 0; i < len(vabHeaderOutput.AudioSizes); i++ {
		rawAudioSize := int(vabHeaderOutput.AudioSizes[i])
		if rawAudioSize == 0 {
			continue
		}

		adpcmData := make([]byte, rawAudioSize*8)
		if err := binary.Read(vabDataReader, binary.LittleEndian, &adpcmData); err != nil {
			return nil, err
		}
		rawADPCMData = append(rawADPCMData, adpcmData)
	}

	vabDataOutput := &VABDataOutput{
		RawADPCMData: rawADPCMData,
	}
	return vabDataOutput, nil
}
