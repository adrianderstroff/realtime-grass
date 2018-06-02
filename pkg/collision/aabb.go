package collision

import "github.com/go-gl/mathgl/mgl32"

type AABB struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

func MakeAABB(min, max mgl32.Vec3) AABB {
	return AABB{
		Min: min,
		Max: max,
	}
}
