package main

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type VertexBuffer struct {
	rendererID uint32 // A private ID for the object (e.g. OpenGL object ID)
}

func NewVertexBuffer(data []float32, size int) VertexBuffer {
	vb := VertexBuffer{}
	gl.GenBuffers(1, &vb.rendererID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.rendererID)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), gl.STATIC_DRAW)
	return vb
}

func (vb *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.rendererID)
}

func (vb *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vb *VertexBuffer) Delete() {
	gl.DeleteBuffers(1, &vb.rendererID)
}
