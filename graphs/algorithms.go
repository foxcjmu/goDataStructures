// algorithms.go: This file contains implementations of various graph algorithms.
//
// author: C. Fox
// version: 11/2013

package graphs

import "containers"
import "errors"

// Perform a recursive depth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func DFS(g Graph, v0 int, visit func(Graph, int, int)) {
	isVisited := make([]bool, g.Vertices())
	var dfs func(int)
	dfs = func(v int) {
		iter, _ := g.NewIterator(v)
		for w, ok := iter.Next(); ok; w, ok = iter.Next() {
			if isVisited[w] {
				continue
			}
			visit(g, v, w)
			isVisited[w] = true
			dfs(w)
		}
	}
	visit(g, -1, v0)
	isVisited[v0] = true
	dfs(v0)
}

// The Edge struct stores edges for use in visiting vertices. The vertex
// visited is w and the edge that got it visited is v.
type Edge struct {
	v, w int // edge from source to taget
}

// Perform a stack-based depth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func StackDFS(g Graph, v0 int, visit func(Graph, int, int)) {
	isVisited := make([]bool, g.Vertices())
	stack := containers.NewLinkedStack()
	stack.Push(Edge{-1, v0})
	for edge, err := stack.Pop(); err == nil; edge, err = stack.Pop() {
		v, w := edge.(Edge).v, edge.(Edge).w
		if isVisited[w] {
			continue
		}
		visit(g, v, w)
		isVisited[w] = true
		iter, _ := g.NewIterator(w)
		for x, ok := iter.Next(); ok; x, ok = iter.Next() {
			if !isVisited[x] {
				stack.Push(Edge{w, x})
			}
		}
	}
}

// Perform a queue-based breadth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func BFS(g Graph, v0 int, visit func(Graph, int, int)) {
	isVisited := make([]bool, g.Vertices())
	queue := containers.NewLinkedQueue()
	queue.Enter(Edge{-1, v0})
	for edge, err := queue.Leave(); err == nil; edge, err = queue.Leave() {
		v, w := edge.(Edge).v, edge.(Edge).w
		if isVisited[w] {
			continue
		}
		visit(g, v, w)
		isVisited[w] = true
		iter, _ := g.NewIterator(w)
		for x, ok := iter.Next(); ok; x, ok = iter.Next() {
			if !isVisited[x] {
				queue.Enter(Edge{w, x})
			}
		}
	}
}

// Return true iff there is a path between v and w in g.
func IsPath(g Graph, v, w int) bool {
	if v < 0 || g.Vertices() <= v {
		return false
	}
	if w < 0 || g.Vertices() <= w {
		return false
	}
	isReached := false
	visit := func(g Graph, v1, v2 int) {
		if w == v2 {
			isReached = true
		}
	}
	DFS(g, v, visit)
	return isReached
}

// Return an int slice with the shortest path between v and w.
// Pre: IsPath(g,v,w)
// Pre violation: Return nil and an error
// Normal return: the path and nil
func ShortestPath(g Graph, v, w int) ([]int, error) {
	if !IsPath(g, v, w) {
		return nil, errors.New("The vertices are not connected")
	}
	toEdge := make([]int, g.Vertices())
	visit := func(g Graph, v1, v2 int) {
		toEdge[v2] = v1
	}
	BFS(g, w, visit)
	result := make([]int, 0, g.Vertices())
	x := v
	for x != w {
		result = append(result, x)
		x = toEdge[x]
	}
	result = append(result, x)
	return result, nil
}

// Return true iff a graph is connected (that is, there is a path between every pair of vertices).
func IsConnected(g Graph) bool {
	vertexCount := 0
	visit := func(g Graph, v1, v2 int) {
		vertexCount++
	}
	DFS(g, 0, visit)
	return vertexCount == g.Vertices()
}

// Return a new linked graph containing a spanning tree for g.
// Pre: g is connected.
// Pre Violation: return nil and false.
// Normal return: the spanning tree and true.
func SpanningTree(g Graph) (Graph, error) {
	if !IsConnected(g) {
		return nil, errors.New("Graph g is not connected")
	}
	result := NewLinkedGraph(g.Vertices())
	visit := func(g Graph, v1, v2 int) {
		result.AddEdge(v1, v2)
	}
	DFS(g, 0, visit)
	return result, nil
}

// Return the maximum number of edges from a vertex in a connected component
// of g containing vertex 0.
func MaxDegree(g Graph) int {
	result := 0
	visit := func(g Graph, v1, v2 int) {
		iter, _ := g.NewIterator(v2)
		count := 0
		for _, ok := iter.Next(); ok; _, ok = iter.Next() {
			count++
		}
		if result < count {
			result = count
		}
	}
	DFS(g, 0, visit)
	return result
}

// Return the number of vertices connected to a vertex v in the connected
// component of g containing vertex v.
func NumConnectedVertices(g Graph, v int) int {
	result := 0
	visit := func(g Graph, v1, v2 int) {
		if v1 != -1 {
			result++
		}
	}
	DFS(g, v, visit)
	return result
}
