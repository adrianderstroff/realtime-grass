// Package collision handles collision checking between an AABB and a view frustum.
package collision

import "github.com/go-gl/mathgl/mgl32"

// Plane is specified by its normal and its distance to the origin.
type Plane struct {
	normal mgl32.Vec3
	d      float32
}

// MakePlane takes the first three components of vec4 and use it as the plane's  normal.
// the fourth component normalized by the length of the three components of the vec4 is the
// distance of the plane from the origin.
func MakePlane(vec mgl32.Vec4) Plane {
	n := vec.Vec3()
	len := n.Len()

	return Plane{
		n.Normalize(),
		vec.W() / len,
	}
}

// Distance returns the distance of a point to this plane.
func (plane *Plane) Distance(point mgl32.Vec3) float32 {
	return plane.normal.Dot(point) + plane.d
}
