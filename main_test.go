package main

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	c := mergeRecursive(asChan(1, 2, 3), asChan(4, 5, 6), asChan(7, 8, 9))
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

var funcs = []struct {
	name string
	f    func(...<-chan int) <-chan int
}{
	{"goroutines", merge},
	{"recursive", mergeRecursive},
}

func BenchmarkMerge(b *testing.B) {
	for _, f := range funcs {
		for n := 1; n <= 1024; n *= 2 {
			chans := make([]<-chan int, n)
			b.Run(fmt.Sprintf("%s/%d", f.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for i := range chans {
						chans[i] = asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
					}
					c := f.f(chans...)
					for range c {
					}
				}
			})

		}
	}

}
