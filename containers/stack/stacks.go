// stacks.go -- implements the containers/stack package
// author: C. Fox
// version: 1/2016
//
// The Stack interface is used to declare any kind of stack.
//
// stack provides two kinds of stack containers:
//  - ArrayStack uses a slice to store elements
//  - LinkedStack stores values in a singly linked list
package stack

import (
	"containers"
	"errors"
	"fmt"
)

// Stack is the interface for stacks in the containers hierarchy.
type Stack interface {
	containers.Container       // include Size, Clear, and Empty
	Top() (interface{}, error) // return the top element of a non-empty stack
	Pop() (interface{}, error) // remove and return top element of a non-empty stack
	Push(e interface{})        // place a new element on the top of the stack
}

// ArrayStack ----------------------------------------------------------------
// A slice is used to store the data and the slice is always exactly as big as
// the stack, so the last element in the slice is the top of the stack.
// Invariant: len(store) == Size()

// ArrayStack is a contiguous implementation of a stack using slices.
type ArrayStack struct {
	store []interface{} // top is always at store[len(store)-1]
}

// Size returns the number of elements stored in the stack.
func (s *ArrayStack) Size() int { return len(s.store) }

// Empty returns true iff the stack is empty.
func (s *ArrayStack) Empty() bool { return len(s.store) == 0 }

// Clear removes all the items from the stack.
func (s *ArrayStack) Clear() { s.store = make([]interface{}, 0, 10) }

// Top returns the top value on the stack without removing it.
// Pre: the stack is not empty.
// Pre violation: return nil and an error indication.
// Normal return: return the top element (which is not removed) and nil.
func (s *ArrayStack) Top() (interface{}, error) {
	if len(s.store) == 0 {
		return nil, errors.New("Top: stack cannot be empty")
	}
	return s.store[len(s.store)-1], nil
}

// Pop removes and returns the top element on the stack.
// Precondition: the stack is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: return the top element (which is removed) and nil.
func (s *ArrayStack) Pop() (interface{}, error) {
	if len(s.store) == 0 {
		return nil, errors.New("Pop: the stack cannot be empty")
	}
	result := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return result, nil
}

// Push adds a new element to the top of the stack.
func (s *ArrayStack) Push(e interface{}) { s.store = append(s.store, e) }

// String makes a report on the container.
func (s *ArrayStack) String() string {
	return fmt.Sprintf("ArrayStack instance:\nstore len: %d\nstore cap: %d\nstore: %v\n",
		len(s.store), cap(s.store), s.store)
}

// LinkedStack --------------------------------------------------------------
// A singly-linked list is used to store the values in a linked stack
// with the top node at the head of the list.
// Invariant topPtr == nil iff Size() == 0
// node is a singly linked list node used for this purpose.
type node struct {
	item interface{}
	next *node
}

// LinkedStack is a linked implementation of a stack.
type LinkedStack struct {
	count  int   // how many elements are present
	topPtr *node // singly-linked list of values
}

// Size returns the number of elements in the stack.
func (s *LinkedStack) Size() int { return s.count }

// Empty returns true iff the stack is empty.
func (s *LinkedStack) Empty() bool { return s.count == 0 }

// Clear makes the stack empty.
func (s *LinkedStack) Clear() {
	s.count = 0
	s.topPtr = nil
}

// Top returns the top value on the stack without removing it.
// Precondition: the stack is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: return the top element (which is not removed) and nil.
func (s *LinkedStack) Top() (interface{}, error) {
	if s.count == 0 {
		return nil, errors.New("Top: the stack cannot be empty")
	}
	return s.topPtr.item, nil
}

// Pop removes and returns the top element on the stack.
// Precondition: the stack is not empty.
// Precondition violation: return nil and an error indication.
// Normal return: return the top element (which is removed) and nil.
func (s *LinkedStack) Pop() (interface{}, error) {
	if s.count == 0 {
		return nil, errors.New("Pop: the stack cannot be empty")
	}
	result := s.topPtr.item
	s.topPtr = s.topPtr.next
	s.count--
	return result, nil
}

// Push adds a new element to the top of the stack.
func (s *LinkedStack) Push(e interface{}) {
	s.topPtr = &node{e, s.topPtr}
	s.count++
}

// String makes a report on the container.
func (s *LinkedStack) String() string {
	var result = fmt.Sprintf("LinkedStack instance:\nsize: %d\ncontents:", s.count)
	for n := s.topPtr; n != nil; n = n.next {
		result += fmt.Sprintf(" %v", n.item)
	}
	return result + "\n"
}
