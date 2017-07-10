// slist.go: Implementation of lists using singly-linked lists.
// author: C. Fox
// version: 10/2012

package list

import (
	"containers"
	"fmt"
)

// A singly-linked list with a cursor is used to store the values. A pointer is kept
// to the head of the list. The cursor sometimes provides faster access in long lists.
// The cursor is undefined if cursorPtr is nil; otherwise cursorPtr point to the
// element with index cursorIdx.
// snode is used for the singly-linked list
type snode struct {
	item interface{}
	next *snode
}

// SinglyLinkedList is a singly-linked implementation of a list.
type SinglyLinkedList struct {
	count     int    // how many elements are in the list
	head      *snode // start of a singly linked list of values
	cursorPtr *snode // snode where the cursor rests
	cursorIdx int    // index where the cursor rests
}

// sLinkedListIterator is the data structure for a SinglyLinkedList external iterator.
type sLinkedListIterator struct {
	list    *SinglyLinkedList // the list that is iterated over
	current *snode            // where we are in the list
}

// Reset prepares an iterator to traverse its associated Collection.
func (iter *sLinkedListIterator) Reset() {
	iter.current = iter.list.head
}

// Done is true iff the iterator has traversed its associated Collection.
func (iter *sLinkedListIterator) Done() bool {
	return iter.current == nil
}

// Next returns an interface{} and an indication of whether iteration is complete.
// Precondition: Iteration is not complete.
// Precondition violation: return nil and false.
// Normal return: the next element in the iteration and true.
func (iter *sLinkedListIterator) Next() (interface{}, bool) {
	if iter.current == nil {
		return nil, false
	}
	result := iter.current.item
	iter.current = iter.current.next
	return result, true
}

// Size indicates how many elements are in the list.
func (list *SinglyLinkedList) Size() int { return list.count }

// Clear removes all elements from the list.
func (list *SinglyLinkedList) Clear() {
	list.count = 0
	list.head = nil
	list.cursorPtr, list.cursorIdx = nil, 0
}

// Empty returns true iff the list has no elements.
func (list *SinglyLinkedList) Empty() bool { return list.count == 0 }

// NewIterator creates and returns an external iterator for a SinglyLinkedList.
func (list *SinglyLinkedList) NewIterator() containers.Iterator {
	result := new(sLinkedListIterator)
	result.list = list
	result.current = list.head
	return result
}

// Apply calls function f on every element in the Collection.
func (list *SinglyLinkedList) Apply(f func(interface{})) {
	for ptr := list.head; ptr != nil; ptr = ptr.next {
		f(ptr.item)
	}
}

// Contains returns true iff element e is in the Collection.
func (list *SinglyLinkedList) Contains(e interface{}) bool {
	for ptr := list.head; ptr != nil; ptr = ptr.next {
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
func (list *SinglyLinkedList) Insert(i int, e interface{}) error {
	if i < 0 || list.count < i {
		return fmt.Errorf("Insert: index out of bounds: %d", i)
	}
	if i == 0 {
		list.head = &snode{e, list.head}
		if list.cursorPtr != nil {
			list.cursorIdx++
		}
	} else {
		list.setCursor(i - 1)
		list.cursorPtr.next = &snode{e, list.cursorPtr.next}
	}
	list.count++
	return nil
}

// Delete removes and returns the element at location i.
// Precondition: 0 <= i < list.count.
// Precondition violation: delete nothing and return an error indication.
// Normal return: delete element i and return it and nil.
func (list *SinglyLinkedList) Delete(i int) (interface{}, error) {
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
		list.head = list.head.next
		if list.cursorIdx > 0 {
			list.cursorIdx--
		} else {
			list.cursorPtr = list.head
		}
	} else {
		list.setCursor(i - 1)
		result = list.cursorPtr.next.item
		list.cursorPtr.next = list.cursorPtr.next.next
	}
	list.count--
	return result, nil
}

// Get returns element i without removing it.
// Precondition: 0 <= i < list.count.
// Precondition violation: return nil and an error indication.
// Normal return: element i and nil.
func (list *SinglyLinkedList) Get(i int) (interface{}, error) {
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
func (list *SinglyLinkedList) Put(i int, e interface{}) error {
	if i < 0 || list.count <= i {
		return fmt.Errorf("Put: index out of bounds: %d", i)
	}
	list.setCursor(i)
	list.cursorPtr.item = e
	return nil
}

// Index returns the location of element e. If e is not present,
// return 0 and false; otherwise return the location and true.
func (list *SinglyLinkedList) Index(e interface{}) (int, bool) {
	for index, ptr := 0, list.head; ptr != nil; index, ptr = index+1, ptr.next {
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
func (list *SinglyLinkedList) Slice(i, j int) (List, error) {
	result := new(SinglyLinkedList)
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
func (list *SinglyLinkedList) Equal(l List) bool {
	if list.count != l.Size() {
		return false
	}
	iter := l.NewIterator()
	v, ok := iter.Next()
	for ptr := list.head; ptr != nil; ptr = ptr.next {
		if !ok || ptr.item != v {
			return false
		}
		v, ok = iter.Next()
	}
	return true
}

// setCursor moves the cursor to a location taking the shortest route.
// Precondition: i is in range.
func (list *SinglyLinkedList) setCursor(index int) {
	if list.cursorPtr == nil || index < list.cursorIdx {
		list.cursorPtr, list.cursorIdx = list.head, 0
	}
	for list.cursorIdx < index {
		list.cursorPtr = list.cursorPtr.next
		list.cursorIdx++
	}
}

// String makes a report on the container.
func (list *SinglyLinkedList) String() string {
	result := fmt.Sprintf("SinglyLinkedList instance:\nsize: %d\n", list.count)
	if list.cursorPtr != nil {
		result += fmt.Sprintf("Cursor: [%d] %v\n", list.cursorIdx, list.cursorPtr.item)
	}
	result += "Values:"
	for ptr := list.head; ptr != nil; ptr = ptr.next {
		result += fmt.Sprintf(" %v ", ptr.item)
	}
	return result + "\n"
}
