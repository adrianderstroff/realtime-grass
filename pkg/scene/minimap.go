package scene

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/collision"
	"github.com/adrianderstroff/realtime-grass/pkg/engine"
)

type Minimap struct {
	camera engine.CameraTrackball
	fbo    engine.FBO
	shader engine.ShaderProgram
	width  int32
	height int32
}

func MakeMinimap(shaderpath string, width, height int32) Minimap {
	// setup camera
	camera := engine.MakeCameraTrackball(int(width), int(height), 1400, mgl32.Vec3{0, 0, 0}, 90, 0.1, 2500.0)
	camera.Rotate(-90.0, 0.0)

	// create fbo
	fbo := engine.MakeFBO(width, height)

	// create wireframe shader for frustum
	wireframeShader, err := engine.MakeProgram(shaderpath+"/wireframe/wireframe.vert", shaderpath+"/wireframe/wireframe.frag")
	if err != nil {
		panic(err)
	}
	frustum := collision.MakeFrustum(camera.Near, camera.Far, camera.Fov)
	wireframeShader.AddRenderable(frustum)

	return Minimap{
		camera,
		fbo,
		wireframeShader,
		width,
		height,
	}
}
func (minimap *Minimap) Update(othercam *engine.CameraFPS) {
	// update camera
	minimap.camera.SetPos(mgl32.Vec3{othercam.Pos.X(), minimap.camera.Target.Y(), othercam.Pos.Z()})
	minimap.camera.Update()
}
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
func (minimap *Minimap) Begin() {
	minimap.fbo.Bind()
	minimap.fbo.Clear()
}
func (minimap *Minimap) End() {
	minimap.fbo.Unbind()
	minimap.fbo.CopyToScreenRegion(
		0,
		0, 0, minimap.width, minimap.height,
		0, 0, minimap.width/4, minimap.height/4,
	)
}
