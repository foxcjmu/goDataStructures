package slice

// Standard sequential search--return (index, true)
// if the key is found, and (-1, false) otherwise
func SequentialSearch(a []int, key int) (int, bool) {
	for i := 0; i < len(a); i++ {
		if key == a[i] {
			return i, true
		}
	}
	return -1, false
}

// Standard recursive binary search
// Pre: the slice is sorted
// Pre violation: undefined behavior (not checked)
// Normal return: (index, true) if key is found, (-1, false) otherwise
func BinarySearchRecursive(a []int, key int) (int, bool) {
	if len(a) == 0 {
		return -1, false
	}
	m := len(a) / 2
	switch {
	case key == a[m]:
		return m, true
	case key < a[m]:
		return BinarySearchRecursive(a[:m], key)
	case key > a[m]:
		return BinarySearchRecursive(a[m+1:], key)
	}
	panic("Unreachable code reached")
}

// Standard non-recursive binary search
// Pre: the slice is sorted
// Pre violation: undefined behavior (not checked)
// Normal return: (index, true) if key is found, (-1, false) otherwise
func BinarySearch(a []int, key int) (int, bool) {
	for lb, ub := 0, len(a)-1; lb <= ub; {
		m := (lb + ub) / 2
		switch {
		case key == a[m]:
			return m, true
		case key < a[m]:
			ub = m - 1
		case key > a[m]:
			lb = m + 1
		}
	}
	return -1, false
}
