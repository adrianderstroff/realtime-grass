package scene

import (
	"math/rand"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
)

type Grass struct {
	shader        engine.ShaderProgram
	buffer        engine.Mesh
	grassAlpha    engine.Texture
	grassDiffuse0 engine.Texture
	grassDiffuse1 engine.Texture
	grassDiffuse2 engine.Texture
	grassDiffuse3 engine.Texture
	bladecount    int32
	height        float32
	viewdist      float32
	time          float32
	windradius    int32
}

func MakeGrass(shaderpath, texpath string, bladecount int, height, viewdist float32, windradius int32) (Grass, error) {
	// make shader
	shader, err := engine.MakeGeomProgram(shaderpath+"/grass/grass.vert", shaderpath+"/grass/grass.geom", shaderpath+"/grass/grass.frag")
	if err != nil {
		return Grass{}, err
	}

	// generate random 2D root positions
	var positions []float32
	for i := 0; i < bladecount; i++ {
		positions = append(positions, rand.Float32(), rand.Float32())
	}

	// generate buffer
	positionBuffer := engine.MakeVBO(positions, 2, gl.STATIC_DRAW)
	positionBuffer.AddVertexAttribute("position", 2, gl.FLOAT)
	vao := engine.MakeVAO(gl.POINTS)
	vao.AddVertexBuffer(&positionBuffer)
	mesh := engine.MakeEmptyMesh(gl.POINTS)
	mesh.SetVAO(vao)
	shader.AddRenderable(mesh)

	// load grass texture
	grassalpha, err := engine.MakeTextureFromPath(texpath + "grassAlpha.png")
	if err != nil {
		return Grass{}, err
	}
	grassdiffuse0, err := engine.MakeTextureFromPath(texpath + "grass0.jpg")
	if err != nil {
		return Grass{}, err
	}
	grassdiffuse1, err := engine.MakeTextureFromPath(texpath + "grass1.jpg")
	if err != nil {
		return Grass{}, err
	}
	grassdiffuse2, err := engine.MakeTextureFromPath(texpath + "grass2.jpg")
	if err != nil {
		return Grass{}, err
	}
	grassdiffuse3, err := engine.MakeTextureFromPath(texpath + "grass3.jpg")
	if err != nil {
		return Grass{}, err
	}
	// generate mipmaps
	grassalpha.GenMipmap()
	grassdiffuse0.GenMipmapNearest()
	grassdiffuse1.GenMipmap()
	grassdiffuse2.GenMipmap()
	grassdiffuse3.GenMipmap()

	return Grass{
		shader,
		mesh,
		grassalpha,
		grassdiffuse0,
		grassdiffuse1,
		grassdiffuse2,
		grassdiffuse3,
		int32(bladecount),
		height,
		viewdist,
		0.0,
		windradius,
	}, nil
}

func (grass *Grass) Render(instancecount int32, tilesize float32, M, V, P mgl32.Mat4, camerapos mgl32.Vec3) {
	lightdir := mgl32.Vec3{10.0, 0.0, 10.0}
	lightcolor := mgl32.Vec3{1.0, 1.0, 0.0}

	grass.grassAlpha.Bind(0)
	grass.grassDiffuse0.Bind(1)
	grass.grassDiffuse1.Bind(2)
	grass.grassDiffuse2.Bind(3)
	grass.grassDiffuse3.Bind(4)

	// render terrain
	grass.shader.Use()
	grass.shader.UpdateMat4("M", M)
	grass.shader.UpdateMat4("V", V)
	grass.shader.UpdateMat4("P", P)
	grass.shader.UpdateFloat32("grassHeight", grass.height)
	grass.shader.UpdateInt32("bladeCount", grass.bladecount)
	grass.shader.UpdateFloat32("tilesize", tilesize)
	grass.shader.UpdateVec3("cameraPos", camerapos)
	grass.shader.UpdateVec3("lightDir", lightdir)
	grass.shader.UpdateVec3("lightColor", lightcolor)
	grass.shader.UpdateFloat32("ambientIntensity", 0.4)
	grass.shader.UpdateFloat32("diffuseIntensity", 0.4)
	grass.shader.UpdateFloat32("d1", grass.viewdist/8)
	grass.shader.UpdateFloat32("d2", grass.viewdist)
	grass.shader.UpdateFloat32("t", grass.time)
	// wind related uniforms
	grass.shader.UpdateInt32("radius", grass.windradius)
	grass.shader.RenderInstanced(instancecount)

	grass.grassAlpha.Unbind()
	grass.grassDiffuse0.Unbind()
	grass.grassDiffuse1.Unbind()
	grass.grassDiffuse2.Unbind()
	grass.grassDiffuse3.Unbind()

	// update time
	grass.time++
}
