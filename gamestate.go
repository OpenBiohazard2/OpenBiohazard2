package main

import (
	"./client"
	"./fileio"
	"./game"
	"./render"
	"./script"
	"fmt"
	"log"
)

type MainGameStateInput struct {
	GameDef        *game.GameDef
	ScriptDef      *script.ScriptDef
	MainGameRender *MainGameRender
}

type MainGameRender struct {
	RenderDef               *render.RenderDef
	RoomOutput              *fileio.RoomImageOutput
	RoomcutBinFilename      string
	RoomcutBinOutput        *fileio.BinOutput
	PlayerModel             PlayerModel
	SpriteTextureIds        [][]uint32
	DebugEntities           []*render.DebugEntity
	CameraSwitchDebugEntity *render.DebugEntity
	ItemEntities            []render.SceneMD1Entity
}

type PlayerModel struct {
	TextureId    uint32
	VertexBuffer []float32
	PLDOutput    *fileio.PLDOutput
}

func NewMainGameStateInput(renderDef *render.RenderDef, gameDef *game.GameDef) *MainGameStateInput {
	return &MainGameStateInput{
		GameDef:        gameDef,
		ScriptDef:      script.NewScriptDef(),
		MainGameRender: NewMainGameRender(renderDef),
	}
}

func NewMainGameRender(renderDef *render.RenderDef) *MainGameRender {
	return &MainGameRender{
		RenderDef:               renderDef,
		RoomOutput:              nil,
		RoomcutBinFilename:      game.ROOMCUT_FILE,
		RoomcutBinOutput:        fileio.LoadBINFile(game.ROOMCUT_FILE),
		PlayerModel:             NewPlayerModel(),
		SpriteTextureIds:        make([][]uint32, 0),
		DebugEntities:           make([]*render.DebugEntity, 0),
		CameraSwitchDebugEntity: nil,
		ItemEntities:            make([]render.SceneMD1Entity, 0),
	}
}

func NewPlayerModel() PlayerModel {
	// Load player model
	pldOutput, err := fileio.LoadPLDFile(game.LEON_MODEL_FILE)
	if err != nil {
		log.Fatal(err)
	}
	modelTexColors := pldOutput.TextureData.ConvertToRenderData()
	playerTextureId := render.BuildTexture(modelTexColors,
		int32(pldOutput.TextureData.ImageWidth), int32(pldOutput.TextureData.ImageHeight))
	playerEntityVertexBuffer := render.BuildEntityComponentVertices(pldOutput.MeshData, pldOutput.TextureData)
	return PlayerModel{
		TextureId:    playerTextureId,
		VertexBuffer: playerEntityVertexBuffer,
		PLDOutput:    pldOutput,
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

	// Load room data from file
	roomFilename := gameDef.GetRoomFilename(game.PLAYER_LEON)
	rdtOutput, err := fileio.LoadRDTFile(roomFilename)
	if err != nil {
		log.Fatal("Error loading RDT file. ", err)
	}
	fmt.Println("Loaded", roomFilename)
	gameDef.LoadNewRoom(rdtOutput)

	// Sprite ids for rendering
	spriteTextureIds := mainGameRender.SpriteTextureIds
	spriteTextureIds = make([][]uint32, 0)
	for i := 0; i < len(gameDef.RenderRoom.SpriteData); i++ {
		spriteFrames := render.BuildSpriteTexture(gameDef.RenderRoom.SpriteData[i])
		spriteTextureIds = append(spriteTextureIds, spriteFrames)
	}

	// Initialize scripts
	scriptDef.Reset()

	// Run initial script once when the room loads
	threadNum := 0
	functionNum := 0
	scriptDef.InitScript(gameDef.GameRoom.InitScriptData, threadNum, functionNum)
	scriptDef.RunScript(gameDef.GameRoom.InitScriptData, 10.0, gameDef)

	// Run the room script in the game loop
	threadNum = 0
	functionNum = 0
	scriptDef.InitScript(gameDef.GameRoom.RoomScriptData, threadNum, functionNum)
	threadNum = 1
	functionNum = 1
	scriptDef.InitScript(gameDef.GameRoom.RoomScriptData, threadNum, functionNum)

	mainGameRender.DebugEntities = render.BuildAllDebugEntities(gameDef)
	mainGameRender.ItemEntities = render.NewItemEntities(gameDef.AotManager.Items, gameDef.RenderRoom.ItemTextureData, gameDef.RenderRoom.ItemModelData)
}

func loadCameraState(mainGameStateInput *MainGameStateInput) {
	gameDef := mainGameStateInput.GameDef
	mainGameRender := mainGameStateInput.MainGameRender
	renderDef := mainGameRender.RenderDef
	roomOutput := mainGameRender.RoomOutput
	roomcutBinFilename := mainGameRender.RoomcutBinFilename
	roomcutBinOutput := mainGameRender.RoomcutBinOutput

	// Update camera position
	cameraPosition := gameDef.GameRoom.CameraPositionData[gameDef.CameraId]
	renderDef.Camera.CameraFrom = cameraPosition.CameraFrom
	renderDef.Camera.CameraTo = cameraPosition.CameraTo
	renderDef.Camera.CameraFov = cameraPosition.CameraFov
	renderDef.ViewMatrix = renderDef.Camera.GetViewMatrix()
	renderDef.SetEnvironmentLight(gameDef.RenderRoom.LightData[gameDef.CameraId])

	backgroundImageNumber := gameDef.GetBackgroundImageNumber()
	roomOutput = fileio.ExtractRoomBackground(roomcutBinFilename, roomcutBinOutput, backgroundImageNumber)

	if roomOutput.BackgroundImage != nil {
		render.GenerateBackgroundImageEntity(renderDef, roomOutput.BackgroundImage.ConvertToRenderData())
		// Camera image mask depends on updated camera position
		render.GenerateCameraImageMaskEntity(renderDef, roomOutput, gameDef.RenderRoom.CameraMaskData[gameDef.CameraId])
	}

	cameraSwitchHandler := gameDef.GameRoom.CameraSwitchHandler
	mainGameRender.CameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId,
		cameraSwitchHandler.CameraSwitches, cameraSwitchHandler.CameraSwitchTransitions)
}

func runGameLoop(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	gameDef := mainGameStateInput.GameDef
	scriptDef := mainGameStateInput.ScriptDef
	mainGameRender := mainGameStateInput.MainGameRender
	renderDef := mainGameRender.RenderDef
	playerModel := mainGameRender.PlayerModel

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	// Only render these entities for debugging
	debugEntitiesRender := render.DebugEntities{
		CameraSwitchDebugEntity: mainGameRender.CameraSwitchDebugEntity,
		DebugEntities:           mainGameRender.DebugEntities,
	}
	// Update screen
	playerEntity := render.NewPlayerEntity(playerModel.TextureId, playerModel.VertexBuffer, playerModel.PLDOutput,
		gameDef.Player, gameDef.Player.PoseNumber)
	spriteEntity := render.SpriteEntity{
		TextureIds: mainGameRender.SpriteTextureIds,
		Sprites:    gameDef.AotManager.Sprites,
	}

	renderDef.RenderFrame(playerEntity, mainGameRender.ItemEntities, debugEntitiesRender, spriteEntity, timeElapsedSeconds)

	handleMainGameInput(gameDef, timeElapsedSeconds, gameDef.GameRoom.CollisionEntities, gameStateManager)
	gameDef.HandleCameraSwitch(gameDef.Player.Position)
	gameDef.HandleRoomSwitch(gameDef.Player.Position)
	aot := gameDef.AotManager.GetAotTriggerNearPlayer(gameDef.Player.Position)
	if aot != nil {
		if aot.Id == game.AOT_EVENT {
			lineData := []byte{1, aot.Data[0], 0, aot.Data[3]}
			scriptDef.ScriptEvtExec(lineData, gameDef.GameRoom.RoomScriptData)
			// Only execute event once
			gameDef.AotManager.RemoveAotTrigger(int(aot.Aot))
		}
	}

	scriptDef.RunScript(gameDef.GameRoom.RoomScriptData, timeElapsedSeconds, gameDef)
}

func handleMainGameInput(gameDef *game.GameDef,
	timeElapsedSeconds float64,
	collisionEntities []fileio.CollisionEntity,
	gameStateManager *GameStateManager) {
	if windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) {
		gameDef.HandlePlayerInputForward(collisionEntities, timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.HandlePlayerInputBackward(collisionEntities, timeElapsedSeconds)
	}

	if !windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) &&
		!windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.PoseNumber = -1
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_LEFT) {
		gameDef.RotatePlayerLeft(timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_RIGHT) {
		gameDef.RotatePlayerRight(timeElapsedSeconds)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_INVENTORY)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}
