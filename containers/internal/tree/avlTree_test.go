// Test AVLTree interface and the AVL Tree implementation
// author: C. Fox
// version: 6/2017

package tree

//import "fmt"
import "testing"

func TestEmptyAVLTree(t *testing.T) {
	var r AVLTree

	// make sure a new AVLTree is empty
	if !r.Empty() || r.Size() != 0 {
		t.Error("AVLTree should be empty when new")
	}

	// test operations on an empty tree
	if r.Contains(KeyValue{4, ""}) {
		t.Error("Empty AVLTree should no contain anything")
	}
	r.Remove(KeyValue{4, ""}) // no panic
	if v, ok := r.Get(KeyValue{4, ""}); ok || v != nil {
		t.Error("Empty AVLTree should not allow a get")
	}
}

func TestNonEmptyAVLTree(t *testing.T) {
	var r AVLTree

	// make a tree
	r.Add(KeyValue{20, "twenty"})
	r.Add(KeyValue{10, "ten"})
	r.Add(KeyValue{5, "five"})
	r.Add(KeyValue{8, "eight"})
	r.Add(KeyValue{7, "seven"})
	r.Add(KeyValue{3, "three"})
	r.Add(KeyValue{15, "fifteen"})
	r.Add(KeyValue{30, "thirty"})
	r.Add(KeyValue{25, "twenty-five"})
	r.Add(KeyValue{30, "thirty"})
	r.Add(KeyValue{27, "twenty-seven"})
	r.Add(KeyValue{15, "fifteen"})
	r.Add(KeyValue{18, "eighteen"})
	r.Add(KeyValue{26, "twenty-six"})
	if r.Empty() {
		t.Error("AVLTree should not be empty")
	}
	if r.Size() != 12 {
		t.Errorf("AVLTree size should be 12 but is %v", r.Size())
	}
	if r.Height() != 4 {
		t.Errorf("AVLTree height should be 4 but is %v", r.Height())
	}
	if v, err := r.RootValue(); err != nil {
		t.Error("AVLTree root should exist")
	} else {
		if p, ok := v.(KeyValue); !ok {
			t.Error("AVLTree root not the right type")
		} else if p.key != 10 {
			t.Errorf("AVLTree root key should be 10 but is %v", p.key)
		}
	}
	if !r.Contains(KeyValue{26, ""}) {
		t.Error("AVLTree should contain 26")
	}
	if r.Contains(KeyValue{13, ""}) {
		t.Error("AVLTree should not contain 13")
	}
	if v, ok := r.Get(KeyValue{27, "glop"}); !ok {
		t.Error("AVLTree should find 27-glop")
	} else if v != (KeyValue{27, "twenty-seven"}) {
		t.Error("AVLTree should return 27-twenty-seven")
	}

	// check insertion order
	inorder := []KeyValue{{3, "three"}, {5, "five"}, {7, "seven"}, {8, "eight"},
		{10, "ten"}, {15, "fifteen"}, {18, "eighteen"}, {20, "twenty"}, {25, "twenty-five"},
		{26, "twenty-six"}, {27, "twenty-seven"}, {30, "thirty"}}
	i := 0
	r.VisitInorder(func(e interface{}) {
		if e != inorder[i] {
			t.Errorf("Inorder internal iterator value is %v should be %v", e, inorder[i])
		}
		i++
	})

	// delete some nodes and check the tree
	r.Remove(KeyValue{30, ""})
	r.Remove(KeyValue{8, ""})
	r.Remove(KeyValue{20, ""})
	r.Remove(KeyValue{3, ""})
	r.Remove(KeyValue{5, ""})
	r.Remove(KeyValue{15, ""})
	r.Remove(KeyValue{25, ""})
	r.Remove(KeyValue{10, ""})
	if r.Size() != 4 {
		t.Errorf("AVLTree size should be 4 but is %v", r.Size())
	}
	if r.Height() != 2 {
		t.Errorf("AVLTree height should be 2 but is %v", r.Height())
	}
	if e, _ := r.RootValue(); e != (KeyValue{26, "twenty-six"}) {
		t.Errorf("AVLTree root value should be 25-twenty-five but is %v", e)
	}
	inorder = []KeyValue{{7, "seven"}, {18, "eighteen"},
		{26, "twenty-six"}, {27, "twenty-seven"}}
	i = 0
	r.VisitInorder(func(e interface{}) {
		if e != inorder[i] {
			t.Errorf("Inorder internal iterator value is %v should be %v", e, inorder[i])
		}
		i++
	})

	// make sure a cleared AVLTree is empty
	r.Clear()
	if !r.Empty() || r.Size() != 0 {
		t.Error("AVLTree should be empty after Clear()")
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
	if r.Height() != 0 {
		t.Errorf("AVLTree height should be 0 but is %v", r.Height())
	}
	if !r.Empty() || r.Size() != 0 {
		t.Error("AVLTree should be empty after deletions")
	}
}
