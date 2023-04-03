package main

import (
	"fmt"
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/script"
	"github.com/OpenBiohazard2/OpenBiohazard2/world"
)

type MainGameStateInput struct {
	GameDef        *game.GameDef
	ScriptDef      *script.ScriptDef
	MainGameRender *MainGameRender
}

type MainGameRender struct {
	RenderDef               *render.RenderDef
	RoomcutBinOutput        *fileio.BinOutput
	RenderRoom              render.RenderRoom
	PlayerEntity            *render.PlayerEntity
	DebugEntities           []*render.DebugEntity
	CameraSwitchDebugEntity *render.DebugEntity
}

func NewMainGameStateInput(renderDef *render.RenderDef, gameDef *game.GameDef) *MainGameStateInput {
	scriptDef := script.NewScriptDef()
	// Set game difficulty (0 is easy, 1 is normal)
	scriptDef.SetBitArray(0, 25, game.DIFFICULTY_EASY)
	// Set camera id
	scriptDef.SetScriptVariable(26, 0)

	return &MainGameStateInput{
		GameDef:        gameDef,
		ScriptDef:      scriptDef,
		MainGameRender: NewMainGameRender(renderDef),
	}
}

func NewMainGameRender(renderDef *render.RenderDef) *MainGameRender {
	// Load player model
	pldOutput, err := fileio.LoadPLDFile(game.LEON_MODEL_FILE)
	if err != nil {
		log.Fatal("Error loading player model: ", err)
	}

	// Core sprite file has sprite ids 0-7
	// All other sprites are loaded based on the room
	fileio.LoadESPFile(game.CORE_SPRITE_FILE)

	return &MainGameRender{
		RenderDef:               renderDef,
		RoomcutBinOutput:        fileio.LoadBINFile(game.ROOMCUT_FILE),
		PlayerEntity:            render.NewPlayerEntity(pldOutput),
		DebugEntities:           make([]*render.DebugEntity, 0),
		CameraSwitchDebugEntity: nil,
	}
}

func handleMainGame(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	gameDef := mainGameStateInput.GameDef

	switch gameDef.StateStatus {
	case game.GAME_LOAD_ROOM:
		loadRoomState(mainGameStateInput)
		gameDef.StateStatus = game.GAME_LOAD_CAMERA
	case game.GAME_LOAD_CAMERA:
		loadCameraState(mainGameStateInput)
		gameDef.StateStatus = game.GAME_LOOP
	case game.GAME_LOOP:
		runGameLoop(mainGameStateInput, gameStateManager)
	}
}

func loadRoomState(mainGameStateInput *MainGameStateInput) {
	gameDef := mainGameStateInput.GameDef
	scriptDef := mainGameStateInput.ScriptDef
	mainGameRender := mainGameStateInput.MainGameRender
	renderDef := mainGameRender.RenderDef

	// Load room data from file
	roomFilename := gameDef.GetRoomFilename(game.PLAYER_LEON)
	rdtOutput, err := fileio.LoadRDTFile(roomFilename)
	if err != nil {
		log.Fatal("Error loading RDT file. ", err)
	}
	fmt.Println("Loaded", roomFilename)
	gameDef.RoomScript = gameDef.NewRoomScript(rdtOutput)
	gameDef.GameWorld.LoadNewRoom(rdtOutput)
	mainGameRender.RenderRoom = render.NewRenderRoom(rdtOutput)

	// Initialize room model objects
	renderDef.ItemGroupEntity.ItemTextureData = mainGameRender.RenderRoom.ItemTextureData
	renderDef.ItemGroupEntity.ItemModelData = mainGameRender.RenderRoom.ItemModelData

	// Initialize sprite textures
	renderDef.SpriteGroupEntity = render.NewSpriteGroupEntity(mainGameRender.RenderRoom.SpriteData)

	initScriptOnRoomLoad(scriptDef, gameDef, renderDef)

	mainGameRender.DebugEntities = render.BuildAllDebugEntities(gameDef.GameWorld)
}

func initScriptOnRoomLoad(scriptDef *script.ScriptDef, gameDef *game.GameDef, renderDef *render.RenderDef) {
	// Reset all state
	scriptDef.Reset()

	gameRoom := gameDef.RoomScript

	// Run initial script once when the room loads
	threadNum := 0
	functionNum := 0
	initScriptData := gameRoom.InitScriptData
	scriptDef.InitScript(initScriptData, threadNum, functionNum)
	scriptDef.RunScript(initScriptData, 10.0, gameDef, renderDef)

	// Initialize the room script to be run in the game loop
	threadNum = 0
	functionNum = 0
	roomScriptData := gameRoom.RoomScriptData
	scriptDef.InitScript(roomScriptData, threadNum, functionNum)
	threadNum = 1
	functionNum = 1
	scriptDef.InitScript(roomScriptData, threadNum, functionNum)
}

func loadCameraState(mainGameStateInput *MainGameStateInput) {
	gameDef := mainGameStateInput.GameDef
	mainGameRender := mainGameStateInput.MainGameRender

	updateCameraView(mainGameRender, gameDef)
	updateRoomBackroundImage(mainGameRender, gameDef)
	updateCameraSwitchZones(mainGameRender, gameDef)
}

func updateCameraView(mainGameRender *MainGameRender, gameDef *game.GameDef) {
	renderDef := mainGameRender.RenderDef

	// Update camera position
	cameraPosition := gameDef.GameWorld.GameRoom.CameraPositionData[gameDef.CameraId]
	renderDef.Camera.Update(cameraPosition.CameraFrom, cameraPosition.CameraTo, cameraPosition.CameraFov)
	renderDef.ViewMatrix = renderDef.Camera.BuildViewMatrix()

	// Update lighting
	renderDef.EnvironmentLight = render.BuildEnvironmentLight(mainGameRender.RenderRoom.LightData[gameDef.CameraId])
}

func updateRoomBackroundImage(mainGameRender *MainGameRender, gameDef *game.GameDef) {
	// Update background image
	roomOutput := fileio.ExtractRoomBackground(game.ROOMCUT_FILE, mainGameRender.RoomcutBinOutput, gameDef.GetBackgroundImageNumber())

	if roomOutput.BackgroundImage != nil {
		renderDef := mainGameRender.RenderDef
		render.UpdateTextureADT(renderDef.BackgroundImageEntity.TextureId, roomOutput.BackgroundImage)
		// Camera image mask depends on updated camera position
		cameraMasks := mainGameRender.RenderRoom.CameraMaskData[gameDef.CameraId]
		renderDef.CameraMaskEntity.UpdateCameraImageMaskEntity(renderDef, roomOutput, cameraMasks)
	}
}

func updateCameraSwitchZones(mainGameRender *MainGameRender, gameDef *game.GameDef) {
	cameraSwitchHandler := gameDef.GameWorld.GameRoom.CameraSwitchHandler
	mainGameRender.CameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId,
		cameraSwitchHandler.CameraSwitches, cameraSwitchHandler.CameraSwitchTransitions)
}

func runGameLoop(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	gameDef := mainGameStateInput.GameDef
	scriptDef := mainGameStateInput.ScriptDef
	mainGameRender := mainGameStateInput.MainGameRender
	renderDef := mainGameRender.RenderDef
	playerEntity := mainGameRender.PlayerEntity

	// Update screen
	playerEntity.UpdatePlayerEntity(gameDef.Player, gameDef.Player.PoseNumber)

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	// Only render these entities for debugging
	debugEntitiesRender := render.DebugEntities{
		CameraSwitchDebugEntity: mainGameRender.CameraSwitchDebugEntity,
		DebugEntities:           mainGameRender.DebugEntities,
	}
	renderDef.RenderFrame(*playerEntity, debugEntitiesRender, timeElapsedSeconds)

	handleMainGameInput(gameDef, timeElapsedSeconds, gameDef.GameWorld, gameStateManager)
	gameDef.HandleCameraSwitch(gameDef.Player.Position)
	gameDef.HandleRoomSwitch(gameDef.Player.Position)
	handleEventTrigger(scriptDef, gameDef)

	scriptDef.RunScript(gameDef.RoomScript.RoomScriptData, timeElapsedSeconds, gameDef, renderDef)
}

func handleEventTrigger(scriptDef *script.ScriptDef, gameDef *game.GameDef) {
	// Handles events like cutscenes
	aot := gameDef.GameWorld.AotManager.GetAotTriggerNearPlayer(gameDef.Player.Position)
	if aot != nil {
		if aot.Header.Id == world.AOT_EVENT {
			threadNum := aot.Data[0]
			eventNum := aot.Data[3]
			lineData := []byte{fileio.OP_EVT_EXEC, threadNum, 0, eventNum}
			scriptDef.ScriptEvtExec(lineData, gameDef.RoomScript.RoomScriptData)
		}
	}
}

func handleMainGameInput(
	gameDef *game.GameDef,
	timeElapsedSeconds float64,
	gameWorld *world.GameWorld,
	gameStateManager *GameStateManager,
) {
	collisionEntities := gameWorld.GameRoom.CollisionEntities

	if windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) {
		gameDef.Player.HandlePlayerInputForward(collisionEntities, timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.HandlePlayerInputBackward(collisionEntities, timeElapsedSeconds)
	}

	if !windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) &&
		!windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.PoseNumber = -1
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_LEFT) {
		gameDef.Player.RotatePlayerLeft(timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_RIGHT) {
		gameDef.Player.RotatePlayerRight(timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameDef.HandlePlayerActionButton(collisionEntities)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_INVENTORY)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}
