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
)

func main() {
	// Run OpenGL code
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize glfw: %v", err))
	}
	defer glfw.Terminate()
	windowHandler = client.NewWindowHandler(WINDOW_WIDTH, WINDOW_HEIGHT, "OpenBiohazard2")

	renderDef := render.InitRenderer(WINDOW_WIDTH, WINDOW_HEIGHT)

	gameDef := game.NewGame(1, 0, 0)
	gameDef.Player = game.NewPlayer(mgl32.Vec3{18781, 0, -2664}, 180)

	// Set game difficulty (0 is easy, 1 is normal)
	gameDef.SetBitArray(0, 25, game.DIFFICULTY_EASY)
	// Set camera id
	gameDef.SetScriptVariable(26, 0)

	gameStateManager := &GameStateManager{
		GameState:           GAME_STATE_MAIN_MENU,
		MainMenuOption:      0,
		LastTimeChangeState: windowHandler.GetCurrentTime(),
	}

	// Initialize main game
	mainGameStateInput := NewMainGameStateInput(renderDef, gameDef)

	// Initialize main menu
	mainMenuStateInput := NewMainMenuStateInput(renderDef)

	// Initialize inventory
	inventoryImages, _ := fileio.LoadTIMImages(game.INVENTORY_FILE)
	inventoryItemImages, _ := fileio.LoadTIMImages(game.ITEMALL_FILE)
	renderDef.GenerateInventoryImageEntity(inventoryImages, inventoryItemImages)

	for !windowHandler.ShouldClose() {
		windowHandler.StartFrame()

		switch gameStateManager.GameState {
		case GAME_STATE_MAIN_MENU:
			handleMainMenu(mainMenuStateInput, gameStateManager)
		case GAME_STATE_MAIN:
			handleMainGame(mainGameStateInput, gameStateManager)
		case GAME_STATE_INVENTORY:
			handleInventory(renderDef, gameStateManager)
		default:
			log.Fatal("Invalid game state: ", gameStateManager.GameState)
		}
	}
}
