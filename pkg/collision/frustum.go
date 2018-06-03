// Package collision handles collision checking between an AABB and a view frustum.
package collision

import (
	"math"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// MakeFrustum creates the mesh of a frustum by providing the near and far plane distance
// as well as the field of view angle in degrees.
func MakeFrustum(near, far, fov float32) engine.Mesh {
	// calculate the half width of the near and far planes
	angle := fov * math.Pi / 180.0
	dnear := float32(math.Tan(float64(angle)/2.0)) * near
	dfar := float32(math.Tan(float64(angle)/2.0)) * far

	// coordinate system for the frustum
	dir := mgl32.Vec3{0, 0, -1}
	right := mgl32.Vec3{1, 0, 0}
	up := mgl32.Vec3{0, 1, 0}

	// create points of the frustum
	v1 := dir.Mul(dnear).Add(right.Mul(-dnear)).Add(up.Mul(dnear))
	v2 := dir.Mul(dnear).Add(right.Mul(-dnear)).Add(up.Mul(-dnear))
	v3 := dir.Mul(dnear).Add(right.Mul(dnear)).Add(up.Mul(dnear))
	v4 := dir.Mul(dnear).Add(right.Mul(dnear)).Add(up.Mul(-dnear))
	v5 := dir.Mul(dfar).Add(right.Mul(-dfar)).Add(up.Mul(dfar))
	v6 := dir.Mul(dfar).Add(right.Mul(-dfar)).Add(up.Mul(-dfar))
	v7 := dir.Mul(dfar).Add(right.Mul(dfar)).Add(up.Mul(dfar))
	v8 := dir.Mul(dfar).Add(right.Mul(dfar)).Add(up.Mul(-dfar))

	// create positions of the frustum
	positions := mathutils.Combine(
		// front
		v1, v2, v3,
		v3, v2, v4,
		// back
		v7, v8, v5,
		v5, v8, v6,
		// left
		v5, v6, v1,
		v1, v6, v2,
		// right
		v3, v4, v7,
		v7, v4, v8,
		// top
		v5, v1, v7,
		v7, v1, v3,
		// bottom
		v8, v4, v6,
		v6, v4, v2,
	)
	// create barycentric coordinates of the frustum used for a wireframe shader.
	barycoords := mathutils.Repeat(
		[]float32{
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
		},
		12,
	)

	// create the mesh for the frustum
	mesh, _ := engine.MakeMeshFromArrays(positions, nil, barycoords, "position", "", "barycoord", 3, 0, 3, gl.TRIANGLES)
	return mesh
}
