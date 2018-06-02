package mathutils

import "github.com/go-gl/mathgl/mgl32"

func Vec3To4(vec mgl32.Vec3) mgl32.Vec4 {
	return mgl32.Vec4{vec.X(), vec.Y(), vec.Z(), 1.0}
}
func Vec4To3(vec mgl32.Vec4) mgl32.Vec3 {
	return mgl32.Vec3{vec.X(), vec.Y(), vec.Z()}
}
