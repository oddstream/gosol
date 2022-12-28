package sol

import (
	"math/rand"
	"sort"
	"testing"
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

const DECKSIZE int = 52

func BenchmarkShuffle(t *testing.B) {
	var cards []int = make([]int, DECKSIZE)
	var dist []int = make([]int, DECKSIZE)

	for i := 0; i < DECKSIZE; i++ {
		cards[i] = i
	}

	// rand.Seed(time.Now().UnixNano())

	for cycles := 0; cycles < 10000; cycles++ {
		sort.Slice(cards, func(a, b int) bool { return a < b })
		for shuffs := 0; shuffs < 1; shuffs++ {
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
