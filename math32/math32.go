// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package math32

import "math"

// MinMax limits the value between a min and a max
func MinMax(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MinMax32 limits the value between a min and a max
func MinMax32(v, min, max int32) int32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func Hypot(x, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

func Logn(x float32) float32 {
	return float32(math.Log(float64(x)))
}

func Log1p(x float32) float32 {
	return float32(math.Log1p(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func RoundF(x float32) float32 {
	if x < 0 {
		return Ceil(x - 0.5)
	}
	return Ceil(x + 0.5)
}
