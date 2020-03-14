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
	RenderDef               *render.RenderDef
	GameDef                 *game.GameDef
	ScriptDef               *script.ScriptDef
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
	// Load player model
	pldOutput, err := fileio.LoadPLDFile(game.LEON_MODEL_FILE)
	if err != nil {
		log.Fatal(err)
	}
	modelTexColors := pldOutput.TextureData.ConvertToRenderData()
	playerTextureId := render.BuildTexture(modelTexColors,
		int32(pldOutput.TextureData.ImageWidth), int32(pldOutput.TextureData.ImageHeight))
	playerEntityVertexBuffer := render.BuildEntityComponentVertices(pldOutput.MeshData, pldOutput.TextureData)
	playerModel := PlayerModel{
		TextureId:    playerTextureId,
		VertexBuffer: playerEntityVertexBuffer,
		PLDOutput:    pldOutput,
	}

	return &MainGameStateInput{
		RenderDef:               renderDef,
		GameDef:                 gameDef,
		ScriptDef:               script.NewScriptDef(),
		RoomOutput:              nil,
		RoomcutBinFilename:      game.ROOMCUT_FILE,
		RoomcutBinOutput:        fileio.LoadBINFile(game.ROOMCUT_FILE),
		PlayerModel:             playerModel,
		SpriteTextureIds:        make([][]uint32, 0),
		DebugEntities:           make([]*render.DebugEntity, 0),
		CameraSwitchDebugEntity: nil,
		ItemEntities:            make([]render.SceneMD1Entity, 0),
	}
}

func NewInventoryStateInput(renderDef *render.RenderDef) *InventoryStateInput {
	inventoryImages, _ := fileio.LoadTIMImages(game.INVENTORY_FILE)
	inventoryItemImages, _ := fileio.LoadTIMImages(game.ITEMALL_FILE)
	return &InventoryStateInput{
		RenderDef:           renderDef,
		InventoryImages:     inventoryImages,
		InventoryItemImages: inventoryItemImages,
	}
}

func handleMainGame(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	renderDef := mainGameStateInput.RenderDef
	gameDef := mainGameStateInput.GameDef
	scriptDef := mainGameStateInput.ScriptDef
	roomOutput := mainGameStateInput.RoomOutput
	roomcutBinFilename := mainGameStateInput.RoomcutBinFilename
	roomcutBinOutput := mainGameStateInput.RoomcutBinOutput
	playerModel := mainGameStateInput.PlayerModel
	spriteTextureIds := mainGameStateInput.SpriteTextureIds

	if !gameDef.IsRoomLoaded {
		roomFilename := gameDef.GetRoomFilename(game.PLAYER_LEON)
		rdtOutput, err := fileio.LoadRDTFile(roomFilename)
		if err != nil {
			log.Fatal("Error loading RDT file. ", err)
		}
		fmt.Println("Loaded", roomFilename)
		gameDef.LoadNewRoom(rdtOutput)

		// Sprite ids for rendering
		spriteTextureIds = make([][]uint32, 0)
		for i := 0; i < len(gameDef.RenderRoom.SpriteData); i++ {
			spriteFrames := render.BuildSpriteTexture(gameDef.RenderRoom.SpriteData[i])
			spriteTextureIds = append(spriteTextureIds, spriteFrames)
		}
		gameDef.IsRoomLoaded = true

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

		mainGameStateInput.DebugEntities = render.BuildAllDebugEntities(gameDef)

		mainGameStateInput.ItemEntities = render.NewItemEntities(gameDef.AotManager.Items, gameDef.RenderRoom.ItemTextureData, gameDef.RenderRoom.ItemModelData)
	}

	if !gameDef.IsCameraLoaded {
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
		mainGameStateInput.CameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId,
			cameraSwitchHandler.CameraSwitches, cameraSwitchHandler.CameraSwitchTransitions)

		gameDef.IsCameraLoaded = true
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	// Only render these entities for debugging
	debugEntitiesRender := render.DebugEntities{
		CameraSwitchDebugEntity: mainGameStateInput.CameraSwitchDebugEntity,
		DebugEntities:           mainGameStateInput.DebugEntities,
	}
	// Update screen
	playerEntity := render.NewPlayerEntity(playerModel.TextureId, playerModel.VertexBuffer, playerModel.PLDOutput,
		gameDef.Player, gameDef.Player.PoseNumber)
	spriteEntity := render.SpriteEntity{
		TextureIds: spriteTextureIds,
		Sprites:    gameDef.AotManager.Sprites,
	}

	renderDef.RenderFrame(playerEntity, mainGameStateInput.ItemEntities, debugEntitiesRender, spriteEntity, timeElapsedSeconds)

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
