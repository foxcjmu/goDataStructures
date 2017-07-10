// Test Set interface and the TreeSet and HashSet data structures.
//
// author: C. Fox
// version: 1/2016

package set

import (
	"fmt"
	"testing"

	"containers"
)

var _ = fmt.Printf // in case we need fmt for debugging

// Make a Hasher and Comparer type
type KeyValue struct {
	key   int
	value string
}

func (p KeyValue) Equal(q interface{}) bool {
	return p.key == q.(KeyValue).key
}

func (p KeyValue) Less(q interface{}) bool {
	return p.key < q.(KeyValue).key
}

func (p KeyValue) Hash(tableSize int) int {
	return p.key % tableSize
}

func TestSets(t *testing.T) {
	testSet(t, new(TreeSet), "TreeSet ")
	testSet(t, new(HashSet), "HashSet ")
}

func testSet(t *testing.T, set Set, name string) {
	// make sure a new Set is empty and that operations work on it
	if !set.Empty() || 0 != set.Size() {
		t.Error(name + "should be empty and size should be 0 when new")
	}

	values := []KeyValue{{2, "two"}, {3, "three"}, {5, "five"}, {10, "ten"}, {20, "twenty"}}

	// add some some data to the set and check that it is there
	set.Insert(values[3])
	set.Insert(values[0])
	set.Insert(values[4])
	set.Insert(values[2])
	set.Insert(values[1])
	if set.Size() != 5 {
		t.Errorf(name+"insertion failure: set should have 5 elements but has %v", set.Size())
	}
	//fmt.Print(set)
	for _, kv := range values {
		if !set.Contains(kv) {
			t.Errorf(name+"set should include %v but does not", kv)
		}
	}

	// try external iterators
	var iter containers.Iterator
	i := 0
	iter = set.NewIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != values[i] {
			t.Errorf("Iterator value should be %v but is %v", i, v)
		}
		i++
	}
	if i < 5 {
		t.Error("Iterator did not iterate all the way through the set")
	}
	if !iter.Done() {
		t.Error("Iterator should be done")
	}

	// try internal iterators
	i = 0
	var evf func(interface{}) = func(e interface{}) {
		if values[i] != e {
			t.Errorf("Expected %v but got %v during internal iteration", values[i], e)
		}
		i++
	}
	set.Apply(evf)
	if i < 5 {
		t.Error("Internal iterator did not iterator all the way through the set")
	}

	// make a new set with the same data and make sure the sets are equal
	var s1 Set
	if name == "HashSet " {
		s1 = new(HashSet)
	} else {
		s1 = new(TreeSet)
	}
	for _, v := range values {
		s1.Insert(v)
	}
	if !set.Equal(s1) {
		t.Error(name + "Set equality test failed")
	}

	// delete some data and insert some data and see that things are in order
	set.Delete(KeyValue{key: 2})
	set.Delete(KeyValue{key: 4})
	set.Delete(KeyValue{key: 3})
	set.Insert(KeyValue{4, "four"})
	set.Insert(KeyValue{14, "fourteen"})
	values[0] = KeyValue{4, "four"}
	values[1] = KeyValue{14, "fourteen"}
	if set.Size() != 5 {
		t.Errorf(name+"deletion failure: set should have 5 elements but has %v", set.Size())
	}
	if set.Contains(KeyValue{key: 2}) {
		t.Error(name + "contains 2-two but it should have been deleted")
	}
	if set.Contains(KeyValue{key: 3}) {
		t.Error(name + "contains 3-three but it should have been deleted")
	}
	for _, kv := range values {
		if !set.Contains(kv) {
			t.Errorf(name+"set should include %v but does not", kv)
		}
	}

	// do some set operations and check them
	values1 := []KeyValue{{2, "two"}, {3, "three"}, {5, "five"}, {10, "ten"}, {20, "twenty"}}
	values1p := []KeyValue{{2, "two"}, {3, "three"}, {6, "six"}, {10, "ten"}, {20, "twenty"}}
	values2 := []KeyValue{{2, "two"}, {4, "four"}, {5, "five"}, {6, "six"}, {7, "seven"},
		{12, "twelve"}, {20, "twenty"}}
	intersection := []KeyValue{{2, "two"}, {5, "five"}, {20, "twenty"}}
	union := []KeyValue{{2, "two"}, {3, "three"}, {4, "four"}, {5, "five"}, {6, "six"}, {7, "seven"},
		{10, "ten"}, {12, "twelve"}, {20, "twenty"}}
	values1minus2 := []KeyValue{{3, "three"}, {10, "ten"}}
	values2minus1 := []KeyValue{{4, "four"}, {6, "six"}, {7, "seven"}, {12, "twelve"}}
	setEmpty := new(HashSet)
	set1 := new(TreeSet)
	set1p := new(HashSet)
	set2 := new(HashSet)
	setIntersection := new(TreeSet)
	setUnion := new(HashSet)
	set1minus2 := new(TreeSet)
	set2minus1 := new(HashSet)
	for _, v := range values1 {
		set1.Insert(v)
	}
	for _, v := range values1p {
		set1p.Insert(v)
	}
	for _, v := range values2 {
		set2.Insert(v)
	}
	for _, v := range intersection {
		setIntersection.Insert(v)
	}
	for _, v := range union {
		setUnion.Insert(v)
	}
	for _, v := range values1minus2 {
		set1minus2.Insert(v)
	}
	for _, v := range values2minus1 {
		set2minus1.Insert(v)
	}

	// check equality and subsets
	if set1.Equal(set2) {
		t.Error("Set1 equal to set2")
	}
	if set2.Equal(set1) {
		t.Error("Set2 equal to set1")
	}
	if set1.Equal(set1p) {
		t.Error("Set1 equal to set1p")
	}
	if set1p.Equal(set1) {
		t.Error("Set1p equal to set1")
	}
	if set1.Subset(set2) {
		t.Error("Set1 subset of set2")
	}
	if set2.Subset(set1) {
		t.Error("Set2 subset of set1")
	}

	// intersection
	if !set1.Intersection(setEmpty).Equal(setEmpty) {
		t.Error("Set1 intersect setEmpty failed")
	}
	if !setEmpty.Intersection(set1).Equal(setEmpty) {
		t.Error("setEmpty intersect set1 failed")
	}
	if !set1.Intersection(set2).Equal(setIntersection) {
		t.Error("Set1 intersect set2 failed")
	}
	if !set2.Intersection(set1).Equal(setIntersection) {
		t.Error("Set2 intersect set1 failed")
	}
	if !setIntersection.Subset(set1) || !setIntersection.Subset(set2) {
		t.Error("Intersection is not a subset")
	}

	// union
	if !set1.Union(setEmpty).Equal(set1) {
		t.Error("Set1 unite setEmpty failed")
	}
	if !setEmpty.Union(set1).Equal(set1) {
		t.Error("setEmpty unite set1 failed")
	}
	if !set1.Union(set2).Equal(setUnion) {
		t.Error("Set1 unite set2 failed")
	}
	if !set2.Union(set1).Equal(setUnion) {
		t.Error("Set2 union set1 failed")
	}
	if !set1.Subset(setUnion) || !set2.Subset(setUnion) {
		t.Error("United sets not a subset of union")
	}

	// complement
	if !set1.Complement(setEmpty).Equal(set1) {
		t.Error("Set1 complement setEmpty failed")
	}
	if !setEmpty.Complement(set1).Equal(setEmpty) {
		t.Error("setEmpty complement set1 failed")
	}
	if !set1.Complement(set2).Equal(set1minus2) {
		t.Error("Set1 complement set2 failed")
	}
	if !set2.Complement(set1).Equal(set2minus1) {
		t.Error("Set2 complement set1 failed")
	}
	if !set1minus2.Subset(set1) || !set2minus1.Subset(set2) {
		t.Error("complement sets not a subset of left set")
	}

	// test clear
	set.Clear()
	if !set.Empty() || 0 != set.Size() {
		t.Error(name + "should be empty and size should be 0 after clear is called")
	}
}
