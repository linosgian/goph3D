package scene

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/linosgian/goph3d/renderer"
)

const (
	FOV  = 45.0
	NEAR = 0.1
	FAR  = 100.0
)

type PointLight struct {
	Position                    mgl32.Vec3
	Ambient, Diffuse, Specular  mgl32.Vec3
	Constant, Linear, Quadratic float32
}
type Scene struct {
	Nodes                []*Node
	Cam                  *Camera
	DeltaTime, LastFrame float64
	Perspective          mgl32.Mat4
	lightPos             mgl32.Vec3
	PointLights          []*PointLight
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

func NewScene(ratio float32, c *Camera, lights []*PointLight) *Scene {
	proj := mgl32.Perspective(mgl32.DegToRad(FOV), ratio, NEAR, FAR)
	return &Scene{
		Nodes:       make([]*Node, 0),
		Cam:         c,
		DeltaTime:   0,
		LastFrame:   0,
		Perspective: proj,
		// lightPos:    lightPos,
		PointLights: lights,
	}
}

func (s *Scene) attach(n *Node) {
	s.Nodes = append(s.Nodes, n)
}

func (s *Scene) InitLights(r *renderer.Renderer) {
	shaderID, err := r.GetProgram("phong")
	if err != nil {
		log.Fatalf("no such shader program: %q\n", err)
	}
	phongShader := r.Programs[shaderID]
	phongShader.Bind()
	// Prepare material's values
	// That depends on what we are rendering!
	// It shouldnt be in scene's initlights
	phongShader.SetVec3f("material.ambient", 1.0, 0.5, 0.31)
	phongShader.SetVec3f("material.diffuse", 1.0, 0.5, 0.31)
	phongShader.SetVec3f("material.specular", 0.5, 0.5, 0.5)
	phongShader.SetFloat("material.shininess", 32.0)

	// Directional light
	phongShader.SetVec3f("dirLight.direction", -0.2, -1.0, -0.3)
	phongShader.SetVec3f("dirLight.ambient", 0.05, 0.05, 0.05)
	phongShader.SetVec3f("dirLight.diffuse", 0.4, 0.4, 0.4)
	phongShader.SetVec3f("dirLight.specular", 0.5, 0.5, 0.5)

	// Spotlights
	for i, light := range s.PointLights {
		phongShader.SetVec3(fmt.Sprintf("pointLights[%d].position", i), light.Position)
		phongShader.SetVec3(fmt.Sprintf("pointLights[%d].ambient", i), light.Ambient)
		phongShader.SetVec3(fmt.Sprintf("pointLights[%d].diffuse", i), light.Diffuse)
		phongShader.SetVec3(fmt.Sprintf("pointLights[%d].specular", i), light.Specular)
		phongShader.SetFloat(fmt.Sprintf("pointLights[%d].constant", i), light.Constant)
		phongShader.SetFloat(fmt.Sprintf("pointLights[%d].linear", i), light.Linear)
		phongShader.SetFloat(fmt.Sprintf("pointLights[%d].quadratic", i), light.Quadratic)
	}

	phongShader.SetVec3("spotLight.position", s.Cam.Position)
	phongShader.SetVec3("spotLight.direction", s.Cam.Front)
	phongShader.SetVec3f("spotLight.ambient", 0.0, 0.0, 0.0)
	phongShader.SetVec3f("spotLight.diffuse", 1.0, 1.0, 1.0)
	phongShader.SetVec3f("spotLight.specular", 1.0, 1.0, 1.0)
	phongShader.SetFloat("spotLight.constant", 1.0)
	phongShader.SetFloat("spotLight.linear", 0.09)
	phongShader.SetFloat("spotLight.quadratic", 0.032)
	phongShader.SetFloat("spotLight.cutOff", float32(math.Cos(float64(mgl32.DegToRad(12.5)))))
	phongShader.SetFloat("spotLight.outerCutOff", float32(math.Cos(float64(mgl32.DegToRad(15)))))
}

// NOTE: This looks meh..
func (s *Scene) Update(r *renderer.Renderer) {
	shaderID, err := r.GetProgram("phong")
	if err != nil {
		log.Fatalf("no such shader program: %q\n", err)
	}

	phongShader := r.Programs[shaderID]
	phongShader.Bind()
	phongShader.SetVec3("viewPos\x00", s.Cam.Position)

	phongShader.SetVec3("spotLight.position", s.Cam.Position)
	phongShader.SetVec3("spotLight.direction", s.Cam.Front)

}
