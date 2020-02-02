package main

import (
	"./client"
	"./fileio"
	"./game"
	"./render"
	"fmt"
	"log"
)

const (
	GAME_STATE_MAIN_MENU = 0
	GAME_STATE_MAIN      = 1
	GAME_STATE_INVENTORY = 2

	STATE_CHANGE_DELAY = 1 // in seconds
)

type GameStateManager struct {
	GameState           int
	MainMenuOption      int
	LastTimeChangeState float64
}

type MainMenuStateInput struct {
	RenderDef                 *render.RenderDef
	MenuBackgroundImageOutput *fileio.ADTOutput
	MenuBackgroundTextImages  []*fileio.TIMOutput
}

type MainGameStateInput struct {
	RenderDef               *render.RenderDef
	GameDef                 *game.GameDef
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

func (gameStateManager *GameStateManager) CanUpdateGameState() bool {
	return windowHandler.GetCurrentTime()-gameStateManager.LastTimeChangeState >= STATE_CHANGE_DELAY
}

func (gameStateManager *GameStateManager) UpdateGameState(newGameState int) {
	gameStateManager.GameState = newGameState
}

func (gameStateManager *GameStateManager) UpdateLastTimeChangeState() {
	gameStateManager.LastTimeChangeState = windowHandler.GetCurrentTime()
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

func handleMainGame(mainGameStateInput *MainGameStateInput, gameStateManager *GameStateManager) {
	renderDef := mainGameStateInput.RenderDef
	gameDef := mainGameStateInput.GameDef
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

		spriteTextureIds = make([][]uint32, 0)
		for i := 0; i < len(gameDef.GameRoom.SpriteData); i++ {
			spriteFrames := render.BuildSpriteTexture(gameDef.GameRoom.SpriteData[i])
			spriteTextureIds = append(spriteTextureIds, spriteFrames)
		}
		gameDef.IsRoomLoaded = true

		debugEntities := make([]*render.DebugEntity, 0)
		debugEntities = append(debugEntities, render.NewDoorTriggerDebugEntity(gameDef.Doors))
		debugEntities = append(debugEntities, render.NewCollisionDebugEntity(gameDef.GameRoom.CollisionEntities))
		debugEntities = append(debugEntities, render.NewSlopedSurfacesDebugEntity(gameDef.GameRoom.CollisionEntities))
		debugEntities = append(debugEntities, render.NewItemTriggerDebugEntity(gameDef.Items))
		mainGameStateInput.DebugEntities = debugEntities

		mainGameStateInput.ItemEntities = render.NewItemEntities(gameDef.Items, gameDef.GameRoom.ItemTextureData, gameDef.GameRoom.ItemModelData)
	}

	if !gameDef.IsCameraLoaded {
		// Update camera position
		cameraPosition := gameDef.GameRoom.CameraPositionData[gameDef.CameraId]
		renderDef.Camera.CameraFrom = cameraPosition.CameraFrom
		renderDef.Camera.CameraTo = cameraPosition.CameraTo
		renderDef.Camera.CameraFov = cameraPosition.CameraFov
		renderDef.ViewMatrix = renderDef.Camera.GetViewMatrix()
		renderDef.SetEnvironmentLight(gameDef.GameRoom.LightData[gameDef.CameraId])

		backgroundImageNumber := gameDef.GetBackgroundImageNumber()
		roomOutput = fileio.ExtractRoomBackground(roomcutBinFilename, roomcutBinOutput, backgroundImageNumber)

		if roomOutput.BackgroundImage != nil {
			render.GenerateBackgroundImageEntity(renderDef, roomOutput.BackgroundImage.ConvertToRenderData())
			// Camera image mask depends on updated camera position
			render.GenerateCameraImageMaskEntity(renderDef, roomOutput, gameDef.GameRoom.CameraMaskData[gameDef.CameraId])
		}

		mainGameStateInput.CameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId, gameDef.GameRoom.CameraSwitches, gameDef.GameRoom.CameraSwitchTransitions)

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
		Sprites:    gameDef.Sprites,
	}

	renderDef.RenderFrame(playerEntity, mainGameStateInput.ItemEntities, debugEntitiesRender, spriteEntity, timeElapsedSeconds)

	handleMainGameInput(gameDef, gameDef.GameRoom.CollisionEntities, gameStateManager)
	gameDef.HandleCameraSwitch(gameDef.Player.Position, gameDef.GameRoom.CameraSwitches, gameDef.GameRoom.CameraSwitchTransitions)
	gameDef.HandleRoomSwitch(gameDef.Player.Position)
}

func handleMainGameInput(gameDef *game.GameDef,
	collisionEntities []fileio.CollisionEntity,
	gameStateManager *GameStateManager) {
	if windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) {
		gameDef.HandlePlayerInputForward(collisionEntities)
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.HandlePlayerInputBackward(collisionEntities)
	}

	if !windowHandler.InputHandler.IsActive(client.PLAYER_FORWARD) &&
		!windowHandler.InputHandler.IsActive(client.PLAYER_BACKWARD) {
		gameDef.Player.PoseNumber = -1
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_LEFT) {
		gameDef.Player.RotationAngle -= 5
		if gameDef.Player.RotationAngle < 0 {
			gameDef.Player.RotationAngle += 360
		}
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_ROTATE_RIGHT) {
		gameDef.Player.RotationAngle += 5
		if gameDef.Player.RotationAngle > 360 {
			gameDef.Player.RotationAngle -= 360
		}
	}

	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_INVENTORY)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}

func handleInventory(renderDef *render.RenderDef, gameStateManager *GameStateManager) {
	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	renderDef.RenderInventory()
}

func NewMainMenuStateInput(renderDef *render.RenderDef) *MainMenuStateInput {
	mainMenuStateInput := MainMenuStateInput{}
	mainMenuStateInput.RenderDef = renderDef
	mainMenuStateInput.MenuBackgroundImageOutput = fileio.LoadADTFile(game.MENU_IMAGE_FILE)
	mainMenuStateInput.MenuBackgroundTextImages, _ = fileio.LoadTIMImages(game.MENU_TEXT_FILE)
	mainMenuStateInput.RenderDef.GenerateMainMenuImageEntity(mainMenuStateInput.MenuBackgroundImageOutput, mainMenuStateInput.MenuBackgroundTextImages)
	return &mainMenuStateInput
}

func handleMainMenu(mainMenuStateInput *MainMenuStateInput, gameStateManager *GameStateManager) {
	renderDef := mainMenuStateInput.RenderDef
	menuBackgroundImageOutput := mainMenuStateInput.MenuBackgroundImageOutput
	menuBackgroundTextOutput := mainMenuStateInput.MenuBackgroundTextImages

	renderDef.RenderMainMenu()

	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		if gameStateManager.MainMenuOption == 1 {
			gameStateManager.UpdateGameState(GAME_STATE_MAIN)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.MainMenuOption--
			if gameStateManager.MainMenuOption < 0 {
				gameStateManager.MainMenuOption = 0
			}
			renderDef.UpdateMainMenu(menuBackgroundImageOutput, menuBackgroundTextOutput, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}

	if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		if gameStateManager.CanUpdateGameState() {
			gameStateManager.MainMenuOption++
			if gameStateManager.MainMenuOption > 2 {
				gameStateManager.MainMenuOption = 2
			}
			renderDef.UpdateMainMenu(menuBackgroundImageOutput, menuBackgroundTextOutput, gameStateManager.MainMenuOption)
			gameStateManager.UpdateLastTimeChangeState()
		}
	}
}
