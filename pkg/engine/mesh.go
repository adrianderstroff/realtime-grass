package engine

import "github.com/go-gl/gl/v4.3-core/gl"

type Mesh struct {
	vao VAO
}

func MakeEmptyMesh(mode uint32) Mesh {
	return Mesh{
		vao: MakeVAO(mode),
	}
}
func MakeMesh(filepath string, mode uint32) (Mesh, error) {
	// load obj
	obj, err := LoadObj(filepath)
	if err != nil {
		return Mesh{}, err
	}

	// extract properties and create buffers
	vertexBuffer := MakeVBO(obj.Vertices, 3, gl.STATIC_DRAW)
	vertexBuffer.AddVertexAttribute("vert", 3, gl.FLOAT)
	texCoordBuffer := MakeVBO(obj.Texcoords, 2, gl.STATIC_DRAW)
	texCoordBuffer.AddVertexAttribute("uv", 2, gl.FLOAT)
	normalBuffer := MakeVBO(obj.Normals, 3, gl.STATIC_DRAW)
	normalBuffer.AddVertexAttribute("normal", 3, gl.FLOAT)

	// combine buffers in vao
	vao := MakeVAO(mode)
	vao.AddVertexBuffer(&vertexBuffer)
	vao.AddVertexBuffer(&texCoordBuffer)
	vao.AddVertexBuffer(&normalBuffer)

	return Mesh{vao}, nil
}
func MakeSimpleMesh(positions []float32, pcount uint32, mode, usage uint32) (Mesh, error) {
	// make vao
	vao := MakeVAO(mode)

	// add position buffer
	vertexBuffer := MakeVBO(positions, pcount, usage)
	vertexBuffer.AddVertexAttribute("position", int32(pcount), gl.FLOAT)
	vao.AddVertexBuffer(&vertexBuffer)

	return Mesh{vao}, nil
}
func MakeMeshFromArrays(positions, normals, texcoords []float32, pname, nname, tname string, pcount, ncount, tcount uint32, mode uint32) (Mesh, error) {
	// make vao
	vao := MakeVAO(mode)

	// extract properties and create buffers
	if len(positions) > 0 {
		vertexBuffer := MakeVBO(positions, pcount, gl.STATIC_DRAW)
		vertexBuffer.AddVertexAttribute(pname, int32(pcount), gl.FLOAT)
		vao.AddVertexBuffer(&vertexBuffer)
	}
	if len(normals) > 0 {
		normalBuffer := MakeVBO(normals, ncount, gl.STATIC_DRAW)
		normalBuffer.AddVertexAttribute(nname, int32(ncount), gl.FLOAT)
		vao.AddVertexBuffer(&normalBuffer)
	}
	if len(texcoords) > 0 {
		texCoordBuffer := MakeVBO(texcoords, tcount, gl.STATIC_DRAW)
		texCoordBuffer.AddVertexAttribute(tname, int32(tcount), gl.FLOAT)
		vao.AddVertexBuffer(&texCoordBuffer)
	}

	return Mesh{vao}, nil
}
func (mesh *Mesh) Delete() {
	mesh.vao.Delete()
}
func (mesh Mesh) Build(shaderProgramHandle uint32) {
	mesh.vao.BuildBuffers(shaderProgramHandle)
}
func (mesh Mesh) Render() {
	mesh.vao.Render()
}
func (mesh Mesh) RenderInstanced(instancecount int32) {
	mesh.vao.RenderInstanced(instancecount)
}
func (mesh *Mesh) GetVAO() *VAO {
	return &mesh.vao
}
func (mesh *Mesh) SetVAO(vao VAO) {
	mesh.vao = vao
}

func combine(slices ...[]float32) []float32 {
	var result []float32
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}
func repeat(slice []float32, number int) []float32 {
	var result []float32
	for i := 0; i < number; i++ {
		result = append(result, slice...)
	}
	return result
}
