// Copyright ©️ 2021 oddstream.games

package util

import (
	"image"
	"log"
	"math"
	"strconv"
)

// InRect returns true if px,py is within Rect returned by function parameter
func InRect(pt image.Point, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return pt.X > x0 && pt.Y > y0 && pt.X < x1 && pt.Y < y1
}

// Lerp see https://en.wikipedia.org/wiki/Linear_interpolation
func Lerp(v0 float64, v1 float64, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// Smoothstep see http://sol.gfxile.net/interpolation/
func Smoothstep(A float64, B float64, v float64) float64 {
	v = (v) * (v) * (3 - 2*(v)) // smoothstep
	// v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10)	// smootherstep
	X := (B * v) + (A * (1.0 - v))
	return X
}

// Smootherstep see http://sol.gfxile.net/interpolation/
func Smootherstep(A float64, B float64, v float64) float64 {
	v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10) // smootherstep
	X := (B * v) + (A * (1.0 - v))
	return X
}

// Normalize is the opposite of lerp. Instead of a range and a factor, we give a range and a value to find out the factor.
func Normalize(start, finish, value float64) float64 {
	return (value - start) / (finish - start)
}

// MapValue converts a value from the scale [fromMin, fromMax] to a value from the scale [toMin, toMax].
// It’s just the normalize and lerp functions working together.
func MapValue(value, fromMin, fromMax, toMin, toMax float64) float64 {
	return Lerp(toMin, toMax, Normalize(fromMin, fromMax, value))
}

// Clamp a value between min and max values
func Clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}

// GetIntFromMap does what it says on the tin
func GetIntFromMap(info map[string]string, key string) int {
	str, exists := info[key]
	if exists {
		i, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(str + " is not an int")
		}
		return i
	}
	return 0
}

// GetStringFromMap does what it says on the tin
func GetStringFromMap(info map[string]string, key string) string {
	str, exists := info[key]
	if exists {
		return str
	}
	return ""
}
