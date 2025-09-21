package render

import (
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/geometry"
	"github.com/OpenBiohazard2/OpenBiohazard2/world"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	RENDER_TYPE_DEBUG = -1
)

// Debug entity colors
var (
	DEBUG_COLOR_RED    = [4]float32{1.0, 0.0, 0.0, 0.3} // Collision entities
	DEBUG_COLOR_GREEN  = [4]float32{0.0, 1.0, 0.0, 0.3} // Camera switches
	DEBUG_COLOR_BLUE   = [4]float32{0.0, 0.0, 1.0, 0.3} // Door triggers
	DEBUG_COLOR_CYAN   = [4]float32{0.0, 1.0, 1.0, 0.3} // Item triggers and AOT triggers
	DEBUG_COLOR_MAGENTA = [4]float32{1.0, 0.0, 1.0, 0.3} // Sloped surfaces
)

type DebugEntity struct {
	Color              [4]float32
	VertexBuffer       []float32
	VertexArrayObject  uint32
	VertexBufferObject uint32
}

func RenderCameraSwitches(r *RenderDef, cameraSwitchDebugEntity *DebugEntity) {
	// Use ShaderSystem method for better performance
	r.ShaderSystem.SetRenderType(RENDER_TYPE_DEBUG)

	RenderDebugEntities(r, []*DebugEntity{cameraSwitchDebugEntity})
}

func RenderDebugEntities(r *RenderDef, debugEntities []*DebugEntity) {
	for _, debugEntity := range debugEntities {
		entityVertexBuffer := debugEntity.VertexBuffer
		if len(entityVertexBuffer) == 0 {
			continue
		}

		// Create render config for debug entity (position only)
		config := r.Renderer.CreateDebugEntityConfig(
			debugEntity.VertexArrayObject,
			debugEntity.VertexBufferObject,
			entityVertexBuffer,
			RENDER_TYPE_DEBUG,
		)

		// Set debug color uniform before rendering using ShaderSystem method
		color := debugEntity.Color
		r.ShaderSystem.SetDebugColor(color)

		// Render the debug entity
		r.Renderer.RenderEntity(config)
	}
}

func BuildAllDebugEntities(gameWorld *world.GameWorld) []*DebugEntity {
	debugEntities := make([]*DebugEntity, 0)
	debugEntities = append(debugEntities, NewDoorTriggerDebugEntity(gameWorld.AotManager.Doors))
	debugEntities = append(debugEntities, NewCollisionDebugEntity(gameWorld.GameRoom.CollisionEntities))
	debugEntities = append(debugEntities, NewSlopedSurfacesDebugEntity(gameWorld.GameRoom.CollisionEntities))
	debugEntities = append(debugEntities, NewItemTriggerDebugEntity(gameWorld.AotManager.Items))
	debugEntities = append(debugEntities, NewAotTriggerDebugEntity(gameWorld.AotManager.AotTriggers))
	return debugEntities
}

// createDebugEntity is a helper function that creates a DebugEntity with the given vertex buffer and color
func createDebugEntity(vertexBuffer []float32, color [4]float32) *DebugEntity {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return &DebugEntity{
		Color:              color,
		VertexBuffer:       vertexBuffer,
		VertexArrayObject:  vao,
		VertexBufferObject: vbo,
	}
}

func NewCollisionDebugEntity(collisionEntities []fileio.CollisionEntity) *DebugEntity {
	vertexBuffer := geometry.NewCollisionDebugEntity(collisionEntities)
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_RED)
}

func NewCameraSwitchDebugEntity(curCameraId int,
	cameraSwitches []fileio.RVDHeader,
	cameraSwitchTransitions map[int][]int) *DebugEntity {
	vertexBuffer := geometry.NewCameraSwitchDebugVertexBuffer(curCameraId, cameraSwitches, cameraSwitchTransitions)
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_GREEN)
}

func NewDoorTriggerDebugEntity(doors []world.AotDoor) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range doors {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_BLUE)
}

func NewItemTriggerDebugEntity(items []world.AotItem) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range items {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_CYAN)
}

func NewAotTriggerDebugEntity(aotTriggers []world.AotObject) *DebugEntity {
	vertexBuffer := make([]float32, 0)
	for _, aot := range aotTriggers {
		vertexBuffer = append(vertexBuffer, aot.Bounds.VertexBuffer...)
	}
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_CYAN)
}

func NewSlopedSurfacesDebugEntity(collisionEntities []fileio.CollisionEntity) *DebugEntity {
	vertexBuffer := geometry.NewSlopedSurfacesDebugVertexBuffer(collisionEntities)
	return createDebugEntity(vertexBuffer, DEBUG_COLOR_MAGENTA)
}
