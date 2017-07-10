// randomizers.go -- implements two kinds of random exit queues:
// - ArrayRandomizer uses a slice to store elements
// - LinkedRandomizer stores values in a singly linked list

// author: C. Fox
// version: 1/2016

package queue

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"containers"
)

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
}

// Randomizer is the interface for randomizers in the containers hierarchy.
type Randomizer interface {
	containers.Container         // include Size, Clear, and Empty
	Leave() (interface{}, error) // remove and return a random element from a non-empty randomizer
	Enter(e interface{})         // place a new element on at the rear of the randomizer
}

// ArrayRandomizer ------------------------------------------------------------
// A slice is used to store the data, and it expands as necessary if the slice
// becomes full. New items are placed from index 0 upwards in the slice. When an
// item is removed, a location in the occupied portion of the slice is chosen at
// random, that element is returned, and the element at the end of the occupied
// poriton of the slice is copied into the returned item's slot.
// Invariant: len(store) >= Size()

// ArrayRandomizer is a contiguous implementation of a randomizer.
type ArrayRandomizer struct {
	count int           // how many elements are in the queue
	store []interface{} // slice for randomizer elements
}

// Size returns the number of elements in the randomizer.
func (r *ArrayRandomizer) Size() int { return r.count }

// Clear makes the randomizer empty.
func (r *ArrayRandomizer) Clear() { r.count = 0 }

// Empty returns true iff the randomizer empty.
func (r *ArrayRandomizer) Empty() bool { return r.count == 0 }

// Leave removes and returns a random element from the randomizer.
// Precondition: the randomizer is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: a random element and nil.
func (r *ArrayRandomizer) Leave() (interface{}, error) {
	if r.count == 0 {
		return nil, errors.New("Leave: the randomizer cannot be empty")
	}
	index := rand.Intn(r.count)
	result := r.store[index]
	r.count--
	r.store[index] = r.store[r.count]
	return result, nil
}

// Enter adds a new element to the randomizer.
func (r *ArrayRandomizer) Enter(e interface{}) {
	if r.count == len(r.store) {
		r.store = append(r.store, e)
	} else {
		r.store[r.count] = e
	}
	r.count++
}

// String makes a report on the container.
func (r *ArrayRandomizer) String() string {
	return fmt.Sprintf("ArrayRandomizer instance:\nsize: %d\nstore len: %d\nstore cap: %d\n"+
		"store: %v\n", r.count, len(r.store), cap(r.store), r.store)
}

// LinkedRandomizer ----------------------------------------------------------
// A singly-linked list is used to store the values with the front node at the
// head of the list. A walk a random distance from the head is used to find
// the element removed from the randomizer.

// LinkedRandomizer is a linked implementation of a randomizer.
type LinkedRandomizer struct {
	count   int   // how many elements are stored in the randomizer
	headPtr *node // head of a singly-linked list of values
}

// Size returns the number of elements in the randomizer.
func (r *LinkedRandomizer) Size() int { return r.count }

// Clear makes the linked randomizer empty.
func (r *LinkedRandomizer) Clear() {
	r.count = 0
	r.headPtr = nil
}

// Empty returns true iff the randomizer is empty.
func (r *LinkedRandomizer) Empty() bool { return r.count == 0 }

// Leave removes and returns a random element from the randomizer.
// Precondition: the randomizer is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: a random element and nil.
func (r *LinkedRandomizer) Leave() (interface{}, error) {
	if r.count == 0 {
		return nil, errors.New("Leave: the randomizer cannot be empty")
	}
	var result interface{}
	index := rand.Intn(r.count)
	if index == 0 {
		result = r.headPtr.item
		r.headPtr = r.headPtr.next
	} else {
		ptr := r.headPtr
		for i := 0; i < index-1; i++ {
			ptr = ptr.next
		}
		result = ptr.next.item
		ptr.next = ptr.next.next
	}
	r.count--
	return result, nil
}

// Enter adds a new element to the randomizer.
func (r *LinkedRandomizer) Enter(e interface{}) {
	r.headPtr = &node{e, r.headPtr}
	r.count++
}

// String makes a report on the contiainer.
func (r *LinkedRandomizer) String() string {
	var result = fmt.Sprintf("LinkedRandomizer instance:\nsize: %d\ncontents:", r.count)
	for ptr := r.headPtr; ptr != nil; ptr = ptr.next {
		result += fmt.Sprintf(" %v", ptr.item)
	}
	return result + "\n"
}
