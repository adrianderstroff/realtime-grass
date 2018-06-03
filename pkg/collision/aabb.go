// Package collision handles collision checking between an AABB and a view frustum.
package collision

import "github.com/go-gl/mathgl/mgl32"

// AABB is an axis-aligned bounding box.
// Min and Max are the two opposite points of the AABB.
type AABB struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

// MakeAABB is a constructor for the AABB specifying the opposite points of the AABB.
func MakeAABB(min, max mgl32.Vec3) AABB {
	return AABB{
		Min: min,
		Max: max,
	}
}
