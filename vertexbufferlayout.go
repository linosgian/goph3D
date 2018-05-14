package main

import "github.com/go-gl/gl/v4.3-core/gl"

type VertexBufferElement struct {
	count      int32
	etype      uint32
	normalized bool
}

type VertexBufferLayout struct {
	Elements []VertexBufferElement
	Stride   int32
}

func (vbl *VertexBufferLayout) PushFloat(count int32) {
	e := VertexBufferElement{
		etype:      gl.FLOAT,
		count:      count,
		normalized: false,
	}
	vbl.Elements = append(vbl.Elements, e)
	vbl.Stride += count * int32(sizes[FLOAT])
}

func (vbl *VertexBufferLayout) PushUint(count int32) {
	e := VertexBufferElement{
		etype:      gl.UNSIGNED_INT,
		count:      count,
		normalized: false,
	}
	vbl.Elements = append(vbl.Elements, e)
	vbl.Stride += count * int32(sizes[UINT32])
}
