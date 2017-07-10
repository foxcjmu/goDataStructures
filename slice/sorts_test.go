package slice

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var _ = fmt.Println

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestSorts(t *testing.T) {
	const M = 1000000

	const smallN = 21631
	small := make([]int, smallN)
	for index := range small {
		small[index] = rand.Int() % M
	}
	smallOracle := make([]int, smallN)
	copy(smallOracle, small)
	sort.IntSlice(smallOracle).Sort()

	const bigN = 2391631
	big := make([]int, bigN)
	for index := range big {
		big[index] = rand.Int() % M
	}
	bigOracle := make([]int, bigN)
	copy(bigOracle, big)
	sort.IntSlice(bigOracle).Sort()

	testSort(t, small, smallOracle, BubbleSort, "Bubble sort")
	testSort(t, small, smallOracle, SelectionSort, "Selection sort")
	testSort(t, small, smallOracle, BinaryInsertionSort, "Binary insertion sort")
	testSort(t, big, bigOracle, ShellSort, "Shell sort")
	testSort(t, big, bigOracle, MergeSort, "Merge sort")
	testSort(t, big, bigOracle, ConcurrentMergeSort, "Concurrent merge sort")
	testSort(t, big, bigOracle, Quicksort, "Basic quicksort")
	testSort(t, big, bigOracle, ConcurrentQuicksort, "Concurrent quicksort")
	testSort(t, big, bigOracle, Qsort, "Improved quicksort")
	testSort(t, big, bigOracle, Heapsort, "Heapsort")
	testSort(t, big, bigOracle, IntrospectiveSort, "Introspective sort")
}

func testSort(t *testing.T, a, oracle []int, sort func([]int), name string) {
	b := make([]int, len(a))
	copy(b, a)
	sort(b)
	for i := range b {
		if oracle[i] != b[i] {
			t.Errorf("%s failed\n", name)
			return
		}
	}
}

func benchmarkSort(b *testing.B, sort func([]int)) {
	a := make([]int, b.N*1000000)
	for index := range a {
		a[index] = rand.Int()
	}
	sort(a)
}

//func BenchmarkQuicksort(b *testing.B)          { benchmarkSort(b, Quicksort) }
//func BenchmarkQsort(b *testing.B)              { benchmarkSort(b, Qsort) }
//func BenchmarkConcurrenQuicksort(b *testing.B) { benchmarkSort(b, ConcurrentQuicksort) }
//func BenchmarkIntrospectiveSort(b *testing.B)  { benchmarkSort(b, IntrospectiveSort) }
func BenchmarkMergeSort(b *testing.B)          { benchmarkSort(b, MergeSort) }
func BenchmarkConcurrenMergeSort(b *testing.B) { benchmarkSort(b, ConcurrentMergeSort) }
