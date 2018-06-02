package mathutils

import "github.com/go-gl/mathgl/mgl32"

func Interpolate(val1, val2, alpha float32) float32 {
	return val1*(1-alpha) + val2*alpha
}

func ExtractRotation(mat4 mgl32.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		mat4.At(0, 0), mat4.At(1, 0), mat4.At(2, 0), 0.0,
		mat4.At(0, 1), mat4.At(1, 1), mat4.At(2, 1), 0.0,
		mat4.At(0, 2), mat4.At(1, 2), mat4.At(2, 2), 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}
