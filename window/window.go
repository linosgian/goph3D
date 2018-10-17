package window

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/linosgian/goph3d/scene"
)

const (
	WIDTH  = 1360
	HEIGHT = 768
)

// Struct embedding so we can add more functionality
type GlWindow struct {
	*glfw.Window
}

func NewWindow(w *glfw.Window) *GlWindow {
	return &GlWindow{w}
}

func Init(c *scene.Camera, mcallback func(w *glfw.Window, xpos, ypos float64)) (*GlWindow, error) {
	gw, err := initGLFW()
	if err != nil {
		return nil, err
	}
	if err := initOpenGL(); err != nil {
		return nil, err
	}
	gw.InitInputs(c, mcallback)
	return gw, nil
}

func initGLFW() (*GlWindow, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// Specify profile and OpenGL version
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w, err := glfw.CreateWindow(WIDTH, HEIGHT, "3D Gamez", nil, nil)
	if err != nil {
		return nil, err
	}
	w.MakeContextCurrent()

	return NewWindow(w), nil
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

func (gw *GlWindow) InitInputs(c *scene.Camera, mcallback func(w *glfw.Window, xpos, ypos float64)) {
	gw.SetCursorPosCallback(mcallback)
	gw.SetInputMode(glfw.CursorMode, glfw.CursorDisabled) // Disable mouse pointer while playing
	gw.SetUserPointer(unsafe.Pointer(c))                  // This is needed for the mouse callback
}

func (gw *GlWindow) Destroy() {
	glfw.Terminate()
}

func (gw *GlWindow) Clear() {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0) // Default scene color.
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
