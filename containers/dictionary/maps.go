// maps.go: Implementation of maps in the container hierarchy
//
// author: C. Fox
// version: 6/2017
//
// The Map interface is for all maps.
//
// map implements two kinds of dictionaries (maps):
//  - TreeMap stores key-value pairs in a search tree by keys
//  - HashMap stores key-value pairs in a hash table by keys
package dictionary

import (
	"containers"
	"containers/internal/hashtbl"
	"containers/internal/tree"
)

// Map is the interface for maps in the container hierarchy.
type Map interface {
	containers.Collection                  // Size, Clear, Empty, Contains, NewIterator, Apply
	Insert(k, v interface{})               // put pair <k,v> in the map; replace <k,w> if any
	Delete(k interface{})                  // remove pair <k,v> from the map, if any
	Get(k interface{}) (interface{}, bool) // retrieve a value by its key
	HasKey(k interface{}) bool             // true iff <k,v> is in the map
	IsEqual(n Map) bool                    // true iff reciever and m have the same pairs
	NewKeyIterator() containers.Iterator   // iterate over keys
}

// Comparable pairs ///////////////////////////////////////////////////////
// cKeyValue holds key-value pairs and implements the Comparer interface so
// pointers to its instances can be used in search trees to implement maps.
type cKeyValue struct {
	key   containers.Comparer
	value interface{}
}

func (kv *cKeyValue) Equal(kw interface{}) bool {
	return kv.key.Equal(kw.(*cKeyValue).key)
}

func (kv *cKeyValue) Less(kw interface{}) bool {
	return kv.key.Less(kw.(*cKeyValue).key)
}

// TreeMap is the data structure for a search-tree-based implementation
// of maps that uses pointers to cKeyValue instances in the nodes.
type TreeMap struct {
	tree tree.AVLTree // holds cKeyValue instances as node values
}

// Size indicates how many items in the tree map.
func (m *TreeMap) Size() int { return m.tree.Size() }

// Clear removes all items from a tree map.
func (m *TreeMap) Clear() { m.tree.Clear() }

// Empty returns true just in case the tree map has no contents.
func (m *TreeMap) Empty() bool { return m.tree.Empty() }

// Contains returns true just in case its argument v is a value
// held in a key-value pair in the tree map.
func (m *TreeMap) Contains(v interface{}) bool {
	iterator := m.NewIterator()
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		if value == v {
			return true
		}
	}
	return false
}

// Apply invokes function f on every value (not key) in the map.
func (m *TreeMap) Apply(f func(interface{})) {
	m.tree.VisitInorder(func(kv interface{}) {
		f(kv.(*cKeyValue).value)
	})
}

// Insert puts the key-value pair <k,v> into a map.
// It replaces the pair <k,w> if it is already there.
func (m *TreeMap) Insert(k, v interface{}) {
	m.tree.Add(&cKeyValue{k.(containers.Comparer), v})
}

// Delete removes a pair <k,v> from a map given the key k.
// Do nothing if it is not there.
func (m *TreeMap) Delete(k interface{}) {
	m.tree.Remove(&cKeyValue{key: k.(containers.Comparer)})
}

// Get retrieves a key-value pair by its key.
// Precondition: The key is in the map.
// Precondition violation: return nil, false.
// Normal return: return the value mapped to the key and true.
func (m *TreeMap) Get(k interface{}) (interface{}, bool) {
	kv := &cKeyValue{key: k.(containers.Comparer)}
	if kw, ok := m.tree.Get(kv); ok {
		return kw.(*cKeyValue).value, true
	}
	return nil, false
}

// HasKey returns true just in case a key-value pair with key
// k is present in the tree map.
func (m *TreeMap) HasKey(k interface{}) bool {
	_, ok := m.tree.Get(&cKeyValue{key: k.(containers.Comparer)})
	return ok
}

// IsEqual returns true jsut in case the receiver map contains
// exactly the same elements as the argument map n.
func (m *TreeMap) IsEqual(n Map) bool {
	if m.Size() != n.Size() {
		return false
	}
	iter := n.NewKeyIterator()
	for k, ok := iter.Next(); ok; k, ok = iter.Next() {
		key := k.(containers.Hasher)
		nValue, _ := n.Get(key)
		mValue, ok := m.Get(key)
		if !ok {
			return false
		}
		if mValue != nValue {
			return false
		}
	}
	return true
}

// TreeMap Value Iterator ////////////////////////////////////////////////
// treeMapValueIterator keeps track of the state of value iteration over a
// search tree whose nodes are pointers to instances of key-value pairs.
type treeMapValueIterator struct {
	treeIter containers.Iterator // iterator over the search tree
}

// Reset prepares for a new iteration.
func (iter *treeMapValueIterator) Reset() { iter.treeIter.Reset() }

// Done returns true iff iteration is complete.
func (iter *treeMapValueIterator) Done() bool { return iter.treeIter.Done() }

// Next returns the next value in the iteration.
// Precondition: Iteration is not complete.
// Precondition violation: return nil and false.
// Normal return: return the value portion of the pair and true.
func (iter *treeMapValueIterator) Next() (interface{}, bool) {
	kv, ok := iter.treeIter.Next()
	if !ok {
		return nil, false
	}
	return kv.(*cKeyValue).value, true
}

// NewIterator creates and returns a new external iterator that
// traverses values (not keys) in the map.
func (m *TreeMap) NewIterator() containers.Iterator {
	result := new(treeMapValueIterator)
	result.treeIter = m.tree.NewInorderIterator()
	return result
}

// TreeMap Key Iterator ////////////////////////////////////////////////
// treeMapKeyIterator keeps track of the state of key iteration over a
// search tree whose nodes are pointers to instances of key-value pairs.
type treeMapKeyIterator struct {
	treeIter containers.Iterator // iterator over the search tree
}

// Reset prepares for a new iteration.
func (iter *treeMapKeyIterator) Reset() { iter.treeIter.Reset() }

// Done returns true iff iteration is complete.
func (iter *treeMapKeyIterator) Done() bool { return iter.treeIter.Done() }

// Next returns the next key in the iteration.
// Precondition: Iteration is not complete.
// Precondition violation: return nil and false.
// Normal return: return the key portion of the pair and true.
func (iter *treeMapKeyIterator) Next() (interface{}, bool) {
	kv, ok := iter.treeIter.Next()
	if !ok {
		return nil, false
	}
	return kv.(*cKeyValue).key, true
}

// NewKeyIterator creates and returns a new external iterator that
// traverses keys (not values) in the map.
func (m *TreeMap) NewKeyIterator() containers.Iterator {
	result := new(treeMapKeyIterator)
	result.treeIter = m.tree.NewInorderIterator()
	return result
}

// HashMap is the data structure for a hash-table-based implementation
// of maps that uses pointers to hKeyValue instances in the table.
type HashMap struct {
	table hashtbl.HashTable // holds hKeyValue instances as node values
}

// Size returns the number of values in the map.
func (m *HashMap) Size() int { return m.table.Size() }

// Clear makes the map empty.
func (m *HashMap) Clear() { m.table.Clear() }

// Empty returns true iff this map is empty.
func (m *HashMap) Empty() bool { return m.table.Empty() }

// Contains returns true just in case its argument v is a value
// held in a key-value pair in the tree map.
func (m *HashMap) Contains(v interface{}) bool {
	iter := m.table.NewIterator()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		if value == v {
			return true
		}
	}
	return false
}

// Apply invokes function f on every value (not key) in the map.
func (m *HashMap) Apply(f func(interface{})) {
	iter := m.table.NewIterator()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		f(value)
	}
}

// Insert puts a pair <k,v> into a hash map. It replaces any pair
// with the same key  <k,w> if it is already there.
func (m *HashMap) Insert(k, v interface{}) {
	m.table.Insert(k.(containers.Hasher), v)
}

// Delete removes a pair <k,v> from a map given its key k.
// It does nothing if the key is is not in the map.
func (m *HashMap) Delete(k interface{}) {
	m.table.Delete(k.(containers.Hasher))
}

// Get retrieves a value by its key.
// Precondition: The key is in the map.
// Precondition violation: return nil, false.
// Normal return: return the value mapped to the key and true
func (m *HashMap) Get(k interface{}) (interface{}, bool) {
	if value, ok := m.table.Get(k.(containers.Hasher)); ok {
		return value, true
	}
	return nil, false
}

// HasKey returns true just in case the hash map contains a
// key-value pair with key k.
func (m *HashMap) HasKey(k interface{}) bool {
	_, ok := m.table.Get(k.(containers.Hasher))
	return ok
}

// IsEqual returns true just in case the receiver map contains
// exactly the same elements as the argument map n.
func (m *HashMap) IsEqual(n Map) bool {
	if m.Size() != n.Size() {
		return false
	}
	iter := n.NewKeyIterator()
	for k, ok := iter.Next(); ok; k, ok = iter.Next() {
		key := k.(containers.Hasher)
		nValue, _ := n.Get(key)
		mValue, ok := m.Get(key)
		if !ok {
			return false
		}
		if mValue != nValue {
			return false
		}
	}
	return true
}

// NewIterator creates and returns a new external iterator that
// traverses values (not keys) in the map.
func (m *HashMap) NewIterator() containers.Iterator {
	return m.table.NewIterator()
}

// NewKeyIterator creates and returns a new external iterator that
// traverses keys (not values) in the map.
func (m *HashMap) NewKeyIterator() containers.Iterator {
	return m.table.NewKeyIterator()
}
