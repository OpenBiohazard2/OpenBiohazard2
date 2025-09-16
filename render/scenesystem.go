package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// SceneSystem manages all scene entities and rendering
type SceneSystem struct {
	SpriteGroupEntity     *SpriteGroupEntity
	BackgroundImageEntity *SceneEntity
	CameraMaskEntity      *SceneEntity
	ItemGroupEntity       *ItemGroupEntity
}

// NewSceneSystem creates a new scene system with all entities
func NewSceneSystem() *SceneSystem {
	return &SceneSystem{
		SpriteGroupEntity:     NewSpriteGroupEntity([]fileio.SpriteData{}),
		BackgroundImageEntity: NewBackgroundImageEntity(),
		CameraMaskEntity:      NewSceneEntity(),
		ItemGroupEntity:       NewItemGroupEntity(),
	}
}

// NewSceneSystemForTesting creates a SceneSystem without OpenGL dependencies for testing
func NewSceneSystemForTesting() *SceneSystem {
	return &SceneSystem{
		SpriteGroupEntity:     nil, // Skip OpenGL-dependent creation for tests
		BackgroundImageEntity: nil,
		CameraMaskEntity:      nil,
		ItemGroupEntity:       nil,
	}
}

// RenderBackground renders the background entities
func (ss *SceneSystem) RenderBackground(renderDef *RenderDef) {
	renderDef.RenderSceneEntity(ss.BackgroundImageEntity, RENDER_GAME_STATE_BACKGROUND_SOLID)
	renderDef.RenderSceneEntity(ss.CameraMaskEntity, RENDER_GAME_STATE_BACKGROUND_TRANSPARENT)
}

// RenderItems renders all item entities
func (ss *SceneSystem) RenderItems(renderDef *RenderDef) {
	for _, itemEntity := range ss.ItemGroupEntity.ModelObjectData {
		renderDef.RenderStaticEntity(*itemEntity, RENDER_TYPE_ITEM)
	}
}

// UpdateCameraMask updates the camera mask entity
func (ss *SceneSystem) UpdateCameraMask(renderDef *RenderDef, roomOutput *fileio.RoomImageOutput, masks []fileio.MaskRectangle) {
	ss.CameraMaskEntity.UpdateCameraImageMaskEntity(renderDef.ViewSystem, roomOutput, masks)
}
