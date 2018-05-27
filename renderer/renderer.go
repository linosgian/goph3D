package renderer

import (
	"fmt"
	"os"
	"path"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	shadersPath = "res/shaders"
)

// This should be given temporarily because of vim-go
var rootPath = os.Getenv("PROJ_PATH") // e.g. /home/lgian/go/src/github.com/linosgian/goph3d

type Renderer struct {
	vaos         []*VertexArray
	textures     []*Texture
	programNames map[string]int
	Programs     []*Shader
}

func NewRenderer() (*Renderer, error) {
	r := &Renderer{
		vaos:         make([]*VertexArray, 0),
		textures:     make([]*Texture, 0),
		programNames: make(map[string]int, 0),
		Programs:     make([]*Shader, 0),
	}

	// Load all default shaders
	if err := r.LoadDefaultPrograms(); err != nil {
		return nil, err
	}
	return r, nil

}

func (r *Renderer) DrawRaw(vaoID, programID, texID int, view, proj, model mgl32.Mat4) error {
	s := r.Programs[programID]
	va := r.vaos[vaoID]

	s.Bind()
	va.Bind()

	// TODO: Improve this by holding an ID for the texture instead of loading
	// it to the first slot all the time
	// Load the right texture for the object
	r.textures[texID].Bind(0)
	s.SetUniform1i("aTexture\x00", 0)

	// Camera
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
func (r *Renderer) LoadTexture(texturePath string, programID int) (int, error) {
	r.Programs[programID].Bind()
	t, err := NewTexture(texturePath)
	if err != nil {
		return 0, err
	}
	objID := len(r.textures)
	r.textures = append(r.textures, t)
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
	r.programNames[progName] = len(r.Programs)
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
	vbl.PushFloat(3) // Normals: a fvec2

	va.Vcount = vbl.Vcount // I DONT LIKE THIS SHIT. Reconsider in the future
	va.DataSize = int32(len(data))
	va.AddBuffer(vb, vbl)

	// State should remain clean after each load
	va.Unbind()
	vb.Unbind()

	objID := len(r.vaos)
	r.vaos = append(r.vaos, va)
	return objID, nil
}

// Loads all default shader programs
// For any new program name added to programNames
// we expect to find <name>_vertex.glsl and <name>_fragment.glsl under the shadersPath
func (r *Renderer) LoadDefaultPrograms() error {
	programNames := []string{"basic", "phong", "lamp"}
	for _, pName := range programNames {
		vsPath := path.Join(rootPath, shadersPath, fmt.Sprintf("%s_vertex.glsl", pName))
		fsPath := path.Join(rootPath, shadersPath, fmt.Sprintf("%s_fragment.glsl", pName))
		if err := r.loadProgram(pName, vsPath, fsPath); err != nil {
			return err
		}
	}
	return nil
}

// Find a program ID by name
// The returned ID is the internal one
func (r *Renderer) GetProgram(progName string) (int, error) {
	if pID, ok := r.programNames[progName]; ok {
		return pID, nil
	}
	return 0, fmt.Errorf("Could not find a program by that name: %q", progName)
}
