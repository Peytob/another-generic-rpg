package isomath

func ClampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func FloorToInt(x float32) int {
	i := int(x)
	if x < 0 && float32(i) != x {
		i--
	}
	return i
}

func CeilToInt(x float32) int {
	i := int(x)
	if x > 0 && float32(i) != x {
		i++
	}
	return i
}
