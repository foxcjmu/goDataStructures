package slice

import "testing"
import "math/rand"
import "time"

//import "fmt"

func TestSearches(t *testing.T) {
	const N = 21631
	const M = 1000000
	var a [N]int
	rand.Seed(int64(time.Now().Nanosecond()))
	for index := range a {
		a[index] = rand.Int() % M
	}
	key := a[0]
	Quicksort(a[:])
	testSearch(t, a[:], key, SequentialSearch, "Sequential search")
	testSearch(t, a[:], key, BinarySearchRecursive, "Recursive binary search")
	testSearch(t, a[:], key, BinarySearch, "Non-recursive binary search")
}

func testSearch(t *testing.T, a []int, key int, search func([]int, int) (int, bool), name string) {
	if _, isFound := search(a, key); !isFound {
		t.Errorf("Search %s failed to find a value\n", name)
	}
	if _, isFound := search(a, a[0]); !isFound {
		t.Errorf("Search %s failed to find a value at 0\n", name)
	}
	if _, isFound := search(a, a[len(a)-1]); !isFound {
		t.Errorf("Search %s failed to find a value at end\n", name)
	}
	if i, isFound := search(a, 12345678); isFound {
		t.Errorf("Search %s thinks it found a value at %v\n", name, i)
	}
}
