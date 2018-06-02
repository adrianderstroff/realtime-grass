package scene

import (
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
	"github.com/go-gl/mathgl/mgl32"
)

type TileFactory struct {
	tilesize      float32
	tilesperblock int32

	heightmap *Heightmap
}

type Tile struct {
	pos  []float32
	data []float32
}

// Creates a tile at position (tx,tz)
// A tile consists of 2 triangles with 6 vertices in total.
// The heights of the vertices are calculated from a heightmap.
// The vertex order is
//    x->
//  ^ 1---3
//  | |  /|
//  z | / |
//    |/  |
//    2---4
// while vertex 2 and 3 are reused for both triangles.
func (tf *TileFactory) MakeTile(tx, tz int32) Tile {
	// get the height values
	h1 := tf.getHeight(tx, tz+1)
	h2 := tf.getHeight(tx, tz)
	h3 := tf.getHeight(tx+1, tz+1)
	h4 := tf.getHeight(tx+1, tz)

	// calc positions of all 4 vertices
	p1 := tf.calcPlanePos(tx, tz+1, h1)
	p2 := tf.calcPlanePos(tx, tz, h2)
	p3 := tf.calcPlanePos(tx+1, tz+1, h3)
	p4 := tf.calcPlanePos(tx+1, tz, h4)

	// solve the plane formula for both triangles
	tri1 := calcPlane(p1, p2, p3)
	tri2 := calcPlane(p3, p2, p4)

	// calculate the tile pos
	pos := mgl32.Vec3{
		(p1.X() + p3.X()) / 2,
		0.0,
		(p1.Z() + p2.Z()) / 2,
	}

	// create tile
	return Tile{
		pos: []float32{pos.X(), 0.0, pos.Z()},
		data: []float32{
			tri1.X(), tri1.Y(), tri1.Z(), tri1.W(), // triangle 1
			tri2.X(), tri2.Y(), tri2.Z(), tri2.W(), // triangle 2
			pos.X(), pos.Z(), // tile position
			3.0, // level of detail
			0.0, // padding to have 12 byte
		},
	}
}

func (tf *TileFactory) getHeight(x, z int32) float32 {
	// calc xz-coordinate relative to the block size
	rx := tf.calcTileBounds(x)
	rz := tf.calcTileBounds(z)
	// map tile coordinate to image pixel position
	ix := mathutils.MapI32(rx, 0, tf.tilesperblock-1, 0, tf.heightmap.GetWidth()-1)
	iz := mathutils.MapI32(rz, 0, tf.tilesperblock-1, 0, tf.heightmap.GetHeight()-1)
	// read height at pixel position (ix,iz)
	return tf.heightmap.GetHeightAt(ix, iz)
}
func (tf *TileFactory) calcTileBounds(x int32) int32 {
	rx := x % tf.tilesperblock
	if rx < 0 {
		rx = tf.tilesperblock + rx
	}
	return rx
}
func (tf *TileFactory) calcPlanePos(x, z int32, height float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(x) * tf.tilesize,
		height,
		float32(z) * tf.tilesize,
	}
}
func calcPlane(v1, v2, v3 mgl32.Vec3) mgl32.Vec4 {
	d1 := v1.Sub(v2)
	d2 := v3.Sub(v2)
	n := d1.Cross(d2)
	len := n.Len()
	D := -n.X()*v2.X() - n.Y()*v2.Y() - n.Z()*v2.Z()

	return mgl32.Vec4{
		n.X() / len,
		n.Y() / len,
		n.Z() / len,
		D / len,
	}
}
