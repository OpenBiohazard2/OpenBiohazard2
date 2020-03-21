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
	RenderRoom              render.RenderRoom
	PlayerEntity            *render.PlayerEntity
	SpriteTextureIds        [][]uint32
	DebugEntities           []*render.DebugEntity
	CameraSwitchDebugEntity *render.DebugEntity
	ItemEntities            []render.SceneMD1Entity
}

func NewMainGameStateInput(renderDef *render.RenderDef, gameDef *game.GameDef) *MainGameStateInput {
	return &MainGameStateInput{
		GameDef:        gameDef,
		ScriptDef:      script.NewScriptDef(),
		MainGameRender: NewMainGameRender(renderDef),
	}
}

func NewMainGameRender(renderDef *render.RenderDef) *MainGameRender {
	// Load player model
	pldOutput, err := fileio.LoadPLDFile(game.LEON_MODEL_FILE)
	if err != nil {
		log.Fatal("Error loading player model: ", err)
	}

	renderDef.SceneEntityMap[render.ENTITY_BACKGROUND_ID] = render.NewSceneEntity()
	// Camera image mask depends on updated camera position
	renderDef.SceneEntityMap[render.ENTITY_CAMERA_MASK_ID] = render.NewSceneEntity()

	return &MainGameRender{
		RenderDef:               renderDef,
		RoomOutput:              nil,
		RoomcutBinFilename:      game.ROOMCUT_FILE,
		RoomcutBinOutput:        fileio.LoadBINFile(game.ROOMCUT_FILE),
		PlayerEntity:            render.NewPlayerEntity(pldOutput),
		SpriteTextureIds:        make([][]uint32, 0),
		DebugEntities:           make([]*render.DebugEntity, 0),
		CameraSwitchDebugEntity: nil,
		ItemEntities:            make([]render.SceneMD1Entity, 0),
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
	gameDef.MaxCamerasInRoom = int(rdtOutput.Header.NumCameras)
	fmt.Println("Max cameras in room = ", gameDef.MaxCamerasInRoom)
	gameDef.GameRoom = gameDef.NewGameRoom(rdtOutput)
	mainGameRender.RenderRoom = render.NewRenderRoom(rdtOutput)

	// Sprite ids for rendering
	spriteTextureIds := mainGameRender.SpriteTextureIds
	spriteTextureIds = make([][]uint32, 0)
	for i := 0; i < len(mainGameRender.RenderRoom.SpriteData); i++ {
		spriteFrames := render.BuildSpriteTexture(mainGameRender.RenderRoom.SpriteData[i])
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
	mainGameRender.ItemEntities = render.NewItemEntities(gameDef.AotManager.Items,
		mainGameRender.RenderRoom.ItemTextureData, mainGameRender.RenderRoom.ItemModelData)
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
	renderDef.EnvironmentLight = render.BuildEnvironmentLight(mainGameRender.RenderRoom.LightData[gameDef.CameraId])

	// Update background image
	backgroundImageNumber := gameDef.GetBackgroundImageNumber()
	roomOutput = fileio.ExtractRoomBackground(roomcutBinFilename, roomcutBinOutput, backgroundImageNumber)

	if roomOutput.BackgroundImage != nil {
		backgroundColors := roomOutput.BackgroundImage.ConvertToRenderData()
		renderDef.SceneEntityMap[render.ENTITY_BACKGROUND_ID].UpdateBackgroundImageEntity(renderDef, backgroundColors)
		// Camera image mask depends on updated camera position
		cameraMasks := mainGameRender.RenderRoom.CameraMaskData[gameDef.CameraId]
		renderDef.SceneEntityMap[render.ENTITY_CAMERA_MASK_ID].UpdateCameraImageMaskEntity(renderDef, roomOutput, cameraMasks)
	}

	// Update camera switch zones
	cameraSwitchHandler := gameDef.GameRoom.CameraSwitchHandler
	mainGameRender.CameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId,
		cameraSwitchHandler.CameraSwitches, cameraSwitchHandler.CameraSwitchTransitions)
}

func runGameLoop(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	gameDef := mainGameStateInput.GameDef
	scriptDef := mainGameStateInput.ScriptDef
	mainGameRender := mainGameStateInput.MainGameRender
	renderDef := mainGameRender.RenderDef
	playerEntity := mainGameRender.PlayerEntity

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	// Only render these entities for debugging
	debugEntitiesRender := render.DebugEntities{
		CameraSwitchDebugEntity: mainGameRender.CameraSwitchDebugEntity,
		DebugEntities:           mainGameRender.DebugEntities,
	}
	// Update screen
	playerEntity.UpdatePlayerEntity(gameDef.Player, gameDef.Player.PoseNumber)
	spriteEntity := render.SpriteEntity{
		TextureIds: mainGameRender.SpriteTextureIds,
		Sprites:    gameDef.AotManager.Sprites,
	}

	renderDef.RenderFrame(*playerEntity, mainGameRender.ItemEntities, debugEntitiesRender, spriteEntity, timeElapsedSeconds)

	handleMainGameInput(gameDef, timeElapsedSeconds, gameDef.GameRoom.CollisionEntities, gameStateManager)
	gameDef.HandleCameraSwitch(gameDef.Player.Position)
	gameDef.HandleRoomSwitch(gameDef.Player.Position)
	aot := gameDef.AotManager.GetAotTriggerNearPlayer(gameDef.Player.Position)
	if aot != nil {
		if aot.Id == game.AOT_EVENT {
			threadNum := aot.Data[0]
			eventNum := aot.Data[3]
			lineData := []byte{fileio.OP_EVT_EXEC, threadNum, 0, eventNum}
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
