// Test HashTable implementation.
//
// author: C. Fox
// version: 8/2012

package hashtbl

import (
	//"fmt"
	"testing"
	//"containers"
)

type Integer int

type KeyValue struct {
	key   Integer
	value string
}

// Make a Hasher type for HashTable values
func (key Integer) Equal(other interface{}) bool {
	return key == other.(Integer)
}

func (key Integer) Hash(tableSize int) int {
	return int(key) % tableSize
}

func TestHashTableCreationSize(t *testing.T) {
	data := []struct {
		input    int
		expected int
	}{{1, DefaultTableSize},
		{2, DefaultTableSize},
		{3, 3},
		{4, 5},
		{5, 5},
		{6, 7},
		{14, 17},
		{155, 157}}

	for _, d := range data {
		size := NewHashTable(d.input).TableSize()
		if size != d.expected {
			t.Errorf("HashTable table size should be %v but is %v", d.expected, size)
		}
	}
}

func TestHashTableCreationProperties(t *testing.T) {
	table := NewHashTable()

	// make sure a new HashTable is empty and operations do not fail on it
	if !table.Empty() || 0 != table.Size() {
		t.Error("HashTable should be empty and size should be 0 when new")
	}
	table.Clear() // no panic
	if _, ok := table.Get(Integer(5)); ok {
		t.Error("Empty Hashtable should not retrieve anything")
	}
	if v, _ := table.Get(Integer(5)); v != nil {
		t.Error("Empty Hashtable should not retrieve anything")
	}
	table.Delete(Integer(5)) // no panic
	for iter := table.NewIterator(); !iter.Done(); iter.Next() {
		t.Error("Empty HashTable should not do value iteration")
	}
	for iter := table.NewKeyIterator(); !iter.Done(); iter.Next() {
		t.Error("Empty HashTable should not do key iteration")
	}
}

func TestNonEmptyHashTable(t *testing.T) {
	table := NewHashTable(6)

	values := []KeyValue{{0, "zero"}, {1, "one"}, {2, "two"}, {3, "three"},
		{4, "four"}, {5, "five"}, {6, "six"}, {7, "seven"},
		{8, "eight"}, {9, "nine"}, {10, "ten"}, {11, "eleven"},
		{12, "twelve"}, {13, "thirteen"}, {14, "fourteen"}, {15, "fifteen"},
		{16, "sixteen"}, {17, "seventeen"}, {18, "eighteen"}, {19, "nineteen"},
		{20, "twenty"}}
	table.Insert(values[3].key, values[3].value)
	table.Insert(values[5].key, values[5].value)
	table.Insert(values[11].key, values[11].value)
	table.Insert(values[12].key, values[12].value)
	table.Insert(values[7].key, values[7].value)
	table.Insert(values[4].key, values[4].value)
	table.Insert(values[13].key, values[13].value)
	table.Insert(values[10].key, values[10].value)
	table.Insert(values[7].key, values[7].value)
	table.Insert(values[16].key, values[16].value)
	table.Insert(values[9].key, values[9].value)
	table.Insert(values[14].key, values[14].value)
	table.Insert(values[8].key, values[8].value)
	table.Insert(values[15].key, values[15].value)
	table.Insert(values[6].key, values[6].value)

	if table.Empty() {
		t.Errorf("Insertion failure: Non empty table considered empty")
	}
	if table.Size() != 14 {
		t.Errorf("Insertion failure: table should have 14 elements but has %v", table.Size())
	}
	//fmt.Print(table)
	for _, kv := range values[3:17] {
		if v, ok := table.Get(kv.key); !ok {
			t.Errorf("Failed to Get %v", kv)
		} else if kv.value != v {
			t.Errorf("Get got the wrong value %v instead of %v", v, kv.value)
		}
	}

	// try external iterators
	found := make([]bool, 14)
	iter := table.NewKeyIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		index := int(v.(Integer))
		if found[index-3] {
			t.Errorf("Key %v returned twice by the iterator", index)
		}
		found[index-3] = true
	}
	for i := range found {
		if !found[i] {
			t.Errorf("Iterator did not enumerate key %v", i+3)
		}
	}
	if !iter.Done() {
		t.Errorf("Iterator should be done")
	}

	// delete some data and insert some data and see that things are in order
	table.Delete(values[2].key)
	table.Delete(values[4].key)
	table.Delete(values[3].key)
	table.Delete(values[16].key)
	table.Delete(values[15].key)
	if table.Size() != 10 {
		t.Errorf("HashTable deletion failure: table should have 10 elements but has %v", table.Size())
	}
	for _, kv := range values[5:15] {
		if v, ok := table.Get(kv.key); !ok {
			t.Errorf("Failed to Get %v", kv)
		} else if kv.value != v {
			t.Errorf("Get got the wrong value %v instead of %v", v, kv)
		}
	}
	table.Insert(values[4].key, values[4].value)
	table.Insert(values[16].key, values[16].value)
	table.Insert(values[3].key, values[3].value)
	table.Insert(values[16].key, values[16].value)
	table.Insert(values[15].key, values[15].value)
	if table.Size() != 14 {
		t.Errorf("HashTable deletion failure: table should have 14 elements but has %v", table.Size())
	}
	for _, kv := range values[3:17] {
		if v, ok := table.Get(kv.key); !ok {
			t.Errorf("Failed to Get %v", kv)
		} else if kv.value != v {
			t.Errorf("Get got the wrong value %v instead of %v", v, kv)
		}
	}

	// test clear
	table.Clear()
	if !table.Empty() || 0 != table.Size() {
		t.Errorf("HashTable should be empty and size should be 0 after clear is called")
	}
}
