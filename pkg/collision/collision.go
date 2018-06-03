// Package collision handles collision checking between an AABB and a view frustum.
package collision

import "github.com/go-gl/mathgl/mgl32"

const (
	OUTSIDE      = -1
	INTERSECTING = 0
	INSIDE       = 1
)

// CheckAABBFrustum performs a collision check between an AABB and a view frustum.
// An AABB  can either be inside, outside the view frustum or it intersects the frustum.
// If the AABB is outside the frustum -1 is returned, while 1 is returned when the AABB
// is completely inside and 0 when it intersects the frustum.
// Source: http://www.lighthouse3d.com/tutorials/view-frustum-culling/geometric-approach-testing-boxes/.
func CheckAABBFrustum(aabb AABB, mvp mgl32.Mat4) int {
	// extract planes
	planes := []Plane{
		// near
		MakePlane(mvp.Row(3).Add(mvp.Row(2))), //A4+A3
		// far
		MakePlane(mvp.Row(3).Sub(mvp.Row(2))), //A4-A3
		// bottom
		MakePlane(mvp.Row(3).Add(mvp.Row(1))), //A4+A2
		// top
		MakePlane(mvp.Row(3).Sub(mvp.Row(1))), //A4-A2
		// left
		MakePlane(mvp.Row(3).Add(mvp.Row(0))), //A4+A1
		// right
		MakePlane(mvp.Row(3).Sub(mvp.Row(0))), //A4-A1
	}

	// transform all AABB points into clip space
	state := INSIDE
	for _, plane := range planes {
		if plane.Distance(getPointP(aabb, plane.normal)) < 0 {
			// AABB is outside of one plane of the frustum and thus outside of the frustum
			return OUTSIDE
		} else if plane.Distance(getPointN(aabb, plane.normal)) < 0 {
			// AABB intersects the plane and thus the frustum
			state = INTERSECTING
		}
	}

	// AABB is inside of all planes. This means that the AABB has to be inside the frustum
	return state
}

// getPointP gets the point of the AABB that is further along the normal's direction.
func getPointP(aabb AABB, normal mgl32.Vec3) mgl32.Vec3 {
	p := aabb.Min
	if normal.X() >= 0 {
		p = mgl32.Vec3{aabb.Max.X(), p.Y(), p.Z()}
	}
	if normal.Y() >= 0 {
		p = mgl32.Vec3{p.X(), aabb.Max.Y(), p.Z()}
	}
	if normal.Z() >= 0 {
		p = mgl32.Vec3{p.X(), p.Y(), aabb.Max.Z()}
	}
	return p
}

// getPointN gets the point of the AABB that is on the opposite of point P.
func getPointN(aabb AABB, normal mgl32.Vec3) mgl32.Vec3 {
	n := aabb.Max
	if normal.X() >= 0 {
		n = mgl32.Vec3{aabb.Min.X(), n.Y(), n.Z()}
	}
	if normal.Y() >= 0 {
		n = mgl32.Vec3{n.X(), aabb.Min.Y(), n.Z()}
	}
	if normal.Z() >= 0 {
		n = mgl32.Vec3{n.X(), n.Y(), aabb.Min.Z()}
	}
	return n
}
