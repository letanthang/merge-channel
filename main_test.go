package main

import "testing"

func TestMerge(t *testing.T) {
	c := merge(asChan(1, 2, 3), asChan(4, 5, 6), asChan(7, 8, 9))
	seen := make(map[int]bool)
	for v := range c {
		if seen[v] {
			t.Errorf("saw %d at least twice", v)
		}
		seen[v] = true
	}

	for i := 1; i <= 9; i++ {
		if !seen[i] {
			t.Errorf("didn't see %d", i)
		}
	}
}
