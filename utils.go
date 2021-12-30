package autotile

import (
	"math/rand"
)

// one chooses one item at random
func one(rng *rand.Rand, items []string) string {
	if items == nil || len(items) == 0 {
		return ""
	}
	return items[rng.Intn(len(items))]
}

// firstFull chooses one 'Full' tile from the first non nil BasicLand
func firstFull(rng *rand.Rand, lnd ...*BasicLand) (string, string) {
	for _, l := range lnd {
		if l == nil {
			continue
		}
		return one(rng, l.Full), l.tag
	}

	return "", Sea
}

// firstTransition chooses one 'Transition' tile from the first non nil BasicLand
func firstTransition(rng *rand.Rand, lnd ...*BasicLand) (string, string) {
	for _, l := range lnd {
		if l == nil {
			continue
		}
		return one(rng, l.Transition), l.tag
	}

	return "", Sea
}
