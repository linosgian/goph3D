package main

import (
	"fmt"
	"path"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	VAOs         []*VertexArray
	Textures     []*Texture
	ProgramNames map[string]int
	Programs     []*Shader
}

func NewRenderer(w *glfw.Window, c *Camera) (*Renderer, error) {
	r := &Renderer{
		VAOs:         make([]*VertexArray, 0),
		Textures:     make([]*Texture, 0),
		ProgramNames: make(map[string]int, 0),
		Programs:     make([]*Shader, 0),
	}
	// Load all default shaders
	if err := r.LoadDefaultPrograms(); err != nil {
		return nil, err
	}

	// Install callbacks for all inputs
	r.InitInputs(w, c)
	return r, nil
}

func (r *Renderer) InitInputs(w *glfw.Window, c *Camera) {
	w.SetCursorPosCallback(MouseCallback)
	w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled) // Disable mouse pointer while playing
	w.SetUserPointer(unsafe.Pointer(c))                  // This is needed for the mouse callback
}

func (r *Renderer) Destroy() {
	glfw.Terminate()
}

func (r *Renderer) Clear() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0) // Default scene color.
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Renderer) DrawRaw(n *Node, cam *Camera, proj mgl32.Mat4) error {
	s := r.Programs[n.programID]
	va := r.VAOs[n.vaoID]

	s.Bind()
	va.Bind()

	// TODO: Improve this by holding an ID for the texture instead of loading
	// it to the first slot all the time
	// Load the right texture for the object
	r.Textures[n.texID].Bind(0)
	s.SetUniform1i("aTexture\x00", 0)

	// Camera
	view := cam.GetViewMatrix()
	s.SetMat4("view\x00", &view[0])

	// Perspective matrix
	s.SetMat4("projection\x00", &proj[0])

	// Model matrix
	s.SetMat4("model\x00", &n.ModelTrans[0])

	gl.DrawArrays(gl.TRIANGLES, 0, va.DataSize/va.Vcount)
	return nil
}

// Loads a texture for a specific program (shader)
// Returns an internal object ID
func (r *Renderer) LoadTexture(texturePath string, programID int) (int, error) {
	r.Programs[programID].Bind()
	t, err := NewTexture(texturePath)
	if err != nil {
		return 0, err
	}
	objID := len(r.Textures)
	r.Textures = append(r.Textures, t)
	r.Programs[programID].Unbind()
	return objID, nil
}

// Loads a program with both the fragment and vertex shaders
// Returns an internal object ID
func (r *Renderer) loadProgram(progName, vsPath, fsPath string) error {
	s, err := NewShader(vsPath, fsPath)
	if err != nil {
		return err
	}
	r.ProgramNames[progName] = len(r.Programs)
	r.Programs = append(r.Programs, s)
	return nil
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
	va.AddBuffer(vb, vbl)

	// State should remain clean after each load
	va.Unbind()
	vb.Unbind()

	objID := len(r.VAOs)
	r.VAOs = append(r.VAOs, va)
	return objID, nil
}

func (r *Renderer) LoadDefaultPrograms() error {
	programNames := []string{"basic"}
	for _, pName := range programNames {
		// TODO: Better string concat handling
		vsPath := path.Join(rootPath, shadersPath, pName+"_vertex.glsl")
		fsPath := path.Join(rootPath, shadersPath, pName+"_fragment.glsl")
		if err := r.loadProgram(pName, vsPath, fsPath); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) GetProgram(progName string) (int, error) {
	if pID, ok := r.ProgramNames[progName]; ok {
		return pID, nil
	}
	return 0, fmt.Errorf("Could not find a program by that name: %q", progName)
}
