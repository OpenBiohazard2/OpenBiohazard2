package main

import (
	"./client"
	"./fileio"
	"./game"
	"./render"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
)

const (
	WINDOW_WIDTH         = 1024
	WINDOW_HEIGHT        = 768
	GAME_STATE_MAIN      = 0
	GAME_STATE_INVENTORY = 1
	STATE_CHANGE_DELAY   = 1 // in seconds
)

var (
	windowHandler           *client.WindowHandler
	gameDef                 *game.GameDef
	gameState               int
	lastTimeChangeState     float64
	debugEntities           []*render.DebugEntity
	cameraSwitchDebugEntity *render.DebugEntity
	itemEntities            []render.SceneMD1Entity
)

type PlayerModel struct {
	TextureId    uint32
	VertexBuffer []float32
	PLDOutput    *fileio.PLDOutput
}

func handleMainGameInput(gameDef *game.GameDef, collisionEntities []fileio.CollisionEntity) {
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
		if windowHandler.GetCurrentTime()-lastTimeChangeState >= STATE_CHANGE_DELAY {
			gameState = GAME_STATE_INVENTORY
			lastTimeChangeState = windowHandler.GetCurrentTime()
		}
	}
}

func handleInventory(renderDef *render.RenderDef) {
	if windowHandler.InputHandler.IsActive(client.PLAYER_VIEW_INVENTORY) {
		if windowHandler.GetCurrentTime()-lastTimeChangeState >= STATE_CHANGE_DELAY {
			gameState = GAME_STATE_MAIN
			lastTimeChangeState = windowHandler.GetCurrentTime()
		}
	}

	renderDef.RenderInventory()
}

func handleMainGame(roomOutput *fileio.RoomImageOutput,
	roomcutBinFilename string,
	roomcutBinOutput *fileio.BinOutput,
	renderDef *render.RenderDef,
	playerModel PlayerModel,
	spriteTextureIds [][]uint32) {
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

		debugEntities = make([]*render.DebugEntity, 0)
		debugEntities = append(debugEntities, render.NewDoorTriggerDebugEntity(gameDef.Doors))
		debugEntities = append(debugEntities, render.NewCollisionDebugEntity(gameDef.GameRoom.CollisionEntities))
		debugEntities = append(debugEntities, render.NewSlopedSurfacesDebugEntity(gameDef.GameRoom.CollisionEntities))
		debugEntities = append(debugEntities, render.NewItemTriggerDebugEntity(gameDef.Items))

		itemEntities = render.NewItemEntities(gameDef.Items, gameDef.GameRoom.ItemTextureData, gameDef.GameRoom.ItemModelData)
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

		cameraSwitchDebugEntity = render.NewCameraSwitchDebugEntity(gameDef.CameraId, gameDef.GameRoom.CameraSwitches, gameDef.GameRoom.CameraSwitchTransitions)

		gameDef.IsCameraLoaded = true
	}

	timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
	// Only render these entities for debugging
	debugEntities := render.DebugEntities{
		CameraSwitchDebugEntity: cameraSwitchDebugEntity,
		DebugEntities:           debugEntities,
	}
	// Update screen
	playerEntity := render.NewPlayerEntity(playerModel.TextureId, playerModel.VertexBuffer, playerModel.PLDOutput,
		gameDef.Player, gameDef.Player.PoseNumber)
	spriteEntity := render.SpriteEntity{
		TextureIds: spriteTextureIds,
		Sprites:    gameDef.Sprites,
	}

	renderDef.RenderFrame(playerEntity, itemEntities, debugEntities, spriteEntity, timeElapsedSeconds)

	handleMainGameInput(gameDef, gameDef.GameRoom.CollisionEntities)
	gameDef.HandleCameraSwitch(gameDef.Player.Position, gameDef.GameRoom.CameraSwitches, gameDef.GameRoom.CameraSwitchTransitions)
	gameDef.HandleRoomSwitch(gameDef.Player.Position)
	/*if gameDef.StageId == 1 && gameDef.RoomId == 0 {
		// for ROOM1000, start at function 1
		gameDef.RunScript(gameDef.GameRoom.RoomScriptData, timeElapsedSeconds, false, 1)
	} else {
		// start at function 0
		gameDef.RunScript(gameDef.GameRoom.RoomScriptData, timeElapsedSeconds, false, 0)
	}*/
}

func main() {
	// Run OpenGL code
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize glfw: %v", err))
	}
	defer glfw.Terminate()
	windowHandler = client.NewWindowHandler(WINDOW_WIDTH, WINDOW_HEIGHT, "OpenBiohazard2")

	renderDef := render.InitRenderer(WINDOW_WIDTH, WINDOW_HEIGHT)

	roomcutBinFilename := game.ROOMCUT_FILE
	roomcutBinOutput := fileio.LoadBINFile(roomcutBinFilename)

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

	gameDef = game.NewGame(1, 0, 0)
	gameDef.Player = game.NewPlayer(mgl32.Vec3{18781, 0, -2664}, 180)

	// Set game difficulty (0 is easy, 1 is normal)
	gameDef.SetBitArray(0, 25, game.DIFFICULTY_EASY)
	// Set camera id
	gameDef.SetScriptVariable(26, 0)

	var roomOutput *fileio.RoomImageOutput
	spriteTextureIds := make([][]uint32, 0)

	gameState = GAME_STATE_MAIN
	lastTimeChangeState = windowHandler.GetCurrentTime()

	inventoryImages, _ := fileio.LoadTIMImages(game.INVENTORY_FILE)
	renderDef.GenerateInventoryImageEntity(inventoryImages)
	for !windowHandler.ShouldClose() {
		windowHandler.StartFrame()

		if gameState == GAME_STATE_MAIN {
			handleMainGame(roomOutput, game.ROOMCUT_FILE, roomcutBinOutput, renderDef, playerModel, spriteTextureIds)
		} else if gameState == GAME_STATE_INVENTORY {
			handleInventory(renderDef)
		}
	}
}
