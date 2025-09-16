package shader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	VertexShader   uint32
	FragmentShader uint32
	ProgramShader  uint32
}

func NewShader(vertexFilePath string, fragmentFilePath string) (*Shader, error) {
	sh := Shader{}

	// compile shaders
	var err error
	sh.VertexShader, err = sh.initVertexShader(vertexFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to compile vertex shader: %w", err)
	}

	sh.FragmentShader, err = sh.initFragmentShader(fragmentFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to compile fragment shader: %w", err)
	}

	programShader := gl.CreateProgram()
	gl.AttachShader(programShader, sh.VertexShader)
	gl.AttachShader(programShader, sh.FragmentShader)
	gl.LinkProgram(programShader)

	// Check if program linked successfully
	var success int32
	gl.GetProgramiv(programShader, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programShader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(programShader, logLength, nil, &log[0])
		return nil, fmt.Errorf("failed to link shader program: %s", string(log))
	}

	sh.ProgramShader = programShader

	// Clean up individual shaders after linking
	gl.DeleteShader(sh.VertexShader)
	gl.DeleteShader(sh.FragmentShader)

	return &sh, nil
}

func (sh *Shader) compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (sh *Shader) initVertexShader(filePath string) (uint32, error) {
	vertexShader, err := sh.compileShader(sh.readShaderCode(filePath), gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("vertex shader compilation failed: %w", err)
	}
	return vertexShader, nil
}

func (sh *Shader) initFragmentShader(filePath string) (uint32, error) {
	fragmentShader, err := sh.compileShader(sh.readShaderCode(filePath), gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("fragment shader compilation failed: %w", err)
	}
	return fragmentShader, nil
}

// Read shader code from file
func (sh *Shader) readShaderCode(filePath string) string {
	var builder strings.Builder
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		builder.WriteString("\n")
		builder.WriteString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	builder.WriteString("\x00")
	return builder.String()
}
