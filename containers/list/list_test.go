// Test List interface and the ArrayList and LinkedList data structure.
//
// author: C. Fox
// version: 5/2012

package list

import (
	"fmt"
	"testing"

	"containers"
)

var _ = fmt.Printf // in case we need fmt for debugging

func TestLists(t *testing.T) {
	testList(t, new(ArrayList), "ArrayList ")
	testList(t, new(LinkedList), "LinkedList ")
	testList(t, new(SinglyLinkedList), "SinglyLinkedList ")
}

func testList(t *testing.T, list List, name string) {
	// make sure a new List is empty
	if !list.Empty() || 0 != list.Size() {
		t.Error(name + "should be empty and size should be 0 when new")
	}

	// make sure get and delete fail on an empty list
	if v, err := list.Delete(0); err == nil {
		t.Errorf(name+"delete operation should fail on an empty list, instead returns %v", v)
	}
	if v, err := list.Get(0); err == nil {
		t.Errorf(name+"get operation should fail on an empty list, instead returns %v", v)
	}

	// add some some data to the list and check that it is there
	if err := list.Insert(0, 3); err != nil {
		t.Errorf(name+"insertion operation failed at index %v", 0)
	}
	if err := list.Insert(0, 1); err != nil {
		t.Errorf(name+"insertion operation failed at index %v", 0)
	}
	if err := list.Insert(2, 4); err != nil {
		t.Errorf(name+"insertion operation failed at index %v", 2)
	}
	if err := list.Insert(1, 2); err != nil {
		t.Errorf(name+"insertion operation failed at index %v", 1)
	}
	if err := list.Insert(0, 0); err != nil {
		t.Errorf(name+"insertion operation failed at index %v", 0)
	}
	if list.Size() != 5 {
		t.Errorf(name+"insertion failure: list should have 4 elements but has %v", list.Size())
	}
	//fmt.Print(list)
	for i := 0; i <= 4; i++ {
		if v, err := list.Get(i); err == nil {
			if v != i {
				t.Errorf(name+"error: value %v should be %v", v, i)
			}
		} else {
			t.Error(name + "error: get operation failure when list should not be empty")
		}
	}

	// try external iterators
	var iter containers.Iterator
	i := 0
	iter = list.NewIterator()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if v != i {
			t.Errorf("Iterator value should be %v but is %v", i, v)
		}
		i++
	}
	if i < 5 {
		t.Error("Iterator did not iterate all the way through the list")
	}
	if !iter.Done() {
		t.Error("Iterator should be done")
	}

	// try internal iterators
	var expectedValue int = 0
	var evf func(interface{}) = func(e interface{}) {
		if expectedValue != e {
			t.Errorf("Expected %v but got %v during internal iteration", expectedValue, e)
		}
		expectedValue++
	}
	list.Apply(evf)
	if expectedValue < 5 {
		t.Error("Internal iterator did not iterator all the way through the list")
	}

	// search the list
	if i, present := list.Index(5); present {
		t.Errorf("Index failure on non-existent element at location %v", i)
	}
	if list.Contains(5) {
		t.Error("Contains failure on non-existent element")
	}
	if i, present := list.Index(3); !present || i != 3 {
		t.Errorf("Index failure locating element at location 3; got %v", i)
	}
	if !list.Contains(3) {
		t.Error("Contains failure: should find 3")
	}

	// delete some data and insert some data and see that things are in order
	if v, err := list.Delete(2); err != nil {
		t.Error(name + "delete failed on non-empty list with legal index")
	} else if v != 2 {
		t.Errorf(name+"delete failed and deleted %v", v)
	}
	if err := list.Insert(2, 2); err != nil {
		t.Error(name + "insertion operation failed at index 2")
	}
	list.Get(0) // force cursor to 0
	if v, err := list.Delete(0); err != nil {
		t.Error(name + "delete failed on non-empty list with legal index")
	} else if v != 0 {
		t.Errorf(name+"delete failed and deleted %v", v)
	}
	if err := list.Insert(0, 0); err != nil {
		t.Error(name + "insertion operation failed at index 2")
	}
	list.Get(1) // force cursor to 1
	if v, err := list.Delete(0); err != nil {
		t.Error(name + "delete failed on non-empty list with legal index")
	} else if v != 0 {
		t.Errorf(name+"delete failed and deleted %v", v)
	}
	if err := list.Insert(0, 0); err != nil {
		t.Error(name + "second insertion operation failed at index 0")
	}
	list.Get(4) // force cursor to the end
	if v, err := list.Delete(4); err != nil {
		t.Error(name + "delete failed on non-empty list with legal index")
	} else if v != 4 {
		t.Errorf(name+"delete failed and deleted %v", v)
	}
	if err := list.Insert(4, 4); err != nil {
		t.Error(name + "insertion operation failed at index 4")
	}
	for i := 0; i < 5; i++ {
		if v, err := list.Get(i); err == nil {
			if v != i {
				t.Errorf(name+"error: value %v should be %v", v, i)
			}
		} else {
			t.Error(name + "error: get operation failure when list should not be empty")
		}
	}

	// check out of bounds errors
	if err := list.Insert(list.Size()+1, 6); err == nil {
		t.Error("Insertion at illegal location")
	}
	if err := list.Insert(-1, -1); err == nil {
		t.Error("Insertion at illegal location")
	}
	if _, err := list.Delete(list.Size()); err == nil {
		t.Error("Deletion at illegal location")
	}
	if _, err := list.Delete(-1); err == nil {
		t.Error("Deletion at illegal location")
	}
	if _, err := list.Get(list.Size()); err == nil {
		t.Error("Get at illegal location")
	}
	if _, err := list.Get(-1); err == nil {
		t.Error("Get at illegal location")
	}
	if err := list.Put(list.Size(), 8); err == nil {
		t.Error("Put at illegal location")
	}
	if err := list.Put(-1, 8); err == nil {
		t.Error("Put at illegal location")
	}
	if s, err := list.Slice(list.Size()+1, list.Size()+1); err == nil {
		t.Errorf("Slice at illegal location; list is %v", s)
	}
	if s, err := list.Slice(0, list.Size()+1); err == nil {
		t.Errorf("Slice at illegal location; list is %v", s)
	}
	if s, err := list.Slice(3, 2); err == nil {
		t.Errorf("Slice at illegal location; list is %v", s)
	}
	if s, err := list.Slice(-1, 2); err == nil {
		t.Errorf("Slice at illegal location; list is %v", s)
	} else {
		if s.Size() != 0 {
			t.Error("Slice error result is not empty")
		}
	}

	// use put to change everything
	for i := 0; i < list.Size(); i++ {
		if err := list.Put(i, list.Size()-i); err != nil {
			t.Errorf("Put failed when it should not on index %v", i)
		}
	}
	iter.Reset()
	i = list.Size()
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		if i != v {
			t.Errorf("Values out of order from put: expecting %v got %v", i, v)
		}
		i--
	}
	list.Put(0, 0)
	list.Put(3, 3)
	list.Put(1, 1)
	list.Put(4, 4)
	list.Put(2, 2)
	for i := 0; i < list.Size(); i++ {
		if v, err := list.Get(i); err != nil {
			t.Errorf(name+"get failed when it should not have on %v", i)
		} else if v != i {
			t.Errorf("Values out of order from random Put: expecting %v got %v", i, v)
		}
	}

	// make some sub-lists and check them
	if s, err := list.Slice(0, list.Size()); err != nil {
		t.Error("Slice failed when it should not")
	} else if !list.Equal(s) {
		t.Error("Slice different from parent list element")
	}
	if s, err := list.Slice(1, 4); err != nil {
		t.Error("Slice failed when it should not")
	} else {
		if s.Size() != 3 {
			t.Errorf("Slice has wrong size %v", s.Size())
		}
		for i := 0; i < s.Size(); i++ {
			if v, err := s.Get(i); err != nil {
				t.Error("Slice get failed when it should be err")
			} else {
				if i+1 != v {
					t.Errorf("Slice value %v should be %v", v, i+1)
				}
			}
		}
	}

	// delete from a list with only one element
	one, _ := list.Slice(0, 1)
	if _, err := one.Delete(0); err != nil {
		t.Error(name + "failed to delete from a list with one element")
	} else {
		if !one.Empty() {
			t.Error(name + "did not remove last element from list")
		}
	}

	// test lists for equality
	other := new(ArrayList)
	other.Insert(0, 0)
	if list.Equal(other) {
		t.Error(name + "should not be equal to a shorter list")
	}
	if other, err := list.Slice(0, list.Size()); err != nil {
		t.Error("Slice failed for entire list")
	} else {
		other.Put(1, 7)
		if list.Equal(other) {
			t.Error(name + "should not be equal to a list with different data")
		}
	}

	// test clear
	list.Clear()
	if !list.Empty() || 0 != list.Size() {
		t.Error(name + "should be empty and size should be 0 after clear is called")
	}
}
