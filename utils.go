package autotile

import (
	"math/rand"
)

// contains returns if `b` contains `a`
func contains(a string, b []string) bool {
	for _, s := range b {
		if s == a {
			return true
		}
	}
	return false
}

// one chooses one item at random
func one(rng *rand.Rand, items []string) string {
	if items == nil || len(items) == 0 {
		return ""
	}
	return items[rng.Intn(len(items))]
}

// firstFull chooses one 'Full' tile from the first non nil BasicLand
func firstFull(rng *rand.Rand, lnd ...*BasicLand) string {
	for _, l := range lnd {
		if l == nil {
			continue
		}
		return one(rng, l.Full)
	}

	return ""
}

// firstTransition chooses one 'Transition' tile from the first non nil BasicLand
func firstTransition(rng *rand.Rand, lnd ...*BasicLand) string {
	for _, l := range lnd {
		if l == nil {
			continue
		}
		return one(rng, l.Transition)
	}
	return ""
}
