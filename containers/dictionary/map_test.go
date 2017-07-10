// Test Map interface and the TreeMap and HashMap data structures.
//
// author: C. Fox
// version: 1/2016

package dictionary

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf // in case we need fmt for debugging

func TestMaps(t *testing.T) {
	testMap(t, new(TreeMap), "TreeMap ")
	testMap(t, new(HashMap), "HashMap ")
}

type Integer int

// Define a Comparer/Hasher key type
func (i Integer) Equal(c interface{}) bool { return int(i) == int(c.(Integer)) }
func (i Integer) Less(c interface{}) bool  { return int(i) < int(c.(Integer)) }
func (i Integer) Hash(tableSize int) int   { return int(i) }

func testMap(t *testing.T, m Map, name string) {
	// make sure a new Map is empty and works with operations
	if !m.Empty() || 0 != m.Size() {
		t.Error(name + "should be empty and size should be 0 when new")
	}
	if m.Contains("abc") {
		t.Error(name + "empty map should not contain anything")
	}
	if m.HasKey(Integer(3)) {
		t.Error(name + "empty map should not have any keys")
	}
	m.Delete(Integer(7)) // no panic
	m.Apply(func(_ interface{}) {
		t.Error(name + "empty map should not apply and functions")
	})
	if v, ok := m.Get(Integer(3)); ok || v != nil {
		t.Error(name + "Get succeeded for an empty map")
	}

	// add some some data to the map and check that it is there
	values := []string{"two", "three", "five", "ten", "twenty"}
	keys := []Integer{Integer(2), Integer(3), Integer(5), Integer(10), Integer(20)}
	m.Insert(Integer(5), "five")
	m.Insert(Integer(10), "ten")
	m.Insert(Integer(2), "two")
	m.Insert(Integer(3), "three")
	m.Insert(Integer(20), "twenty")
	if m.Size() != 5 {
		t.Errorf(name+"insertion failure: map should have 5 elements but has %v", m.Size())
	}
	//fmt.Print(m)
	for _, v := range values {
		if !m.Contains(v) {
			t.Errorf(name+"map should include value %v but does not", v)
		}
	}
	for _, k := range keys {
		if !m.HasKey(k) {
			t.Errorf(name+"map should include key %v but does not", k)
		}
	}

	// try external iterators
	i := 0
	iter := m.NewIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != values[i] {
			t.Errorf("Value iterator value should be %v but is %v", i, v)
		}
		i++
	}
	i = 0
	iter.Reset()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != values[i] {
			t.Errorf("Value iterator value should be %v but is %v", i, v)
		}
		i++
	}
	if i < 5 {
		t.Error("Value iterator did not iterate all the way through the map")
	}
	if !iter.Done() {
		t.Error("Value iterator should be done")
	}

	i = 0
	iter = m.NewKeyIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != keys[i] {
			t.Errorf("Key tterator value should be %v but is %v", i, v)
		}
		i++
	}
	i = 0
	iter.Reset()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != keys[i] {
			t.Errorf("Key tterator value should be %v but is %v", i, v)
		}
		i++
	}
	if i < 5 {
		t.Error("Key iterator did not iterate all the way through the map")
	}
	if !iter.Done() {
		t.Error("Key iterator should be done")
	}

	// try internal iterators
	i = 0
	var evf func(interface{}) = func(e interface{}) {
		if values[i] != e {
			t.Errorf("Expected %v but got %v during internal iteration", values[i], e)
		}
		i++
	}
	m.Apply(evf)
	if i < 5 {
		t.Error("Internal iterator did not iterator all the way through the map")
	}

	// delete some data and insert some data and see that things are in order
	m.Delete(Integer(2))
	m.Delete(Integer(4))
	m.Delete(Integer(3))
	m.Insert(Integer(4), "four")
	m.Insert(Integer(14), "fourteen")
	keys[0], values[0] = Integer(4), "four"
	keys[1], values[1] = Integer(14), "fourteen"
	if m.Size() != 5 {
		t.Errorf(name+"deletion failure: map should have 5 elements but has %v", m.Size())
	}
	if m.HasKey(Integer(2)) {
		t.Error(name + "contains 2 but it should have been deleted")
	}
	if m.HasKey(Integer(3)) {
		t.Error(name + "contains 3 but it should have been deleted")
	}
	for _, v := range values {
		if !m.Contains(v) {
			t.Errorf(name+"map should include value %v but does not", v)
		}
	}
	for _, k := range keys {
		if !m.HasKey(k) {
			t.Errorf(name+"map should include key %v but does not", k)
		}
	}

	// check Get and copy-over insertion
	if v, ok := m.Get(Integer(4)); !ok {
		t.Error(name + "should contain key 4 but does not")
	} else if v != "four" {
		t.Errorf(name+"has wrong value for key 4: %v", v)
	}
	if v, ok := m.Get(Integer(20)); !ok {
		t.Error(name + "should contain key 20 but does not")
	} else if v != "twenty" {
		t.Errorf(name+"has wrong value for key 20: %v", v)
	}
	m.Insert(Integer(20), "Twenty")
	if v, ok := m.Get(Integer(20)); !ok {
		t.Error(name + "should contain key 20 but does not")
	} else if v != "Twenty" {
		t.Errorf(name+"has wrong value for key 20: %v", v)
	}

	// test clear
	m.Clear()
	if !m.Empty() || 0 != m.Size() {
		t.Error(name + "should be empty and size should be 0 after clear is called")
	}

	// check equality testing
	m0 := new(TreeMap)
	if !m.IsEqual(m0) {
		t.Error(name + "fails equality test on empty maps")
	}
	for i := range values {
		m.Insert(keys[i], values[i])
		m0.Insert(keys[i], values[i])
	}
	if !m.IsEqual(m0) {
		t.Error(name + "fails equality test on non-empty maps")
	}
	m0.Insert(Integer(50), "fifty")
	if m.IsEqual(m0) {
		t.Error(name + "fails equality test on unequal-size maps maps")
	}
	m0.Delete(keys[0])
	if m.IsEqual(m0) {
		t.Error(name + "fails equality test on equal-size maps with different keys")
	}
	m0.Delete(Integer(5))
	m0.Insert(keys[0], "blah")
	if m.IsEqual(m0) {
		t.Error(name + "fails equality test on equal-size maps with same keys but different values")
	}
}
