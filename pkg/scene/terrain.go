package scene

import (
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/adrianderstroff/realtime-grass/pkg/collision"
	"github.com/adrianderstroff/realtime-grass/pkg/engine"
	"github.com/adrianderstroff/realtime-grass/pkg/mathutils"
)

type Terrain struct {
	// rendering
	shader        engine.ShaderProgram
	buffer        *engine.Mesh
	terrainbuffer engine.SSBO
	grass         Grass
	wind          Wind
	// factories
	cf *ChunkFactory
	tf *TileFactory
	// chunk
	chunks    map[string]Chunk
	chunksize float32
	// tile
	tilesize      float32
	tilesperchunk int32
	tilecount     int32
	// loading
	loaddist   float32
	unloaddist float32
}

func MakeTerrain(shaderpath, texpath string, blocksize float32, blockresolution, chunkresolution int32, terrainheight float32, bladecount int, grassheight, viewdist float32, windradius int32, windinfluence float32) (Terrain, error) {
	// setup shaderprogram
	shader, err := engine.MakeGeomProgram(shaderpath+"/terrain/terrain.vert", shaderpath+"/terrain/terrain.geom", shaderpath+"/terrain/terrain.frag")
	if err != nil {
		return Terrain{}, err
	}

	// get heightmap data
	heightmap, err := MakeHeightmap(texpath+"heightmap.png", terrainheight)
	if err != nil {
		return Terrain{}, err
	}

	// setup vertex buffer
	positionsbuffer, err := engine.MakeSimpleMesh(nil, 3, gl.POINTS, gl.STREAM_DRAW)
	if err != nil {
		return Terrain{}, err
	}
	shader.AddRenderable(&positionsbuffer)
	// setup ssbo
	tilebytesize := 12 * 4 // vec4 + vec4 + vec2 + float + float
	terrainbuffer := engine.MakeEmptySSBO(tilebytesize)

	// setup factories
	chunksize := blocksize / float32(blockresolution)
	tilesize := chunksize / float32(chunkresolution)
	tf := TileFactory{
		tilesize:      tilesize,
		tilesperblock: blockresolution * chunkresolution,
		heightmap:     &heightmap,
	}
	cf := ChunkFactory{
		chunksize:       chunksize,
		chunkheight:     terrainheight + grassheight,
		terrainheight:   terrainheight,
		chunkresolution: chunkresolution,
		tf:              &tf,
	}

	// setup grass
	grass, err := MakeGrass(shaderpath, texpath, bladecount, grassheight, viewdist, windradius)
	if err != nil {
		return Terrain{}, err
	}

	// setup wind
	wind, err := MakeWind(shaderpath, int(windradius), windinfluence, tilesize)
	if err != nil {
		return Terrain{}, err
	}

	// create terrain
	return Terrain{
		// rendering
		shader:        shader,
		buffer:        &positionsbuffer,
		terrainbuffer: terrainbuffer,
		grass:         grass,
		wind:          wind,
		// factories
		cf: &cf,
		tf: &tf,
		// chunk
		chunks:    map[string]Chunk{},
		chunksize: chunksize,
		// tile
		tilesize:      tilesize,
		tilesperchunk: chunkresolution * chunkresolution,
		tilecount:     0,
		// loading
		loaddist:   viewdist + chunksize,
		unloaddist: viewdist + chunksize,
	}, nil
}
func (terrain *Terrain) Update(pos, cameradelta mgl32.Vec3, mvp mgl32.Mat4) {
	// update wind
	terrain.wind.Update(pos, cameradelta)

	// update chunks
	terrain.unload(pos)
	terrain.load(pos)

	// collect terrain data
	poss := []float32{}
	data := []float32{}
	var tilecount int32 = 0
	for _, chunk := range terrain.chunks {
		if collision.CheckAABBFrustum(chunk.aabb, mvp) != collision.OUTSIDE {
			poss = append(poss, chunk.poss...)
			data = append(data, chunk.data...)
			tilecount += terrain.tilesperchunk
		}
	}
	terrain.tilecount = tilecount

	// early return to prevent error
	if tilecount == 0 {
		return
	}

	// update terrain tile positions
	terrain.buffer.GetVAO().GetVertexBuffer(0).UpdateData(poss)

	// update terrain data
	terrain.terrainbuffer.Resize(int(tilecount))
	terrain.terrainbuffer.UploadArray(data)
}
func (terrain *Terrain) Render(M, V, P mgl32.Mat4, camerapos mgl32.Vec3) {
	lightdir := mgl32.Vec3{2.0, 2.0, 0.0}
	lightcolor := mgl32.Vec3{0.0, 1.0, 0.0}

	terrain.terrainbuffer.Bind(0)
	terrain.wind.velocityfield.Bind(1)

	// render terrain
	terrain.shader.Use()
	terrain.shader.UpdateMat4("M", M)
	terrain.shader.UpdateMat4("V", V)
	terrain.shader.UpdateMat4("P", P)
	terrain.shader.UpdateFloat32("tilesize", terrain.tilesize)
	terrain.shader.UpdateVec3("cameraPos", camerapos)
	terrain.shader.UpdateVec3("lightDir", lightdir)
	terrain.shader.UpdateVec3("lightColor", lightcolor)
	terrain.shader.UpdateFloat32("ambientIntensity", 0.4)
	terrain.shader.UpdateFloat32("diffuseIntensity", 0.4)
	terrain.shader.UpdateFloat32("d1", (terrain.unloaddist-200)/8)
	terrain.shader.UpdateFloat32("d2", terrain.unloaddist-200)
	terrain.shader.Render()

	// render grass
	terrain.grass.Render(int32(terrain.tilecount), terrain.tilesize, M, V, P, camerapos)

	terrain.terrainbuffer.Unbind()
	terrain.wind.velocityfield.Unbind()
}
func (terrain *Terrain) GetHeight(pos mgl32.Vec3) float32 {
	x, z := terrain.getChunkPos(pos.X(), pos.Z())

	chunkresolution := int32(mathutils.RoundF32(terrain.chunksize / terrain.tilesize))

	var height float32 = 0.0
	if chunk, ok := terrain.chunks[makeKey(x, z)]; ok {
		// get position within chunk
		dx := pos.X() - float32(x)*terrain.chunksize
		dz := pos.Z() - float32(z)*terrain.chunksize
		if dx < 0 {
			dx = terrain.chunksize + dx
		}
		if dz < 0 {
			dz = terrain.chunksize + dz
		}
		tx := int32(dx / terrain.tilesize)
		tz := int32(dz / terrain.tilesize)
		tidx := tz*chunkresolution + tx

		tileposx := chunk.data[tidx*12+8]
		tileposz := chunk.data[tidx*12+9]
		if mathutils.AbsF32(tileposx-pos.X()) > terrain.tilesize/2 ||
			mathutils.AbsF32(tileposz-pos.Z()) > terrain.tilesize {
			fmt.Printf("Wrong tile (%v,%v)  (%v,%v : %v,%v)\n", pos.X(), pos.Z(), x, z, tx, tz)
			panic("Err")
		}

		if tidx >= 0 && tidx < terrain.tilesperchunk {
			rx := dx - float32(tx)*terrain.tilesize
			rz := dz - float32(tz)*terrain.tilesize
			if rx < rz {
				a := chunk.data[tidx*12]
				b := chunk.data[tidx*12+1]
				c := chunk.data[tidx*12+2]
				d := chunk.data[tidx*12+3]
				height = -(d + a*pos.X() + c*pos.Z()) / b
			} else {
				a := chunk.data[tidx*12+4]
				b := chunk.data[tidx*12+5]
				c := chunk.data[tidx*12+6]
				d := chunk.data[tidx*12+7]
				height = -(d + a*pos.X() + c*pos.Z()) / b
			}
		}
	}
	return height
}
func (terrain *Terrain) load(pos mgl32.Vec3) {
	centerx, centerz := terrain.getChunkPos(pos.X(), pos.Z())

	chunkrad := int32(mathutils.RoundF32(terrain.loaddist / terrain.chunksize))

	// check if chunk position it not in chunks map and create a new one
	for z := -chunkrad; z <= chunkrad; z++ {
		for x := -chunkrad; x <= chunkrad; x++ {
			cx := centerx + x
			cz := centerz + z
			// is chunk not yet present?
			if _, ok := terrain.chunks[makeKey(cx, cz)]; !ok {
				// is chunk in load radius?
				chunkpos := mgl32.Vec3{
					float32(cx)*terrain.chunksize + terrain.chunksize/2,
					0.0,
					float32(cz)*terrain.chunksize + terrain.chunksize/2,
				}
				if distxz(pos, chunkpos) < terrain.loaddist {
					// create a new chunk
					terrain.chunks[makeKey(cx, cz)] = terrain.cf.MakeChunk(cx, cz)
				}
			}
		}
	}
}
func (terrain *Terrain) unload(pos mgl32.Vec3) {
	for key, chunk := range terrain.chunks {
		if distxz(chunk.pos, pos) > terrain.unloaddist {
			delete(terrain.chunks, key)
		}
	}
}
func (terrain *Terrain) getChunkPos(x, z float32) (int32, int32) {
	cx := int32(x / terrain.chunksize)
	cz := int32(z / terrain.chunksize)
	if x < 0 {
		cx -= 1
	}
	if z < 0 {
		cz -= 1
	}
	return cx, cz
}
func makeKey(x, z int32) string {
	return fmt.Sprint(x, "-", z)
}
func distxz(posa, posb mgl32.Vec3) float32 {
	dx := mathutils.AbsF32(posa.X() - posb.X())
	dz := mathutils.AbsF32(posa.Z() - posb.Z())
	return mathutils.SqrtF32(dx*dx + dz*dz)
}
