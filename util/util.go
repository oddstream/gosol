// Copyright ©️ 2021 oddstream.games

package util

import (
	"image"
	"math"
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

// https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
// https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
// func drawCircle(img *ebiten.Image, x0, y0, r int, c color.Color) {
// 	x, y, dx, dy := r-1, 0, 1, 1
// 	err := dx - (r * 2)

// 	for x > y {
// 		img.Set(x0+x, y0+y, c)
// 		img.Set(x0+y, y0+x, c)
// 		img.Set(x0-y, y0+x, c)
// 		img.Set(x0-x, y0+y, c)
// 		img.Set(x0-x, y0-y, c)
// 		img.Set(x0-y, y0-x, c)
// 		img.Set(x0+y, y0-x, c)
// 		img.Set(x0+x, y0-y, c)

// 		if err <= 0 {
// 			y++
// 			err += dy
// 			dy += 2
// 		}
// 		if err > 0 {
// 			x--
// 			dx += 2
// 			err += dx - (r * 2)
// 		}
// 	}
// }

// Forward returns the direction (0-3)
func Forward(dir int) int {
	return dir
}

// Backward returns the direction (0-3)
func Backward(dir int) int {
	d := [4]int{2, 3, 0, 1}
	return d[dir]
}

// Leftward returns the direction (0-3)
func Leftward(dir int) int {
	d := [4]int{3, 0, 1, 2}
	return d[dir]
}

// Rightward returns the direction (0-3)
func Rightward(dir int) int {
	d := [4]int{1, 2, 3, 0}
	return d[dir]
}
