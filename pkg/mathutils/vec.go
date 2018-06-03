// Package mathutils provides utility functions for scalar and vectorial math.
package mathutils

import "github.com/go-gl/mathgl/mgl32"

// Vec3To4 takes a vec3 and turns it into a vec4 with the fourth component being 1.
func Vec3To4(vec mgl32.Vec3) mgl32.Vec4 {
	return mgl32.Vec4{vec.X(), vec.Y(), vec.Z(), 1.0}
}

// Vec4To3 extracts the first three components of a vec4 and returns them as a vec3.
func Vec4To3(vec mgl32.Vec4) mgl32.Vec3 {
	return mgl32.Vec3{vec.X(), vec.Y(), vec.Z()}
}
