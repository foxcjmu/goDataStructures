package slice

import "testing"
import "math/rand"
import "time"

//import "fmt"

func TestKthLargest(t *testing.T) {
	// test non-overlapping slices
	const N = 20
	a1 := make([]int, N)
	a2 := make([]int, N)
	a3 := make([]int, N)

	rand.Seed(int64(time.Now().Nanosecond()))

	const M = 1000
	for index := range a1 {
		a1[index] = rand.Int() % M
	}
	for index := range a2 {
		a2[index] = M + rand.Int()%M
	}
	copy(a3, a1)
	a3 = append(a3, a2...)

	IntrospectiveSort(a1)
	IntrospectiveSort(a2)
	IntrospectiveSort(a3)

	// test non-overlappig slices
	if a3[0] != kthLargest(a1, a2, 1) {
		t.Errorf("Failure to find smallest element in non-overlapping slices")
	}
	if a3[0] != kthLargest(a2, a1, 1) {
		t.Errorf("Failure to find smallest element in non-overlapping slices")
	}
	if a3[len(a1)-1] != kthLargest(a1, a2, len(a1)) {
		t.Errorf("Failure to find largest element in first slice in non-overlapping slices")
	}
	if a3[len(a1)-1] != kthLargest(a2, a1, len(a1)) {
		t.Errorf("Failure to find largest element in first slice in non-overlapping slices")
	}
	if a3[len(a1)] != kthLargest(a1, a2, len(a1)+1) {
		t.Errorf("Failure to find smallest element in second slice in non-overlapping slices")
	}
	if a3[len(a1)] != kthLargest(a2, a1, len(a1)+1) {
		t.Errorf("Failure to find smallest element in second slice in non-overlapping slices")
	}
	if a3[len(a3)-1] != kthLargest(a1, a2, len(a1)+len(a2)) {
		t.Errorf("Failure to find largest element in second slice in non-overlapping slices")
	}
	if a3[len(a3)-1] != kthLargest(a2, a1, len(a1)+len(a2)) {
		t.Errorf("Failure to find largest element in second slice in non-overlapping slices")
	}

	// test overlapping slices
	const NUM_TESTS = 5000
	for i := 0; i < NUM_TESTS; i++ {
		a1 = make([]int, rand.Intn(1000)+1)
		for index := range a1 {
			a1[index] = rand.Int() % M
		}
		a2 = make([]int, rand.Intn(1000)+1)
		for index := range a2 {
			a2[index] = rand.Int() % M
		}

		a3 = make([]int, 0)
		a3 = append(a3, a1...)
		a3 = append(a3, a2...)

		IntrospectiveSort(a1)
		IntrospectiveSort(a2)
		IntrospectiveSort(a3)
		k := rand.Intn(len(a3)-1) + 1
		if a3[k-1] != kthLargest(a2, a1, k) {
			t.Errorf("Failure: k = %d", k)
		}
	}
}
