// Test Stack interface and the ArrayStack and LinkedStack data structures.
// author: C. Fox
// version: 1/2016

package stack

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf // in case we need fmt for debugging

type tstruct struct { // for test values
	k int
	r float64
}

func TestStacks(t *testing.T) {
	var s Stack = new(ArrayStack)
	testStack(t, s)
	s = new(LinkedStack)
	testStack(t, s)
}

func testStack(t *testing.T, s Stack) {

	// make sure a new Stack is empty
	if !s.Empty() || 0 != s.Size() {
		t.Error("Stack should be empty and size should be 0 when new")
	}

	// make sure top and pop fail on an empty stack
	if v, err := s.Top(); err == nil {
		t.Errorf("Stack top operation should fail on an empty stack, instead returns %v", v)
	}
	if v, err := s.Pop(); err == nil {
		t.Errorf("Stack pop operation should fail on an empty stack, instead returns %v", v)
	}

	// push some data and check that everything works
	for i := 1; i <= 10; i++ {
		x := tstruct{i, float64(i) * 8.2}
		s.Push(x)
	}
	//fmt.Print(s)
	if s.Size() != 10 {
		t.Errorf("Stack push failure: stack should have 10 elements but has %v", s.Size())
	}
	for i := 10; 0 < i; i-- {
		if v, err := s.Top(); err == nil {
			if v != (tstruct{i, float64(i) * 8.2}) {
				t.Errorf("Stack Top error: value %v should be %v", v, tstruct{i, float64(i) * 8.2})
			}
		} else {
			t.Error("Stack error: Top operation failure when stack should not be empty")
		}
		if v, err := s.Pop(); err == nil {
			if v != (tstruct{i, float64(i) * 8.2}) {
				t.Errorf("Stack Pop error: value %v should be %v", v, tstruct{i, float64(i) * 8.2})
			}
		} else {
			t.Error("Stack error: Pop operation failure when stack should not be empty")
		}
	}
	if !s.Empty() || 0 != s.Size() {
		t.Error("Stack should be empty and size should be 0 when all data is popped")
	}

	// check that the clear operation works
	for i := float64(1); i <= 10; i++ {
		s.Push(i)
	}
	s.Clear()
	if !s.Empty() || 0 != s.Size() {
		t.Error("Stack should be empty and size should be 0 after clear is called")
	}
}
