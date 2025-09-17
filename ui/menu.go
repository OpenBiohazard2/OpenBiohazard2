package ui

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/client"
)

type Menu struct {
	CurrentOption    int
	MaxOptions       int
	IsNewOption      bool // switches current option to a different option
	IsOptionSelected bool
}

func NewMenu(maxOptions int) *Menu {
	return &Menu{
		CurrentOption:    0,
		MaxOptions:       maxOptions,
		IsNewOption:      false,
		IsOptionSelected: false,
	}
}

func (menu *Menu) HandleMenuEvent(windowHandler *client.WindowHandler) {
	if windowHandler.InputHandler.IsActive(client.ACTION_BUTTON) {
		menu.IsOptionSelected = true
	} else {
		menu.IsOptionSelected = false
	}

	if windowHandler.InputHandler.IsActive(client.MENU_UP_BUTTON) {
		menu.CurrentOption--
		if menu.CurrentOption < 0 {
			menu.CurrentOption = 0
		}
		menu.IsNewOption = true
	} else if windowHandler.InputHandler.IsActive(client.MENU_DOWN_BUTTON) {
		menu.CurrentOption++
		if menu.CurrentOption >= menu.MaxOptions {
			menu.CurrentOption = menu.MaxOptions - 1
		}
		menu.IsNewOption = true
	} else {
		menu.IsNewOption = false
	}
}
