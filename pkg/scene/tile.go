// Package scene contains all main entities for rendering and/or interaction with the user.
package scene

import (
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
	"github.com/go-gl/mathgl/mgl32"
)

// TileFactory is creating single Tiles.
// The TileFactory knows of the size of a tile and how many tiles are in one block.
// Additionally it has a reference to a height-map used to grab the height of the four points in a Tile.
type TileFactory struct {
	tilesize      float32
	tilesperblock int32

	heightmap *Heightmap
}

// Tile contains its position and the plane data of the two triangles that make up a Tile.
// A Tile has an upper left and a lower right triangle where the points 2 and 3 are shared between both triangles.
// To have the normals both pointing in the right direction the point orders of both triangles is 1-2-3 and 3-2-4.
// 1-------3
// |     / |
// |   /   |
// | /     |
// 2-------4
type Tile struct {
	pos  []float32
	data []float32
}

// MakeTile constructs a Tile at position (tx,tz)
// A Tile consists of 2 triangles with 6 vertices in total.
// The heights of the vertices are calculated from a height-map.
// The Tile coordinates go from left to right for x and bottom to top for z.
//    x------->
//  ^ 1-------3
//  | |     / |
//  | |   /   |
//  | | /     |
//  z 2-------4
func (tf *TileFactory) MakeTile(tx, tz int32) Tile {
	// get the height values for the height-map
	h1 := tf.getHeight(tx, tz+1)
	h2 := tf.getHeight(tx, tz)
	h3 := tf.getHeight(tx+1, tz+1)
	h4 := tf.getHeight(tx+1, tz)

	// calculate positions of all 4 points
	p1 := tf.calcPlanePos(tx, tz+1, h1)
	p2 := tf.calcPlanePos(tx, tz, h2)
	p3 := tf.calcPlanePos(tx+1, tz+1, h3)
	p4 := tf.calcPlanePos(tx+1, tz, h4)

	// solve the plane formula for both triangles
	tri1 := calcPlane(p1, p2, p3)
	tri2 := calcPlane(p3, p2, p4)

	// calculate the Tile position
	pos := mgl32.Vec3{
		(p1.X() + p3.X()) / 2,
		0.0,
		(p1.Z() + p2.Z()) / 2,
	}

	// construct the Tile from the position and plane data
	return Tile{
		pos: []float32{pos.X(), 0.0, pos.Z()}, // position of the Tile's center
		data: []float32{
			tri1.X(), tri1.Y(), tri1.Z(), tri1.W(), // triangle 1
			tri2.X(), tri2.Y(), tri2.Z(), tri2.W(), // triangle 2
			pos.X(), pos.Z(), // tile position
			3.0, // level of detail
			0.0, // padding to have 12 byte
		},
	}
}

// getHeight grabs the height from the height-map at position (x,z).
// This position gets normalized and repeated in all directions to have a continuous repeating heightmap.
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

// calcTileBounds repeats a position p to be relative to the block size.
// The returned value is between 0 and the side length of a block.
func (tf *TileFactory) calcTileBounds(x int32) int32 {
	rx := x % tf.tilesperblock
	if rx < 0 {
		rx = tf.tilesperblock + rx
	}
	return rx
}

// clacPlanePos returns the position of the Tile center as vec3.
// The y component is specified by height.
func (tf *TileFactory) calcPlanePos(x, z int32, height float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(x) * tf.tilesize,
		height,
		float32(z) * tf.tilesize,
	}
}

// calcPlane calculates the components A,B,C,D for the plane equation Ax+By+Cz+D=0.
// The points v1,v2,v3 have to be specified with v2 being the center point and the points have to be defined counter-clockwise.
func calcPlane(v1, v2, v3 mgl32.Vec3) mgl32.Vec4 {
	// normal n = (v1-v2) x (v3-v2)
	d1 := v1.Sub(v2)
	d2 := v3.Sub(v2)
	n := d1.Cross(d2)

	// D = -Ax -By -Cz
	D := -n.X()*v2.X() - n.Y()*v2.Y() - n.Z()*v2.Z()

	// normalize all components
	len := n.Len()
	// D/||(x y z)|| is the distance of the plane to the origin
	return mgl32.Vec4{
		n.X() / len,
		n.Y() / len,
		n.Z() / len,
		D / len,
	}
}
