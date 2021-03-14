// Copyright ©️ 2021 oddstream.games

package util

import (
	"fmt"
	"math"
)

// InRect returns true if px,py is within Rect returned by function parameter
func InRect(x, y int, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return x > x0 && y > y0 && x < x1 && y < y1
}

// RectEmpty returns true if rect is empty
func RectEmpty(x0, y0, x1, y1 int) bool {
	return x0 == x1 || y0 == y1
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

// Clamp a value between min and max values
func ClampInt(value, min, max int) int {
	return Min(Max(value, min), max)
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Max returns the largest of of it's two int parameters
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the largest of of it's two int parameters
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Pow returns x ** y
func Pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// DistanceFloat64 finds the length of the hypotenuse between two points.
// Formula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func DistanceFloat64(x1, y1, x2, y2 float64) float64 {
	first := math.Pow(x2-x1, 2)
	second := math.Pow(y2-y1, 2)
	return math.Sqrt(first + second)
}

// Distance finds the length of the hypotenuse between two points.
// Formula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func DistanceInt(x1, y1, x2, y2 int) int {
	first := math.Pow(float64(x2)-float64(x1), 2)
	second := math.Pow(float64(y2)-float64(y1), 2)
	return int(math.Sqrt(first + second))
}

// OverlapArea returns the intersection of two rectangles
func OverlapArea(x1, y1, x2, y2, X1, Y1, X2, Y2 int) int {
	xOverlap := Max(0, Min(x2, X2)-Max(x1, X1))
	yOverlap := Max(0, Min(y2, Y2)-Max(y1, Y1))
	return xOverlap * yOverlap
}

// OverlapAreaFloat64 returns the intersection of two rectangles
func OverlapAreaFloat64(x1, y1, x2, y2, X1, Y1, X2, Y2 float64) float64 {
	xOverlap := math.Max(0, math.Min(x2, X2)-math.Max(x1, X1))
	yOverlap := math.Max(0, math.Min(y2, Y2)-math.Max(y1, Y1))
	return xOverlap * yOverlap
}

// OrdinalToShortString converts an ordinal (1..13) to a single(ish) character (A .. K)
func OrdinalToShortString(ord int) string {
	var chars = []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	return chars[ord]
}

// Pluralize returns a string containing an attempt at a plural form of the word
func Pluralize(word string, n int) string {
	if n == 0 {
		return fmt.Sprintf("no %ss", word)
	}
	if n == 1 {
		return fmt.Sprintf("one %s", word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// Contains tells whether a contains x.
// func SearchStrings(a []string, x string) int
// assumes the input slice is sorted
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
