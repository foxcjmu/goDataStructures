// PA3: A program comparing the performance of various sorts.
// author: C. Fox
// version: 1/2016

package main

import (
	"fmt"
	"math/rand"
	"time"

	"slice"
)

// timeSort copies a slice and then times how long it takes to sort it using
// the sorter function argument.
func timeSort(a []int, sorter func([]int)) float64 {
	b := make([]int, len(a))
	copy(b, a)
	start := time.Now().UnixNano()
	sorter(b)
	end := time.Now().UnixNano()
	return float64(end-start) / 1e9
}

// Report on three comparions:
// - All sorts on random data in slices of sizes ranging from 10,000 to 80,000 ints
// - Fast sorts on random data in slices of sizes ranging from 100,000 to 12,800,000 ints
// - Sorts fast for sorted lists except basic quicksort, on ordered slices ranging from
//   10,000 to 160,000 ints (shows how quicksort blows up on sorted lists)
func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	fmt.Println("       N    QConc2  QConc3  QConc4  QConc5  QConc6  QConc7")
	for n := 100000; n < 13000000; n *= 2 {
		a := make([]int, n)
		for index := range a {
			a[index] = rand.Int()
		}
		fmt.Printf("%10d  ", n)
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort2))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort3))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort4))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort5))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort6))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort7))
		fmt.Println()
	}
}
