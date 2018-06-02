package collision

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	OUTSIDE      = -1
	INTERSECTING = 0
	INSIDE       = 1
)

func CheckAABBFrustum(aabb AABB, mvp mgl32.Mat4) int {
	// extract planes
	planes := []Plane{
		// near
		extractPlane(mvp.Row(3).Add(mvp.Row(2))), //A4+A3
		// far
		extractPlane(mvp.Row(3).Sub(mvp.Row(2))), //A4-A3
		// bottom
		extractPlane(mvp.Row(3).Add(mvp.Row(1))), //A4+A2
		// top
		extractPlane(mvp.Row(3).Sub(mvp.Row(1))), //A4-A2
		// left
		extractPlane(mvp.Row(3).Add(mvp.Row(0))), //A4+A1
		// right
		extractPlane(mvp.Row(3).Sub(mvp.Row(0))), //A4-A1
	}

	// transform aabb points in clip space
	state := INSIDE
	for _, plane := range planes {
		if planeDistance(plane, getPointP(aabb, plane.normal)) < 0 {
			// aabb outside
			return OUTSIDE
		} else if planeDistance(plane, getPointN(aabb, plane.normal)) < 0 {
			// aabb intersects frustum
			state = INTERSECTING
		}
	}

	// completely inside
	return state
}

func extractPlane(vec mgl32.Vec4) Plane {
	n := vec.Vec3()
	len := n.Len()

	return Plane{
		n.Normalize(),
		vec.W() / len,
	}
}

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

func planeDistance(plane Plane, point mgl32.Vec3) float32 {
	return plane.normal.Dot(point) + plane.d
}
