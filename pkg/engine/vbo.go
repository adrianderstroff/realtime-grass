package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type VertexAttribute struct {
	name   string
	count  int32
	glType uint32
}

type VBO struct {
	handle     uint32
	count      int32
	stride     uint32
	usage      uint32
	attributes []VertexAttribute
}

func MakeVBO(data []float32, elementsPerVertex uint32, usage uint32) VBO {
	vbo := VBO{
		handle:     0,
		count:      int32(len(data)) / int32(elementsPerVertex),
		stride:     elementsPerVertex,
		usage:      usage,
		attributes: nil,
	}

	gl.GenBuffers(1, &vbo.handle)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
	if len(data) != 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), usage)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, 4, gl.Ptr([]float32{0.0}), usage)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vbo
}
func (vbo *VBO) UpdateData(data []float32) {
	// update buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), vbo.usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// update size
	vbo.count = int32(len(data)) / int32(vbo.stride)
}
func (vbo *VBO) Delete() {
	vbo.count = 0
	vbo.stride = 0
	vbo.attributes = nil
	gl.DeleteBuffers(1, &vbo.handle)
}

func (vbo *VBO) AddVertexAttribute(name string, count int32, glType uint32) {
	vbo.attributes = append(vbo.attributes, VertexAttribute{name, count, glType})
}
func (vbo *VBO) BuildVertexAttributes(shaderProgramHandle uint32) {
	// specify all vertex attributes
	var offset int = 0
	for _, attrib := range vbo.attributes {
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
		location := gl.GetAttribLocation(shaderProgramHandle, gl.Str(attrib.name+"\x00"))
		if location != -1 {
			gl.EnableVertexAttribArray(uint32(location))
			gl.VertexAttribPointer(uint32(location), attrib.count, attrib.glType, false, int32(vbo.stride*4), gl.PtrOffset(offset*4))
		}
		offset += int(attrib.count)
	}

	// unbind vbo to prevent overwrites
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
