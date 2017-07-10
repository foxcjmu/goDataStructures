// Test graph algorithms.
// author: C. Fox
// version: 8/2012

package graphs

import "fmt"
import "testing"
import "math/rand"
import "time"

func TestAlgorithms(t *testing.T) {
	testSearch(t, "ArrayGraph", NewArrayGraph(20))
	testSearch(t, "LinkedGraph", NewLinkedGraph(20))
	testAlgorithms(t, "ArrayGraph", NewArrayGraph(10))
	testAlgorithms(t, "LinkedGraph", NewLinkedGraph(10))
	s := fmt.Sprintf("Be quiet about fmt imported but not used already")
	s += "!"
}

func testSearch(t *testing.T, name string, g Graph) {
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := 0; i < 100; i++ {
		g.AddEdge(rand.Intn(20), rand.Intn(20))
	}
	/* Add this code to guarantee that the graph is connected. Usually it will be by chance.
	for i := 0; i < g.Vertices(); i++ {
		g.AddEdge(i,(i+12)%20)
	}
	*/
	counts := make([]int, g.Vertices())
	f := func(g Graph, v, w int) {
		counts[w]++
	}
	g.DFS(0, f)
	for i := 0; i < g.Vertices(); i++ {
		if counts[i] == 0 {
			t.Errorf(name+": DFS did not visit vertex %v", i)
		} else if counts[i] > 1 {
			t.Errorf(name+": DFS visited vertex %v %v times", i, counts[i])
		}
	}

	counts = make([]int, g.Vertices())
	g.StackDFS(0, f)
	for i := 0; i < g.Vertices(); i++ {
		if counts[i] == 0 {
			t.Errorf(name+": StackDFS did not visit vertex %v", i)
		} else if counts[i] > 1 {
			t.Errorf(name+": StackDFS visited vertex %v %v times", i, counts[i])
		}
	}

	counts = make([]int, g.Vertices())
	g.BFS(0, f)
	for i := 0; i < g.Vertices(); i++ {
		if counts[i] == 0 {
			t.Errorf(name+": BFS did not visit vertex %v", i)
		} else if counts[i] > 1 {
			t.Errorf(name+": BFS visited vertex %v %v times", i, counts[i])
		}
	}
}

func testAlgorithms(t *testing.T, name string, g Graph) {
	g.AddEdge(0, 1)
	g.AddEdge(0, 3)
	g.AddEdge(2, 3)
	g.AddEdge(4, 3)
	g.AddEdge(3, 5)
	g.AddEdge(4, 5)
	g.AddEdge(6, 7)
	g.AddEdge(8, 7)
	g.AddEdge(8, 9)
	if g.IsPath(3, 6) || g.IsPath(6, 3) {
		t.Errorf(name + ": there is no path from 3 to 6")
	}
	if !g.IsPath(5, 1) && !g.IsPath(1, 5) {
		t.Errorf(name + ": there is a path from 5 to 1")
	}
	if g.IsConnected() {
		t.Errorf(name + ": graph is not connected, but IsConnected says it is")
	}
	g.AddEdge(2, 8)
	g.AddEdge(3, 6)
	if !g.IsConnected() {
		t.Errorf(name + ": graph is connected, but IsConnected says it is not")
	}

	// test shortest path generation
	path, _ := g.ShortestPath(0, 9)
	if !samePath(path, []int{0, 3, 2, 8, 9}) {
		t.Errorf(name + ": failed to find the shortes path from 0 to 9")
	}
	path, _ = g.ShortestPath(5, 7)
	if !samePath(path, []int{5, 3, 6, 7}) {
		t.Errorf(name + ": failed to find the shortes path from 5 to 7")
	}

	// test spanning tree generation
	h, err := g.SpanningTree()
	if err != nil {
		t.Errorf(name + ": spanning tree generation failed")
	}
	if g.Vertices() != h.Vertices() || !h.IsConnected() || h.Edges() != h.Vertices()-1 {
		t.Errorf(name + ": spanning tree generation failed with a bad spanning tree")
	}
}

func samePath(p, q []int) bool {
	if len(p) != len(q) {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}
