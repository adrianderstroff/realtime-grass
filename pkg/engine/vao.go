package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type VAO struct {
	handle        uint32
	mode          uint32
	vertexBuffers []*VBO
	indexBuffer   *IBO
}

// Creates a new VAO.
// 'mode' specified the drawing mode used.
// Some modes would be TRIANGLE, TRIANGLE_STRIP, TRIANGLE_FAN
func MakeVAO(mode uint32) VAO {
	vao := VAO{0, mode, nil, nil}
	gl.GenVertexArrays(1, &vao.handle)
	return vao
}
func (vao *VAO) Delete() {
	// delete buffers
	if vao.vertexBuffers != nil {
		for _, vertBuf := range vao.vertexBuffers {
			vertBuf.Delete()
		}
	}
	vao.indexBuffer.Delete()

	// delete vertex array
	gl.DeleteVertexArrays(1, &vao.handle)
}
func (vao *VAO) Render() {
	gl.BindVertexArray(vao.handle)
	if vao.indexBuffer != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.indexBuffer.handle)
		gl.DrawElements(vao.mode, vao.indexBuffer.count, gl.UNSIGNED_SHORT, nil)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		gl.DrawArrays(vao.mode, 0, vao.vertexBuffers[0].count)
	}
	gl.BindVertexArray(0)
}
func (vao *VAO) RenderInstanced(instancecount int32) {
	gl.BindVertexArray(vao.handle)
	if vao.indexBuffer != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vao.indexBuffer.handle)
		gl.DrawElementsInstanced(vao.mode, vao.indexBuffer.count, gl.UNSIGNED_SHORT, nil, instancecount)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		gl.DrawArraysInstanced(vao.mode, 0, vao.vertexBuffers[0].count, instancecount)
	}
	gl.BindVertexArray(0)
}

func (vao *VAO) AddVertexBuffer(vbo *VBO) {
	vao.vertexBuffers = append(vao.vertexBuffers, vbo)
}
func (vao *VAO) AddIndexBuffer(ibo *IBO) {
	vao.indexBuffer = ibo
}
func (vao *VAO) GetVertexBuffer(idx int) *VBO {
	return vao.vertexBuffers[idx]
}
func (vao *VAO) GetIndexBuffer() *IBO {
	return vao.indexBuffer
}
func (vao *VAO) BuildBuffers(shaderProgramHandle uint32) {
	gl.BindVertexArray(vao.handle)
	for _, vbo := range vao.vertexBuffers {
		vbo.BuildVertexAttributes(shaderProgramHandle)
	}
	gl.BindVertexArray(0)
}
