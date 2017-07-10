// Test BinarySearchTree interface and the Binary Search Tree implementation
// author: C. Fox
// version: 1/2013

package tree

import (
	"strconv"
	"testing"
)

///////////////////////////////////////////////////////
// Make a Comparer Stringer value for insertion
type KeyValue struct {
	key   int
	value string
}

func (p KeyValue) Equal(v interface{}) bool {
	q := v.(KeyValue)
	return p.key == q.key
}

func (p KeyValue) Less(v interface{}) bool {
	q := v.(KeyValue)
	return p.key < q.key
}

func (p KeyValue) String() string {
	return "(" + strconv.Itoa(p.key) + ", " + p.value + ")"
}

///////////////////////////////////////////////////////

func TestEmptyBinarySearchTree(t *testing.T) {
	var r BinarySearchTree

	// make sure a new BinarySearchTree is empty
	if !r.Empty() || r.Size() != 0 {
		t.Error("BinarySearchTree should be empty when new")
	}

	// test operations on an empty tree
	if r.Contains(KeyValue{4, ""}) {
		t.Error("Empty BinarySearchTree should no contain anything")
	}
	r.Remove(KeyValue{4, ""}) // no panic
	if v, ok := r.Get(KeyValue{4, ""}); ok || v != nil {
		t.Error("Empty BinarySearchTree should not allow a get")
	}
}

func TestNonEmptyBinarySearchTree(t *testing.T) {
	var r BinarySearchTree

	// make a tree
	r.Add(KeyValue{20, "twenty"})
	r.Add(KeyValue{10, "ten"})
	r.Add(KeyValue{30, "thirty"})
	r.Add(KeyValue{5, "five"})
	r.Add(KeyValue{15, "fifteen"})
	r.Add(KeyValue{25, "twenty-five"})
	r.Add(KeyValue{30, "thirty"})
	r.Add(KeyValue{27, "twenty-seven"})
	r.Add(KeyValue{3, "three"})
	r.Add(KeyValue{15, "fifteen"})
	r.Add(KeyValue{18, "eighteen"})
	r.Add(KeyValue{26, "twenty-six"})
	if r.Empty() {
		t.Error("BinarySearchTree should not be empty")
	}
	if r.Size() != 10 {
		t.Errorf("BinarySearchTree size should be 9 but is %v", r.Size())
	}
	if v, err := r.RootValue(); err != nil {
		t.Error("BinarySearchTree root should exist")
	} else {
		if p, ok := v.(KeyValue); !ok {
			t.Error("BinarySearchTree root not the right type")
		} else if p.key != 20 {
			t.Errorf("BinarySearchTree root key should be 20 but is %v", p.key)
		}
	}
	if !r.Contains(KeyValue{26, ""}) {
		t.Error("BinarySearchTree should contain 26")
	}
	if r.Contains(KeyValue{13, ""}) {
		t.Error("BinarySearchTree should not contain 13")
	}
	if v, ok := r.Get(KeyValue{27, "glop"}); !ok {
		t.Error("BinarySearchTree should find 27-glop")
	} else if v != (KeyValue{27, "twenty-seven"}) {
		t.Error("BinarySearchTree should return 27-twenty-seven")
	}

	// check insertion order
	inorder := []KeyValue{{3, "three"}, {5, "five"}, {10, "ten"},
		{15, "fifteen"}, {18, "eighteen"}, {20, "twenty"}, {25, "twenty-five"},
		{26, "twenty-six"}, {27, "twenty-seven"}, {30, "thirty"}}
	i := 0
	r.VisitInorder(func(e interface{}) {
		if e != inorder[i] {
			t.Errorf("Inorder internal iterator value is %v should be %v", e, inorder[i])
		}
		i++
	})

	// delete some nodes and check the tree
	r.Remove(KeyValue{8, ""})
	r.Remove(KeyValue{15, ""})
	r.Remove(KeyValue{3, ""})
	r.Remove(KeyValue{25, ""})
	r.Remove(KeyValue{20, ""})
	r.Remove(KeyValue{30, ""})
	r.Remove(KeyValue{10, ""})
	if r.Size() != 4 {
		t.Errorf("BinarySearchTree size should be 4 but is %v", r.Size())
	}
	if r.Height() != 2 {
		t.Errorf("BinarySearchTree height should be 2 but is %v", r.Height())
	}
	if e, _ := r.RootValue(); e != (KeyValue{26, "twenty-six"}) {
		t.Errorf("BinarySearchTree root value should be 25-twenty-five but is %v", e)
	}
	inorder = []KeyValue{{5, "five"}, {18, "eighteen"},
		{26, "twenty-six"}, {27, "twenty-seven"}, {30, "thirty"}}
	i = 0
	r.VisitInorder(func(e interface{}) {
		if e != inorder[i] {
			t.Errorf("Inorder internal iterator value is %v should be %v", e, inorder[i])
		}
		i++
	})

	// make sure a cleared BinarySearchTree is empty
	r.Clear()
	if !r.Empty() || r.Size() != 0 {
		t.Error("BinarySearchTree should be empty after Clear()")
	}

	// special case: empty a tree by deleting
	r.Add(KeyValue{10, "ten"})
	r.Add(KeyValue{25, "twenty-five"})
	r.Add(KeyValue{27, "twenty-seven"})
	r.Add(KeyValue{3, "three"})
	r.Remove(KeyValue{10, ""})
	r.Remove(KeyValue{3, ""})
	r.Remove(KeyValue{27, ""})
	r.Remove(KeyValue{25, ""})
	if !r.Empty() || r.Size() != 0 {
		t.Error("BinarySearchTree should be empty after deletions")
	}
}
