package renderer

import (
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

const (
	FLOAT  = gl.FLOAT
	UINT32 = gl.UNSIGNED_INT
)

// Calculate the byte-size of all types at runtime
var sizes map[int]int = map[int]int{
	FLOAT:  int(unsafe.Sizeof(float32(0))),
	UINT32: int(unsafe.Sizeof(uint32(0))),
}
