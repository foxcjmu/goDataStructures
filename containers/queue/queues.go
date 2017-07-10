// queues.go -- implementation of the ArrayQueue and LinkedQueue parts of
// the containers/queue package
// author: C. Fox
// version: 1/2016
//
// The Queue interface is for all queues.
//
// queue implements four kinds of queues:
//  - ArrayQueue uses a slice to store elements
//  - LinkedQueue stores values in a singly linked list
//  - ArrayRandomizer stores entered values in a slice and
//    releases them in random order
//  - LinkedRandomizer stores values in a singly linked list
//    and releases them in random order
package queue

import (
	"errors"
	"fmt"

	"containers"
)

// Queue is the interface for queues in the container hierarchy.
type Queue interface {
	containers.Container         // include Size, Clear, and Empty
	Front() (interface{}, error) // return the front element of a non-empty queue
	Leave() (interface{}, error) // remove and return the front element of a non-empty queue
	Enter(e interface{})         // place a new element on at the rear of the queue
}

// ArrayQueue -----------------------------------------------------------------------
// A slice is used to store the data, and it expands as necessary if the queue
// becomes full. The front location is recorded, and the next open spot at the
// rear is at store[(frontIndex+count)%len(store)]. The built-in append operation is
// used to add a copy of elements at the end of the store because it will reallocate
// the underlying array automatically.
// Invariant: len(store) >= Size()

// ArrayQueue is a contiguous implementation of a queue.
type ArrayQueue struct {
	count      int           // how many elements are in the queue
	frontIndex int           // store[frontIndex] comes out next
	store      []interface{} // circular buffer for queue elements
}

// Size returns the number of elements inthe queue.
func (q *ArrayQueue) Size() int { return q.count }

// Clear makes the queue empty.
func (q *ArrayQueue) Clear() { q.count = 0 }

// Empty returns true iff the queue is empty.
func (q *ArrayQueue) Empty() bool { return q.count == 0 }

// Front returns the front element in the queue without removing it.
// Precondition: the queue is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: the front element and nil.
func (q *ArrayQueue) Front() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("Front: the queue cannot be empty")
	}
	return q.store[q.frontIndex], nil
}

// Leave removes and returns the front element on the queue.
// Precondition: the queue is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: the front element and nil.
func (q *ArrayQueue) Leave() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("Leave: the queue cannot be empty")
	}
	result := q.store[q.frontIndex]
	q.frontIndex++
	q.count--
	return result, nil
}

// Enter adds a new element to the rear of the queue.
func (q *ArrayQueue) Enter(e interface{}) {
	if len(q.store) == 0 {
		q.store = make([]interface{}, 4)
	}
	if q.count == len(q.store) {
		q.store = append(q.store, q.store...)
	}
	q.store[(q.frontIndex+q.count)%len(q.store)] = e
	q.count++
}

// String makes a report on the container.
func (q *ArrayQueue) String() string {
	return fmt.Sprintf("ArrayQueue instance:\nsize: %d\nfrontIndex: %d\nstore len: %d\nstore cap: %d\n"+
		"store: %v\n", q.count, q.frontIndex, len(q.store), cap(q.store), q.store)
}

// LinkedQueue ----------------------------------------------------------------
// A singly-linked list is used to store the values with the front node at the
// head of the list. A pointer to the tail of the list is kept in rearPtr. This
// makes both adding and removing elements from the queue efficient.
type node struct {
	item interface{}
	next *node
}

// LinkedQueue is a linked implementation of a queue.
type LinkedQueue struct {
	count    int   // how many elements are in the queue
	frontPtr *node // head of a singly-linked list of values
	rearPtr  *node // last element in the singly-linked list
}

// Size returns the number of elements in the queue.
func (q *LinkedQueue) Size() int { return q.count }

// Clear makes the linked queue empty.
func (q *LinkedQueue) Clear() {
	q.count = 0
	q.frontPtr, q.rearPtr = nil, nil
}

// Empty returns true iff the queue is empty.
func (q *LinkedQueue) Empty() bool { return q.count == 0 }

// Front returns the front element in the queue without removing it.
// Precondition: the queue is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: the front element and nil.
func (q *LinkedQueue) Front() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("Front: the queue cannot be empty")
	}
	return q.frontPtr.item, nil
}

// Leave removes and returns the front element on the queue.
// Precondition: the queue is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: the front element and nil.
func (q *LinkedQueue) Leave() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("Leave: the queue cannot be empty")
	}
	result := q.frontPtr.item
	q.frontPtr = q.frontPtr.next
	if q.frontPtr == nil {
		q.rearPtr = nil
	}
	q.count--
	return result, nil
}

// Enter adds a new element to the rear of the queue.
func (q *LinkedQueue) Enter(e interface{}) {
	if q.count == 0 {
		q.frontPtr = &node{e, nil}
		q.rearPtr = q.frontPtr
	} else {
		q.rearPtr.next = &node{e, nil}
		q.rearPtr = q.rearPtr.next
	}
	q.count++
}

// String makes a report on the container.
func (q *LinkedQueue) String() string {
	var result = fmt.Sprintf("LinkedQueue instance:\nsize: %d\ncontents:", q.count)
	for n := q.frontPtr; n != nil; n = n.next {
		result += fmt.Sprintf(" %v", n.item)
	}
	return result + "\n"
}
