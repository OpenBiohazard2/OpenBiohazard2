package script

import (
	"bytes"
	"encoding/binary"

	"github.com/samuelyuan/openbiohazard2/fileio"
)

// PLC commands are used for 3D model animation

func (scriptDef *ScriptDef) ScriptPlcMotion(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcMotion{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

  // TODO: implement

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcDest(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcDest{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

  // TODO: implement

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcNeck(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcNeck{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

  // TODO: implement

	return 1
}
