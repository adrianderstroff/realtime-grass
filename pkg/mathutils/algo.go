// Package mathutils provides utility functions for scalar and vectorial math.
package mathutils

import "github.com/go-gl/mathgl/mgl32"

// Interpolate returns the interpolated value between val1 and val2 specified by alpha.
// An alpha value of 0 returns val1 while a value of 1 return val2.
func Interpolate(val1, val2, alpha float32) float32 {
	return val1*(1-alpha) + val2*alpha
}

// ExtractRotation takes the upper-left 3x3 rotation marix of a view matrix and returns it as a 4x4 matrix.
func ExtractRotation(viewmat *mgl32.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		viewmat.At(0, 0), viewmat.At(1, 0), viewmat.At(2, 0), 0.0,
		viewmat.At(0, 1), viewmat.At(1, 1), viewmat.At(2, 1), 0.0,
		viewmat.At(0, 2), viewmat.At(1, 2), viewmat.At(2, 2), 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

// MapI32 specifies a function that maps a value val from the source domain between srcstart to srcend
// onto a value in the destination domain between dststart to dstend.
// The return value gets truncated casted to int32.
func MapI32(val, srcstart, srcend, dststart, dstend int32) int32 {
	return int32(MapF32(float32(val), float32(srcstart), float32(srcend), float32(dststart), float32(dstend)))
}

// MapI32 specifies a function that maps a value val from the source domain between srcstart to srcend
// onto a value in the destination domain between dststart to dstend.
func MapF32(val, srcstart, srcend, dststart, dstend float32) float32 {
	x := (val - srcstart) / (srcend - srcstart)
	return x*(dstend-dststart) + dststart
}

// Combine concatenates multiple slices into one slice.
func Combine(slices ...mgl32.Vec3) []float32 {
	var result []float32
	for _, s := range slices {
		result = append(result, s.X(), s.Y(), s.Z())
	}
	return result
}

// Repeat concatenates the same slices several times to one slice.
func Repeat(slice []float32, number int) []float32 {
	var result []float32
	for i := 0; i < number; i++ {
		result = append(result, slice...)
	}
	return result
}
