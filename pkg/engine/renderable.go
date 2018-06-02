package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type Renderable interface {
	Build(shaderProgramHandle uint32)
	Render()
	RenderInstanced(instancecount int32)
}

type SquarePyramid struct {
	vao VAO
}

func MakeSquarePyramid(width, depth, height float32) SquarePyramid {
	// vertices
	v1 := []float32{-1.0, 0.0, -1.0}
	v2 := []float32{-1.0, 0.0, 1.0}
	v3 := []float32{1.0, 0.0, 1.0}
	v4 := []float32{1.0, 0.0, -1.0}
	v5 := []float32{0.0, 1.0, 0.0}
	vertices := combine(
		// bottom
		v1, v2, v3,
		v1, v3, v4,
		// left
		v1, v5, v2,
		// front
		v1, v4, v5,
		// right
		v4, v3, v5,
		// back
		v3, v2, v5,
	)
	vertexBuffer := MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vertexBuffer.AddVertexAttribute("vert", 3, gl.FLOAT)

	// colors
	black := []float32{0.0, 0.0, 0.0}
	colors := repeat(black, 18)
	colorBuffer := MakeVBO(colors, 3, gl.STATIC_DRAW)
	colorBuffer.AddVertexAttribute("color", 3, gl.FLOAT)

	vao := MakeVAO(gl.TRIANGLES)
	vao.AddVertexBuffer(&vertexBuffer)
	vao.AddVertexBuffer(&colorBuffer)

	return SquarePyramid{vao}
}
func (squarePyramid *SquarePyramid) Delete() {
	squarePyramid.vao.Delete()
}
func (squarePyramid SquarePyramid) Build(shaderProgramHandle uint32) {
	squarePyramid.vao.BuildBuffers(shaderProgramHandle)
}
func (squarePyramid SquarePyramid) Render() {
	squarePyramid.vao.Render()
}
func (squarePyramid SquarePyramid) RenderInstanced(instancecount int32) {
	squarePyramid.vao.RenderInstanced(instancecount)
}

type Cube struct {
	vao VAO
}

func MakeCube(halfWidth, halfHeight, halfDepth float32) Cube {
	// vertex positions
	v1 := []float32{-halfWidth, halfHeight, halfDepth}
	v2 := []float32{-halfWidth, -halfHeight, halfDepth}
	v3 := []float32{halfWidth, halfHeight, halfDepth}
	v4 := []float32{halfWidth, -halfHeight, halfDepth}
	v5 := []float32{-halfWidth, halfHeight, -halfDepth}
	v6 := []float32{-halfWidth, -halfHeight, -halfDepth}
	v7 := []float32{halfWidth, halfHeight, -halfDepth}
	v8 := []float32{halfWidth, -halfHeight, -halfDepth}
	vertices := combine(
		// front
		v1, v2, v3,
		v3, v2, v4,
		// back
		v7, v8, v5,
		v5, v8, v6,
		// left
		v5, v6, v1,
		v1, v6, v2,
		// right
		v3, v4, v7,
		v7, v4, v8,
		// top
		v5, v1, v7,
		v7, v1, v3,
		// bottom
		v2, v6, v4,
		v4, v6, v8,
	)
	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := repeat(combine(t1, t2, t3, t3, t2, t4), 6)
	// normals
	right := []float32{1.0, 0.0, 0.0}
	left := []float32{-1.0, 0.0, 0.0}
	top := []float32{0.0, 1.0, 0.0}
	bottom := []float32{0.0, -1.0, 0.0}
	front := []float32{0.0, 0.0, -1.0}
	back := []float32{0.0, 0.0, 1.0}
	normals := combine(
		repeat(bottom, 6),
		repeat(top, 6),
		repeat(left, 6),
		repeat(right, 6),
		repeat(front, 6),
		repeat(back, 6),
	)

	vertexBuffer := MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vertexBuffer.AddVertexAttribute("vert", 3, gl.FLOAT)
	uvBuffer := MakeVBO(uvs, 2, gl.STATIC_DRAW)
	uvBuffer.AddVertexAttribute("uv", 2, gl.FLOAT)
	normalBuffer := MakeVBO(normals, 3, gl.STATIC_DRAW)
	normalBuffer.AddVertexAttribute("normal", 3, gl.FLOAT)

	vao := MakeVAO(gl.TRIANGLES)
	vao.AddVertexBuffer(&vertexBuffer)
	vao.AddVertexBuffer(&uvBuffer)
	vao.AddVertexBuffer(&normalBuffer)

	return Cube{vao}
}
func (cube *Cube) Delete() {
	cube.vao.Delete()
}
func (cube Cube) Build(shaderProgramHandle uint32) {
	cube.vao.BuildBuffers(shaderProgramHandle)
}
func (cube Cube) Render() {
	cube.vao.Render()
}
func (cube Cube) RenderInstanced(instancecount int32) {
	cube.vao.RenderInstanced(instancecount)
}

type Skybox struct {
	vao     VAO
	texture Texture
}

func MakeSkybox(cubeTexture Texture) Skybox {
	var size float32 = 1.0
	v0 := []float32{-size, -size, -size}
	v1 := []float32{-size, -size, size}
	v2 := []float32{size, -size, size}
	v3 := []float32{size, -size, -size}
	v4 := []float32{-size, size, -size}
	v5 := []float32{-size, size, size}
	v6 := []float32{size, size, size}
	v7 := []float32{size, size, -size}
	vertices := combine(
		// right face
		v2, v7, v3, v2, v6, v7,
		// left face
		v0, v4, v5, v0, v5, v1,
		// top face
		v7, v6, v5, v7, v5, v4,
		// bottom face
		v0, v1, v2, v0, v2, v3,
		// back face
		v0, v7, v4, v0, v3, v7,
		// front face
		v6, v2, v5, v5, v2, v1,
	)
	vbo := MakeVBO(vertices, 3, gl.STATIC_DRAW)
	vbo.AddVertexAttribute("vert", 3, gl.FLOAT)
	vao := MakeVAO(gl.TRIANGLES)
	vao.AddVertexBuffer(&vbo)

	return Skybox{vao, cubeTexture}
}
func (skybox *Skybox) Delete() {
	skybox.vao.Delete()
	skybox.texture.Delete()
}
func (skybox Skybox) Build(shaderProgramHandle uint32) {
	skybox.vao.BuildBuffers(shaderProgramHandle)
}
func (skybox Skybox) Render() {
	gl.DepthMask(false)
	skybox.texture.Bind(0)
	skybox.vao.Render()
	skybox.texture.Unbind()
	gl.DepthMask(true)
}
func (skybox Skybox) RenderInstanced(instancecount int32) {
	gl.DepthMask(false)
	skybox.texture.Bind(0)
	skybox.vao.RenderInstanced(instancecount)
	skybox.texture.Unbind()
	gl.DepthMask(true)
}
