package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/samuelyuan/openbiohazard2/client"
	"github.com/samuelyuan/openbiohazard2/game"
	"github.com/samuelyuan/openbiohazard2/render"
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
	gameDef.Player = game.NewPlayer(game.DebugLocations[game.RoomMapKey{gameDef.StageId, gameDef.RoomId}], 180)

	// Set game difficulty (0 is easy, 1 is normal)
	gameDef.SetBitArray(0, 25, game.DIFFICULTY_EASY)
	// Set camera id
	gameDef.SetScriptVariable(26, 0)

	gameStateManager := NewGameStateManager()

	// Initialize main game
	mainGameStateInput := NewMainGameStateInput(renderDef, gameDef)

	// Initialize main menu
	mainMenuStateInput := &MainMenuStateInput{
		RenderDef: renderDef,
	}

	// Initialize inventory
	inventoryStateInput := NewInventoryStateInput(renderDef)

	for !windowHandler.ShouldClose() {
		windowHandler.StartFrame()

		switch gameStateManager.GameState {
		case GAME_STATE_MAIN_MENU:
			handleMainMenu(mainMenuStateInput, gameStateManager)
		case GAME_STATE_MAIN_GAME:
			handleMainGame(mainGameStateInput, gameStateManager)
		case GAME_STATE_INVENTORY:
			handleInventory(inventoryStateInput, gameStateManager)
		case GAME_STATE_LOAD_SAVE:
			handleLoadSave(renderDef, gameStateManager)
		case GAME_STATE_SPECIAL_MENU:
			handleSpecialMenu(mainMenuStateInput, gameStateManager)
		default:
			log.Fatal("Invalid game state: ", gameStateManager.GameState)
		}
	}
}
