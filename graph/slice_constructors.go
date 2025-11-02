package graph

func IntLinearArray(min, max int) []float64 {
	if max == min {
		return nil
	}

	result := make([]float64, max-min)

	for i := range result {
		result[i] = float64(min + i)
	}

	return result
}

func UniformArray(start float64, step float64, length int) []float64 {
	if length <= 0 {
		return nil
	}

	result := make([]float64, length)
	for i := 0; i < length; i++ {
		result[i] = start + float64(i)*step
	}
	return result
}

func LinearArray(min, max float64, length int) []float64 {
	if max == min || length <= 0 {
		return nil
	}

	if length == 1 {
		return []float64{min}
	}

	norm := make([]float64, length)
	step := (max - min) / float64(length-1)
	for i := 0; i < length; i++ {
		norm[i] = min + step*float64(i)
	}

	return norm
}

func ToFloatSlice(ints []int) []float64 {
	floats := make([]float64, len(ints))
	for i, v := range ints {
		floats[i] = float64(v)
	}
	return floats
}

func ScaleArray(arr []float64, scale func(float64) float64) []float64 {
	scaled := make([]float64, len(arr))
	for i, v := range arr {
		scaled[i] = scale(v)
	}
	return scaled
}
