package engine

import (
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type Texture struct {
	handle uint32
	target uint32
	texPos uint32 // e.g. gl.TEXTURE0
}

func MakeEmptyTexture() Texture {
	return Texture{0, gl.TEXTURE_2D, 0}
}
func MakeTexture(width, height, internalformat int32, format, pixelType uint32, data unsafe.Pointer, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t)

	// specify a texture image
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalformat, width, height, 0, format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}
func MakeColorTexture(width, height int32) Texture {
	return MakeTexture(width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}
func MakeDepthTexture(width, height int32) Texture {
	tex := MakeTexture(width, height, gl.DEPTH_COMPONENT, gl.DEPTH_COMPONENT, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
	return tex
}
func MakeCubeMapTexture(right, left, top, bottom, front, back string) (Texture, error) {
	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// load images
	imagePaths := []string{right, left, top, bottom, front, back}
	for i, path := range imagePaths {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)
		image, err := MakeImage(path)
		if err != nil {
			return Texture{}, err
		}
		gl.TexImage2D(target, 0, image.internalFormat, image.width, image.height,
			0, image.format, image.pixelType, image.data)
	}

	// format texture
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// unset active texture
	tex.Unbind()

	return tex, nil
}
func MakeTextureFromPath(path string) (Texture, error) {
	image, err := MakeImage(path)
	if err != nil {
		return Texture{}, err
	}

	return MakeTexture(image.width, image.height, image.internalFormat, image.format,
		image.pixelType, image.data, gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

func MakeMultisampleTexture(width, height, samples int32, format uint32, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D_MULTISAMPLE, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	/* gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_T, t) */

	// specify a texture image
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, samples, format, width, height, false)

	// unbind texture
	texture.Unbind()

	return texture
}
func MakeColorMultisampleTexture(width, height, samples int32) Texture {
	return MakeMultisampleTexture(width, height, samples, gl.RGBA,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}
func MakeDepthMultisampleTexture(width, height, samples int32) Texture {
	return MakeMultisampleTexture(width, height, samples, gl.DEPTH_COMPONENT,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

func (tex *Texture) Delete() {
	gl.DeleteTextures(1, &tex.handle)
}

func (tex *Texture) GenMipmap() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}
func (tex *Texture) GenMipmapNearest() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

func (tex *Texture) Bind(index uint32) {
	tex.texPos = gl.TEXTURE0 + index
	gl.ActiveTexture(tex.texPos)
	gl.BindTexture(tex.target, tex.handle)
}
func (tex *Texture) Unbind() {
	tex.texPos = 0
	gl.BindTexture(tex.target, 0)
}
