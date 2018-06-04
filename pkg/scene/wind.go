// Package scene contains all main entities for rendering and/or interaction with the user.
package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/engine"
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
)

// Wind has two grids of the same size for the wind velocity and acceleration.
// The movement of the camera creates acceleration around the camera.
// The acceleration intensity is a bell-curve with strong acceleration near the center that degrades to the outside.
type Wind struct {
	shader            *engine.ShaderProgram
	velocityfield     *engine.SSBO
	accelerationfield *engine.SSBO
	groupcount        uint32
	griddimension     int32
	fieldsize         int32
	cellsize          float32
	prevcenterx       int32
	prevcenterz       int32
	dt                float32
	t                 float32
}

// MakeWind constructs the Wind struct.
// The side of a wind grid is 2*radius+1 making each grid an odd number of cells.
// The influence specifies the relation between grid radius and bell-curve spread.
// The bell-curve spread is calculated by radius/influence.
// This means that higher values of influence yield a stronger contracted spread.
// The cellsize has to equal the tilesize of the terrain.
func MakeWind(shaderpath string, radius int, influence, cellsize float32) (Wind, error) {
	griddim := 2*radius + 1
	fieldsize := griddim * griddim

	// both grids store a vec4 (4 x float32) per cell
	valuecount := 4
	bytesize := 4 * valuecount

	// create velocityfield
	velocityfield := engine.MakeSSBO(bytesize, fieldsize)
	velocityfield.UploadValue([]float32{0, 0, 0, 0})

	// create accelerationfield
	accelerationfield := engine.MakeSSBO(bytesize, fieldsize)
	afielddata := make([]float32, fieldsize*valuecount)
	var max float32 = 0.0
	for z := 0; z < griddim; z++ {
		dz := 4 * float64(z-radius) / float64(radius) * float64(influence)
		for x := 0; x < griddim; x++ {
			dx := 4 * float64(x-radius) / float64(radius) * float64(influence)
			// calculate acceleration vector
			e := math.Pow(math.E, -(dx*dx)-(dz*dz))
			idx := (z*griddim + x) * valuecount
			afielddata[idx] = float32(dx) * float32(e)
			afielddata[idx+1] = float32(dz) * float32(e)
			afielddata[idx+2] = 0.0
			afielddata[idx+3] = 0.0
			// calculate max magnitude of all acceleration vectors
			max = mathutils.MaxF32(max, afielddata[idx])
			max = mathutils.MaxF32(max, afielddata[idx+1])
		}
	}
	// normalize values in such a way that the magnitudes of all acceleration vectors are between 0 and 1
	for z := 0; z < griddim; z++ {
		for x := 0; x < griddim; x++ {
			idx := (z*griddim + x) * valuecount
			afielddata[idx] /= max
			afielddata[idx+1] /= max
		}
	}
	accelerationfield.UploadArray(afielddata)

	// create wind compute shader
	shader, err := engine.MakeComputeProgram(shaderpath + "wind/wind.comp")
	if err != nil {
		return Wind{}, err
	}

	// calculate number of work groups necessary
	groupcount := uint32(mathutils.CeilF32(float32(fieldsize) / 16.0))

	return Wind{
		shader:            &shader,
		velocityfield:     &velocityfield,
		accelerationfield: &accelerationfield,
		groupcount:        groupcount,
		griddimension:     int32(griddim),
		fieldsize:         int32(fieldsize),
		cellsize:          cellsize,
		prevcenterx:       0,
		prevcenterz:       0,
		dt:                0.1,
		t:                 0.0,
	}, nil
}

// Update performs the wind simulation.
// pos is the new camera position and is used to specify the new center position of the wind grid.
// cameradelta is the vector from the previous to the current camera position.
func (wind *Wind) Update(pos, cameradelta mgl32.Vec3) {
	// get cell in which the camera is in
	centerx := int32(pos.X() / wind.cellsize)
	centerz := int32(pos.Z() / wind.cellsize)
	if pos.X() < 0 {
		centerx -= 1
	}
	if pos.Z() < 0 {
		centerz -= 1
	}

	// calculate wind offset
	dx := centerx - wind.prevcenterx
	dz := centerz - wind.prevcenterz

	// bind buffers
	wind.velocityfield.Bind(0)
	wind.accelerationfield.Bind(1)

	// get direction
	distance := cameradelta.Len()
	speed := distance / wind.dt
	dir := cameradelta.Normalize()

	// update wind simulation
	wind.shader.Use()
	wind.shader.UpdateVec2("viewDir", mgl32.Vec2{dir.X(), dir.Z()})
	wind.shader.UpdateFloat32("speed", speed)
	wind.shader.UpdateInt32("size", wind.fieldsize)
	wind.shader.UpdateInt32("dim", wind.griddimension)
	wind.shader.UpdateInt32("dx", dx)
	wind.shader.UpdateInt32("dz", dz)
	wind.shader.UpdateFloat32("dt", wind.dt)
	wind.shader.UpdateFloat32("t", wind.t)
	wind.shader.Compute(wind.groupcount, 1, 1)

	// unbind buffers
	wind.velocityfield.Unbind()
	wind.accelerationfield.Unbind()

	// save last position
	wind.prevcenterx = centerx
	wind.prevcenterz = centerz

	// update time
	wind.t++
}
