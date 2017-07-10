// Test TwoThreeTree interface and implementation
// author: C. Fox
// version: 6/2017

package tree

import (
	//"fmt"
	"testing"
)

func makeTestTree() TwoThreeTree {
	var t TwoThreeTree
	t.Add(KeyValue{30, "30"})
	t.Add(KeyValue{40, "40"})
	t.Add(KeyValue{25, "25"})
	t.Add(KeyValue{25, "25"})
	t.Add(KeyValue{50, "50"})
	t.Add(KeyValue{20, "20"})
	t.Add(KeyValue{22, "22"})
	t.Add(KeyValue{27, "27"})
	t.Add(KeyValue{24, "24"})
	t.Add(KeyValue{35, "35"})
	t.Add(KeyValue{10, "10"})
	t.Add(KeyValue{45, "45"})
	return t
}

func TestEmptyTree(t *testing.T) {
	var r TwoThreeTree
	if !r.Empty() || r.Size() != 0 || r.Height() != 0 {
		t.Error("2-3 tree should be empty when new")
	}

	r.Clear()
	if !r.Empty() || r.Size() != 0 || r.Height() != 0 {
		t.Error("Cleared 2-3 tree be empty when new")
	}

	r.Remove(KeyValue{4, ""}) // no panic
	if v, ok := r.Get(KeyValue{4, ""}); ok || v != nil {
		t.Error("Empty 2-3 tree should not allow a get")
	}
}

func TestAdd(t *testing.T) {
	r := makeTestTree()
	if r.Empty() {
		t.Error("In TestAdd, tree should not be empty")
	}
	if v := r.Size(); v != 11 {
		t.Errorf("In TestAdd, size should be 11 but is %v", v)
	}
	if v := r.Height(); v != 2 {
		t.Errorf("In TestAdd, height should be 2 but is %v", v)
	}
}

func TestContainsAndGet(t *testing.T) {
	r := makeTestTree()
	if v, ok := r.Get(KeyValue{25, "g"}); !ok {
		t.Error("Tree should contain 25")
	} else {
		x := v.(KeyValue)
		if x.key != 25 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{10, "g"}); !ok {
		t.Error("Tree should contain 10")
	} else {
		x := v.(KeyValue)
		if x.key != 10 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{20, "g"}); !ok {
		t.Error("Tree should contain 20")
	} else {
		x := v.(KeyValue)
		if x.key != 20 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{22, "g"}); !ok {
		t.Error("Tree should contain 22")
	} else {
		x := v.(KeyValue)
		if x.key != 22 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{24, "g"}); !ok {
		t.Error("Tree should contain 24")
	} else {
		x := v.(KeyValue)
		if x.key != 24 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{30, "g"}); !ok {
		t.Error("Tree should contain 30")
	} else {
		x := v.(KeyValue)
		if x.key != 30 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{40, "g"}); !ok {
		t.Error("Tree should contain 40")
	} else {
		x := v.(KeyValue)
		if x.key != 40 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if v, ok := r.Get(KeyValue{45, "g"}); !ok {
		t.Error("Tree should contain 45")
	} else {
		x := v.(KeyValue)
		if x.key != 45 {
			t.Errorf("Wrong key value: %v", x)
		}
	}
	if _, ok := r.Get(KeyValue{23, "g"}); ok {
		t.Error("Tree should not contain 23")
	}
	if !r.Contains(KeyValue{25, "g"}) {
		t.Error("Tree should contain 25")
	}
	if !r.Contains(KeyValue{10, "g"}) {
		t.Error("Tree should contain 10")
	}
	if !r.Contains(KeyValue{20, "g"}) {
		t.Error("Tree should contain 20")
	}
	if !r.Contains(KeyValue{22, "g"}) {
		t.Error("Tree should contain 22")
	}
	if !r.Contains(KeyValue{24, "g"}) {
		t.Error("Tree should contain 24")
	}
	if !r.Contains(KeyValue{30, "g"}) {
		t.Error("Tree should contain 30")
	}
	if !r.Contains(KeyValue{40, "g"}) {
		t.Error("Tree should contain 40")
	}
	if !r.Contains(KeyValue{45, "g"}) {
		t.Error("Tree should contain 45")
	}
	if !r.Contains(KeyValue{35, "g"}) {
		t.Error("Tree should contain 35")
	}
	if !r.Contains(KeyValue{27, "g"}) {
		t.Error("Tree should contain 27")
	}
	if r.Contains(KeyValue{5, "g"}) {
		t.Error("Tree should not contain 5")
	}
}

var listing string

func stringer(v interface{}) {
	kv := v.(KeyValue)
	listing += kv.value
}

func shapeTest(t *testing.T, r TwoThreeTree, in, pre string, h, s int) {
	listing = ""
	r.Visit(stringer)
	if listing != in {
		t.Errorf("Expected %v got %v", in, listing)
	}
	listing = ""
	r.VisitPreorder(stringer)
	if listing != pre {
		t.Errorf("Expected %v got %v", pre, listing)
	}
	if v := r.Height(); v != h {
		t.Errorf("Expected height of %v got %v", h, v)
	}
	if v := r.Size(); v != s {
		t.Errorf("Expected size of %v got %v", s, v)
	}
}

func TestDelete(t *testing.T) {
	r := makeTestTree()
	r.Remove(KeyValue{45, ""})
	shapeTest(t, r, "10202224252730354050", "25221020243040273550", 2, 10)
	r.Remove(KeyValue{50, ""})
	shapeTest(t, r, "102022242527303540", "252210202430273540", 2, 9)
	r.Remove(KeyValue{50, ""})
	shapeTest(t, r, "102022242527303540", "252210202430273540", 2, 9)
	r.Remove(KeyValue{40, ""})
	shapeTest(t, r, "1020222425273035", "2522102024302735", 2, 8)
	r.Remove(KeyValue{35, ""})
	shapeTest(t, r, "10202224252730", "22251020242730", 1, 7)
	r.Remove(KeyValue{24, ""})
	shapeTest(t, r, "102022252730", "202510222730", 1, 6)
	r.Remove(KeyValue{22, ""})
	shapeTest(t, r, "1020252730", "2027102530", 1, 5)
	r.Remove(KeyValue{25, ""})
	shapeTest(t, r, "10202730", "20102730", 1, 4)
	r.Remove(KeyValue{10, ""})
	shapeTest(t, r, "202730", "272030", 1, 3)
	r.Remove(KeyValue{20, ""})
	shapeTest(t, r, "2730", "2730", 0, 2)
	r.Remove(KeyValue{27, ""})
	shapeTest(t, r, "30", "30", 0, 1)
	r.Remove(KeyValue{30, ""})
	shapeTest(t, r, "", "", 0, 0)
	r.Remove(KeyValue{30, ""})
	shapeTest(t, r, "", "", 0, 0)
	r.Add(KeyValue{30, "30"})
	r.Add(KeyValue{40, "40"})
	r.Add(KeyValue{50, "50"})
	r.Add(KeyValue{10, "10"})
	r.Add(KeyValue{60, "60"})
	r.Add(KeyValue{70, "70"})
	r.Remove(KeyValue{70, "30"})
	shapeTest(t, r, "1030405060", "3050104060", 1, 5)
	r.Add(KeyValue{45, "45"})
	r.Remove(KeyValue{60, ""})
	r.Remove(KeyValue{10, ""})
	shapeTest(t, r, "30404550", "45304050", 1, 4)
	r.Add(KeyValue{10, "10"})
	r.Add(KeyValue{20, "20"})
	r.Add(KeyValue{60, "60"})
	r.Add(KeyValue{15, "15"})
	r.Add(KeyValue{25, "25"})
	r.Add(KeyValue{35, "35"})
	r.Remove(KeyValue{30, "10"})
	shapeTest(t, r, "101520253540455060", "351510202545405060", 2, 9)
	r.Add(KeyValue{5, "5"})
	r.Add(KeyValue{8, "8"})
	r.Add(KeyValue{12, "12"})
	r.Add(KeyValue{18, "18"})
	r.Add(KeyValue{70, "70"})
	r.Remove(KeyValue{35, "5"})
	shapeTest(t, r, "581012151820254045506070", "154085101220182560455070", 2, 13)
	r.Add(KeyValue{55, "55"})
	r.Remove(KeyValue{12, "55"})
	shapeTest(t, r, "581015182025404550556070", "154085102018255060455570", 2, 13)
	r.Remove(KeyValue{8, ""})
	shapeTest(t, r, "51015182025404550556070", "20501551018402545605570", 2, 12)
}

func TestInorderVisitor(t *testing.T) {
	var r TwoThreeTree
	listing = ""
	r.Visit(stringer)
	if listing != "" {
		t.Errorf("Expected nothing got %v", listing)
	}
	r.Add(KeyValue{3, "3"})
	listing = ""
	r.Visit(stringer)
	if listing != "3" {
		t.Errorf("Expected 3 got %v", listing)
	}
	r = makeTestTree()
	listing = ""
	r.Visit(stringer)
	if listing != "1020222425273035404550" {
		t.Errorf("Expected 1020222425273035404550 got %v", listing)
	}
	//fmt.Print(r)
}

func TestPreorderVisitor(t *testing.T) {
	var r TwoThreeTree
	listing = ""
	r.VisitPreorder(stringer)
	if listing != "" {
		t.Errorf("Expected nothing got %v", listing)
	}
	r.Add(KeyValue{3, "3"})
	listing = ""
	r.VisitPreorder(stringer)
	if listing != "3" {
		t.Errorf("Expected 3 got %v", listing)
	}
	r = makeTestTree()
	listing = ""
	r.VisitPreorder(stringer)
	if listing != "2522102024304027354550" {
		t.Errorf("Expected 2522102024304027354550 got %v", listing)
	}
}

func TestPostorderVisitor(t *testing.T) {
	var r TwoThreeTree
	listing = ""
	r.VisitPostorder(stringer)
	if listing != "" {
		t.Errorf("Expected nothing got %v", listing)
	}
	r.Add(KeyValue{3, "3"})
	listing = ""
	r.VisitPostorder(stringer)
	if listing != "3" {
		t.Errorf("Expected 3 got %v", listing)
	}
	r = makeTestTree()
	listing = ""
	r.VisitPostorder(stringer)
	if listing != "1020242227354550304025" {
		t.Errorf("Expected 1020242227354550304025 got %v", listing)
	}
}

func TestInorderIterator(t *testing.T) {
	var r TwoThreeTree
	listing := ""
	iter := r.NewIterator()
	for !iter.Done() {
		if _, ok := iter.Next(); !ok {
			t.Error("Attempting to iterate over an empty tree")
		}
	}
	if listing != "" {
		t.Error("Empty tree iteration is broken")
	}

	listing = ""
	r.Add(KeyValue{30, "30"})
	iter = r.NewIterator()
	for !iter.Done() {
		if v, ok := iter.Next(); !ok {
			t.Error("Failure during iteration over singleton")
		} else {
			listing += v.(KeyValue).value
		}
	}
	if listing != "30" {
		t.Error("Singelton tree iteration is broken")
	}

	listing = ""
	r = makeTestTree()
	iter = r.NewIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		listing += v.(KeyValue).value
	}
	if listing != "1020222425273035404550" {
		t.Error("Big tree iteration is broken")
	}
}
