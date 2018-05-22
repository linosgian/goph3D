package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraMovement uint8

// These are just enums in go
const (
	FORWARD CameraMovement = iota
	BACKWARD
	LEFT
	RIGHT
)

const (
	YAW         float32 = -90
	PITCH       float32 = 0
	SPEED       float32 = 2.5
	SENSITIVITY float32 = 0.1
	ZOOM        float32 = 45
)

type Camera struct {
	// Euler values
	Yaw, Pitch float32
	// Camera options
	MovementSpeed, Zoom, MouseSensitivity float32
	// Camera attributes
	Position, Front, Up, Right, WorldUp mgl32.Vec3
}

func NewCamera(pos mgl32.Vec3) *Camera {
	c := &Camera{
		Yaw:              YAW,
		Pitch:            PITCH,
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVITY,
		Front:            mgl32.Vec3{0, 0, -1},
		Position:         pos,
		WorldUp:          mgl32.Vec3{0, 1, 0},
	}
	c.updateCameraVectors()
	return c
}

// Returns a matrix with the camera view for OpenGL
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessMouseMovement(xoffset, yoffset float32) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch -= yoffset // Add offset for inverse mode

	// Constrain pitch values
	if c.Pitch > 89 {
		c.Pitch = 89
	}
	if c.Pitch < -89 {
		c.Pitch = -89
	}

	c.updateCameraVectors()
}

func (c *Camera) ProcessKeyboard(direction CameraMovement, deltaTime float64) {
	// deltaTime between frames is used so that
	// players with better hardware don't move faster
	velocity := c.MovementSpeed * float32(deltaTime)

	switch direction {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}

func (c *Camera) updateCameraVectors() {
	// NOTE: Too many conversions. Maybe wrap this once in a library (create a math32 library)
	yaw, pitch := float64(mgl32.DegToRad(c.Yaw)), float64(mgl32.DegToRad(c.Pitch))
	frontX := float32(math.Cos(yaw) * math.Cos(pitch))
	frontY := float32(math.Sin(pitch))
	frontZ := float32(math.Sin(yaw) * math.Cos(pitch))

	c.Front = mgl32.Vec3{frontX, frontY, frontZ}.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
