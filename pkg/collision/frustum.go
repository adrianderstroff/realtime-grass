package collision

import (
	"math"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func MakeFrustum(near, far, fov float32) engine.Mesh {
	angle := fov * math.Pi / 180.0
	dnear := float32(math.Tan(float64(angle)/2.0)) * near
	dfar := float32(math.Tan(float64(angle)/2.0)) * far

	// direction
	dir := mgl32.Vec3{0, 0, -1}
	right := mgl32.Vec3{1, 0, 0}
	up := mgl32.Vec3{0, 1, 0}

	// create points
	v1 := dir.Mul(dnear).Add(right.Mul(-dnear)).Add(up.Mul(dnear))
	v2 := dir.Mul(dnear).Add(right.Mul(-dnear)).Add(up.Mul(-dnear))
	v3 := dir.Mul(dnear).Add(right.Mul(dnear)).Add(up.Mul(dnear))
	v4 := dir.Mul(dnear).Add(right.Mul(dnear)).Add(up.Mul(-dnear))

	v5 := dir.Mul(dfar).Add(right.Mul(-dfar)).Add(up.Mul(dfar))
	v6 := dir.Mul(dfar).Add(right.Mul(-dfar)).Add(up.Mul(-dfar))
	v7 := dir.Mul(dfar).Add(right.Mul(dfar)).Add(up.Mul(dfar))
	v8 := dir.Mul(dfar).Add(right.Mul(dfar)).Add(up.Mul(-dfar))

	// create position buffer
	positions := combine(
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
	barycoords := repeat(
		[]float32{
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
		},
		12,
	)

	mesh, _ := engine.MakeMeshFromArrays(positions, nil, barycoords, "position", "", "barycoord", 3, 0, 3, gl.TRIANGLES)
	return mesh
}

func combine(slices ...mgl32.Vec3) []float32 {
	var result []float32
	for _, s := range slices {
		result = append(result, s.X(), s.Y(), s.Z())
	}
	return result
}

func repeat(slice []float32, number int) []float32 {
	var result []float32
	for i := 0; i < number; i++ {
		result = append(result, slice...)
	}
	return result
}
