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
	"github.com/go-gl/mathgl/mgl32"
)

const (
	projectRoot  = "src/github.com/linosgian/glfw-test2"
	fragmentPath = "res/shaders/fragment.glsl"
	vertexPath   = "res/shaders/vertex.glsl"
	texturePath  = "res/textures/container.jpg"
)

const (
	FLOAT  = gl.FLOAT
	UINT32 = gl.UNSIGNED_INT
	WIDTH  = 800
	HEIGHT = 600
)

var (
	lastX      float64 = WIDTH / 2.0
	lastY      float64 = HEIGHT / 2.0
	firstMouse bool    = true
	deltaTime  float64 = 0
	lastFrame  float64 = 0
)

// Calculate the byte-size of all types at runtime
var sizes map[int]int = map[int]int{
	FLOAT:  int(unsafe.Sizeof(float32(0))),
	UINT32: int(unsafe.Sizeof(uint32(0))),
}

var cube = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
}

var tri = []float32{
	// X, Y / R, G, B / S, T
	-0.5, -0.5, 1, 0, 0, 0, 0,
	0.5, -0.5, 0, 1, 0, 1, 0,
	0.5, 0.5, 0, 0, 1, 1, 1,
	-0.5, 0.5, 1, 1, 0, 0, 1,
}

// Indices must be uint32 instead of uint
// in order to match gl.UNSIGNED_INT
var triIndices = []uint32{
	0, 1, 2,
	2, 3, 0,
}

//Initialize camera object
var cam = NewCamera(mgl32.Vec3{0, 0, 7})

func main() {
	rootPath := path.Join(os.Getenv("GOPATH"), projectRoot)
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()

	window, err := initGLFW()
	if err != nil {
		log.Fatalf("[GLFW error]: %q", err)
	}
	defer glfw.Terminate()

	if err := initOpenGL(); err != nil {
		log.Fatalf("OpenGL could not be initialized: %v\n", err)
	}

	// Set Input callbacks
	window.SetCursorPosCallback(MouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Instantiate all data needed for rendering
	vb := NewVertexBuffer(cube, len(cube)*sizes[FLOAT])
	va := NewVertexArray()

	vbl := new(VertexBufferLayout)
	vbl.PushFloat(3) // position: a fvec2
	// vbl.PushFloat(3) // color: a fvec3
	vbl.PushFloat(2) // texture: a fvec2

	va.AddBuffer(&vb, vbl)

	// ib := NewIndexBuffer(cubeIndices)

	// Create a new Shader
	shader, err := NewShader(path.Join(rootPath, vertexPath), path.Join(rootPath, fragmentPath))
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

	// Initialize renderer
	r := Renderer{}

	// Clear all state before game loop
	va.Unbind()
	vb.Unbind()
	shader.Unbind()
	// ib.Unbind()

	gl.Enable(gl.DEPTH_TEST)
	for !window.ShouldClose() {
		// Per-frame time
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		r.Clear()

		processInput(window)
		// r.Draw(va, ib, shader)
		r.DrawRaw(va, cam, shader)

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

	window, err := glfw.CreateWindow(WIDTH, HEIGHT, "3D Gamez", nil, nil)
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

func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
	if w.GetKey(glfw.KeyW) == glfw.Press {
		cam.ProcessKeyboard(FORWARD, deltaTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		cam.ProcessKeyboard(BACKWARD, deltaTime)
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		cam.ProcessKeyboard(LEFT, deltaTime)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		cam.ProcessKeyboard(RIGHT, deltaTime)
	}
}
func MouseCallback(w *glfw.Window, xpos, ypos float64) {
	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := float32(xpos - lastX)
	yoffset := float32(ypos - lastY)

	lastX = xpos
	lastY = ypos
	cam.ProcessMouseMovement(xoffset, yoffset)
}
