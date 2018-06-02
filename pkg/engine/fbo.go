package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type FBO struct {
	handle        uint32
	isBound       bool
	ColorTextures []*Texture
	DepthTexture  *Texture
	textureType   uint32
}

func MakeEmptyFBO() FBO {
	fbo := FBO{0, false, nil, nil, gl.TEXTURE_2D}
	gl.GenFramebuffers(1, &fbo.handle)
	return fbo
}
func MakeFBO(width, height int32) FBO {
	fbo := FBO{0, false, nil, nil, gl.TEXTURE_2D}
	gl.GenFramebuffers(1, &fbo.handle)
	color := MakeColorTexture(width, height)
	depth := MakeDepthTexture(width, height)
	fbo.AttachColorTexture(color, 0)
	fbo.AttachDepthTexture(depth)
	return fbo
}
func MakeMultisampleFBO() FBO {
	fbo := FBO{0, false, nil, nil, gl.TEXTURE_2D_MULTISAMPLE}
	gl.GenFramebuffers(1, &fbo.handle)
	return fbo
}
func (fbo *FBO) Delete() {
	// delete textures
	if fbo.ColorTextures != nil {
		for _, colTex := range fbo.ColorTextures {
			if colTex != nil {
				colTex.Delete()
			}
		}
	}
	if fbo.DepthTexture != nil {
		fbo.DepthTexture.Delete()
	}

	// unbind fbo
	if fbo.isBound {
		fbo.Unbind()
	}

	// delete buffer
	gl.DeleteFramebuffers(1, &fbo.handle)
}
func (fbo *FBO) Clear() {
	if fbo.isBound {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
}
func (fbo *FBO) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo.handle)
	fbo.isBound = true
}
func (fbo *FBO) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	fbo.isBound = false
}

func (fbo *FBO) AttachColorTexture(texture Texture, index uint32) {
	fbo.Bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+index, fbo.textureType, texture.handle, 0)
	drawBuffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &drawBuffers[0])
	fbo.Unbind()
	// add handle
	fbo.ColorTextures = append(fbo.ColorTextures, &texture)
}
func (fbo *FBO) AttachDepthTexture(texture Texture) {
	fbo.Bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, fbo.textureType, texture.handle, 0)
	fbo.Unbind()
	// add handle
	fbo.DepthTexture = &texture
}

// Checks if the framebuffer is complete
func (fbo *FBO) IsComplete() bool {
	fbo.Bind()
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	fbo.Unbind()
	return status == gl.FRAMEBUFFER_COMPLETE
}

func (fbo *FBO) CopyToScreen(index uint32, x, y, width, height int32) {
	fbo.CopyToScreenRegion(index, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyToScreenRegion(index uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.DrawBuffer(gl.BACK)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyToFBORegion(other, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyColorToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyColorToFBORegion(other, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyColorToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyColorToFBOSmooth(other *FBO, x, y, width, height int32) {
	fbo.CopyColorToFBORegionSmooth(other, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyColorToFBORegionSmooth(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.LINEAR,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyAttachmentColorToFBO(other *FBO, index1, index2 uint32, x, y, width, height int32) {
	fbo.CopyColorAttachmentToFBORegion(other, index1, index2, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyColorAttachmentToFBORegion(other *FBO, index1, index2 uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index1)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0 + index2)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyAttachmentColorToFBOSmooth(other *FBO, index1, index2 uint32, x, y, width, height int32) {
	fbo.CopyAttachmentColorToFBORegionSmooth(other, index1, index2, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyAttachmentColorToFBORegionSmooth(other *FBO, index1, index2 uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index1)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0 + index2)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.LINEAR,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo *FBO) CopyDepthToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyDepthToFBORegion(other, x, y, width, height, x, y, width, height)
}
func (fbo *FBO) CopyDepthToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
