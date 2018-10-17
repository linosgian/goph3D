package main

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kpelelis/go-engine/objloader"
	"github.com/linosgian/goph3d/renderer"
	"github.com/linosgian/goph3d/scene"
	"github.com/linosgian/goph3d/window"
)

const (
	metalPath  = "res/textures/wood.png"
	marblePath = "res/textures/marble.jpg"

	FOV = 55.0
)

// This should be given temporarily because of vim-go
var rootPath = os.Getenv("PROJ_PATH") // e.g. /home/lgian/go/src/github.com/linosgian/goph3d

func init() {
	runtime.LockOSThread() // This is needed to arrange that main() runs on main thread.
}

func main() {
	//Initialize camera object at a certain position
	cam := scene.NewCamera(mgl32.Vec3{0, 2, 17}, window.WIDTH/2.0, window.HEIGHT/2.0)

	// Initializes Graphics API and GLFW
	// NOTE: A pointer to a camera is needed in order for the mouse callback to control it
	w, err := window.Init(cam, MouseCallback)
	if err != nil {
		log.Fatalf("Window initialization failed: %q\n", err)
	}

	// Temp: Light position
	lightPos := mgl32.Vec3{1.2, 2.3, 2}

	aspectRatio := float32(window.WIDTH) / window.HEIGHT

	// Create all Point lights
	lights := make([]*scene.PointLight, 0)
	for _, lpos := range PointLightPositions {
		l := &scene.PointLight{
			Position:  lpos,
			Ambient:   mgl32.Vec3{0.05, 0.05, 0.05},
			Diffuse:   mgl32.Vec3{0.8, 0.8, 0.8},
			Specular:  mgl32.Vec3{1, 1, 1},
			Constant:  1,
			Linear:    0.09,
			Quadratic: 0.032,
		}
		lights = append(lights, l)
	}
	sc := scene.NewScene(aspectRatio, cam, lights)

	r, err := renderer.NewRenderer()
	if err != nil {
		log.Fatalf("could not create renderer: %q\n", err)
	}

	// Instantiate all scene nodes and set their model matrices
	// -----------------------------
	reader, err := objloader.New("/home/lgian/Desktop/obj.obj")
	if err != nil {
		log.Fatalf("Could not read file: %q\n", err)
	}
	reader.Read()
	reader.Close()
	cube := reader.ExportToFloat32Array()[0]
	if err := sc.NewNodes(r, "crate", true, cube, path.Join(rootPath, marblePath), "phong", cubePositions); err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}

	if err := sc.NewNode(r, "lamp", true, cube, path.Join(rootPath, marblePath), "lamp", lightPos); err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}

	if err := sc.NewNode(r, "plane", true, planeVertices, path.Join(rootPath, metalPath), "phong", mgl32.Vec3{0, 0, 0}); err != nil {
		log.Fatalf("Could not create node: %q\n", err)
	}
	// ----------------------------

	gl.Enable(gl.DEPTH_TEST)

	sc.InitLights(r)
	for !w.ShouldClose() {
		// Per-frame time. Used for speed normalization
		currentFrame := glfw.GetTime()
		sc.DeltaTime = currentFrame - sc.LastFrame
		sc.LastFrame = currentFrame

		w.Clear()
		processInput(w, sc)

		// Update everything per-frame
		sc.Update(r)

		// Rotate all crates according to current time
		view := sc.Cam.GetViewMatrix()
		rot := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0, 1, 0})
		for _, n := range sc.Nodes {
			if n.Name == "crate" {
				translate := mgl32.Translate3D(n.Position.X(), n.Position.Y(), n.Position.Z())
				n.SetModelMatrix(translate.Mul4(rot))
			}
			r.DrawRaw(n.VaoID, n.ProgramID, n.TexID, view, sc.Perspective, n.ModelMatrix)
		}

		w.SwapBuffers()
		glfw.PollEvents()
	}
	w.Destroy()
}

func processInput(w *window.GlWindow, sc *scene.Scene) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
	if w.GetKey(glfw.KeyW) == glfw.Press {
		sc.Cam.ProcessKeyboard(scene.FORWARD, sc.DeltaTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		sc.Cam.ProcessKeyboard(scene.BACKWARD, sc.DeltaTime)
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		sc.Cam.ProcessKeyboard(scene.LEFT, sc.DeltaTime)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		sc.Cam.ProcessKeyboard(scene.RIGHT, sc.DeltaTime)
	}
	if w.GetKey(glfw.KeySpace) == glfw.Press {
		w.SetShouldClose(true)
	}
}

func MouseCallback(w *glfw.Window, xpos, ypos float64) {
	// This is needed so we don't have a global camera variable
	cam := (*scene.Camera)(w.GetUserPointer())

	// This solves the issue when the mouse enters the scene
	// and the camera immediately turns to that point instantly
	if cam.FirstMouse {
		cam.LastX = xpos
		cam.LastY = ypos
		cam.FirstMouse = false
	}

	xoffset := float32(xpos - cam.LastX)
	yoffset := float32(ypos - cam.LastY)

	cam.LastX = xpos
	cam.LastY = ypos

	cam.ProcessMouseMovement(xoffset, yoffset)
}

var cubePositions = []mgl32.Vec3{
	mgl32.Vec3{1, 1, 1},
	mgl32.Vec3{-4, 2, 3},
	mgl32.Vec3{-3, 3, 7},
	mgl32.Vec3{4, 7, 2},
	mgl32.Vec3{9, 4, 1},
}

var PointLightPositions = []mgl32.Vec3{
	mgl32.Vec3{0.7, 0.2, 2},
	mgl32.Vec3{2.3, 10, 4},
	mgl32.Vec3{4, 3, 5},
	mgl32.Vec3{1, 2, 12},
}

// var cube = []float32{
// 	// X,Y,Z / U,V / Nx, Ny, Nz
// 	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, 0.0, -1.0,
// 	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, -1.0,
// 	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, -1.0,
// 	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, -1.0,
// 	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, -1.0,
// 	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, 0.0, -1.0,

// 	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 1.0,
// 	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
// 	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 1.0,
// 	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 1.0,
// 	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
// 	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 1.0,

// 	-0.5, 0.5, 0.5, 1.0, 0.0, -1.0, 0.0, 0.0,
// 	-0.5, 0.5, -0.5, 1.0, 1.0, -1.0, 0.0, 0.0,
// 	-0.5, -0.5, -0.5, 0.0, 1.0, -1.0, 0.0, 0.0,
// 	-0.5, -0.5, -0.5, 0.0, 1.0, -1.0, 0.0, 0.0,
// 	-0.5, -0.5, 0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
// 	-0.5, 0.5, 0.5, 1.0, 0.0, -1.0, 0.0, 0.0,

// 	0.5, 0.5, 0.5, 1.0, 0.0, 1.0, 0.0, 0.0,
// 	0.5, 0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
// 	0.5, -0.5, -0.5, 0.0, 1.0, 1.0, 0.0, 0.0,
// 	0.5, -0.5, -0.5, 0.0, 1.0, 1.0, 0.0, 0.0,
// 	0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
// 	0.5, 0.5, 0.5, 1.0, 0.0, 1.0, 0.0, 0.0,

// 	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, -1.0, 0.0,
// 	0.5, -0.5, -0.5, 1.0, 1.0, 0.0, -1.0, 0.0,
// 	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, -1.0, 0.0,
// 	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, -1.0, 0.0,
// 	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, -1.0, 0.0,
// 	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, -1.0, 0.0,

// 	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
// 	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 1.0, 0.0,
// 	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
// 	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
// 	-0.5, 0.5, 0.5, 0.0, 0.0, 0.0, 1.0, 0.0,
// 	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
// }

var planeVertices = []float32{
	10.0, -0.5, 10.0, 10.0, 0.0, 0.0, 1.0, 0.0,
	-10.0, -0.5, 10.0, 0.0, 0.0, 0.0, 1.0, 0.0,
	-10.0, -0.5, -10.0, 0.0, 10.0, 0.0, 1.0, 0.0,

	10.0, -0.5, 10.0, 10.0, 0.0, 0.0, 1.0, 0.0,
	-10.0, -0.5, -10.0, 0.0, 10.0, 0.0, 1.0, 0.0,
	10.0, -0.5, -10.0, 10.0, 10.0, 0.0, 1.0, 0.0,
}
