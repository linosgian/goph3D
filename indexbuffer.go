package main

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type IndexBuffer struct {
	rendererID uint32 // A private ID for the object (e.g. OpenGL object ID)
	count      int    // Count of indices
}

func NewIndexBuffer(indices []uint32) *IndexBuffer {
	ib := IndexBuffer{count: len(indices)}
	gl.GenBuffers(1, &ib.rendererID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.rendererID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, ib.count*sizes[UINT32], gl.Ptr(indices), gl.STATIC_DRAW)
	return &ib
}

func (ib *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.rendererID)
}

func (ib *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBuffer) Delete() {
	gl.DeleteBuffers(1, &ib.rendererID)
}
