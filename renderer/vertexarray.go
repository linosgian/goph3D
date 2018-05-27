package renderer

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type VertexArray struct {
	rendererID uint32
	Vcount     int32 // vertex counter for draw call
	DataSize   int32 // Size of input data
}

func NewVertexArray() *VertexArray {
	va := VertexArray{}
	gl.GenVertexArrays(1, &va.rendererID)
	return &va
}

func (va *VertexArray) DeleteVertexArray() {
	gl.DeleteVertexArrays(1, &va.rendererID)
}

func (va *VertexArray) Bind() {
	gl.BindVertexArray(va.rendererID)
}

func (va *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

func (va *VertexArray) AddBuffer(vb *VertexBuffer, vbl *VertexBufferLayout) {
	va.Bind()
	vb.Bind()

	offset := 0
	for i, e := range vbl.Elements {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointer(
			uint32(i),
			e.count,
			uint32(e.etype),
			e.normalized,
			vbl.Stride,
			gl.PtrOffset(offset),
		)
		offset += int(e.count) * sizes[int(e.etype)]
	}
}
