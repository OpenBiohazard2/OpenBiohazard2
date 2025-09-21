package render

const (
	// OpenGL Constants
	FLOAT_SIZE_BYTES       = 4
	VERTEX_ATTRIB_POSITION = 0
	VERTEX_ATTRIB_TEXTURE  = 1
	VERTEX_ATTRIB_NORMAL   = 2
	TEXTURE_UNIT_0         = 0

	// Render Game State Constants
	RENDER_GAME_STATE_MAIN                   = 0
	RENDER_GAME_STATE_BACKGROUND_SOLID       = 1
	RENDER_GAME_STATE_BACKGROUND_TRANSPARENT = 2
	RENDER_TYPE_ITEM                         = 5

	// Camera Constants
	DEFAULT_FOV_DEGREES = 60.0
	NEAR_PLANE          = 16.0
	FAR_PLANE           = 45000.0
	ASPECT_RATIO        = 4.0 / 3.0
)
