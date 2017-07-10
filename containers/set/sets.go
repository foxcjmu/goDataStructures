// sets.go: Implementation of sets in the container hierarchy
//
// author: C. Fox
// version: 7/2012
//
// The Set interface is for all sets.
//
// set implements two kinds of sets:
//  - HashSet stores values in a hash table
//  - TreeSet stores values in a binary search tree
package set

import (
	"containers"
	"containers/internal/hashtbl"
	"containers/internal/tree"
)

// Set is the interface for sets in the containers hierarchy.
type Set interface {
	containers.Collection     // Size, Clear, Empty, Contains, NewIterator, Apply
	Subset(set Set) bool      // Say whether the receiver is contained in another set
	Insert(e interface{})     // Put e into a set--replace the value if it is already there
	Delete(e interface{})     // Remove e from a set--do nothing it is not there
	Intersection(set Set) Set // Create the intersection of the receiver and set
	Union(set Set) Set        // Create the union of the receiver and set
	Complement(set Set) Set   // Create the relative complemenh of the receiver and set
	Equal(set Set) bool       // true iff set is identical to the receiver
}

// TreeSet ////////////////////////////////////////////////////////////
// TreeSet is the data structure for a search-tree-based implementation
// of sets that uses values that implement the Comparer interface.
type TreeSet struct {
	tree tree.AVLTree // holds comparable set members as node values
}

// Size returns the number of values in the set.
func (s *TreeSet) Size() int { return s.tree.Size() }

// Clear makes the set empty.
func (s *TreeSet) Clear() { s.tree.Clear() }

// Empty returns true iff this set is empty.
func (s *TreeSet) Empty() bool { return s.tree.Empty() }

// Contains returns true iff this set includes value e.
// Get() in a binary search tree is much faster than Contains in a binary
// tree, so we convert e to a Comparer value and use Get.
func (s *TreeSet) Contains(e interface{}) bool {
	return s.tree.Contains(e)
}

// NewIterator creates and returns a new external iterator value.
func (s *TreeSet) NewIterator() containers.Iterator {
	return s.tree.NewInorderIterator()
}

// Apply invokes function f on every value in the set.
func (s *TreeSet) Apply(f func(interface{})) { s.tree.VisitInorder(f) }

// Equal returns true iff the receiver contains the same elements as set.
func (s *TreeSet) Equal(set Set) bool {
	if s.Size() != set.Size() {
		return false
	}
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !set.Contains(e) {
			return false
		}
	}
	return true
}

// Subset returns true iff the receiver is contained in another set.
func (s *TreeSet) Subset(set Set) bool {
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !set.Contains(e) {
			return false
		}
	}
	return true
}

// Insert puts e into the receiver or replaces e if it is already there.
func (s *TreeSet) Insert(e interface{}) { s.tree.Add(e.(containers.Comparer)) }

// Delete removes e from the receiver, or does nothing if it is not there.
func (s *TreeSet) Delete(e interface{}) { s.tree.Remove(e.(containers.Comparer)) }

// Intersection returns the intersection of the receiver and set.
func (s *TreeSet) Intersection(set Set) Set {
	result := new(TreeSet)
	iter := s.tree.NewPreorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if set.Contains(e) {
			result.Insert(e)
		}
	}
	return result
}

// Union returns the union of the receiver and set.
func (s *TreeSet) Union(set Set) Set {
	result := new(TreeSet)
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		result.Insert(e)
	}
	iter = set.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		result.Insert(e)
	}
	return result
}

// Complement returns the relative complement of the receiver and set.
func (s *TreeSet) Complement(set Set) Set {
	result := new(TreeSet)
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !set.Contains(e) {
			result.Insert(e)
		}
	}
	return result
}

// HashSet ////////////////////////////////////////////////////////////
// HashSet is the data structure for a hash-table-based implementation
// of sets that uses values that implement the Hasher interface.
type HashSet struct {
	table hashtbl.HashTable // holds hashed set members as keys and values
}

// Size returns the number of values in the set.
func (s *HashSet) Size() int { return s.table.Size() }

// Clear makes the set empty.
func (s *HashSet) Clear() { s.table.Clear() }

// Empty returns true iff this set is empty.
func (s *HashSet) Empty() bool { return s.table.Empty() }

// Contains returns true iff this set includes value e.
func (s *HashSet) Contains(e interface{}) bool {
	if _, ok := s.table.Get(e.(containers.Hasher)); ok {
		return true
	}
	return false
}

// NewIterator creates and returns a new external iterator value.
func (s *HashSet) NewIterator() containers.Iterator { return s.table.NewIterator() }

// Apply invokes function f on every value in the set.
func (s *HashSet) Apply(f func(interface{})) {
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		f(e)
	}
}

// Equal returns true iff the receiver contains the same elements as set.
func (s *HashSet) Equal(set Set) bool {
	if s.Size() != set.Size() {
		return false
	}
	iter := set.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !s.Contains(e.(containers.Hasher)) {
			return false
		}
	}
	return true
}

// Subset returns true iff the receiver is contained in another set.
func (s *HashSet) Subset(set Set) bool {
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !set.Contains(e.(containers.Hasher)) {
			return false
		}
	}
	return true
}

// Insert puts e into the receiver or replaces e if it is already there.
func (s *HashSet) Insert(e interface{}) { s.table.Insert(e.(containers.Hasher), e) }

// Delete removes e from the receiver, or does nothing if it is not there.
func (s *HashSet) Delete(e interface{}) { s.table.Delete(e.(containers.Hasher)) }

// Intersection returns the intersection of the receiver and set.
func (s *HashSet) Intersection(set Set) Set {
	result := new(HashSet)
	iter := set.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if s.Contains(e.(containers.Hasher)) {
			result.Insert(e)
		}
	}
	return result
}

// Union returns the union of the receiver and set.
func (s *HashSet) Union(set Set) Set {
	result := new(HashSet)
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		result.Insert(e)
	}
	iter = set.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		result.Insert(e)
	}
	return result
}

// Complement returns the relative complement of the receiver and set.
func (s *HashSet) Complement(set Set) Set {
	result := new(HashSet)
	iter := s.NewIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if !set.Contains(e) {
			result.Insert(e)
		}
	}
	return result
}
