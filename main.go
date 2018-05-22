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
	fragmentPath = "res/shaders/fragment.glsl"
	vertexPath   = "res/shaders/vertex.glsl"
	metalPath    = "res/textures/metal.png"
	marblePath   = "res/textures/marble.jpg"
	shadersPath  = "res/shaders"

	FLOAT  = gl.FLOAT
	UINT32 = gl.UNSIGNED_INT

	WIDTH  = 800 // These will be in a config file eventually
	HEIGHT = 600 // And won't be contants
)

// Calculate the byte-size of all types at runtime
var sizes map[int]int = map[int]int{
	FLOAT:  int(unsafe.Sizeof(float32(0))),
	UINT32: int(unsafe.Sizeof(uint32(0))),
}

var (
	lastX      float64 = WIDTH / 2.0 // In the middle of the screen
	lastY      float64 = HEIGHT / 2.0
	firstMouse bool    = true
	deltaTime  float64 = 0
	lastFrame  float64 = 0
)

// This should be given temporarily because of vim-go
var rootPath = os.Getenv("PROJ_PATH") // e.g. /home/lgian/go/src/github.com/linosgian/glfw-test2

// This runs before main
func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()
}

func main() {
	// Initialize GLFW and OpenGL
	r := NewRenderer()
	window := r.Init()

	//Initialize camera object at a certain position
	cam := NewCamera(mgl32.Vec3{0, 0, 7})

	// NOTE: These are glfw specific - Maybe an issue (?)
	window.SetCursorPosCallback(MouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled) // Disable mouse pointer while playing
	window.SetUserPointer(unsafe.Pointer(cam))                // This is needed for the mouse callback

	// Instantiate all data needed for rendering
	cubeID, err := r.LoadData(cube)
	if err != nil {
		log.Fatalf("could not load data: %v\n", err)
	}
	planeID, err := r.LoadData(planeVertices)
	if err != nil {
		log.Fatalf("could not load data: %v\n", err)
	}

	// Load all default shaders
	if err := r.LoadDefaultPrograms(); err != nil {
		log.Fatalf("could not shaders: %v\n", err)
	}

	// Find the program by name
	shaderID, err := r.GetProgram("basic")
	if err != nil {
		log.Fatal(err)
	}

	// Prepare textures
	cubeTextureID, err := r.LoadTexture(path.Join(rootPath, marblePath), shaderID)
	if err != nil {
		log.Fatalf("could not create texture: %v\n", err)
	}

	floorTextureID, err := r.LoadTexture(path.Join(rootPath, metalPath), shaderID)
	if err != nil {
		log.Fatalf("could not create texture: %v\n", err)
	}

	gl.Enable(gl.DEPTH_TEST)

	// Scene perspective
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), float32(WIDTH)/HEIGHT, 0.1, 100)

	for !window.ShouldClose() {
		// Per-frame time. Used for speed normalization
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		r.Clear()
		processInput(window, cam)

		model := mgl32.Translate3D(-1, 0, -1) // These transformations will be on Node.Update()
		r.DrawRaw(cubeID, cubeTextureID, shaderID, cam, proj, model)
		model = mgl32.Translate3D(2, 0, 0)
		r.DrawRaw(cubeID, cubeTextureID, shaderID, cam, proj, model)
		model = mgl32.Ident4()
		r.DrawRaw(planeID, floorTextureID, shaderID, cam, proj, model)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	r.Destroy()
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

func processInput(w *glfw.Window, c *Camera) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.ProcessKeyboard(FORWARD, deltaTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.ProcessKeyboard(BACKWARD, deltaTime)
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		c.ProcessKeyboard(LEFT, deltaTime)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		c.ProcessKeyboard(RIGHT, deltaTime)
	}
}
func MouseCallback(w *glfw.Window, xpos, ypos float64) {
	// This solves the issue when the mouse enters the scene
	// and the camera immediately turns to that point instantly
	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := float32(xpos - lastX)
	yoffset := float32(ypos - lastY)

	lastX = xpos
	lastY = ypos

	// This is needed so we don't have a global camera variable
	cam := (*Camera)(w.GetUserPointer())
	cam.ProcessMouseMovement(xoffset, yoffset)
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

var planeVertices = []float32{
	5.0, -0.5, 5.0, 2.0, 0.0,
	-5.0, -0.5, 5.0, 0.0, 0.0,
	-5.0, -0.5, -5.0, 0.0, 2.0,

	5.0, -0.5, 5.0, 2.0, 0.0,
	-5.0, -0.5, -5.0, 0.0, 2.0,
	5.0, -0.5, -5.0, 2.0, 2.0,
}
