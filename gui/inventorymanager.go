package gui

// InventoryItem represents an item in the player's inventory
type InventoryItem struct {
	Id   int
	Num  int
	Size int
}

// InventoryManager manages inventory state and timing
type InventoryManager struct {
	totalInventoryTime        float64
	updateInventoryCursorTime float64 // milliseconds
	playerInventoryItems      []InventoryItem
}

// NewInventoryManager creates a new inventory manager
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		totalInventoryTime:        0,
		updateInventoryCursorTime: 30, // milliseconds
		playerInventoryItems:      initializeInventoryItems(),
	}
}

// initializeInventoryItems creates the initial inventory items
func initializeInventoryItems() []InventoryItem {
	items := make([]InventoryItem, 11)
	items[0] = InventoryItem{Id: 2, Num: 18, Size: 1}                  // hand gun
	items[1] = InventoryItem{Id: 1, Num: 1, Size: 0}                   // knife
	items[RESERVED_ITEM_SLOT] = InventoryItem{Id: 47, Num: 1, Size: 0} // lighter
	return items
}

// UpdateInventoryTime updates the total inventory time
func (im *InventoryManager) UpdateInventoryTime(timeElapsedSeconds float64) {
	im.totalInventoryTime += timeElapsedSeconds * 1000
}

// ShouldUpdateCursor returns true if the cursor should be updated
func (im *InventoryManager) ShouldUpdateCursor() bool {
	return im.totalInventoryTime >= im.updateInventoryCursorTime
}

// ResetInventoryTime resets the inventory time counter
func (im *InventoryManager) ResetInventoryTime() {
	im.totalInventoryTime = 0
}

// GetPlayerInventoryItems returns a copy of the player's inventory items
func (im *InventoryManager) GetPlayerInventoryItems() []InventoryItem {
	// Return a copy to prevent external modification
	items := make([]InventoryItem, len(im.playerInventoryItems))
	copy(items, im.playerInventoryItems)
	return items
}
