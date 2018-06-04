// Package scene contains all main entities for rendering and/or interaction with the user.
package scene

import "github.com/adrianderstroff/realtime-grass/pkg/engine"

// Heightmap holds the normalized image data of a height texture as well as the maximum height.
// Values read range from 0 to maximum height.
type Heightmap struct {
	data      *engine.RawImageData
	maxheight float32
}

// MakeHeightmap creates a Heightmap for the image of the given path and a maximum height.
func MakeHeightmap(path string, maxheight float32) (Heightmap, error) {
	// load image data
	data, err := engine.MakeRawImageData(path)
	if err != nil {
		return Heightmap{}, err
	}

	// construct height-map
	heightmap := Heightmap{
		data:      &data,
		maxheight: maxheight,
	}

	// preprocess the image data
	heightmap.preprocessing()

	return heightmap, nil
}

// GetWidth returns the width of the underlying image.
func (heightmap *Heightmap) GetWidth() int32 {
	return heightmap.data.GetWidth()
}

// GetHeight returns the height of the underlying image.
func (heightmap *Heightmap) GetHeight() int32 {
	return heightmap.data.GetHeight()
}

// GetHeight returns the height value at pixel (x,y) within the image.
// x and y have to be in bounds of the image dimensions.
func (heightmap *Heightmap) GetHeightAt(x, y int32) float32 {
	// height is between 0 and 1 thus scale with the maximum height
	height := heightmap.getHeightValue(x, y)
	return height * heightmap.maxheight
}

// getHeightValue returns the height value at pixel (x,z)
func (heightmap *Heightmap) getHeightValue(x, z int32) float32 {
	return float32(heightmap.data.GetR(x, z)) / 255.0
}

// setHeightValue set the red channel of the image at pixel (x,z) to the value val.
func (heightmap *Heightmap) setHeightValue(x, z int32, val float32) {
	heightmap.data.SetR(x, z, uint8(val*255.0))
}

// preprocessing evens out the height values at the borders of the image
func (heightmap *Heightmap) preprocessing() {
	// extract image dimensions
	var width int32 = heightmap.data.GetWidth()
	var height int32 = heightmap.data.GetHeight()

	// top and bottom
	var x int32
	for x = 0; x < width-1; x++ {
		h1 := heightmap.getHeightValue(x, 0)
		h2 := heightmap.getHeightValue(x, height-1)
		avg := (h1 + h2) / 2.0
		heightmap.setHeightValue(x, 0, avg)
		heightmap.setHeightValue(x, height-1, avg)
	}

	// left and right
	var z int32
	for z = 0; z < height-1; z++ {
		h1 := heightmap.getHeightValue(0, z)
		h2 := heightmap.getHeightValue(width-1, z)
		avg := (h1 + h2) / 2.0
		heightmap.setHeightValue(0, z, avg)
		heightmap.setHeightValue(width-1, z, avg)
	}
}
