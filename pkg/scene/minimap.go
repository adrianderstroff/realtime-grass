// Package scene contains all main entities for rendering and/or interaction with the user.
package scene

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/collision"
	"github.com/adrianderstroff/realtime-grass/pkg/engine"
)

// Minimap has a camera that renders the scene from the top perspective using an orthographic projection.
// In addition will the view frustum by visualized.
type Minimap struct {
	camera engine.CameraTrackball
	fbo    engine.FBO
	shader engine.ShaderProgram
	width  int32
	height int32
}

// MakeMinimap constructs a Minimap with the viewport specified by width and height.
func MakeMinimap(shaderpath string, width, height int32) Minimap {
	// setup camera
	camera := engine.MakeCameraTrackball(int(width), int(height), 1400, mgl32.Vec3{0, 0, 0}, 90, 0.1, 2500.0)
	camera.Rotate(-90.0, 0.0)

	// create fbo
	fbo := engine.MakeFBO(width, height)

	// create wireframe shader for the view frustum of the other camera used in the scene
	wireframeshader, err := engine.MakeProgram(shaderpath+"/wireframe/wireframe.vert", shaderpath+"/wireframe/wireframe.frag")
	if err != nil {
		panic(err)
	}
	frustum := collision.MakeFrustum(camera.Near, camera.Far, camera.Fov)
	wireframeshader.AddRenderable(frustum)

	return Minimap{
		camera: camera,
		fbo:    fbo,
		shader: wireframeshader,
		width:  width,
		height: height,
	}
}

// Update sets the minimap-camera to the x-z position of othercam.
func (minimap *Minimap) Update(othercam *engine.CameraFPS) {
	// update camera
	minimap.camera.SetPos(mgl32.Vec3{othercam.Pos.X(), minimap.camera.Target.Y(), othercam.Pos.Z()})
	minimap.camera.Update()
}

// Render displays the othercam's view frustum as a wireframe.
func (minimap *Minimap) Render(othercam *engine.CameraFPS) {
	M := othercam.GetView().Inv()
	V := minimap.camera.GetView()
	P := minimap.camera.GetOrtho()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	minimap.shader.Use()
	minimap.shader.UpdateMat4("M", M)
	minimap.shader.UpdateMat4("V", V)
	minimap.shader.UpdateMat4("P", P)
	minimap.shader.UpdateFloat32("width", 5)
	minimap.shader.Render()
	gl.Disable(gl.BLEND)
}

// Begin binds the Minimap's FBO.
// The next entity that is being rendered will appear in the Minimap.
func (minimap *Minimap) Begin() {
	minimap.fbo.Bind()
	minimap.fbo.Clear()
}

// End unbinds the FBO and writes it's contents to the screen.
func (minimap *Minimap) End() {
	minimap.fbo.Unbind()
	minimap.fbo.CopyToScreenRegion(
		0,
		0, 0, minimap.width, minimap.height,
		0, 0, minimap.width/4, minimap.height/4,
	)
}
