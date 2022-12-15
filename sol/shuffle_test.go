package sol

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func BenchmarkShuffle(t *testing.B) {
	var cards []int = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	var dist []int = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	rand.Seed(time.Now().UnixNano())

	for cycles := 0; cycles < 1000000; cycles++ {
		sort.Slice(cards, func(a, b int) bool { return a < b })
		for shuffs := 0; shuffs < 6; shuffs++ {
			rand.Shuffle(len(cards), func(a, b int) { cards[a], cards[b] = cards[b], cards[a] })
		}
		for i, j := range cards {
			if i == j {
				dist[i] = dist[i] + 1
			}
		}
	}
	min, max := MinMax(dist)
	t.Logf("min=%d, max=%d, diff=%d", min, max, max-min)
}
