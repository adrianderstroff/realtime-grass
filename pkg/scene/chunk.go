package scene

import (
	"github.com/adrianderstroff/realtime-grass/pkg/collision"
	"github.com/go-gl/mathgl/mgl32"
)

type ChunkFactory struct {
	chunksize       float32
	chunkheight     float32
	terrainheight   float32
	chunkresolution int32

	tf *TileFactory
}

type Chunk struct {
	pos  mgl32.Vec3
	aabb collision.AABB
	poss []float32
	data []float32
}

func (cf *ChunkFactory) MakeChunk(cx, cz int32) Chunk {
	pos := makeCenteredVec(cx, cz, cf.chunkheight/2, cf.chunksize)
	dir := makeCenteredVec(0, 0, cf.chunkheight/2, cf.chunksize)

	// create aabb
	aabb := collision.MakeAABB(pos.Sub(dir), pos.Add(dir))

	// calc tile start for this chunk
	acx := cx * cf.chunkresolution
	acz := cz * cf.chunkresolution

	// create all tiles
	var (
		poss []float32 // center positions of all tiles
		data []float32 // tile data of all tiles
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
		pos,
		aabb,
		poss,
		data,
	}
}
func makeCenteredVec(x, z int32, height, size float32) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(x)*size + size/2,
		height,
		float32(z)*size + size/2,
	}
}
