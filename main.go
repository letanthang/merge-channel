package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

func main() {
	a := asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	b := asChan(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	c := asChan(20, 21, 22, 23, 24, 25, 26, 27, 28, 29)

	for i := range mergeRecursive(a, b, c) {
		fmt.Println("mergeRecursive", i)
	}

	// var forever chan int
	// <-forever
}

func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		// defer close(out)
		for _, c := range chans {
			go func(c <-chan int) {
				for v := range c {
					out <- v
				}
				// fmt.Println("done routine")
				wg.Done()
			}(c)
		}
		wg.Wait()
		// fmt.Println("close channel!")
		close(out)
	}()
	return out
}

func mergeReflect(chans ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface().(int)
		}
	}()
	return out
}

func mergeRecursive(chans ...<-chan int) <-chan int {
	switch len(chans) {
	case 0:
		c := make(chan int)
		close(c)
		return c
	case 1:
		return chans[0]
	default:
		m := len(chans) / 2
		return mergeTwo(
			mergeRecursive(chans[:m]...),
			mergeRecursive(chans[m:]...))
	}
}

func mergeTwo(ch1, ch2 <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					fmt.Println("ch1 is done")
					continue
				}
				c <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					fmt.Println("ch2 is done")
					continue
				}
				c <- v
			}

		}
	}()
	return c
}

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
	}()
	return c
}
