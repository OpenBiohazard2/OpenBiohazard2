package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

type EnemyEntity struct {
	EnemyType    uint8
	ModelType    uint8
	Status       uint8
	Motion       uint16
	Position     mgl32.Vec3
	RotationY    float32
	EMDOutput    *fileio.EMDOutput
	DebugEntity  *DebugEntity
}

func NewEnemyEntity(emdOutput *fileio.EMDOutput) *EnemyEntity {
	return &EnemyEntity{
		EMDOutput: emdOutput,
	}
}

func (enemy *EnemyEntity) SetEnemyData(instruction fileio.ScriptInstrSceEmSet) {
	enemy.EnemyType = instruction.Type
	enemy.ModelType = instruction.ModelType
	enemy.Status = instruction.Status
	enemy.Motion = instruction.Motion
	enemy.Position = mgl32.Vec3{
		float32(instruction.X),
		float32(instruction.Y),
		float32(instruction.Z),
	}
	enemy.RotationY = float32(instruction.DirY) * (180.0 / 32768.0) // Convert to degrees
	
	enemy.DebugEntity = NewEnemyDebugEntity(enemy.Position, enemy.RotationY)
}

func (enemy *EnemyEntity) GetModelMatrix() mgl32.Mat4 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(enemy.Position.X(), enemy.Position.Y(), enemy.Position.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(enemy.RotationY)))
	return modelMatrix
}

