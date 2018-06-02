package collision

import "github.com/go-gl/mathgl/mgl32"

type Plane struct {
	normal mgl32.Vec3
	d      float32
}
