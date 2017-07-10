// Test the Queue interface and the ArrayQueue and LinkedQueue data structures.
// author: C. Fox
// version: 5/2012

package queue

import (
	"fmt"
	"math"
	"testing"
)

var _ = fmt.Printf // in case we need fmt for debugging

func TestQueues(t *testing.T) {
	var q Queue = new(ArrayQueue)
	testQueue(t, q)
	q = new(LinkedQueue)
	testQueue(t, q)
}

func testQueue(t *testing.T, q Queue) {

	// make sure a new Queue is empty
	if !q.Empty() || 0 != q.Size() {
		t.Error("Queue should be empty and size should be 0 when new")
	}

	// make sure front and leave fail on an empty queue
	if v, err := q.Front(); err == nil {
		t.Errorf("Queue front operation should fail on an empty queue, instead returns %v", v)
	}
	if v, err := q.Leave(); err == nil {
		t.Errorf("Queue leave operation should fail on an empty queue, instead returns %v", v)
	}

	// enter some data and check that everything works
	for i := 1; i <= 10; i++ {
		q.Enter(i)
	}
	//fmt.Print(q)
	if q.Size() != 10 {
		t.Errorf("Queue Enter failure: queue should have 10 elements but has %v", q.Size())
	}
	for i := 1; i <= 10; i++ {
		if v, err := q.Front(); err == nil {
			if v != i {
				t.Errorf("Queue Front error: value", v, "should be", i)
			}
		} else {
			t.Error("Queue error: Front operation failure when queue should not be empty")
		}
		if v, err := q.Leave(); err == nil {
			if v != i {
				t.Errorf("Queue Leave error: value %v should be %v", v, i)
			}
		} else {
			t.Error("Queue error: Leave operation failure when queue should not be empty")
		}
	}
	if !q.Empty() || 0 != q.Size() {
		t.Error("Queue should be empty and size should be 0 when all data is removed")
	}

	// check that the clear operation works
	for i := 1; i <= 10; i++ {
		q.Enter(math.Sqrt(float64(i)))
	}
	q.Clear()
	if !q.Empty() || 0 != q.Size() {
		t.Error("Queue should be empty and size should be 0 after Clear is called")
	}
}
