// Package util provides general-purpose utility functions for package sol
package util

import (
	"fmt"
	"log"
	"math"
	"time"
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
	v = (v) * (v) * (3 - 2*(v))
	X := (B * v) + (A * (1.0 - v))
	return X
}

// Smootherstep see http://sol.gfxile.net/interpolation/
func Smootherstep(A float64, B float64, v float64) float64 {
	v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10)
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

// ClampInt a value between min and max values
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

// DistanceInt finds the length of the hypotenuse between two points.
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
	var chars = [14]string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	return chars[ord]
}

// RuneToOrdinal convert a single rune to an ordinal (1..13)
func RuneToOrdinal(r rune) int {
	var runes = [14]rune{'0', 'A', '2', '3', '4', '5', '6', '7', '8', '9', 'X', 'J', 'Q', 'K'}
	for idx, r2 := range runes {
		if r == r2 {
			return idx
		}
	}
	return 0
}

// OrdinalToLongString converts an ordinal (1..13) to a single(ish) character (A .. K)
func OrdinalToLongString(ord int) string {
	var cardValueEnglish [14]string = [14]string{"", "Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	return cardValueEnglish[ord]
}

// StringToOrdinal converts a string to an int 1..13
func StringToOrdinal(str string) int {
	switch str {
	case "1", "A", "Ace":
		return 1
	case "2", "Two":
		return 2
	case "3", "Three":
		return 3
	case "4", "Four":
		return 4
	case "5", "Five":
		return 5
	case "6", "Six":
		return 6
	case "7", "Seven":
		return 7
	case "8", "Eight":
		return 8
	case "9", "Nine":
		return 9
	case "X", "Ten":
		return 10
	case "J", "Jack":
		return 11
	case "Q", "Queen":
		return 12
	case "K", "King":
		return 13
	}
	log.Panicf("Unknown input to StringToOrdinal '%s'", str)
	return 0
}

func RuneToSuit(r rune) int {
	switch r {
	case '♣', 'C':
		return 1 //CLUB
	case '♥', 'H':
		return 2 //HEART
	case '♦', 'D':
		return 3 //DIAMOND
	case '♠', 'S':
		return 4 //SPADE
	default:
		log.Panic("Unknown suit rune", r)
	}
	return 0
}

// Pluralize returns a string containing an attempt at a plural form of the word
func Pluralize(word string, n int) string {
	if n == 0 {
		return fmt.Sprintf("no %ss", word)
	}
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// Contains tells whether a contains x.
// func SearchStrings(a []string, x string) int
// assumes the input slice is sorted; func Contains does not
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Duration of a func call
// Arguments to a defer statement is immediately evaluated and stored.
// The deferred function receives the pre-evaluated values when its invoked.
// usage: defer util.Duration(time.Now(), "IntFactorial")
func Duration(invocation time.Time, name string) {
	elapsed := time.Since(invocation)
	if elapsed.Milliseconds() > 0 {
		log.Printf("%s took %s", name, elapsed)
	}
}
