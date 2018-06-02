package scene

import "github.com/adrianderstroff/realtime-grass/pkg/engine"

type Heightmap struct {
	data      *engine.RawImageData
	maxheight float32
}

func MakeHeightmap(path string, maxheight float32) (Heightmap, error) {
	// load data
	data, err := engine.MakeRawImageData(path)
	if err != nil {
		return Heightmap{}, err
	}

	return Heightmap{
		&data,
		maxheight,
	}, nil
}

func (heightmap *Heightmap) GetWidth() int32 {
	return heightmap.data.GetWidth()
}
func (heightmap *Heightmap) GetHeight() int32 {
	return heightmap.data.GetHeight()
}
func (heightmap *Heightmap) GetHeightAt(x, y int32) float32 {
	// extract image dimensions
	var width int32 = heightmap.data.GetWidth()
	var height int32 = heightmap.data.GetHeight()

	// collect all height values
	heights := []float32{}
	heights = append(heights, getHeightValue(x, y, heightmap.data))
	if x == 0 {
		heights = append(heights, getHeightValue(width-1, y, heightmap.data))
	} else if x == width-1 {
		heights = append(heights, getHeightValue(0, y, heightmap.data))
	}
	if y == 0 {
		heights = append(heights, getHeightValue(x, height-1, heightmap.data))
	} else if y == height-1 {
		heights = append(heights, getHeightValue(x, 0, heightmap.data))
	}

	// even out all values
	len := float32(len(heights))
	var averageheight float32 = 0.0
	for _, height := range heights {
		averageheight += height
	}
	averageheight /= len

	return averageheight * heightmap.maxheight
}
func getHeightValue(x, z int32, image *engine.RawImageData) float32 {
	return float32(image.GetR(x, z)) / 255.0
}
