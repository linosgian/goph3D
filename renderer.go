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
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
func (r *Renderer) Draw(va *VertexArray, ib *IndexBuffer, s *Shader) {
	s.Bind()
	va.Bind()
	ib.Bind()
	rotate := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 0, 1})
	trans := mgl32.Translate3D(0.5, -0.5, 0)
	final := trans.Mul4(rotate) // Rotation -> translation! order of transformation is reversed

	loc, _ := s.GetUniformLocation("transform" + "\x00")
	gl.UniformMatrix4fv(loc, 1, false, &final[0])
	gl.DrawElements(gl.TRIANGLES, int32(ib.count), gl.UNSIGNED_INT, nil)

	scaleNum := float32(math.Sin(glfw.GetTime()))
	scale := mgl32.Scale3D(scaleNum, scaleNum, scaleNum)
	trans = mgl32.Translate3D(-0.5, 0.5, 0)
	final = trans.Mul4(scale) // Rotation -> translation! order of transformation is reversed

	loc, _ = s.GetUniformLocation("transform" + "\x00")
	gl.UniformMatrix4fv(loc, 1, false, &final[0])
	gl.DrawElements(gl.TRIANGLES, int32(ib.count), gl.UNSIGNED_INT, nil)
}
