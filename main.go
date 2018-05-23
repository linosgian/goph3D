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
	metalPath   = "res/textures/metal.png"
	marblePath  = "res/textures/marble.jpg"
	shadersPath = "res/shaders"

	FLOAT  = gl.FLOAT
	UINT32 = gl.UNSIGNED_INT

	WIDTH  = 800 // These will be in a config file eventually
	HEIGHT = 600 // And won't be contants
	FOV    = 45.0
)

// Calculate the byte-size of all types at runtime
var sizes map[int]int = map[int]int{
	FLOAT:  int(unsafe.Sizeof(float32(0))),
	UINT32: int(unsafe.Sizeof(uint32(0))),
}

// This should be given temporarily because of vim-go
var rootPath = os.Getenv("PROJ_PATH") // e.g. /home/lgian/go/src/github.com/linosgian/glfw-test2

func init() {
	runtime.LockOSThread() // This is needed to arrange that main() runs on main thread.
}

func main() {
	// Initializes Graphics API and GLFW
	window := Init()

	//Initialize camera object at a certain position
	cam := NewCamera(mgl32.Vec3{0, 0, 7}, WIDTH/2.0, HEIGHT/2.0)
	// Scene perspective
	proj := mgl32.Perspective(mgl32.DegToRad(FOV), float32(WIDTH)/HEIGHT, 0.1, 100)

	sc := NewScene(proj, cam)

	r, err := NewRenderer(window, sc.Cam)
	if err != nil {
		log.Fatalf("could not create renderer: %q\n", err)
	}

	// Instantiate all scene nodes and set their model matrices
	Ncube, err := sc.NewNode(r, true, cube, path.Join(rootPath, marblePath), "basic")
	if err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}
	Ncube.SetModelMatrix(mgl32.Translate3D(-1, 0, -1))

	Ncube2, err := sc.NewNode(r, true, cube, path.Join(rootPath, marblePath), "basic")
	if err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}
	Ncube2.SetModelMatrix(mgl32.Translate3D(2, 0, 0))

	plane, err := sc.NewNode(r, true, planeVertices, path.Join(rootPath, metalPath), "basic")
	if err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}
	plane.SetModelMatrix(mgl32.Ident4())

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		// Per-frame time. Used for speed normalization
		currentFrame := glfw.GetTime()
		sc.deltaTime = currentFrame - sc.lastFrame
		sc.lastFrame = currentFrame

		r.Clear()
		processInput(window, sc)

		for _, n := range sc.Nodes {
			r.DrawRaw(n, sc.Cam, sc.Perspective)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
	r.Destroy()
}

func Init() *glfw.Window {
	window, err := initGLFW()
	if err != nil {
		// These could be propagated to the main function but we're gonna halt anyway
		log.Fatalf("[GLFW error]: %q", err)
	}
	if err := initOpenGL(); err != nil {
		log.Fatalf("OpenGL could not be initialized: %v\n", err)
	}
	return window
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

func processInput(w *glfw.Window, sc *Scene) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
	if w.GetKey(glfw.KeyW) == glfw.Press {
		sc.Cam.ProcessKeyboard(FORWARD, sc.deltaTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		sc.Cam.ProcessKeyboard(BACKWARD, sc.deltaTime)
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		sc.Cam.ProcessKeyboard(LEFT, sc.deltaTime)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		sc.Cam.ProcessKeyboard(RIGHT, sc.deltaTime)
	}
	if w.GetKey(glfw.KeySpace) == glfw.Press {
		w.SetShouldClose(true)
	}
}

func MouseCallback(w *glfw.Window, xpos, ypos float64) {
	// This is needed so we don't have a global camera variable
	cam := (*Camera)(w.GetUserPointer())

	// This solves the issue when the mouse enters the scene
	// and the camera immediately turns to that point instantly
	if cam.firstMouse {
		cam.lastX = xpos
		cam.lastY = ypos
		cam.firstMouse = false
	}

	xoffset := float32(xpos - cam.lastX)
	yoffset := float32(ypos - cam.lastY)

	cam.lastX = xpos
	cam.lastY = ypos

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
