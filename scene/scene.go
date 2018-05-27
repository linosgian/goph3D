package scene

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/linosgian/goph3d/renderer"
)

const (
	FOV  = 45.0
	NEAR = 0.1
	FAR  = 100.0
)

type Scene struct {
	Nodes                []*Node
	Cam                  *Camera
	DeltaTime, LastFrame float64
	Perspective          mgl32.Mat4
	lightPos             mgl32.Vec3
}

// Holds all internal IDs for the VAO, Texture and Shader program
type Node struct {
	VaoID, TexID, ProgramID int
	ModelMatrix             mgl32.Mat4
	Renderable              bool
	Position                mgl32.Vec3
	Name                    string
}

// Creates a Node based on the data, texture and shader program
// Receives a renderer, the data(VBO) the path to the texture we want
// and a shader program by name (e.g. "basic", "phong")
// Returns a pointer to the created Node
// TODO: Textures are loaded multiple times!
func (s *Scene) NewNode(r *renderer.Renderer, name string, renderable bool, data []float32, texturePath, programName string, modelPos mgl32.Vec3) error {
	programID, err := r.GetProgram(programName)
	if err != nil {
		return err
	}

	vaoID, err := r.LoadData(data)
	if err != nil {
		return err
	}

	texID, err := r.LoadTexture(texturePath, programID)
	if err != nil {
		return err
	}
	n := &Node{
		Renderable: renderable,
		VaoID:      vaoID,
		TexID:      texID,
		ProgramID:  programID,
		Position:   modelPos,
		Name:       name,
	}
	n.SetModelMatrix(mgl32.Translate3D(modelPos.X(), modelPos.Y(), modelPos.Z()))
	s.attach(n)
	return nil
}

// Take the same arguments as NewNode alongside with all the different model positions
// This is useful when we having a single VAO with multiple transformations
func (s *Scene) NewNodes(r *renderer.Renderer, name string, renderable bool, data []float32, texturePath, programName string, modelPositions []mgl32.Vec3) error {
	programID, err := r.GetProgram(programName)
	if err != nil {
		return err
	}

	vaoID, err := r.LoadData(data)
	if err != nil {
		return err
	}

	texID, err := r.LoadTexture(texturePath, programID)
	if err != nil {
		return err
	}
	for _, pos := range modelPositions {
		node := &Node{
			Renderable: renderable,
			VaoID:      vaoID,
			TexID:      texID,
			ProgramID:  programID,
			Position:   pos,
			Name:       name,
		}
		node.ModelMatrix = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
		s.attach(node)
	}
	return nil
}

func (n *Node) SetModelMatrix(model mgl32.Mat4) {
	n.ModelMatrix = model
}

func NewScene(ratio float32, c *Camera, lightPos mgl32.Vec3) *Scene {
	proj := mgl32.Perspective(mgl32.DegToRad(FOV), ratio, NEAR, FAR)
	return &Scene{
		Nodes:       make([]*Node, 0),
		Cam:         c,
		DeltaTime:   0,
		LastFrame:   0,
		Perspective: proj,
		lightPos:    lightPos,
	}
}

func (s *Scene) attach(n *Node) {
	s.Nodes = append(s.Nodes, n)
}

// NOTE: This looks meh..
func (s *Scene) Update(r *renderer.Renderer) {
	shaderID, err := r.GetProgram("phong")
	if err != nil {
		log.Fatalf("no such shader program: %q\n", err)
	}

	phongShader := r.Programs[shaderID]
	phongShader.Bind()
	phongShader.SetVec3("lightColor\x00", mgl32.Vec3{1, 1, 1})
	phongShader.SetVec3("lightPos\x00", s.lightPos)
	phongShader.SetVec3("viewPos\x00", s.Cam.Position)

}
