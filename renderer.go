package main

import (
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	Vas      []*VertexArray
	Textures []*Texture
	Shaders  []*Shader
}

func NewRenderer() *Renderer {
	return &Renderer{
		Vas:      make([]*VertexArray, 0),
		Textures: make([]*Texture, 0),
		Shaders:  make([]*Shader, 0),
	}
}

func (r *Renderer) Init() *glfw.Window {
	window, err := initGLFW()
	if err != nil {
		// These could be propagated to the main function but we're gonna halt anyway
		log.Fatalf("[GLFW error]: %q", err)
	}
	if err := initOpenGL(); err != nil {
		log.Fatalf("OpenGL could not be initialized: %v\n", err)
	}
	return window
}

func (r *Renderer) Destroy() {
	glfw.Terminate()
}

func (r *Renderer) Clear() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0) // Default scene color.
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Renderer) DrawRaw(vaID, tID, sID int, cam *Camera, proj, model mgl32.Mat4) error {
	s := r.Shaders[sID]
	va := r.Vas[vaID]

	s.Bind()
	va.Bind()

	// TODO: Improve this by holding an ID for the texture instead of loading
	// it to the first slot all the time
	// Load the right texture for the object
	r.Textures[tID].Bind(0)
	s.SetUniform1i("texture1\x00", 0)

	// Camera
	view := cam.GetViewMatrix()
	s.SetMat4("view\x00", &view[0])

	// Perspective matrix
	s.SetMat4("projection\x00", &proj[0])

	// Model matrix
	s.SetMat4("model\x00", &model[0])

	gl.DrawArrays(gl.TRIANGLES, 0, va.DataSize/va.Vcount)
	return nil
}

// Loads a texture for a specific program (shader)
// Returns an internal object ID
func (r *Renderer) LoadTexture(texturePath string, shaderID int) (int, error) {
	r.Shaders[shaderID].Bind()
	t, err := NewTexture(texturePath)
	if err != nil {
		return 0, err
	}
	objID := len(r.Textures)
	r.Textures = append(r.Textures, t)
	r.Shaders[shaderID].Unbind()
	return objID, nil
}

// Loads a program with both the fragment and vertex shaders
// Returns an internal object ID
func (r *Renderer) LoadProgram(vsPath, fsPath string) (int, error) {
	s, err := NewShader(vsPath, fsPath)
	if err != nil {
		return 0, err
	}
	objID := len(r.Shaders)
	r.Shaders = append(r.Shaders, s)
	return objID, nil
}

// Loads a vertex buffer
// Returns an internal object ID
func (r *Renderer) LoadData(data []float32) (int, error) {
	vb := NewVertexBuffer(data, len(data)*sizes[FLOAT])
	va := NewVertexArray()

	vbl := new(VertexBufferLayout)
	vbl.PushFloat(3) // position: a fvec3
	vbl.PushFloat(2) // texture: a fvec2

	va.Vcount = vbl.Vcount // I DONT LIKE THIS SHIT. Reconsider in the future
	va.DataSize = int32(len(data))
	va.AddBuffer(&vb, vbl)

	// State should remain clean after each load
	va.Unbind()
	vb.Unbind()

	objID := len(r.Vas)
	r.Vas = append(r.Vas, va)
	return objID, nil
}
