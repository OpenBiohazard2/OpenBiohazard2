package script

import (
	"fmt"
	"log"

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

func (scriptDef *ScriptDef) ScriptDebugFunction(threadNum int, functionIds []int, lineBytes []byte) {
	if !scriptDef.DebugEnabled {
		return
	}

	currentFunctionId := functionIds[len(functionIds)-1]
	functionData := fmt.Sprintf("[Thread %d][Function %d] %s%s",
		threadNum, currentFunctionId, getFunctionNameFromOpcode(lineBytes[0]), showParameters(lineBytes))
	log.Printf("SCRIPT-DEBUG: %s", functionData)
}

func (scriptDef *ScriptDef) ScriptDebugLine(line string) {
	if !scriptDef.DebugEnabled {
		return
	}

	log.Printf("SCRIPT-DEBUG: %s", line)
}

func getFunctionNameFromOpcode(opcode byte) string {
	return FunctionName[opcode]
}

func showParameters(lineBytes []byte) string {
	return GetOpcodeSignature(lineBytes)
}
