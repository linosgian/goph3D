package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	Nodes                []*Node
	Cam                  *Camera
	deltaTime, lastFrame float64
	Perspective          mgl32.Mat4
}

// Holds all internal IDs for the VAO, Texture and Shader program
//
type Node struct {
	vaoID, texID, programID int
	ModelTrans              mgl32.Mat4
	Renderable              bool
}

// Creates a Node based on the data, texture and shader program
func (s *Scene) NewNode(r *Renderer, renderable bool, data []float32, texturePath, programName string) (*Node, error) {
	programID, err := r.GetProgram(programName)
	if err != nil {
		return nil, err
	}

	vaoID, err := r.LoadData(data)
	if err != nil {
		return nil, err
	}

	texID, err := r.LoadTexture(texturePath, programID)
	if err != nil {
		return nil, err
	}
	n := &Node{
		Renderable: renderable,
		vaoID:      vaoID,
		texID:      texID,
		programID:  programID,
	}
	s.attach(n)
	return n, nil
}

func (n *Node) SetModelMatrix(model mgl32.Mat4) {
	n.ModelTrans = model
}

func NewScene(p mgl32.Mat4, c *Camera) *Scene {
	return &Scene{
		Nodes:       make([]*Node, 0),
		Cam:         c,
		deltaTime:   0,
		lastFrame:   0,
		Perspective: p,
	}
}

func (s *Scene) attach(n *Node) {
	s.Nodes = append(s.Nodes, n)
}
