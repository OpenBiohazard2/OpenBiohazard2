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
		ege.EnemyEntities = append(ege.EnemyEntities[:index], ege.EnemyEntities[index+1:]...)
	}
}

func (ege *EnemyGroupEntity) ClearEnemies() {
	ege.EnemyEntities = make([]*EnemyEntity, 0)
}

func (ege *EnemyGroupEntity) GetEnemyCount() int {
	return len(ege.EnemyEntities)
}
