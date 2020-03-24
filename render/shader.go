package render

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

func NewShader(vertexFilePath string, fragmentFilePath string) *Shader {
	sh := Shader{}

	// compile shaders
	sh.VertexShader = sh.initVertexShader(vertexFilePath)
	sh.FragmentShader = sh.initFragmentShader(fragmentFilePath)

	programShader := gl.CreateProgram()
	gl.AttachShader(programShader, sh.VertexShader)
	gl.AttachShader(programShader, sh.FragmentShader)
	gl.LinkProgram(programShader)
	sh.ProgramShader = programShader

	return &sh
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

		return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (sh *Shader) initVertexShader(filePath string) uint32 {
	vertexShader, err := sh.compileShader(sh.readShaderCode(filePath), gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	return vertexShader
}

func (sh *Shader) initFragmentShader(filePath string) uint32 {
	fragmentShader, err := sh.compileShader(sh.readShaderCode(filePath), gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	return fragmentShader
}

// Read shader code from file
func (sh *Shader) readShaderCode(filePath string) string {
	code := ""
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		code += "\n" + scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	code += "\x00"
	return code
}
