// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package engine

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

// Obj contains vertices, normals and texture coordinates.
type Obj struct {
	Vertices  []float32
	Normals   []float32
	Texcoords []float32
}

// Face specifies three or more vertices with vertex position, normal and texture coordinate.
type Face struct {
	Vertices  []int
	Normals   []int
	Texcoords []int
}

// LoadObj load an Obj from the specified filepath.
func LoadObj(filepath string) (Obj, error) {
	// opening the file
	file, err := os.Open(filepath)
	if err != nil {
		return Obj{}, err
	}
	defer file.Close()

	// setup temp variables
	tempVertices := []float32{}
	tempNormals := []float32{}
	tempTexcoords := []float32{}

	// faces
	faces := []Face{}
	neighborhood := map[int][]int{}

	// setup obj properties
	vertices := []float32{}
	normals := []float32{}
	texcoords := []float32{}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		switch tokens[0] {
		case "v":
			for i := 1; i < len(tokens); i++ {
				v, err := strconv.ParseFloat(tokens[i], 32)
				if err == nil {
					tempVertices = append(tempVertices, float32(v))
				}
			}
		case "vn":
			for i := 1; i < len(tokens); i++ {
				v, err := strconv.ParseFloat(tokens[i], 32)
				if err == nil {
					tempNormals = append(tempNormals, float32(v))
				}
			}
		case "vt":
			for i := 1; i < len(tokens); i++ {
				v, err := strconv.ParseFloat(tokens[i], 32)
				if err == nil {
					tempTexcoords = append(tempTexcoords, float32(v))
				}
			}
		case "f":
			face := Face{}
			faceIdx := len(faces)
			// iterate over vertices
			for i := 1; i < len(tokens); i++ {
				faceTokens := strings.Split(tokens[i], "/")
				vertexProps := []int{-1, -1, -1}
				for faceIdx, faceToken := range faceTokens {
					idx, err := strconv.Atoi(faceToken)
					if err == nil {
						vertexProps[faceIdx] = idx
					}
				}
				if vertexProps[0] == -1 {
					continue
				}
				if vertexProps[1] == -1 {
					vertexProps[1] = vertexProps[0]
				}
				if vertexProps[2] == -1 {
					vertexProps[2] = vertexProps[0]
				}

				// added vertex infos
				face.Vertices = append(face.Vertices, vertexProps[0])
				face.Texcoords = append(face.Texcoords, vertexProps[1])
				face.Normals = append(face.Normals, vertexProps[2])

				// save neighbor faces for current vertex
				vertexIdx := vertexProps[0]
				if _, ok := neighborhood[vertexIdx]; !ok {
					neighborhood[vertexIdx] = []int{faceIdx}
				} else {
					neighbors := neighborhood[vertexIdx]
					neighbors = append(neighbors, faceIdx)
				}
			}
			faces = append(faces, face)
		}
	}

	// was reading from file was successful ?
	if err := scanner.Err(); err != nil {
		return Obj{}, err
	}

	// build structure from faces
	for _, face := range faces {
		// iterate over all triangles if face has more than one triangle
		for fIdx := 2; fIdx < len(face.Vertices); fIdx++ {
			if len(tempVertices) > 0 {
				i1 := face.Vertices[0] - 1
				i2 := face.Vertices[fIdx-1] - 1
				i3 := face.Vertices[fIdx] - 1
				extractVertexProperty(&vertices, &tempVertices, i1, i2, i3, 3)
			}
			if len(tempNormals) > 0 {
				i1 := face.Normals[0] - 1
				i2 := face.Normals[fIdx-1] - 1
				i3 := face.Normals[fIdx] - 1
				extractVertexProperty(&normals, &tempNormals, i1, i2, i3, 3)
			} else {
				// for all three vertices
				for _, vIdx := range face.Vertices {
					var (
						n1 float32 = 0.0
						n2 float32 = 0.0
						n3 float32 = 0.0
					)
					// iterate all neighboring faces and calc face normals
					// then add all face normals together and normalize
					for _, faceIdx := range neighborhood[vIdx] {
						neighborFace := faces[faceIdx]
						idx1 := neighborFace.Vertices[0]
						idx2 := neighborFace.Vertices[fIdx-1]
						idx3 := neighborFace.Vertices[fIdx]
						tempN1, tempN2, tempN3 := cross(&tempVertices, idx1, idx2, idx3)
						n1 += tempN1
						n2 += tempN2
						n3 += tempN3
					}
					n1, n2, n3 = normalize(n1, n2, n3)
					normals = append(normals, n1)
					normals = append(normals, n2)
					normals = append(normals, n3)
				}
			}
			if len(tempTexcoords) > 0 {
				i1 := face.Texcoords[0] - 1
				i2 := face.Texcoords[fIdx-1] - 1
				i3 := face.Texcoords[fIdx] - 1
				extractVertexProperty(&texcoords, &tempTexcoords, i1, i2, i3, 2)
			} else {
				texcoords = append(texcoords, 0, 0)
				texcoords = append(texcoords, 1, 0)
				texcoords = append(texcoords, 1, 1)
			}
		}
	}

	// calc center of gravity
	vertices = center(vertices)

	return Obj{vertices, normals, texcoords}, nil
}

// extractVertexProperty copies the vertex property of the triangle specifed by the indices idx1, idx2, idx3
// offset times to the outslice.
func extractVertexProperty(outSlice, inSlice *[]float32, idx1, idx2, idx3, offset int) {
	idxs := []int{idx1, idx2, idx3}
	for i := 0; i < 3; i++ {
		fpIdx := idxs[i]
		for o := 0; o < offset; o++ {
			(*outSlice) = append((*outSlice), (*inSlice)[(fpIdx)*offset+o])
		}
	}
}

// cross calculates the normal of the three vertices specified by the indices.
func cross(vertices *[]float32, idx1, idx2, idx3 int) (float32, float32, float32) {
	v1 := (*vertices)[(idx1-1)*3+0]
	v2 := (*vertices)[(idx1-1)*3+1]
	v3 := (*vertices)[(idx1-1)*3+2]
	v4 := (*vertices)[(idx2-1)*3+0]
	v5 := (*vertices)[(idx2-1)*3+1]
	v6 := (*vertices)[(idx2-1)*3+2]
	v7 := (*vertices)[(idx3-1)*3+0]
	v8 := (*vertices)[(idx3-1)*3+1]
	v9 := (*vertices)[(idx3-1)*3+2]

	// calc directions
	a1 := v1 - v4
	a2 := v2 - v5
	a3 := v3 - v6
	b1 := v7 - v4
	b2 := v8 - v5
	b3 := v9 - v6

	// calc cross product
	n1 := a2*b3 - a3*b2
	n2 := a3*b1 - a1*b3
	n3 := a1*b2 - a2*b1

	return normalize(n1, n2, n3)
}

// normalize normalizes the vector (v1, v2,v3).
func normalize(v1, v2, v3 float32) (float32, float32, float32) {
	norm := float32(math.Sqrt(float64(v1*v1 + v2*v2 + v3*v3)))
	if norm != 0.0 {
		return v1 / norm, v2 / norm, v3 / norm
	}
	return v1, v2, v3
}

// center moves all vertex positions to be relative to the center of gravitation.
func center(vertices []float32) []float32 {
	var (
		x float32 = 0.0
		y float32 = 0.0
		z float32 = 0.0

		minX float64 = math.Inf(1)
		maxX float64 = math.Inf(-1)
		minY float64 = math.Inf(1)
		maxY float64 = math.Inf(-1)
		minZ float64 = math.Inf(1)
		maxZ float64 = math.Inf(-1)
	)
	vertexCount := len(vertices) / 3
	for i := 0; i < vertexCount; i++ {
		posX := vertices[i*3+0]
		posY := vertices[i*3+1]
		posZ := vertices[i*3+2]
		x += posX
		y += posY
		z += posZ
		minX = math.Min(minX, float64(posX))
		minY = math.Min(minY, float64(posY))
		minZ = math.Min(minZ, float64(posZ))
		maxX = math.Max(maxX, float64(posX))
		maxY = math.Max(maxY, float64(posY))
		maxZ = math.Max(maxZ, float64(posZ))
	}
	x /= float32(vertexCount)
	y /= float32(vertexCount)
	z /= float32(vertexCount)

	// center vertices
	for i := 0; i < vertexCount; i++ {
		diffX := (maxX - minX) / 2
		diffY := (maxY - minY) / 2
		diffZ := (maxZ - minZ) / 2
		diff := float32(math.Max(diffX, math.Max(diffY, diffZ)))
		vertices[i*3+0] = (vertices[i*3+0] - x) / diff
		vertices[i*3+1] = (vertices[i*3+1] - y) / diff
		vertices[i*3+2] = (vertices[i*3+2] - z) / diff
	}

	return vertices
}
