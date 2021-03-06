package renderer

import (
	"fmt"
	"strings"

	"io/ioutil"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	rendererID           uint32           // A private ID for the object (e.g. OpenGL object ID)
	uniformLocationCache map[string]int32 // Cache for uniform locations
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	s := Shader{
		rendererID:           gl.CreateProgram(),
		uniformLocationCache: make(map[string]int32),
	}

	// Compile shaders
	vsf, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return &s, fmt.Errorf("Could not read shader file for vertex shader: %q", err)
	}
	// TODO: Concat strings more effeciently and make a wrapper for the null character
	vsSource := fmt.Sprintf("%s%s", vsf, "\x00") // Make it a C-Style null-terminated string
	vs, err := s.compileShader(gl.VERTEX_SHADER, vsSource)
	if err != nil {
		return &s, err
	}

	fsf, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return &s, fmt.Errorf("Could not read shader file for fragment shader: %q", err)
	}
	fsSource := fmt.Sprintf("%s%s", fsf, "\x00")

	fs, err := s.compileShader(gl.FRAGMENT_SHADER, fsSource)
	if err != nil {
		return &s, err
	}

	gl.AttachShader(s.rendererID, vs)
	gl.AttachShader(s.rendererID, fs)
	gl.LinkProgram(s.rendererID)

	// TODO: Abstract error handling
	var success int32
	gl.GetProgramiv(s.rendererID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(s.rendererID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.rendererID, logLength, nil, gl.Str(log))
		return &s, fmt.Errorf("failed to link: %q", log)
	}

	// Delete shaders as they are linked already
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	return &s, nil
}

func (s *Shader) compileShader(shaderType uint32, source string) (uint32, error) {
	id := gl.CreateShader(shaderType)

	src, free := gl.Strs(source)
	gl.ShaderSource(id, 1, src, nil)
	free()

	gl.CompileShader(id)

	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)

	// If an error occured, grab the info
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(id, logLength, nil, gl.Str(log))

		gl.DeleteShader(id)
		return 0, fmt.Errorf("failed to compile\n%v: %v", source, log)
	}
	return id, nil
}

func (s *Shader) Bind() {
	gl.UseProgram(s.rendererID)
}

func (s *Shader) Unbind() {
	gl.UseProgram(0)
}

func (s *Shader) SetUniform1i(name string, val int32) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.Uniform1i(location, val)
	return nil
}

func (s *Shader) SetMat4(name string, v *float32) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.UniformMatrix4fv(location, 1, false, v)
	return nil
}

func (s *Shader) SetFloat(name string, val float32) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.Uniform1f(location, val)
	return nil
}
func (s *Shader) SetVec3f(name string, x, y, z float32) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.Uniform3f(location, x, y, z)
	return nil
}

func (s *Shader) SetVec3(name string, vec3 mgl32.Vec3) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.Uniform3f(location, vec3.X(), vec3.Y(), vec3.Z())
	return nil
}

func (s *Shader) SetUniform4f(name string, v0, v1, v2, v3 float32) error {
	nullTermString := fmt.Sprintf("%s\x00", name)
	location, err := s.GetUniformLocation(nullTermString)
	if err != nil {
		return err
	}
	gl.Uniform4f(location, v0, v1, v2, v3)
	return nil
}

func (s *Shader) GetUniformLocation(name string) (int32, error) {
	nullTermString := fmt.Sprintf("%s\x00", name)
	// Check if it is cached
	if val, ok := s.uniformLocationCache[name]; ok {
		return val, nil
	}
	location := gl.GetUniformLocation(s.rendererID, gl.Str(nullTermString))
	s.uniformLocationCache[name] = location
	if location == -1 {
		return 0, fmt.Errorf("Uniform variable location not found: %v", name)
	}
	return location, nil
}
