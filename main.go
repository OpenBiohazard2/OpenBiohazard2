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
	WINDOW_WIDTH  = 1024
	WINDOW_HEIGHT = 768
)

var (
	windowHandler *client.WindowHandler
	gameDef       *game.GameDef
)

func handleInput(gameDef *game.GameDef, collisionEntities []fileio.CollisionEntity) {
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
	playerEntityVertexBuffer := render.BuildEntityComponentVertices(pldOutput)

	gameDef = game.NewGame(1, 0, 0)
	gameDef.Player = game.NewPlayer(mgl32.Vec3{18781, 0, -2664}, 180)

	// Set game difficulty (0 is easy, 1 is normal)
	gameDef.SetBitArray(0, 25, game.DIFFICULTY_EASY)
	// Set camera id
	gameDef.SetScriptVariable(26, 0)
	// Fire animation for ROOM1000
	gameDef.SetBitArray(5, 5, 0)
	gameDef.SetBitArray(5, 6, 0)
	gameDef.SetBitArray(5, 7, 0)

	// Unknown
	gameDef.SetBitArray(5, 12, 1)
	gameDef.SetBitArray(5, 14, 1)
	gameDef.SetBitArray(6, 35, 1)
	gameDef.SetBitArray(6, 36, 1)
	gameDef.SetBitArray(6, 37, 1)
	gameDef.SetBitArray(6, 38, 1)
	gameDef.SetBitArray(6, 39, 1)
	gameDef.SetBitArray(6, 40, 1)

	var roomOutput *fileio.RoomImageOutput
	spriteTextureIds := make([][]uint32, 0)

	for !windowHandler.ShouldClose() {
		windowHandler.StartFrame()

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
		}

		if !gameDef.IsCameraLoaded {
			// Update camera position
			cameraPosition := gameDef.GameRoom.CameraPositionData[gameDef.CameraId]
			renderDef.Camera.CameraFrom = cameraPosition.CameraFrom
			renderDef.Camera.CameraTo = cameraPosition.CameraTo
			renderDef.ViewMatrix = renderDef.Camera.GetViewMatrix()
			renderDef.SetEnvironmentLight(gameDef.GameRoom.LightData[gameDef.CameraId])

			backgroundImageNumber := gameDef.GetBackgroundImageNumber()
			roomOutput = fileio.ExtractRoomBackground(roomcutBinFilename, roomcutBinOutput, backgroundImageNumber)

			if roomOutput.BackgroundImage != nil {
				render.GenerateBackgroundImageEntity(renderDef, roomOutput.BackgroundImage.ConvertToRenderData())
				// Camera image mask depends on updated camera position
				render.GenerateCameraImageMaskEntity(renderDef, roomOutput, gameDef.GameRoom.CameraMaskData[gameDef.CameraId])
			}

			gameDef.IsCameraLoaded = true
		}

		timeElapsedSeconds := windowHandler.GetTimeSinceLastFrame()
		// Only render these entities for debugging
		debugEntities := render.DebugEntities{
			CameraId:                gameDef.CameraId,
			CameraSwitches:          gameDef.GameRoom.CameraSwitches,
			CameraSwitchTransitions: gameDef.GameRoom.CameraSwitchTransitions,
			CollisionEntities:       gameDef.GameRoom.CollisionEntities,
			Doors:                   gameDef.Doors,
		}
		// Update screen
		playerEntity := render.PlayerEntity{
			TextureId:           playerTextureId,
			VertexBuffer:        playerEntityVertexBuffer,
			PLDOutput:           pldOutput,
			Player:              gameDef.Player,
			AnimationPoseNumber: gameDef.Player.PoseNumber,
		}

		spriteEntity := render.SpriteEntity{
			TextureIds: spriteTextureIds,
			Sprites:    gameDef.Sprites,
		}

		renderDef.RenderFrame(playerEntity, debugEntities, spriteEntity, timeElapsedSeconds)

		handleInput(gameDef, gameDef.GameRoom.CollisionEntities)
		gameDef.HandleCameraSwitch(gameDef.Player.Position, gameDef.GameRoom.CameraSwitches, gameDef.GameRoom.CameraSwitchTransitions)
		gameDef.HandleRoomSwitch(gameDef.Player.Position)
		if gameDef.StageId == 1 && gameDef.RoomId == 0 {
			// for ROOM1000, start at function 1
			gameDef.RunScript(gameDef.GameRoom.RoomScriptData, timeElapsedSeconds, false, 1)
		} else {
			// start at function 0
			gameDef.RunScript(gameDef.GameRoom.RoomScriptData, timeElapsedSeconds, false, 0)
		}
	}
}
