package render

func RenderEnemyEntity(r *RenderDef, enemyEntity EnemyEntity, timeElapsedSeconds float64) {
	if enemyEntity.EMDOutput == nil {
		return
	}

	// Only render debug placeholder
	if enemyEntity.DebugEntity != nil {
		RenderDebugEntities(r, []*DebugEntity{enemyEntity.DebugEntity})
	}
}

