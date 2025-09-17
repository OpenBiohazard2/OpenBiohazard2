package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/OpenBiohazard2/OpenBiohazard2/client"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/state"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui"
	"github.com/OpenBiohazard2/OpenBiohazard2/ui_render"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	WINDOW_WIDTH  = 1024
	WINDOW_HEIGHT = 768
)

func main() {
	fmt.Println("Validating game folders exist...")
	if err := game.ValidateFilesExist(); err != nil {
		log.Fatal("File validation failed: ", err)
	}
	fmt.Println("Validated game folders exist")

	// Run OpenGL code
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}
	defer glfw.Terminate()
	windowHandler := client.NewWindowHandler(WINDOW_WIDTH, WINDOW_HEIGHT, "OpenBiohazard2")

	// Initialize game components
	renderDef, gameDef, gameStateManager := initializeGame()

	// Create all state inputs
	stateInputs := createStateInputs(renderDef, gameDef)

	// Run the main game loop
	runMainGameLoop(windowHandler, gameStateManager, stateInputs, renderDef)
}

// initializeGame sets up the core game components
func initializeGame() (*render.RenderDef, *game.GameDef, *state.GameStateManager) {
	renderDef := render.InitRenderer(WINDOW_WIDTH, WINDOW_HEIGHT)

	gameDef := game.NewGame(1, 0, 0)
	gameDef.Player = game.NewPlayer(game.DebugLocations[game.RoomMapKey{StageId: gameDef.StageId, RoomId: gameDef.RoomId}], 180)

	gameStateManager := state.NewGameStateManager()

	return renderDef, gameDef, gameStateManager
}

// createStateInputs initializes all game state input handlers
func createStateInputs(renderDef *render.RenderDef, gameDef *game.GameDef) map[string]interface{} {
	return map[string]interface{}{
		"mainGame": state.NewMainGameStateInput(renderDef, gameDef),
		"mainMenu": &state.MainMenuStateInput{
			RenderDef:  renderDef,
			UIRenderer: ui_render.NewUIRenderer(renderDef),
			Menu:       ui.NewMenu(4),
		},
		"specialMenu": &state.MainMenuStateInput{
			RenderDef:  renderDef,
			UIRenderer: ui_render.NewUIRenderer(renderDef),
			Menu:       ui.NewMenu(2),
		},
		"inventory": state.NewInventoryStateInput(renderDef),
	}
}

// runMainGameLoop handles the main game loop and state management
func runMainGameLoop(windowHandler *client.WindowHandler, gameStateManager *state.GameStateManager, stateInputs map[string]interface{}, renderDef *render.RenderDef) {
	for !windowHandler.ShouldClose() {
		windowHandler.StartFrame()

		switch gameStateManager.GameState {
		case state.GAME_STATE_MAIN_MENU:
			state.HandleMainMenu(stateInputs["mainMenu"].(*state.MainMenuStateInput), gameStateManager, windowHandler)
		case state.GAME_STATE_MAIN_GAME:
			state.HandleMainGame(stateInputs["mainGame"].(*state.MainGameStateInput), gameStateManager, windowHandler)
		case state.GAME_STATE_INVENTORY:
			state.HandleInventory(stateInputs["inventory"].(*state.InventoryStateInput), gameStateManager, windowHandler)
		case state.GAME_STATE_LOAD_SAVE:
			state.HandleLoadSave(renderDef, gameStateManager, windowHandler)
		case state.GAME_STATE_SPECIAL_MENU:
			state.HandleSpecialMenu(stateInputs["specialMenu"].(*state.MainMenuStateInput), gameStateManager, windowHandler)
		default:
			log.Fatal("Invalid game state: ", gameStateManager.GameState)
		}
	}
}
