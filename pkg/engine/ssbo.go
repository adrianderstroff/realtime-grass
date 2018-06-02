package engine

import (
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type SSBO struct {
	handle   uint32
	typesize int
	len      int
	pos      int32
}

func MakeSSBO(typesize, len int) SSBO {
	ssbo := SSBO{
		0,
		typesize,
		len,
		-1,
	}
	ssbo.init()
	return ssbo
}
func MakeEmptySSBO(typesize int) SSBO {
	ssbo := SSBO{
		0,
		typesize,
		0,
		-1,
	}
	ssbo.init()
	return ssbo
}

func (ssbo *SSBO) Delete() {
	// unbind if not done yet
	if ssbo.pos != -1 {
		ssbo.Unbind()
	}

	// delete buffer
	gl.DeleteBuffers(1, &ssbo.handle)
	ssbo.handle = 0
	ssbo.len = 0
}

func (ssbo *SSBO) Len() int {
	return ssbo.len
}
func (ssbo *SSBO) GetHandle() uint32 {
	return ssbo.handle
}

func (ssbo *SSBO) Bind(pos int32) {
	ssbo.pos = pos
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(ssbo.pos), ssbo.handle)
}
func (ssbo *SSBO) Unbind() {
	if ssbo.pos == -1 {
		return
	}

	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, uint32(ssbo.pos), 0)
	ssbo.pos = -1
}

func (ssbo *SSBO) Resize(len int) {
	// early return if size stayed the same
	if ssbo.len == len {
		return
	}

	// get the smaller size of old and new buffer
	minsize := int(math.Min(float64(ssbo.len), float64(len)))
	ssbo.len = len

	// create new buffer of destination size
	var newhandle uint32 = 0
	gl.GenBuffers(1, &newhandle)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, newhandle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*ssbo.len, nil, gl.DYNAMIC_COPY)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	// copy data from old to new buffer
	gl.BindBuffer(gl.COPY_READ_BUFFER, ssbo.handle)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, newhandle)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER, 0, 0, ssbo.typesize*minsize)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)

	// delete old buffer
	gl.DeleteBuffers(1, &ssbo.handle)

	// assign new buffer handle to this buffers
	ssbo.handle = newhandle
}

func (ssbo *SSBO) UploadValue(value []float32) {
	// fill array with value
	values := []float32{}
	for i := 0; i < ssbo.len; i++ {
		values = append(values, value...)
	}

	// upload array content to ssbo
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.handle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*ssbo.len, gl.Ptr(values), gl.DYNAMIC_COPY)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
func (ssbo *SSBO) UploadValueInRange(value []float32, start, len int) {
	// fill array with value
	values := []float32{}
	for i := 0; i < len; i++ {
		values = append(values, value...)
	}

	// upload array content to part of ssbo
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.handle)
	gl.BufferSubData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*start, ssbo.typesize*len, gl.Ptr(values))
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
func (ssbo *SSBO) UploadArray(values []float32) {
	// upload array content to ssbo
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.handle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*ssbo.len, gl.Ptr(values), gl.DYNAMIC_COPY)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
func (ssbo *SSBO) UploadArrayInRange(values []float32, start, len int) {
	// upload array content to part of ssbo
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.handle)
	gl.BufferSubData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*start, ssbo.typesize*len, gl.Ptr(values))
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}

func (ssbo *SSBO) Download() []float32 {
	// create slice of the right size
	values := make([]float32, ssbo.len)

	// copy data to array
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	for i := 0; i < ssbo.len; i++ {
		values[i] = *(*float32)(ptr)
	}
	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	return values
}

func (ssbo *SSBO) init() {
	// buffer must be at least of length 1
	bufferlen := int(math.Max(float64(ssbo.len), 1.0))

	// create an empty buffer of the right size
	gl.GenBuffers(1, &ssbo.handle)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.handle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, ssbo.typesize*bufferlen, nil, gl.DYNAMIC_COPY)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}
