package engine

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type CameraFPS struct {
	width  int
	height int
	theta  float32
	phi    float32
	dir    mgl32.Vec3
	speed  float32

	Pos    mgl32.Vec3
	Target mgl32.Vec3
	Up     mgl32.Vec3
	Right  mgl32.Vec3
	Fov    float32
	Near   float32
	Far    float32
}

func MakeCameraFPS(width, height int, pos mgl32.Vec3, speed, fov, near, far float32) CameraFPS {
	dir := mgl32.Vec3{0.0, 0.0, 1.0}
	camera := CameraFPS{
		width:  width,
		height: height,
		theta:  90.0,
		phi:    0.0,
		dir:    dir,
		speed:  speed,

		Pos:    pos,
		Target: pos.Add(dir),
		Up:     mgl32.Vec3{0, 1, 0},
		Right:  mgl32.Vec3{1, 0, 0},
		Fov:    fov,
		Near:   near,
		Far:    far,
	}
	camera.Update()

	return camera
}

func (camera *CameraFPS) Update() {
	theta := mgl32.DegToRad(camera.theta)
	phi := mgl32.DegToRad(camera.phi)

	// sphere coordinates with inverse y
	btheta := float64(theta)
	bphi := float64(phi)
	camera.dir = mgl32.Vec3{
		float32(math.Sin(btheta) * math.Cos(bphi)),
		-float32(math.Cos(btheta)),
		float32(math.Sin(btheta) * math.Sin(bphi)),
	}
	camera.dir = camera.dir.Normalize()

	// set target
	camera.Target = camera.Pos.Add(camera.dir)

	// calculate up vector
	look := camera.dir.Mul(-1)
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
	camera.Right = worldUp.Cross(look).Normalize()
	camera.Up = look.Cross(camera.Right)
}
func (camera *CameraFPS) Rotate(theta, phi float32) {
	camera.theta += theta
	camera.phi += phi

	// limit angles
	camera.theta = float32(math.Max(math.Min(float64(camera.theta), 179.9), 0.01))
	if camera.phi < 0 {
		camera.phi = 360 + camera.phi
	} else if camera.phi >= 360 {
		camera.phi = camera.phi - 360
	}
}
func (camera *CameraFPS) Zoom(distance float32) {}

func (camera *CameraFPS) GetView() mgl32.Mat4 {
	return mgl32.LookAtV(camera.Pos, camera.Target, camera.Up)
}
func (camera *CameraFPS) GetPerspective() mgl32.Mat4 {
	fov := mgl32.DegToRad(camera.Fov)
	aspect := float32(camera.width) / float32(camera.height)
	return mgl32.Perspective(fov, aspect, camera.Near, camera.Far)
}
func (camera *CameraFPS) GetOrtho() mgl32.Mat4 {
	angle := camera.Fov * math.Pi / 180.0
	dfar := float32(math.Tan(float64(angle/2.0))) * camera.Far
	d := dfar
	return mgl32.Ortho(-d, d, -d, d, camera.Near, camera.Far)
}
func (camera *CameraFPS) GetViewPerspective() mgl32.Mat4 {
	return camera.GetPerspective().Mul4(camera.GetView())
}
func (camera *CameraFPS) SetPos(pos mgl32.Vec3) {
	camera.Pos = pos
	camera.Target = camera.Pos.Add(camera.dir)
}

func (camera *CameraFPS) OnCursorPosMove(x, y, dx, dy float64) bool {
	dPhi := float32(-dx) / 2.0
	dTheta := float32(-dy) / 2.0
	camera.Rotate(dTheta, -dPhi)
	return false
}
func (camera *CameraFPS) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}
func (camera *CameraFPS) OnMouseScroll(x, y float64) bool {
	return false
}
func (camera *CameraFPS) OnKeyPress(key, action, mods int) bool {
	dir := camera.dir.Mul(camera.speed)
	right := camera.Right.Mul(camera.speed)
	if key == int(glfw.KeyW) {
		camera.Pos = camera.Pos.Add(dir)
	} else if key == int(glfw.KeyS) {
		camera.Pos = camera.Pos.Sub(dir)
	} else if key == int(glfw.KeyA) {
		camera.Pos = camera.Pos.Sub(right)
	} else if key == int(glfw.KeyD) {
		camera.Pos = camera.Pos.Add(right)
	}
	return false
}
