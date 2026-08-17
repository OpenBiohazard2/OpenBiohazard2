package render

type EnemyGroupEntity struct {
	EnemyEntities []*EnemyEntity
}

func NewEnemyGroupEntity() *EnemyGroupEntity {
	return &EnemyGroupEntity{
		EnemyEntities: make([]*EnemyEntity, 0),
	}
}

func (ege *EnemyGroupEntity) AddEnemy(enemy *EnemyEntity) {
	ege.EnemyEntities = append(ege.EnemyEntities, enemy)
}

func (ege *EnemyGroupEntity) RemoveEnemy(index int) {
	if index >= 0 && index < len(ege.EnemyEntities) {
		deleteEnemyGPUResources(ege.EnemyEntities[index])
		ege.EnemyEntities = append(ege.EnemyEntities[:index], ege.EnemyEntities[index+1:]...)
	}
}

// ClearEnemies removes all enemies, releasing each one's GPU resources first.
func (ege *EnemyGroupEntity) ClearEnemies() {
	for _, enemy := range ege.EnemyEntities {
		deleteEnemyGPUResources(enemy)
	}
	ege.EnemyEntities = make([]*EnemyEntity, 0)
}

func deleteEnemyGPUResources(enemy *EnemyEntity) {
	if enemy != nil && enemy.DebugEntity != nil {
		enemy.DebugEntity.Delete()
	}
}

func (ege *EnemyGroupEntity) GetEnemyCount() int {
	return len(ege.EnemyEntities)
}
