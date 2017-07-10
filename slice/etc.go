package slice

// Find the kth largest value from two sorted slices in O(lg k) time.
// pre: 1 <= k, len(a1) > 0, len(a2) > 0, k <= len(a1) + len(a2)
// pre (unchecked): a1 and a2 are sorted
// pre violation: panic
// normal return: the kth largest value in both a1 and a2
// Strategy: Make sure that at least one element comes from slice a1,
// then do a binary search of a1 looking for element m1 such that
// m2 == k-m1-2 (in other words, the elements from a1 up to m1 and
// a2 up to m2 total k elements) and such that a1[m1] <= a2[m2+1]
// and a2[m2] <= a1[m1+1]. In such a case, the larger of a1[m1] and
// a2[m2] is the kth element. Though in principle this is very simple
// there are lots of special cases, so the code is a bit fussy.
func kthLargest(a1, a2 []int, k int) int {
	// check preconditions
	if k <= 0 {
		panic("k must be at least 1")
	}
	if len(a1) == 0 || len(a2) == 0 {
		panic("slices cannot be empty")
	}
	if len(a1)+len(a2) < k {
		panic("k is too large")
	}

	// trivial case: the slices do not overlap
	if a1[len(a1)-1] <= a2[0] {
		if k <= len(a1) {
			return a1[k-1]
		} else {
			return a2[k-len(a1)-1]
		}
	}
	if a2[len(a2)-1] <= a1[0] {
		if k <= len(a2) {
			return a2[k-1]
		} else {
			return a1[k-len(a2)-1]
		}
	}

	// binary search case: the slices overlap
	if a2[0] < a1[0] {
		a1, a2 = a2, a1
	}
	if k == 1 {
		return a1[0]
	}
	lo, hi := 0, k-1
	if len(a1) <= hi {
		hi = len(a1) - 1
	}
	for {
		m1 := (lo + hi) / 2
		m2 := k - m1 - 2
		if len(a2) <= m2 {
			m2 = len(a2) - 1
			m1 = k - m2 - 2
		}
		switch {
		case m1 == len(a1)-1:
			return max(a1[m1], a2[m2])
		case m2 == len(a2)-1:
			if a2[m2] <= a1[m1+1] {
				return max(a1[m1], a2[m2])
			} else {
				lo = m1 + 1
			}
		case m2 == -1:
			if a1[m1] <= a2[0] {
				return a1[m1]
			} else {
				hi = m1 - 1
			}
		case a1[m1] <= a2[m2+1] && a2[m2] <= a1[m1+1]:
			return max(a1[m1], a2[m2])
		case a2[m2+1] < a1[m1]:
			hi = m1 - 1
		case a1[m1+1] < a2[m2]:
			lo = m1 + 1
		default:
			panic("Unreachable")
		}
	}
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
