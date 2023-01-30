// Package util provides general-purpose utility functions for package sol
package util

import (
	"image"
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
func Lerp(v0, v1, t float64) float64 {
	if t > 1.0 {
		t = 1.0
	}
	return (1-t)*v0 + t*v1
}

// Smoothstep see http://sol.gfxile.net/interpolation/
func Smoothstep(A, B, v float64) float64 {
	// cards that have spun off the screen may have -ve position
	if A < 0.0 {
		A = 0.0
	}
	if B < 0.0 {
		B = 0.0
	}
	if v > 1.0 {
		v = 1.0
	}
	v = (v) * (v) * (3 - 2*(v))
	X := (B * v) + (A * (1.0 - v))
	return X
}

// Smootherstep see http://sol.gfxile.net/interpolation/
func Smootherstep(A, B, v float64) float64 {
	// cards that have spun off the screen may have -ve position
	if A < 0.0 {
		A = 0.0
	}
	if B < 0.0 {
		B = 0.0
	}
	if v > 1.0 {
		v = 1.0
	}
	v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10)
	X := (B * v) + (A * (1.0 - v))
	return X
}

// func EaseInSine(A, B, v float64) float64 {
// 	v = 1.0 - math.Cos((v*math.Pi)/2.0) // easings.net
// 	return (B * v) + (A * (1.0 - v))
// }

// func EaseInCubic(A, B, v float64) float64 {
// 	v = v * v * v
// 	return (B * v) + (A * (1.0 - v))
// }

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

// Abs returns the absolute value of x
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Max returns the largest of it's two int parameters
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smallest of it's two int parameters
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

// Distance finds the length of the hypotenuse between two points.
func Distance(p1, p2 image.Point) float64 {
	first := math.Pow(float64(p2.X-p1.X), 2)
	second := math.Pow(float64(p2.Y-p1.Y), 2)
	return math.Sqrt(first + second)
}

// DistanceFloat64 finds the length of the hypotenuse between two points.
// Formula is the square root of (x2 - x1)^2 + (y2 - y1)^2
// func DistanceFloat64(x1, y1, x2, y2 float64) float64 {
// 	first := math.Pow(x2-x1, 2)
// 	second := math.Pow(y2-y1, 2)
// 	return math.Sqrt(first + second)
// }

// DistanceInt finds the length of the hypotenuse between two points.
// Formula is the square root of (x2 - x1)^2 + (y2 - y1)^2
// func DistanceInt(x1, y1, x2, y2 int) int {
// 	first := math.Pow(float64(x2)-float64(x1), 2)
// 	second := math.Pow(float64(y2)-float64(y1), 2)
// 	return int(math.Sqrt(first + second))
// }

// OverlapArea returns the intersection of two rectangles
// func OverlapArea(x1, y1, x2, y2, X1, Y1, X2, Y2 int) int {
// 	xOverlap := Max(0, Min(x2, X2)-Max(x1, X1))
// 	yOverlap := Max(0, Min(y2, Y2)-Max(y1, Y1))
// 	return xOverlap * yOverlap
// }

// OverlapAreaFloat64 returns the intersection of two rectangles
// func OverlapAreaFloat64(x1, y1, x2, y2, X1, Y1, X2, Y2 float64) float64 {
// 	xOverlap := math.Max(0, math.Min(x2, X2)-math.Max(x1, X1))
// 	yOverlap := math.Max(0, math.Min(y2, Y2)-math.Max(y1, Y1))
// 	return xOverlap * yOverlap
// }

// OrdinalToShortString converts an ordinal (1..13) to a single(ish) character (A .. K)
func OrdinalToShortString(ord int) string {
	var chars = [14]string{"?", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	return chars[ord]
}

// ShortStringToOrdinal converts a foundation label (probably) into an integer ordinal
// func ShortStringToOrdinal(str string) int {
// 	var chars = [14]string{"?", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
// 	for i, s := range chars {
// 		if str == s {
// 			return i
// 		}
// 	}
// 	log.Panicf("Unknown ordinal short string '%s'", str)
// 	return 0
// }

// RuneToOrdinal convert a single rune to an ordinal (1..13)
// func RuneToOrdinal(r rune) int {
// 	var runes = [14]rune{'?', 'A', '2', '3', '4', '5', '6', '7', '8', '9', 'X', 'J', 'Q', 'K'}
// 	for idx, r2 := range runes {
// 		if r == r2 {
// 			return idx
// 		}
// 	}
// 	return 99 // accept no card
// }

// ParseRunesCard parses short form of card (eg in Deal attribute)
// func ParseRunesCard(runes []rune) (ordinal int, suit int, prone bool) {
// 	// "AC" or "ac" or "A" or "a"
// 	// the suit is optional (think: a tableaux will accept any King)
// 	if len(runes) == 0 {
// 		return // default to 0, 0, false
// 	}
// 	ordinal = RuneToOrdinal(runes[0])
// 	if len(runes) > 1 {
// 		suit = RuneToSuit(runes[1])
// 		prone = unicode.IsLower(runes[1])
// 	}
// 	return
// }

// OrdinalToLongString converts an ordinal (1..13) to a single(ish) character (A .. K)
// func OrdinalToLongString(ord int) string {
// 	var cardValueEnglish [14]string = [14]string{"", "Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
// 	return cardValueEnglish[ord]
// }

// StringToOrdinal converts a short ordinal string to it's longer equivalent
func ShortOrdinalToLongOrdinal(str string) string {
	switch str {
	case "1", "A":
		return "Ace"
	case "J":
		return "Jack"
	case "Q":
		return "Queen"
	case "K":
		return "King"
	default:
		return str
	}
}

// func RuneToSuit(r rune) int {
// 	switch r {
// 	case '♣', 'C', 'c':
// 		return 1 //CLUB
// 	case '♥', 'H', 'h':
// 		return 2 //HEART
// 	case '♦', 'D', 'd':
// 		return 3 //DIAMOND
// 	case '♠', 'S', 's':
// 		return 4 //SPADE
// 	default:
// 		log.Panic("Unknown suit rune", r)
// 	}
// 	return 0
// }

// func SuitToRune(s int) rune {
// 	switch s {
// 	case 1:
// 		return '♣'
// 	case 2:
// 		return '♥'
// 	case 3:
// 		return '♦'
// 	case 4:
// 		return '♠'
// 	default:
// 		log.Panic("Unknown suit", s)
// 	}
// 	return '?'
// }

// Pluralize returns a string containing an attempt at a plural form of the word
// func Pluralize(word string, n int) string {
// 	if n == 1 {
// 		return fmt.Sprintf("%d %s", n, word)
// 	}
// 	return fmt.Sprintf("%d %ss", n, word)
// }

// Duration of a func call
// Arguments to a defer statement are immediately evaluated and stored.
// The deferred function receives the pre-evaluated values when its invoked.
// usage: defer util.Duration(time.Now(), "IntFactorial")
func Duration(invocation time.Time, name string) {
	elapsed := time.Since(invocation)
	if elapsed.Milliseconds() > 0 {
		log.Printf("%s %s", elapsed, name)
	}
}

// DeepCopy deepcopies src to dst using json marshaling
// beware: can turn an int struct member into a string
// func DeepCopy(dst, src interface{}) {
// byt, _ := json.Marshal(src)
// json.Unmarshal(byt, dst)
// }

// Clone deep-copies src to dst
// func Clone(dst, src interface{}) {

// 	buff := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buff)
// 	dec := gob.NewDecoder(buff)
// 	if err := enc.Encode(src); err != nil {
// 		log.Panic("Clone Encode error")
// 	}
// 	if err := dec.Decode(dst); err != nil {
// 		log.Panic("Clone Decode error")
// 	}
// }

// func Contains[T comparable](s []T, e T) bool {
// 	for _, v := range s {
// 		if v == e {
// 			return true
// 		}
// 	}
// 	return false
// }
