package script

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

var (
	FunctionName = map[byte]string{
		fileio.OP_NO_OP:            "NoOp",
		fileio.OP_EVT_END:          "EvtEnd",
		fileio.OP_EVT_NEXT:         "EvtNext",
		fileio.OP_EVT_CHAIN:        "EvtChain",
		fileio.OP_EVT_EXEC:         "EvtExec",
		fileio.OP_EVT_KILL:         "EvtKill",
		fileio.OP_IF_START:         "IfStart",
		fileio.OP_ELSE_START:       "ElseStart",
		fileio.OP_END_IF:           "EndIf",
		fileio.OP_SLEEP:            "Sleep",
		fileio.OP_SLEEPING:         "Sleeping",
		fileio.OP_WSLEEP:           "Wsleep",
		fileio.OP_WSLEEPING:        "Wsleeping",
		fileio.OP_FOR:              "ForStart",
		fileio.OP_FOR_END:          "ForEnd",
		fileio.OP_WHILE_START:      "WhileStart",
		fileio.OP_WHILE_END:        "WhileEnd",
		fileio.OP_DO_START:         "DoStart",
		fileio.OP_DO_END:           "DoEnd",
		fileio.OP_SWITCH:           "Switch",
		fileio.OP_CASE:             "Case",
		fileio.OP_DEFAULT:          "Default",
		fileio.OP_END_SWITCH:       "EndSwitch",
		fileio.OP_GOTO:             "Goto",
		fileio.OP_GOSUB:            "Gosub",
		fileio.OP_GOSUB_RETURN:     "GosubReturn",
		fileio.OP_BREAK:            "Break",
		fileio.OP_WORK_COPY:        "WorkCopy",
		fileio.OP_NO_OP2:           "NoOp2",
		fileio.OP_CHECK:            "CheckBit",
		fileio.OP_SET_BIT:          "SetBit",
		fileio.OP_COMPARE:          "Compare",
		fileio.OP_SAVE:             "Save",
		fileio.OP_COPY:             "Copy",
		fileio.OP_CALC:             "Calc",
		fileio.OP_CALC2:            "Calc2",
		fileio.OP_SCE_RND:          "SceRnd",
		fileio.OP_CUT_CHG:          "CutChg",
		fileio.OP_CUT_OLD:          "CutOld",
		fileio.OP_MESSAGE_ON:       "MessageOn",
		fileio.OP_AOT_SET:          "AotSet",
		fileio.OP_OBJ_MODEL_SET:    "ObjModelSet",
		fileio.OP_WORK_SET:         "WorkSet",
		fileio.OP_SPEED_SET:        "SpeedSet",
		fileio.OP_ADD_SPEED:        "AddSpeed",
		fileio.OP_ADD_ASPEED:       "AddAspeed",
		fileio.OP_POS_SET:          "PosSet",
		fileio.OP_DIR_SET:          "DirSet",
		fileio.OP_MEMBER_SET:       "MemberSet",
		fileio.OP_MEMBER_SET2:      "MemberSet2",
		fileio.OP_SE_ON:            "SeOn",
		fileio.OP_SCA_ID_SET:       "ScaIdSet",
		fileio.OP_DIR_CK:           "DirCk",
		fileio.OP_SCE_ESPR_ON:      "SceEsprOn",
		fileio.OP_DOOR_AOT_SET:     "DoorAotSet",
		fileio.OP_CUT_AUTO:         "CutAuto",
		fileio.OP_MEMBER_COPY:      "MemberCopy",
		fileio.OP_MEMBER_CMP:       "MemberCmp",
		fileio.OP_PLC_MOTION:       "PlcMotion",
		fileio.OP_PLC_DEST:         "PlcDest",
		fileio.OP_PLC_NECK:         "PlcNeck",
		fileio.OP_PLC_RET:          "PlcRet",
		fileio.OP_PLC_FLAG:         "PlcFlag",
		fileio.OP_SCE_EM_SET:       "SceEmSet",
		fileio.OP_AOT_RESET:        "AotReset",
		fileio.OP_AOT_ON:           "AotOn",
		fileio.OP_SUPER_SET:        "SuperSet",
		fileio.OP_CUT_REPLACE:      "CutReplace",
		fileio.OP_SCE_ESPR_KILL:    "SceEsprKill",
		fileio.OP_DOOR_MODEL_SET:   "DoorModelSet",
		fileio.OP_ITEM_AOT_SET:     "ItemAotSet",
		fileio.OP_SCE_TRG_CK:       "SceTrgCk",
		fileio.OP_SCE_BGM_CONTROL:  "SceBgmControl",
		fileio.OP_SCE_ESPR_CONTROL: "SceEsprControl",
		fileio.OP_SCE_FADE_SET:     "SceFadeSet",
		fileio.OP_SCE_ESPR3D_ON:    "SceEspr3dOn",
		fileio.OP_SCE_BGMTBL_SET:   "SceBgmTblSet",
		fileio.OP_PLC_ROT:          "PlcRot",
		fileio.OP_XA_ON:            "XaOn",
		fileio.OP_WEAPON_CHG:       "WeaponChg",
		fileio.OP_PLC_CNT:          "PlcCnt",
		fileio.OP_SCE_SHAKE_ON:     "SceShakeOn",
		fileio.OP_MIZU_DIV_SET:     "MizuDivSet",
		fileio.OP_KEEP_ITEM_CK:     "KeepItemCk",
		fileio.OP_XA_VOL:           "XaVol",
		fileio.OP_KAGE_SET:         "KageSet",
		fileio.OP_CUT_BE_SET:       "CutBeSet",
		fileio.OP_SCE_ITEM_LOST:    "SceItemLost",
		fileio.OP_PLC_GUN_EFF:      "PlcGunEff",
		fileio.OP_SCE_ESPR_ON2:     "SceEsprOn2",
		fileio.OP_SCE_ESPR_KILL2:   "SceEsprKill2",
		fileio.OP_PLC_STOP:         "PlcStop",
		fileio.OP_AOT_SET_4P:       "AotSet4P",
		fileio.OP_DOOR_AOT_SET_4P:  "DoorAotSet4P",
		fileio.OP_ITEM_AOT_SET_4P:  "ItemAotSet4P",
		fileio.OP_LIGHT_POS_SET:    "LightPosSet",
		fileio.OP_LIGHT_KIDO_SET:   "LightKidoSet",
		fileio.OP_RBJ_RESET:        "RbjReset",
		fileio.OP_SCE_SCR_MOVE:     "SceScrMove",
		fileio.OP_PARTS_SET:        "PartsSet",
		fileio.OP_MOVIE_ON:         "MovieOn",
		fileio.OP_SCE_PARTS_BOMB:   "ScePartsBomb",
		fileio.OP_SCE_PARTS_DOWN:   "ScePartsDown",
	}
)

func (scriptDef *ScriptDef) ScriptDebugFunction(threadNum int, lineBytes []byte) {
	if !scriptDebugEnabled {
		return
	}

	functionData := fmt.Sprintf("[ScriptThread %v] %s%s",
		threadNum, getFunctionNameFromOpcode(lineBytes[0]), showParameters(lineBytes))
	fmt.Println(functionData)
}

func (scriptDef *ScriptDef) ScriptDebugLine(line string) {
	if !scriptDebugEnabled {
		return
	}

	fmt.Println(line)
}

func getFunctionNameFromOpcode(opcode byte) string {
	return FunctionName[opcode]
}

func showParameters(lineBytes []byte) string {
	opcode := lineBytes[0]
	parameterString := "("
	switch opcode {
	case fileio.OP_GOSUB: // 0x18
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrGoSub{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Event=%d", instruction.Event)
	case fileio.OP_CHECK: // 0x21
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrCheckBitTest{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("BitArray=%d, ", instruction.BitArray)
		parameterString += fmt.Sprintf("BitNumber=%d, ", instruction.BitNumber)
		parameterString += fmt.Sprintf("Value=%d", instruction.Value)
	case fileio.OP_SET_BIT: // 0x22
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrSetBit{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("BitArray=%d, ", instruction.BitArray)
		parameterString += fmt.Sprintf("BitNumber=%d, ", instruction.BitNumber)
		parameterString += fmt.Sprintf("Operation=%d", instruction.Operation)
	case fileio.OP_CUT_CHG: // 0x29
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrCutChg{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("CameraId=%d", instruction.CameraId)
	case fileio.OP_AOT_SET: // 0x2c
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrAotSet{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Aot=%d, ", instruction.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Floor=%d, ", instruction.Floor)
		parameterString += fmt.Sprintf("Super=%d, ", instruction.Super)
		parameterString += fmt.Sprintf("X=%d, Z=%d, ", instruction.X, instruction.Z)
		parameterString += fmt.Sprintf("Width=%d, Depth=%d, ", instruction.Width, instruction.Depth)
		parameterString += fmt.Sprintf("Data=[%d,%d,%d,%d,%d,%d]", instruction.Data[0], instruction.Data[1], instruction.Data[2],
			instruction.Data[3], instruction.Data[4], instruction.Data[5])
	case fileio.OP_OBJ_MODEL_SET: // 0x2d
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrObjModelSet{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("ObjectIndex=%d, ", instruction.ObjectIndex)
		parameterString += fmt.Sprintf("ObjectId=%d, ", instruction.ObjectId)
		parameterString += fmt.Sprintf("Counter=%d, ", instruction.Counter)
		parameterString += fmt.Sprintf("Wait=%d, ", instruction.Wait)
		parameterString += fmt.Sprintf("Num=%d, ", instruction.Num)
		parameterString += fmt.Sprintf("Floor=%d, ", instruction.Floor)
		parameterString += fmt.Sprintf("Flag0=%d, ", instruction.Flag0)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Flag1=%d, ", instruction.Flag1)
		parameterString += fmt.Sprintf("Attribute=%d, ", instruction.Attribute)
		parameterString += fmt.Sprintf("Position=[%d, %d, %d], ", instruction.Position[0], instruction.Position[1], instruction.Position[2])
		parameterString += fmt.Sprintf("Direction=[%d, %d, %d], ", instruction.Direction[0], instruction.Direction[1], instruction.Direction[2])
		parameterString += fmt.Sprintf("Offset=[%d, %d, %d], ", instruction.Offset[0], instruction.Offset[1], instruction.Offset[2])
		parameterString += fmt.Sprintf("Dimensions=[%d, %d, %d]", instruction.Dimensions[0], instruction.Dimensions[1], instruction.Dimensions[2])
	case fileio.OP_POS_SET: // 0x32
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrPosSet{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Dummy=%d, ", instruction.Dummy)
		parameterString += fmt.Sprintf("X=%d, Y=%d, Z=%d", instruction.X, instruction.Y, instruction.Z)
	case fileio.OP_SCE_ESPR_ON: // 0x3a
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrSceEsprOn{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Dummy=%d, ", instruction.Dummy)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Work=%d, ", instruction.Work)
		parameterString += fmt.Sprintf("Unknown1=%d, ", instruction.Unknown1)
		parameterString += fmt.Sprintf("X=%d, Y=%d, Z=%d, ", instruction.X, instruction.Y, instruction.Z)
		parameterString += fmt.Sprintf("DirY=%d", instruction.DirY)
	case fileio.OP_DOOR_AOT_SET: // 0x3b
		byteArr := bytes.NewBuffer(lineBytes)
		door := fileio.ScriptInstrDoorAotSet{}
		binary.Read(byteArr, binary.LittleEndian, &door)
		parameterString += fmt.Sprintf("Aot=%d, ", door.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", door.Id)
		parameterString += fmt.Sprintf("Type=%d, ", door.Type)
		parameterString += fmt.Sprintf("Floor=%d, ", door.Floor)
		parameterString += fmt.Sprintf("Super=%d, ", door.Super)
		parameterString += fmt.Sprintf("X=%d, Z=%d, ", door.X, door.Z)
		parameterString += fmt.Sprintf("Width=%d, Depth=%d, ", door.Width, door.Depth)
		parameterString += fmt.Sprintf("NextX=%d, NextY=%d, ", door.NextX, door.NextY)
		parameterString += fmt.Sprintf("NextZ=%d, NextDir=%d, ", door.NextZ, door.NextDir)
		parameterString += fmt.Sprintf("Stage=%d, Room=%d, Camera=%d, ", door.Stage, door.Room, door.Camera)
		parameterString += fmt.Sprintf("NextFloor=%d, ", door.NextFloor)
		parameterString += fmt.Sprintf("TextureType=%d, ", door.TextureType)
		parameterString += fmt.Sprintf("DoorType=%d, ", door.DoorType)
		parameterString += fmt.Sprintf("KnockType=%d, ", door.KnockType)
		parameterString += fmt.Sprintf("KeyId=%d, ", door.KeyId)
		parameterString += fmt.Sprintf("KeyType=%d, ", door.KeyType)
		parameterString += fmt.Sprintf("Free=%d", door.Free)
	case fileio.OP_PLC_NECK: // 0x41
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrPlcNeck{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Operation=%d, ", instruction.Operation)
		parameterString += fmt.Sprintf("NeckX=%d, ", instruction.NeckX)
		parameterString += fmt.Sprintf("NeckY=%d, ", instruction.NeckY)
		parameterString += fmt.Sprintf("NeckZ=%d, ", instruction.NeckZ)
		parameterString += fmt.Sprintf("Unknown=[%d, %d]", instruction.Unknown[0], instruction.Unknown[1])
	case fileio.OP_SCE_EM_SET: // 0x44
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrSceEmSet{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Dummy=%d, ", instruction.Dummy)
		parameterString += fmt.Sprintf("Aot=%d, ", instruction.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Status=%d, ", instruction.Status)
		parameterString += fmt.Sprintf("Floor=%d, ", instruction.Floor)
		parameterString += fmt.Sprintf("SoundFlag=%d, ", instruction.SoundFlag)
		parameterString += fmt.Sprintf("ModelType=%d, ", instruction.ModelType)
		parameterString += fmt.Sprintf("EmSetFlag=%d, ", instruction.EmSetFlag)
		parameterString += fmt.Sprintf("X=%d, Y=%d, Z=%d, ", instruction.X, instruction.Y, instruction.Z)
		parameterString += fmt.Sprintf("DirY=%d, ", instruction.DirY)
		parameterString += fmt.Sprintf("Motion=%d, ", instruction.Motion)
		parameterString += fmt.Sprintf("CtrFlag=%d", instruction.CtrFlag)
	case fileio.OP_AOT_RESET: // 0x46
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrAotReset{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Aot=%d, ", instruction.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Data=[%d,%d,%d,%d,%d,%d]", instruction.Data[0], instruction.Data[1], instruction.Data[2],
			instruction.Data[3], instruction.Data[4], instruction.Data[5])
	case fileio.OP_ITEM_AOT_SET: // 0x4e
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrItemAotSet{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Aot=%d, ", instruction.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Floor=%d, ", instruction.Floor)
		parameterString += fmt.Sprintf("Super=%d, ", instruction.Super)
		parameterString += fmt.Sprintf("X=%d, Z=%d, ", instruction.X, instruction.Z)
		parameterString += fmt.Sprintf("Width=%d, Depth=%d, ", instruction.Width, instruction.Depth)
		parameterString += fmt.Sprintf("ItemId=%d, ", instruction.ItemId)
		parameterString += fmt.Sprintf("Amount=%d, ", instruction.Amount)
		parameterString += fmt.Sprintf("ItemPickedIndex=%d, ", instruction.ItemPickedIndex)
		parameterString += fmt.Sprintf("Md1ModelId=%d, ", instruction.Md1ModelId)
		parameterString += fmt.Sprintf("Act=%d", instruction.Act)
	case fileio.OP_SCE_BGM_CONTROL: // 0x51
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrSceBgmControl{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Operation=%d, ", instruction.Operation)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("LeftVolume=%d, ", instruction.LeftVolume)
		parameterString += fmt.Sprintf("RightVolume=%d", instruction.RightVolume)
	case fileio.OP_AOT_SET_4P: // 0x67
		byteArr := bytes.NewBuffer(lineBytes)
		instruction := fileio.ScriptInstrAotSet4p{}
		binary.Read(byteArr, binary.LittleEndian, &instruction)
		parameterString += fmt.Sprintf("Aot=%d, ", instruction.Aot)
		parameterString += fmt.Sprintf("Id=%d, ", instruction.Id)
		parameterString += fmt.Sprintf("Type=%d, ", instruction.Type)
		parameterString += fmt.Sprintf("Floor=%d, ", instruction.Floor)
		parameterString += fmt.Sprintf("Super=%d, ", instruction.Super)
		parameterString += fmt.Sprintf("X1=%d, Z1=%d, ", instruction.X1, instruction.Z1)
		parameterString += fmt.Sprintf("X2=%d, Z2=%d, ", instruction.X2, instruction.Z2)
		parameterString += fmt.Sprintf("X3=%d, Z3=%d, ", instruction.X3, instruction.Z3)
		parameterString += fmt.Sprintf("X4=%d, Z4=%d, ", instruction.X4, instruction.Z4)
		parameterString += fmt.Sprintf("Data=[%d,%d,%d,%d,%d,%d]", instruction.Data[0], instruction.Data[1], instruction.Data[2],
			instruction.Data[3], instruction.Data[4], instruction.Data[5])
	default:
		// Log each byte as its own parameter
		for i := 1; i < len(lineBytes); i++ {
			parameterString += fmt.Sprintf("%d, ", lineBytes[i])
		}
	}
	parameterString += ");"
	return parameterString
}
