package mathutils

import "math"

func MaxF32(vala, valb float32) float32 {
	return float32(math.Max(float64(vala), float64(valb)))
}

func AbsF32(val float32) float32 {
	return float32(math.Abs(float64(val)))
}

func RoundF32(val float32) float32 {
	return float32(math.Round(float64(val)))
}

func CeilF32(val float32) float32 {
	return float32(math.Ceil(float64(val)))
}

func SqrtF32(val float32) float32 {
	return float32(math.Sqrt(float64(val)))
}

func MapI32(val, srcstart, srcend, dststart, dstend int32) int32 {
	return int32(MapF32(float32(val), float32(srcstart), float32(srcend), float32(dststart), float32(dstend)))
}
func MapF32(val, srcstart, srcend, dststart, dstend float32) float32 {
	x := (val - srcstart) / (srcend - srcstart)
	return x*(dstend-dststart) + dststart
}
