// Package scene contains all main entities for rendering and/or interaction with the user.
package scene

import (
	"github.com/adrianderstroff/realtime-grass/pkg/collision"
	"github.com/go-gl/mathgl/mgl32"
)

// ChunkFactory is creating single Chunks.
// Each Chunk consists of multiple chunks specified by the chunkresolution.
// To create the Tiles it has a reference to the TileFactory.
// The chunkheight is the terrain height plus the maximum height of the grass.
type ChunkFactory struct {
	chunksize       float32
	chunkheight     float32
	terrainheight   float32
	chunkresolution int32

	tf *TileFactory
}

// Chunk is a collection of Tiles.
// In addition a chunk has a position and an AABB.
// poss are all Tile positions while data is the Tile data of all Tiles.
type Chunk struct {
	pos  mgl32.Vec3
	aabb collision.AABB
	poss []float32
	data []float32
}

// MakeChunk creates a single Chunk at position (cx,cz).
func (cf *ChunkFactory) MakeChunk(cx, cz int32) Chunk {
	pos := makeCenteredVec(cx, cz, cf.chunkheight/2, cf.chunksize)
	dir := makeCenteredVec(0, 0, cf.chunkheight/2, cf.chunksize)

	// create AABB around the Chunk
	aabb := collision.MakeAABB(pos.Sub(dir), pos.Add(dir))

	// calc tile start for this Chunk
	acx := cx * cf.chunkresolution
	acz := cz * cf.chunkresolution

	// create all Tiles and collect all positions and Tile data
	var (
		poss []float32 // center positions of all Tiles
		data []float32 // Tile data of all Tiles
	)
	var tx, tz int32
	for tz = 0; tz < cf.chunkresolution; tz++ {
		for tx = 0; tx < cf.chunkresolution; tx++ {
			tile := cf.tf.MakeTile(acx+tx, acz+tz)
			poss = append(poss, tile.pos...)
			data = append(data, tile.data...)
		}
	}

	return Chunk{
		pos:  pos,
		aabb: aabb,
		poss: poss,
		data: data,
	}
}

// makeCenteredVec creates a Chunk centered vec3 with the provided height.
func makeCenteredVec(x, z int32, height, size float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(x)*size + size/2,
		height,
		float32(z)*size + size/2,
	}
}
