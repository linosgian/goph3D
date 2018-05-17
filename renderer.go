package main

import (
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
}

func (r *Renderer) Clear() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
func (r *Renderer) DrawRawFloor(va *VertexArray, t *Texture, cam *Camera, s *Shader) error {
	s.Bind()
	va.Bind()
	t.Bind(0)
	s.SetUniform1i("texture1\x00", 0)

	view := cam.GetViewMatrix()
	s.SetMat4("view\x00", &view[0])

	// Perspective matrix
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WIDTH)/HEIGHT, 0.1, 100)
	s.SetMat4("projection\x00", &proj[0])

	// 1st rectangle
	model := mgl32.Ident4()
	s.SetMat4("model\x00", &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	return nil
}
func (r *Renderer) DrawRaw(va *VertexArray, t *Texture, cam *Camera, s *Shader) error {
	s.Bind()
	va.Bind()
	t.Bind(0)
	s.SetUniform1i("texture1\x00", 0)

	view := cam.GetViewMatrix()
	s.SetMat4("view\x00", &view[0])

	// Perspective matrix
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WIDTH)/HEIGHT, 0.1, 100)
	s.SetMat4("projection\x00", &proj[0])

	model := mgl32.Translate3D(-1, 0, -1)
	s.SetMat4("model\x00", &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	model = mgl32.Translate3D(2, 0, 0)
	s.SetMat4("model\x00", &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	return nil
}
func (r *Renderer) Draw(va *VertexArray, ib *IndexBuffer, s *Shader) error {
	s.Bind()
	va.Bind()

	var radius float32 = 10
	camX := float32(math.Sin(glfw.GetTime())) * radius
	camZ := float32(math.Cos(glfw.GetTime())) * radius
	view := mgl32.LookAtV(mgl32.Vec3{camX, 0, camZ}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// Push scene backwards so that everything is visible (or move camera backwards)
	// view := mgl32.Translate3D(0, 0, -3)
	s.SetMat4("view\x00", &view[0])

	// Perspective matrix
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WIDTH)/HEIGHT, 0.1, 100)
	s.SetMat4("projection\x00", &proj[0])

	// 1st rectangle
	rotate := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(0.5, -0.5, 0)
	model := trans.Mul4(rotate) // Rotation -> translation! order of transformation is reversed
	s.SetMat4("model\x00", &model[0])
	gl.DrawElements(gl.TRIANGLES, int32(ib.count), gl.UNSIGNED_INT, nil)

	// // 2nd rectangle
	// scaleNum := float32(math.Sin(glfw.GetTime()))
	// scale := mgl32.Scale3D(scaleNum, scaleNum, scaleNum)
	// trans = mgl32.Translate3D(-0.5, 0.5, 0)
	// model = trans.Mul4(scale) // Rotation -> translation! order of transformation is reversed
	// s.SetMat4("model\x00", &model[0])
	// gl.DrawElements(gl.TRIANGLES, int32(ib.count), gl.UNSIGNED_INT, nil)

	return nil
}
