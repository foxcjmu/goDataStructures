// Test Randomizer interface and the ArrayRandomizer and LinkedRandomizer data structures.
// author: C. Fox
// version: 11/2012

package queue

import (
	"fmt"
	"math"
	"testing"

	"containers/list"
)

var _ = fmt.Printf // in case we need fmt for debugging

func TestRandomizers(t *testing.T) {
	var q Randomizer = new(ArrayRandomizer)
	testRandomizer(t, q)
	q = new(LinkedRandomizer)
	testRandomizer(t, q)
}

func testRandomizer(t *testing.T, q Randomizer) {

	// make sure a new Randomizer is empty
	if !q.Empty() || 0 != q.Size() {
		t.Errorf("Randomizer should be empty and size should be 0 when new")
	}

	// make sure leave fails on an empty randomizer
	if v, err := q.Leave(); err == nil {
		t.Errorf("Randomizer leave operation should fail on an empty randomizer, instead returns %v", v)
	}

	// enter some data and check that everything works
	values := new(list.LinkedList)
	for i := 1; i <= 10; i++ {
		q.Enter(i)
		values.Insert(0, i)
	}
	if q.Size() != 10 {
		t.Errorf("Randomizer enter failure: randomizer should have 10 elements but has %v", q.Size())
	}
	for i := 1; i <= 10; i++ {
		if v, err := q.Leave(); err == nil {
			if v.(int) < 1 || 10 < v.(int) {
				t.Errorf("Randomizer error: value %v is out of range", v)
			}
			if index, success := values.Index(v); !success {
				t.Errorf("Duplicate value %v delivered by randomizer", v)
			} else {
				values.Delete(index)
			}
		} else {
			t.Errorf("Randomizer error: leave operation failure when randomizer should not be empty")
		}
	}
	if !q.Empty() || 0 != q.Size() {
		t.Errorf("Randomizer should be empty and size should be 0 when all data is removed")
	}

	// check that the clear operation works
	for i := 1; i <= 10; i++ {
		q.Enter(math.Sqrt(float64(i)))
	}
	q.Clear()
	if !q.Empty() || 0 != q.Size() {
		t.Errorf("Randomizer should be empty and size should be 0 after clear is called")
	}

}
