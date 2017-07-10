// Test Graph interface and the ArrayGraph and LinkedGraph data structures.
// author: C. Fox
// version: 8/2012

package graphs

//import "fmt"
import "testing"

func TestGraphs(t *testing.T) {
	testGraph(t, "ArrayGraph", NewArrayGraph(20))
	testGraph(t, "LinkedGraph", NewLinkedGraph(20))
}

func testGraph(t *testing.T, name string, g Graph) {

	// make sure a new Graph is the right size
	if g.Vertices() != 20 || 0 != g.Edges() {
		t.Errorf(name+" should have 20 vertices and no edges but has %v verties and %v edges", g.Vertices(), g.Edges())
	}

	// add some illegal edges
	if err := g.AddEdge(2, 2); err == nil {
		t.Errorf(name + ": Illegal edge 2-2 detected")
	}
	if err := g.AddEdge(2, 200); err == nil {
		t.Errorf(name + ": Illegal edge 2-200 detected")
	}
	if err := g.AddEdge(200, 2); err == nil {
		t.Errorf(name + ": Illegal edge 200-2 detected")
	}

	// add some legal edges
	a := []int{0, 4, 7, 11, 15, 19}
	for _, w := range a {
		if err := g.AddEdge(2, w); err != nil {
			t.Errorf(name+": Legal edge 2-%v flagged", w)
		}
	}
	if g.Edges() != len(a) {
		t.Errorf(name+": Edge count should be %v but is %v", len(a), g.Edges())
	}
	for _, w := range a {
		if !g.IsEdge(2, w) {
			t.Errorf(name+": Edge 2-%v missing", w)
		}
	}

	// test vertex iteration
	if _, err := g.NewIterator(-1); err == nil {
		t.Errorf(name + ": Failed to detect illegal vertex -1")
	}
	if name == "ArrayGraph" {
		i := 0
		iter, _ := g.NewIterator(2)
		for v, ok := iter.Next(); ok; v, ok = iter.Next() {
			if v != a[i] {
				t.Errorf(name+": Iterator value should be %v but is %v", a[i], v)
			}
			i++
		}
		if i != len(a) {
			t.Errorf(name+": Iteration ended prematurely after %v values", i+1)
		}
	} else {
		i := len(a) - 1
		iter, _ := g.NewIterator(2)
		for v, ok := iter.Next(); ok; v, ok = iter.Next() {
			if v != a[i] {
				t.Errorf(name+": Iterator value should be %v but is %v", a[i], v)
			}
			i--
		}
		if i != -1 {
			t.Errorf(name+": Iteration ended prematurely after %v values", i+1)
		}
	}
}
