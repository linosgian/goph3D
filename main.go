package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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

const (
	projectRoot      = "src/github.com/linosgian/glfw-test2"
	fragmentFilepath = "res/shaders/fragment.glsl"
	vertexFilepath   = "res/shaders/vertex.glsl"
	texturePath      = "res/textures/face.png"
)

var tri = []float32{
	// X, Y / R, G, B / S, T
	-0.5, -0.5, 1, 0, 0, 0, 0,
	0.5, -0.5, 0, 1, 0, 1, 0,
	0.5, 0.5, 0, 0, 1, 1, 1,
	-0.5, 0.5, 1, 1, 0, 0, 1,
}

// Indices must be uint32 instead of uint
// in order to match gl.UNSIGNED_INT
var indices = []uint32{
	0, 1, 2,
	2, 3, 0,
}

func main() {
	rootPath := path.Join(os.Getenv("GOPATH"), projectRoot)
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()

	window, err := initGLFW()
	if err != nil {
		log.Fatalf("[GLFW error]: %q", err)
	}
	defer glfw.Terminate()

	if err := initOpenGL(); err != nil {
		log.Fatalf("OpenGL could not be initialized: %v\n", err)
	}

	// Instantiate all data needed for rendering
	vb := NewVertexBuffer(tri, len(tri)*sizes[FLOAT])

	va := NewVertexArray()

	vbl := &VertexBufferLayout{Stride: 0}
	vbl.PushFloat(2) // position: a fvec2
	vbl.PushFloat(3) // color: a fvec3
	vbl.PushFloat(2) // texture: a fvec2
	va.AddBuffer(&vb, vbl)

	ib := NewIndexBuffer(indices)

	// Create a new Shader
	shader, err := NewShader(path.Join(rootPath, vertexFilepath), path.Join(rootPath, fragmentFilepath))
	if err != nil {
		log.Fatalf("could not create shader program: %v\n", err)
	}
	shader.Bind()

	// Prepare textures
	t, err := NewTexture(path.Join(rootPath, texturePath))
	if err != nil {
		log.Fatalf("could not create texture: %v\n", err)
	}
	t.Bind(0)
	shader.SetUniform1i("u_Texture\x00", 0)

	r := Renderer{}
	// Clear all state before game loop
	va.Unbind()
	vb.Unbind()
	shader.Unbind()
	ib.Unbind()

	for !window.ShouldClose() {
		r.Clear()

		r.Draw(va, ib, shader)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGLFW() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// Specify profile and OpenGL version
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(1024, 680, "3D Gamez", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	return window, nil
}

func initOpenGL() error {
	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// Enable debugging and hook callback
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(Debug, nil)
	return nil
}

func Debug(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer,
) {
	fmt.Printf("[OpenGL Error]: %q\n", message)
}
