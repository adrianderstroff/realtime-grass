package scene

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
)

type Postprocessing struct {
	width            float32
	height           float32
	blurradius       float32
	luminocityshader *engine.ShaderProgram
	gaussianshader   *engine.ShaderProgram
	bloomshader      *engine.ShaderProgram
	dofshader        *engine.ShaderProgram
	fogshader        *engine.ShaderProgram
	fboa             *engine.FBO
	fbob             *engine.FBO
	fbosmalla        *engine.FBO
	fbosmallb        *engine.FBO
	scale            float32
}

func MakePostprocessing(shaderpath string, width, height int32) (Postprocessing, error) {
	// setup shaders
	luminocityshader, err := engine.MakeProgram(shaderpath+"postprocessing/pass.vert", shaderpath+"postprocessing/luminocity.frag")
	if err != nil {
		return Postprocessing{}, err
	}
	gaussianshader, err := engine.MakeProgram(shaderpath+"postprocessing/pass.vert", shaderpath+"postprocessing/gaussian.frag")
	if err != nil {
		return Postprocessing{}, err
	}
	bloomshader, err := engine.MakeProgram(shaderpath+"postprocessing/pass.vert", shaderpath+"postprocessing/bloom.frag")
	if err != nil {
		return Postprocessing{}, err
	}
	dofshader, err := engine.MakeProgram(shaderpath+"postprocessing/pass.vert", shaderpath+"postprocessing/dof.frag")
	if err != nil {
		return Postprocessing{}, err
	}
	fogshader, err := engine.MakeProgram(shaderpath+"postprocessing/pass.vert", shaderpath+"postprocessing/fog.frag")
	if err != nil {
		return Postprocessing{}, err
	}

	// add cube to all shaders
	cube := engine.MakeCube(1, 1, 1)
	luminocityshader.AddRenderable(cube)
	gaussianshader.AddRenderable(cube)
	bloomshader.AddRenderable(cube)
	dofshader.AddRenderable(cube)
	fogshader.AddRenderable(cube)

	// setup fbos
	var scale float32 = 8.0
	fboa := engine.MakeFBO(width, height)
	fbob := engine.MakeFBO(width, height)
	fbosmalla := engine.MakeFBO(int32(float32(width)/scale), int32(float32(height)/scale))
	fbosmallb := engine.MakeFBO(int32(float32(width)/scale), int32(float32(height)/scale))

	return Postprocessing{
		width:            float32(width),
		height:           float32(height),
		blurradius:       1.0,
		luminocityshader: &luminocityshader,
		gaussianshader:   &gaussianshader,
		bloomshader:      &bloomshader,
		dofshader:        &dofshader,
		fogshader:        &fogshader,
		fboa:             &fboa,
		fbob:             &fbob,
		fbosmalla:        &fbosmalla,
		fbosmallb:        &fbosmallb,
		scale:            scale,
	}, nil
}

func (pp *Postprocessing) Bloom(fbo *engine.FBO) {
	// luminocity calculating
	pp.fboa.Bind()
	pp.fboa.Clear()
	fbo.ColorTextures[0].Bind(0)
	pp.luminocityshader.Use()
	pp.luminocityshader.UpdateFloat32("threshold", 0.85)
	pp.luminocityshader.Render()
	fbo.ColorTextures[0].Unbind()
	pp.fboa.Unbind()

	// gaussian blur in x and y
	pp.gaussian(pp.fboa, pp.fbob, pp.blurradius, pp.width, mgl32.Vec2{1, 0})
	pp.gaussian(pp.fbob, pp.fboa, pp.blurradius, pp.height, mgl32.Vec2{0, 1})

	// final bloom
	pp.fbob.Bind()
	pp.fbob.Clear()
	fbo.ColorTextures[0].Bind(0)
	pp.fboa.ColorTextures[0].Bind(1)
	pp.bloomshader.Use()
	pp.bloomshader.Render()
	fbo.ColorTextures[0].Unbind()
	pp.fboa.ColorTextures[0].Unbind()
	pp.fbob.Unbind()

	// copy tex to output
	pp.fbob.CopyColorToFBO(fbo, 0, 0, int32(pp.width), int32(pp.height))
}
func (pp *Postprocessing) Fog(fbo *engine.FBO, camera *engine.CameraFPS) {
	// get inverse view projection matrix
	invviewproj := camera.GetViewPerspective().Inv()

	// fog calculating
	pp.fboa.Bind()
	pp.fboa.Clear()
	fbo.ColorTextures[0].Bind(0)
	fbo.DepthTexture.Bind(1)
	pp.fogshader.Use()
	pp.fogshader.UpdateFloat32("zNear", camera.Near)
	pp.fogshader.UpdateFloat32("zFar", camera.Far)
	pp.fogshader.UpdateVec3("cameraPos", camera.Pos)
	pp.fogshader.UpdateMat4("InvViewProj", invviewproj)
	pp.fogshader.UpdateVec3("lightDir", mgl32.Vec3{0.6, -0.8, 2.0})
	pp.fogshader.UpdateFloat32("lightIntensity", 12.0)
	pp.fogshader.UpdateFloat32("fogDensity", 1.5)
	pp.fogshader.Render()
	fbo.ColorTextures[0].Unbind()
	fbo.DepthTexture.Unbind()
	pp.fboa.Unbind()

	// copy result
	pp.fboa.CopyColorToFBO(fbo, 0, 0, int32(pp.width), int32(pp.height))
}
func (pp *Postprocessing) DOF(fbo *engine.FBO) {
	// gaussian blur
	pp.gaussian2d(fbo, pp.fbob, pp.width, pp.height)

	// depth of field
	pp.fboa.Bind()
	pp.fboa.Clear()
	fbo.ColorTextures[0].Bind(0)
	fbo.DepthTexture.Bind(1)
	pp.fbob.ColorTextures[0].Bind(2)
	pp.dofshader.Use()
	pp.dofshader.UpdateFloat32("zNear", 0.1)
	pp.dofshader.UpdateFloat32("zFar", 4000.0)
	pp.dofshader.UpdateFloat32("focal", 0.0)
	pp.dofshader.UpdateFloat32("range", 3000.0)
	pp.dofshader.Render()
	fbo.ColorTextures[0].Unbind()
	fbo.DepthTexture.Unbind()
	pp.fbob.ColorTextures[0].Unbind()
	pp.fboa.Unbind()

	// copy tex to output
	pp.fboa.CopyColorToFBO(fbo, 0, 0, int32(pp.width), int32(pp.height))
}

func (pp *Postprocessing) gaussian2d(in, out *engine.FBO, width, height float32) {
	w := int32(width)
	h := int32(height)
	wsmall := int32(width / pp.scale)
	hsmall := int32(height / pp.scale)

	// downscale
	in.CopyColorToFBORegionSmooth(pp.fbosmalla, 0, 0, w, h, 0, 0, wsmall, hsmall)
	gl.Viewport(0, 0, wsmall, hsmall)
	// euler in x
	pp.gaussian(pp.fbosmalla, pp.fbosmallb, 1.0, float32(wsmall), mgl32.Vec2{1, 0})
	pp.gaussian(pp.fbosmallb, pp.fbosmalla, 1.0, float32(hsmall), mgl32.Vec2{0, 1})
	// upscale
	gl.Viewport(0, 0, w, h)
	pp.fbosmalla.CopyColorToFBORegionSmooth(out, 0, 0, wsmall, hsmall, 0, 0, w, h)
}
func (pp *Postprocessing) gaussian(in, out *engine.FBO, radius, dimension float32, dir mgl32.Vec2) {
	out.Bind()
	out.Clear()
	in.ColorTextures[0].Bind(0)
	pp.gaussianshader.Use()
	pp.gaussianshader.UpdateFloat32("resolution", dimension)
	pp.gaussianshader.UpdateFloat32("radius", radius)
	pp.gaussianshader.UpdateVec2("dir", dir)
	pp.gaussianshader.Render()
	in.ColorTextures[0].Unbind()
	out.Unbind()
}
