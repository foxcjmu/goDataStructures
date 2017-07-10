// This package is home to implementations of standard searching and
// sorting algorithms. All algorithms are for int slices because they
// are meant to illustrate the algorithms. Generalizations for all
// types can be done as in the standard Go sort package, if desired.
//
// author:  C. Fox
// version: 11/2013

package slice

import (
	"math"
)

// Standard stupid bubble sort with no improvements
func BubbleSort(a []int) {
	for j := len(a) - 1; 0 < j; j-- {
		for i := 0; i < j; i++ {
			if a[i+1] < a[i] {
				a[i], a[i+1] = a[i+1], a[i]
			}
		}
	}
}

// Standard selection sort finding the minimum on each pass
func SelectionSort(a []int) {
	for j := 0; j < len(a)-1; j++ {
		minIndex := j
		for i := j + 1; i < len(a); i++ {
			if a[i] < a[minIndex] {
				minIndex = i
			}
		}
		a[j], a[minIndex] = a[minIndex], a[j]
	}
}

// Standard linear insertion sort
func InsertionSort(a []int) {
	for j := 1; j < len(a); j++ {
		element := a[j]
		var i int
		for i = j; 0 < i && element < a[i-1]; i-- {
			a[i] = a[i-1]
		}
		a[i] = element
	}
}

// Binary insertion sort
func BinaryInsertionSort(a []int) {
	for j := 1; j < len(a); j++ {
		element := a[j]
		lo, hi := 0, j-1
		for lo < hi {
			mid := (lo + hi) / 2
			switch {
			case element == a[mid]:
				lo, hi = mid, mid
			case element < a[mid]:
				hi = mid - 1
			case a[mid] < element:
				lo = mid + 1
			}
		}
		if a[lo] < element {
			lo++
		}
		copy(a[lo+1:j+1], a[lo:j])
		a[lo] = element
	}
}

// Shell sort using powers of three as the spacing increment
func ShellSort(a []int) {
	// compute the starting value of h
	h := 1
	for h < len(a)/9 {
		h = 3*h + 1
	}

	// insertion sort using decreasing values of h
	for 0 < h {
		for j := h; j < len(a); j++ {
			element := a[j]
			var i int
			for i = j; h <= i && element < a[i-h]; i -= h {
				a[i] = a[i-h]
			}
			a[i] = element
		}
		h /= 3
	}
}

// Mergesort using an auxiliary slice of size len(a)
func MergeSort(a []int) {
	var mergeInto func([]int, []int)

	// merge sub-lists upward
	mergeInto = func(dst []int, src []int) {
		if len(dst) < 2 {
			return
		}
		m := len(dst) / 2
		mergeInto(src[:m], dst[:m])
		mergeInto(src[m:], dst[m:])
		j, k := 0, m
		for i := 0; i < len(dst); i++ {
			if j < m && k < len(src) {
				if src[j] < src[k] {
					dst[i], j = src[j], j+1
				} else {
					dst[i], k = src[k], k+1
				}
			} else if j < m {
				dst[i], j = src[j], j+1
			} else {
				dst[i], k = src[k], k+1
			}
		}
	}

	auxiliary := make([]int, len(a))
	copy(auxiliary, a)
	mergeInto(a, auxiliary)
}

// ConcurrentMergesort using an auxiliary slice of size len(a) taht sorts sub-lists
// in goroutines if the sub-lists are bigger than the goThreshold.
func ConcurrentMergeSort(a []int) {
	const goThreshold = 60000 // merge in a goroutine for lists bigger than this
	var mergeInto func([]int, []int, chan bool)

	// merge sub-lists upward
	mergeInto = func(dst []int, src []int, done chan bool) {
		if len(dst) < 2 {
			return
		}
		m := len(dst) / 2
		if m < goThreshold {
			mergeInto(src[:m], dst[:m], nil)
			mergeInto(src[m:], dst[m:], nil)
		} else {
			done := make(chan bool)
			go mergeInto(src[:m], dst[:m], done)
			go mergeInto(src[m:], dst[m:], done)
			<-done
			<-done
		}
		j, k := 0, m
		for i := 0; i < len(dst); i++ {
			if j < m && k < len(src) {
				if src[j] < src[k] {
					dst[i], j = src[j], j+1
				} else {
					dst[i], k = src[k], k+1
				}
			} else if j < m {
				dst[i], j = src[j], j+1
			} else {
				dst[i], k = src[k], k+1
			}
		}
		if done != nil {
			done <- true
		}
	}

	auxiliary := make([]int, len(a))
	copy(auxiliary, a)
	mergeInto(a, auxiliary, nil)
}

// Quicksort with no improvements
func Quicksort(a []int) {
	if len(a) < 2 {
		return
	}

	// use the last element as the pivot
	ub := len(a) - 1
	pivot := a[ub]

	// partition the list
	i, j := -1, ub
	for i < j {
		for i++; a[i] < pivot; i++ {
		}
		for j--; 0 < j && a[j] > pivot; j-- {
		}
		a[i], a[j] = a[j], a[i]
	}
	a[j], a[i], a[ub] = a[i], pivot, a[j]

	// recursively sort the sublists
	Quicksort(a[:i])
	Quicksort(a[i+1:])
}

// Concurrent quicksort: add concurrency to basic quicksort with no other improvement.
// Making every recursive call of quicksort into a goroutine actually slows the
// sort down a lot. It appears that goroutine overhead is only worth it for
// larger lists. After some profiling, that size appears to be about 75000.
func ConcurrentQuicksort(a []int) {

	// cqs runs as a goroutine on sub-lists only if the sub-list size does not
	// exceeded the goThreshold.
	const goThreshold = 75000
	var cqs func([]int, chan bool)

	cqs = func(a []int, done chan bool) {
		if len(a) < 2 {
			return
		}

		// use the last element as the pivot
		ub := len(a) - 1
		pivot := a[ub]

		// partition the list
		i, j := -1, ub
		for i < j {
			for i++; a[i] < pivot; i++ {
			}
			for j--; 0 < j && a[j] > pivot; j-- {
			}
			a[i], a[j] = a[j], a[i]
		}
		a[j], a[i], a[ub] = a[i], pivot, a[j]

		// recursively sort the sublists
		if goThreshold < len(a) {
			done := make(chan bool)
			go cqs(a[:i], done)
			go cqs(a[i+1:], done)
			<-done
			<-done
		} else {
			cqs(a[:i], nil)
			cqs(a[i+1:], nil)
		}
		if done != nil {
			done <- true
		}
	}

	cqs(a, nil)
}

// Quicksort with the median-of-three improvement.
func Qsort(a []int) {
	if len(a) < 2 {
		return
	}

	// find sentinels for the list ends and the median for the pivot
	m, ub := len(a)/2, len(a)-1
	if a[m] < a[0] {
		a[m], a[0] = a[0], a[m]
	}
	if a[ub] < a[m] {
		a[ub], a[m] = a[m], a[ub]
	}
	if a[m] < a[0] {
		a[m], a[0] = a[0], a[m]
	}

	// a list of length 2 or 3 is now sorted
	if len(a) < 4 {
		return
	}

	// put the pivot just shy of the end of the list
	pivot := a[m]
	a[m], a[ub-1] = a[ub-1], a[m]

	// partition the list
	i, j := 0, ub-1
	for i < j {
		for i++; a[i] < pivot; i++ {
		}
		for j--; a[j] > pivot; j-- {
		}
		a[i], a[j] = a[j], a[i]
	}
	a[j], a[i], a[ub-1] = a[i], pivot, a[j]

	// recursively sort the sublists
	Qsort(a[:i])
	Qsort(a[i+1:])
}

// Standard heapsort
func Heapsort(a []int) {

	// siftDown makes a from i to maxIndex into a heap
	siftDown := func(a []int, i, maxIndex int) {
		tmp := a[i]
		for j := 2*i + 1; j <= maxIndex; j = 2*i + 1 {
			if j < maxIndex && a[j] < a[j+1] {
				j++
			}
			if a[j] <= tmp {
				break
			}
			a[i], i = a[j], j
		}
		a[i] = tmp
	}

	// we are done when there is only one thing in the list
	if len(a) < 2 {
		return
	}

	// make the entire slice into a heap
	maxIndex := len(a) - 1
	for i := (maxIndex - 1) / 2; 0 <= i; i-- {
		siftDown(a, i, maxIndex)
	}

	// repeatedly remove the root and remake the heap
	for {
		a[0], a[maxIndex] = a[maxIndex], a[0]
		maxIndex--
		if maxIndex <= 0 {
			break
		}
		siftDown(a, 0, maxIndex)
	}
}

// Introspective sort is an improved quicksort that uses heapsort to sort
// in the worst case and insertion sort for small sublists.
// This version also uses concurrency, so it is as fast as it can be.
// Note: The alternative sort is used if the recursion depth exceeds
// DepthThreshold. This alternative is currently Heapsort but it could
// be any O(n lg n) sort (this guarantees that overall performance is
// O(n lg n)). This version also includes concurrency.
func IntrospectiveSort(a []int) {
	const (
		smallThreshold = 16    // insertion sort lists smaller than this
		goThreshold    = 75000 // make goroutine recursive calls for sub-list bigger than this
	)
	var (
		altThreshold = 2 * int(math.Log(float64(len(a)))) // use alternate sort at this depth
		ispectSort   func([]int, int, chan bool)          // recursive helper
	)
	altSort := Heapsort

	// ispectSort does the real work
	ispectSort = func(a []int, recursionCount int, done chan bool) {
		// insertion sort small lists at the end
		if len(a) < smallThreshold {
			InsertionSort(a)
			return
		}

		// find sentinels for the list ends and the median for the pivot
		m, ub := len(a)/2, len(a)-1
		if a[m] < a[0] {
			a[m], a[0] = a[0], a[m]
		}
		if a[ub] < a[m] {
			a[ub], a[m] = a[m], a[ub]
		}
		if a[m] < a[0] {
			a[m], a[0] = a[0], a[m]
		}

		// put the pivot just shy of the end of the list
		pivot := a[m]
		a[m], a[ub-1] = a[ub-1], a[m]

		// partition the list
		i, j := 0, ub-1
		for i < j {
			for i++; a[i] < pivot; i++ {
			}
			for j--; a[j] > pivot; j-- {
			}
			a[i], a[j] = a[j], a[i]
		}
		a[j], a[i], a[ub-1] = a[i], pivot, a[j]

		// depending on depth, either recursively ispecSort or altSort the sublists
		if 0 < recursionCount {
			if goThreshold < len(a) {
				done := make(chan bool)
				go ispectSort(a[:i], recursionCount-1, done)
				go ispectSort(a[i+1:], recursionCount-1, done)
				<-done
				<-done
			} else {
				ispectSort(a[:i], recursionCount-1, nil)
				ispectSort(a[i+1:], recursionCount-1, nil)
			}
		} else {
			altSort(a[:i])
			altSort(a[i+1:])
		}
		if done != nil {
			done <- true
		}
	}

	ispectSort(a, altThreshold, nil)
}

// IsSorted tests to see whether a slice is sorted
func IsSorted(a []int) bool {
	for i := 0; i < len(a)-1; i++ {
		if a[i+1] < a[i] {
			return false
		}
	}
	return true
}
