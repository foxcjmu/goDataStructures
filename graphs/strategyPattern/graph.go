// graph.go: This file contains an implementation of unidrected graphs for the
// graphs package. In particular, it includes the Graph and VertexIterator interfaces,
// the unidirectedGraph type, and implementations of all the operations in the Graph
// interface, which is its method set.
//
// author: C. Fox
// version: 8/2012

// Package graphs implements basic undirected graphs using both the adjacency
// matrix and adjacency list representations. A graph provides methods for
// determing paths, connectedness, and spanning trees.

package graphs

import "containers" // use a linked list in the linked graph representation
import "errors"     // for illegal vertices and like errors
import "fmt"        // for the String function

// Graph is the interface for undirected graphs.
type Graph interface {
	// basic graph builders and accessors
	Edges() int                          // return the number of edges in the graph
	Vertices() int                       // return the number of vertices in the graph
	AddEdge(v, w int) error              // add an edge between vertices v and w
	IsEdge(v, w int) bool                // true iff there is an edge between v and w
	NewIterator(v int) (Iterator, error) // make an iterator over edges adjacent to v

	// derived graph algorithms
	DFS(v int, f func(Graph, int, int))      // resursive depth-first search from vertex v
	StackDFS(v int, f func(Graph, int, int)) // stack-based depth-first search from vertex v
	BFS(v int, f func(Graph, int, int))      // breadth-first search from vertex v
	IsPath(v, w int) bool                    // true ff there is a path between v and w
	ShortestPath(v, w int) ([]int, error)    // the vertices on the shortest path from v to w
	IsConnected() bool                       // true iff the graph is connected
	SpanningTree() (Graph, error)            // construct a spanning tree for this graph
}

// Iterator is the interface for external iterators over vertices
type Iterator interface {
	Reset()            // prepare for another iteration
	IsDone() bool      // return true iff this iterator is finished
	Next() (int, bool) // return the next vertex and ok indication
}

// graphRepresentation is the interface for the methods essential for an internal
// representations of a graph: this is the method set for graph representations.
type graphRepresentation interface {
	edges() int                          // return the number of edges in the graph
	vertices() int                       // return the number of vertices in the graph
	addEdge(v, w int) error              // add an edge between vertices v and w
	isEdge(v, w int) bool                // true iff there is an edge between v and w
	newIterator(v int) (Iterator, error) // make an iterator over edges adjacent to v
}

/////////////////////////////////////////////////////////////////////////////
// undirectedGraph holds a graph representation and has the Graph method set.
type undirectedGraph struct {
	representation graphRepresentation // adjacencyMatrix or adjacencyLists
}

// NewArrayGraph returns a pointer to a graph represented by an adjacency matrix.
// Pre: n > 0
// Pre violation: return a graph with 1 vertex.
// Normal return: return a graph with n vertices.
func NewArrayGraph(n int) *undirectedGraph {
	result := new(undirectedGraph)
	matrix := new(adjacencyMatrix)
	if n < 0 {
		n = 1
	}
	matrix.adjacent = make([][]bool, n)
	for i := 0; i < n; i++ {
		matrix.adjacent[i] = make([]bool, n)
	}
	result.representation = matrix
	return result
}

// NewLinkedGraph returns a pointer to a graph represented by adjacency lists.
// Pre: n > 0
// Pre violation: return a graph with 1 vertex.
// Normal return: return a graph with n vertices.
func NewLinkedGraph(n int) *undirectedGraph {
	result := new(undirectedGraph)
	lists := new(adjacencyLists)
	if n < 0 {
		n = 1
	}
	lists.adjacent = make([]containers.List, n)
	for i := 0; i < n; i++ {
		lists.adjacent[i] = containers.NewLinkedList()
	}
	result.representation = lists
	return result
}

// Edges returns the number of edges in the graph g.
func (g *undirectedGraph) Edges() int {
	return g.representation.edges()
}

// Vertices returns the number of vertices in the graph g.
func (g *undirectedGraph) Vertices() int {
	return g.representation.vertices()
}

// AddEdge puts a new edge in graph g between v and w; it does nothing if the
// edge already exists.
// Pre: v != w and 0 <= v, w < g.vertices()
// Pre violation: return an error indication.
// Normal return: add the edge and return nil.
func (g *undirectedGraph) AddEdge(v, w int) error {
	return g.representation.addEdge(v, w)
}

// IsEdge determines whether the graph g contains edge {v,w}.
func (g *undirectedGraph) IsEdge(v, w int) bool {
	return g.representation.isEdge(v, w)
}

// newIterator returns an iterator over the vertices adjacent to v.
// Pre: 0 <= v < g.Vertices()
// Pre violation: return nil and an error indication.
// Normal return: return a new iterator and nil.
func (g *undirectedGraph) NewIterator(v int) (Iterator, error) {
	return g.representation.newIterator(v)
}

// String returns a string representation of a graph.
func (g *undirectedGraph) String() string {
	result := ""
	for i := 0; i < g.Vertices(); i++ {
		result += fmt.Sprintf("%d:", i)
		iter, _ := g.NewIterator(i)
		for v, ok := iter.Next(); ok; v, ok = iter.Next() {
			result += fmt.Sprintf(" %d", v)
		}
		result += "\n"
	}
	return result
}

// DFS performs a recursive depth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func (g *undirectedGraph) DFS(v0 int, visit func(Graph, int, int)) {
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

// StackDFS performs a stack-based depth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func (g *undirectedGraph) StackDFS(v0 int, visit func(Graph, int, int)) {
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

// BFS performs a queue-based breadth-first search of g starting at v0 and
// applying the visit function to every vertex as it is visited.
// Pre: v0 is in g
// Pre violation: panic
// Normal return: all vertices in g connected to v0 are visited once
func (g *undirectedGraph) BFS(v0 int, visit func(Graph, int, int)) {
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

// IsPath returns true iff there is a path between v and w in g
func (g *undirectedGraph) IsPath(v, w int) bool {
	if v < 0 || g.Vertices() < v {
		return false
	}
	if w < 0 || g.Vertices() < w {
		return false
	}
	isReached := false
	visit := func(g Graph, v1, v2 int) {
		if w == v2 {
			isReached = true
		}
	}
	g.DFS(v, visit)
	return isReached
}

// ShortestPath returns an int slice with the shortes path between v and w.
// Pre: IsPath(g,v,w)
// Pre violation: Return nil and an error
// Normal return: the path and nil
func (g *undirectedGraph) ShortestPath(v, w int) ([]int, error) {
	if !g.IsPath(v, w) {
		return nil, errors.New("The vertices are not connected")
	}
	toEdge := make([]int, g.Vertices())
	visit := func(g Graph, v1, v2 int) {
		toEdge[v2] = v1
	}
	g.BFS(w, visit)
	result := make([]int, 0, g.Vertices())
	x := v
	for x != w {
		result = append(result, x)
		x = toEdge[x]
	}
	result = append(result, x)
	return result, nil
}

// IsConnected returns true iff a graph is connected (that is, there is a
// path between every pair of vertices).
func (g *undirectedGraph) IsConnected() bool {
	isVisited := make([]bool, g.Vertices())
	visit := func(g Graph, v1, v2 int) {
		isVisited[v2] = true
	}
	g.DFS(0, visit)
	result := true
	for i := 0; i < len(isVisited); i++ {
		result = result && isVisited[i]
	}
	return result
}

// SpanningTree returns a new linked graph containing a spanning tree for g.
// Pre: g is connected.
// Pre Violation: return nil and false.
// Normal return: the spanning tree and true.
func (g *undirectedGraph) SpanningTree() (Graph, error) {
	if !g.IsConnected() {
		return nil, errors.New("Graph g is not connected")
	}
	result := NewLinkedGraph(g.Vertices())
	visit := func(g Graph, v1, v2 int) {
		result.AddEdge(v1, v2)
	}
	g.DFS(0, visit)
	return result, nil
}
