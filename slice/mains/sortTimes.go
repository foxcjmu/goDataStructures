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

	fmt.Println("Comparison of all sorts on small slices of random data.")
	fmt.Println("  N    Bubble  Select  Insert   Shell   Merge  QBasic  QImprove  Heap  Inspct")
	for n := 10000; n < 81000; n *= 2 {
		a := make([]int, n)
		for index := range a {
			a[index] = rand.Int() % 1000000
		}
		fmt.Printf("%d  ", n)
		fmt.Printf("%6.3f  ", timeSort(a, slice.BubbleSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.SelectionSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.InsertionSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ShellSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.MergeSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Quicksort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Qsort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Heapsort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.IntrospectiveSort))
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Comparison of fast sorts on large slices of random data.")
	fmt.Println("    N      Shell   Merge   MConc  QBasic   QConc  QImprove  Heap  Inspct")
	for n := 100000; n < 26000000; n *= 2 {
		a := make([]int, n)
		for index := range a {
			a[index] = rand.Int() % 1000000
		}
		fmt.Printf("%8d  ", n)
		fmt.Printf("%6.3f  ", timeSort(a, slice.ShellSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.MergeSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentMergeSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Quicksort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ConcurrentQuicksort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Qsort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Heapsort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.IntrospectiveSort))
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Comparison of selected sorts on small slices of ordered data.")
	fmt.Println("   N    Insert   Shell  QBasic QImprove Inspct")
	for n := 10000; n < 170000; n *= 2 {
		a := make([]int, n)
		for index := range a {
			a[index] = index
		}
		fmt.Printf("%6d  ", n)
		fmt.Printf("%6.3f  ", timeSort(a, slice.InsertionSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.ShellSort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Quicksort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.Qsort))
		fmt.Printf("%6.3f  ", timeSort(a, slice.IntrospectiveSort))
		fmt.Println()
	}
}
