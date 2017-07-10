// lists.go -- implementations of ArrayList and LinkedList
// author: C. Fox
// version: 1/2016
//
// The List interface is for all lists.
//
// list implements three kinds of lists:
//  - ArrayList uses slices to store consecutive values
//  - LinedList uses a doubley linked list with a cursor
//  - SinglyLinkedList uses a singly linked list with a cursor
package list

import (
	"containers"
	"fmt"
)

// List is the interface for lists in the container hierarchy.
type List interface {
	containers.Collection              // includes Size, Clear, Empty, NewIterator, and Contains
	Insert(i int, e interface{}) error // insert e at i; pre: 0 <= i <= Size()
	Delete(i int) (interface{}, error) // remove and return element at i; pre: 0 <= i < Size()
	Get(i int) (interface{}, error)    // return element at i; pre: 0 <= i < Size()
	Put(i int, e interface{}) error    // replace element at i; pre: 0 <= i < Size()
	Index(e interface{}) (int, bool)   // return index of e, true, or 0, false if e not present
	Slice(i, j int) (List, error)      // return a duplicate list from i to j-1; pre: 0 <= i <= j <= Size()
	Equal(l List) bool                 // true iff l is identical to the receiver
}

// ArrayList is a contiguous implementation of a list.
type ArrayList struct {
	count int           // how many items are in the list
	store []interface{} // slice holding the list data
}

// An ArrayList is like a slice except that the user does not have to
// worry about the underlying array, so ArrayLists are just slightly more
// abstract than slices. The implementation copies data up and down in the
// slice as one would expect, and reallocates the underlying array as needed.

// arrayListIterator is the data structure for an ArrayList external iterator.
type arrayListIterator struct {
	list *ArrayList // the list that is iterated over
	next int        // which element is next
}

// Reset prepares an iterator to traverse its associated Collection.
func (iter *arrayListIterator) Reset() { iter.next = 0 }

// Done is true iff the iterator has traversed its associated Collection.
func (iter *arrayListIterator) Done() bool { return iter.list.Size() <= iter.next }

// Next returns a value and an indication of whether iteration is complete..
// Precondition: Iteration is not complete.
// Precondition violation: return nil and false.
// Normal return: the next element in the iteration and true.
func (iter *arrayListIterator) Next() (interface{}, bool) {
	if iter.next < iter.list.Size() {
		result, _ := iter.list.Get(iter.next)
		iter.next++
		return result, true
	}
	return nil, false
}

// NewIterator creates and returns an external iterator for an ArrayList.
func (list *ArrayList) NewIterator() containers.Iterator {
	result := new(arrayListIterator)
	result.list = list
	return result
}

// Size indicates how many elements are in the list.
func (list *ArrayList) Size() int { return list.count }

// Clear makes the list empty.
func (list *ArrayList) Clear() { list.count = 0 }

// Empty returns true iff the list has no elements.
func (list *ArrayList) Empty() bool { return list.count == 0 }

// Contains returns true iff element e is in the list.
func (list *ArrayList) Contains(e interface{}) bool {
	for index := 0; index < list.count; index++ {
		if list.store[index] == e {
			return true
		}
	}
	return false
}

// Apply calls function f on every element in the list.
func (list *ArrayList) Apply(f func(interface{})) {
	for index := 0; index < list.count; index++ {
		f(list.store[index])
	}
}

// Insert puts element e into the list at location i.
// Precondition: 0 <= i <= list.count.
// Precondition violation: No insertion and an error indication returned.
// Normal return: e is inserted at i and nil is returned.
func (list *ArrayList) Insert(i int, e interface{}) error {
	if i < 0 || list.count < i {
		return fmt.Errorf("Insert: index out of bounds: %d", i)
	}
	if len(list.store) <= list.count {
		list.store = append(list.store, 0)
	}
	copy(list.store[i+1:], list.store[i:list.count])
	list.store[i] = e
	list.count++
	return nil
}

// Delete removes and returns the element at location i.
// Precondition: 0 <= i < list.count.
// Precondition violation: delete nothing and return an error indication.
// Normal return: delete element i and return it and nil.
func (list *ArrayList) Delete(i int) (interface{}, error) {
	if i < 0 || list.count <= i {
		return nil, fmt.Errorf("Delete: index out of bounds: %d", i)
	}
	result := list.store[i]
	copy(list.store[i:], list.store[i+1:list.count])
	list.count--
	return result, nil
}

// Get returns element i without removing it.
// Precondition: 0 <= i < list.count.
// Precondition violation: return nil and error.
// Normal return: element i and nil.
func (list *ArrayList) Get(i int) (interface{}, error) {
	if i < 0 || list.count <= i {
		return nil, fmt.Errorf("Get: index out of bounds: %d", i)
	}
	return list.store[i], nil
}

// Put changes element i.
// Precondition: 0 <= i < list.count.
// Precondition violation: change nothing and return an error indication.
// Normal return: change the value at i and return nil.
func (list *ArrayList) Put(i int, e interface{}) error {
	if i < 0 || list.count <= i {
		return fmt.Errorf("Put: index out of bounds: %d", i)
	}
	list.store[i] = e
	return nil
}

// Index returns the location of element e. If e is not present,
// return 0 and false; otherwise return the location and true.
func (list *ArrayList) Index(e interface{}) (int, bool) {
	for index := 0; index < list.count; index++ {
		if list.store[index] == e {
			return index, true
		}
	}
	return 0, false
}

// Slice makes a new list duplicating part of this list.
// Precondition: 0 <= i <= j <= list.count.
// Precondition violation: return an empty slice and an error indication.
// Normal return: create a new list and fill it with the items between
// location i and j-1; return this new list and nil.
func (list *ArrayList) Slice(i, j int) (List, error) {
	result := new(ArrayList)
	if i < 0 || j < i || list.count < j {
		return result, fmt.Errorf("Slice: illegal indices: %d %d", i, j)
	}
	result.count = j - i
	result.store = make([]interface{}, result.count)
	copy(result.store[0:], list.store[i:j])
	return result, nil
}

// Equal determines whether another List is identical to this one.
// Two List are identical if they are the same size and have the same
// elements in the same order.
// Precondition: element can be compared using ==.
// Precondition violation: panic.
// Normal return: true iff both lists have the same elements in the same order.
func (list *ArrayList) Equal(l List) bool {
	if list.count != l.Size() {
		return false
	}
	iter := l.NewIterator()
	v, ok := iter.Next()
	for index := 0; index < list.count; index++ {
		if !ok || list.store[index] != v {
			return false
		}
		v, ok = iter.Next()
	}
	return true
}

// String makes a report on the data structure.
func (list *ArrayList) String() string {
	return fmt.Sprintf("ArrayList instance:\nsize: %d\nstore len: %d\nstore cap: %d\nstore: %v\n",
		list.count, len(list.store), cap(list.store), list.store)
}

// LinkedList -----------------------------------------------------------
// A doubly-linked list with a cursor is used to store the values.
// A pointer is kept to the head of the list and to the tail of the list.
// The cursor provides faster access in long lists. The cursor is indefined
// if cursorPtr is nil; otherwise cursorPtr point to the element with
// index cursorIdx.

type node struct {
	item interface{} // data at this node
	pred *node       // the predecessor node
	succ *node       // the successor node
}

// LinkedList is a linked implementation of a list.
type LinkedList struct {
	count     int   // how many elements are in the list
	head      *node // start of a doubly-linked list of values
	tail      *node // end of a doubly-linked list of values
	cursorPtr *node // node where the cursor rests
	cursorIdx int   // index where the cursor rests
}

// linkedListIterator is the data structure for a LinkedList external iterator.
type linkedListIterator struct {
	list    *LinkedList // the list that is iterated over
	current *node       // where we are in the list
}

// Reset prepares an iterator to traverse its associated Collection.
func (iter *linkedListIterator) Reset() { iter.current = iter.list.head }

// Done is true iff the iterator has traversed its associated Collection.
func (iter *linkedListIterator) Done() bool { return iter.current == nil }

// Next returns an interface{} and an indication of whether iteration is complete.
// Precondition: Iteration is not complete.
// Precondition violation: return nil and false.
// Normal return: the next element in the iteration and true.
func (iter *linkedListIterator) Next() (interface{}, bool) {
	if iter.current == nil {
		return nil, false
	}
	result := iter.current.item
	iter.current = iter.current.succ
	return result, true
}

// Size indicates how many elements are in the list.
func (list *LinkedList) Size() int { return list.count }

// Clear removes all elements from the list.
func (list *LinkedList) Clear() {
	list.count = 0
	list.head, list.tail = nil, nil
	list.cursorPtr, list.cursorIdx = nil, 0
}

// Empty returns true just in case the list has not elements.
func (list *LinkedList) Empty() bool { return list.count == 0 }

// NewIterator creates and returns an external iterator for a LinkedList.
func (list *LinkedList) NewIterator() containers.Iterator {
	result := new(linkedListIterator)
	result.list = list
	result.current = list.head
	return result
}

// Apply calls function f on every element in the Collection.
func (list *LinkedList) Apply(f func(interface{})) {
	for ptr := list.head; ptr != nil; ptr = ptr.succ {
		f(ptr.item)
	}
}

// Contains returns true iff element e is in the Collection.
func (list *LinkedList) Contains(e interface{}) bool {
	for ptr := list.head; ptr != nil; ptr = ptr.succ {
		if ptr.item == e {
			return true
		}
	}
	return false
}

// Insert puts element e into the list at location i.
// Precondition: 0 <= i <= list.count.
// Precondition violation: No insertion and an error indication is returned.
// Normal return: e is inserted at i and nil returned.
func (list *LinkedList) Insert(i int, e interface{}) error {
	if i < 0 || list.count < i {
		return fmt.Errorf("Insert: index out of bounds: %d", i)
	}
	if i == 0 {
		list.head = &node{e, nil, list.head}
		if list.head.succ != nil {
			list.head.succ.pred = list.head
		}
		if list.tail == nil {
			list.tail = list.head
		}
		if list.cursorPtr != nil {
			list.cursorIdx++
		}
	} else if i == list.count {
		list.tail.succ = &node{e, list.tail, nil}
		list.tail = list.tail.succ
	} else {
		list.setCursor(i)
		newNode := &node{e, list.cursorPtr.pred, list.cursorPtr}
		list.cursorPtr.pred.succ = newNode
		list.cursorPtr.pred = newNode
		list.cursorPtr = newNode
	}
	list.count++
	return nil
}

// Delete removes and returns the element at location i.
// Precondition: 0 <= i < list.count.
// Precondition violation: delete nothing and return an error indication.
// Normal return: delete element i and return it and nil.
func (list *LinkedList) Delete(i int) (interface{}, error) {
	if i < 0 || list.count <= i {
		return nil, fmt.Errorf("Delete: index out of bounds: %d", i)
	}
	var result interface{}
	if list.count == 1 {
		result = list.head.item
		list.Clear()
		return result, nil
	}
	if i == 0 {
		result = list.head.item
		list.head.succ.pred = nil
		list.head = list.head.succ
		if list.cursorIdx > 0 {
			list.cursorIdx--
		} else {
			list.cursorPtr = list.head
		}
	} else if i == list.count-1 {
		result = list.tail.item
		list.tail = list.tail.pred
		list.tail.succ = nil
		if list.cursorIdx == list.count-1 {
			list.cursorIdx, list.cursorPtr = list.count-2, list.tail
		}
	} else {
		list.setCursor(i)
		result = list.cursorPtr.item
		list.cursorPtr.pred.succ = list.cursorPtr.succ
		list.cursorPtr.succ.pred = list.cursorPtr.pred
		list.cursorPtr = list.cursorPtr.succ
	}
	list.count--
	return result, nil
}

// Get returns element i without removing it.
// Precondition: 0 <= i < list.count.
// Precondition violation: return nil and an error indication.
// Normal return: element i and nil.
func (list *LinkedList) Get(i int) (interface{}, error) {
	if i < 0 || list.count <= i {
		return nil, fmt.Errorf("Get: index out of bounds: %d", i)
	}
	list.setCursor(i)
	return list.cursorPtr.item, nil
}

// Put changes element i.
// Precondition: 0 <= i < list.count.
// Precondition violation: change nothing and return an error indication.
// Normal return: change the value at i and return nil.
func (list *LinkedList) Put(i int, e interface{}) error {
	if i < 0 || list.count <= i {
		return fmt.Errorf("Put: index out of bounds: %d", i)
	}
	list.setCursor(i)
	list.cursorPtr.item = e
	return nil
}

// Index returns the location of element e. If e is not present,
// return 0 and false; otherwise return the location and true.
func (list *LinkedList) Index(e interface{}) (int, bool) {
	for index, ptr := 0, list.head; ptr != nil; index, ptr = index+1, ptr.succ {
		if ptr.item == e {
			return index, true
		}
	}
	return 0, false
}

// Slice makes a new list duplicating part of this list.
// Precondition: 0 <= i <= j <= list.count.
// Precondition violation: return nil and an error indication.
// Normal return: create a new list and fill it with the items between
// location i and j-1; return this new list and nil.
func (list *LinkedList) Slice(i, j int) (List, error) {
	result := new(LinkedList)
	if i < 0 || j < i || list.count < j {
		return result, fmt.Errorf("Slice: illegal indices: %d %d", i, j)
	}
	for srcIndex, dstIndex := i, 0; srcIndex < j; srcIndex, dstIndex = srcIndex+1, dstIndex+1 {
		v, _ := list.Get(srcIndex)
		result.Insert(dstIndex, v)
	}
	return result, nil
}

// Equal determines whether another List is identical to this one.
// Two Lists are identical if they are the same size and have the same
// elements in the same order.
// Precondition: element can be compared using ==.
// Precondition violation: panic.
// Normal return: true iff both lists have the same elements in the same order.
func (list *LinkedList) Equal(l List) bool {
	if list.count != l.Size() {
		return false
	}
	iter := l.NewIterator()
	v, ok := iter.Next()
	for ptr := list.head; ptr != nil; ptr = ptr.succ {
		if !ok || ptr.item != v {
			return false
		}
		v, ok = iter.Next()
	}
	return true
}

// abs returns the absolute value of an integer (used by setCursor).
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// setCursor moves the cursor to a location taking the shortest route.
// Precondition: i is in range.
func (list *LinkedList) setCursor(index int) {
	if list.cursorPtr == nil || index < abs(list.cursorIdx-index) {
		list.cursorIdx, list.cursorPtr = 0, list.head
	}
	if list.count-index < abs(list.cursorIdx-index) {
		list.cursorIdx, list.cursorPtr = list.count-1, list.tail
	}
	for index < list.cursorIdx { // go backwards
		list.cursorIdx, list.cursorPtr = list.cursorIdx-1, list.cursorPtr.pred
	}
	for list.cursorIdx < index { // go forwards
		list.cursorIdx, list.cursorPtr = list.cursorIdx+1, list.cursorPtr.succ
	}
}

// String makes a report on the data structure.
func (list *LinkedList) String() string {
	result := fmt.Sprintf("LinkedList instance:\nsize: %d\n", list.count)
	if list.cursorPtr != nil {
		result += fmt.Sprintf("Cursor: [%d] %v\n", list.cursorIdx, list.cursorPtr.item)
	}
	result += "Values:"
	for ptr := list.head; ptr != nil; ptr = ptr.succ {
		result += fmt.Sprintf(" %v ", ptr.item)
	}
	return result + "\n"
}
