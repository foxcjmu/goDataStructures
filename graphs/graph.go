// graph.go: This file contains the fundamental declarations for the graphs
// package. In particular, it includes the Graph and Iterator interfaces,
// and the arrayGraph and linkedGraph types as receivers that implement the
// adjacency matrix and adjacency list representations of undirected graphs,
// respectively.
//
// author: C. Fox
// version: 11/2013

// Package graphs implements basic undirected graphs using both the adjacency
// matrix and adjacency list representations.

package graphs

import "containers" // use a linked list in the linked graph representation
import "errors"     // for illegal vertices and like errors
import "fmt"        // for the String function

// Graph is the interface for undirected graphs.
type Graph interface {
	Edges() int                          // return the number of items in the container
	Vertices() int                       // return the number of items in the container
	AddEdge(v, w int) error              // add an edge between vertices v and w
	IsEdge(v, w int) bool                // true iff there is an edge between v and w
	NewIterator(v int) (Iterator, error) // make an iterator over edges adjacent to v
}

// Iterator is the interface for all external iterators over vertices
type Iterator interface {
	Reset()            // prepare for another iteration
	IsDone() bool      // return true iff this iterator is finished
	Next() (int, bool) // return the next vertex and ok indication
}

///////////////////////////////////////////////////////////////////////////////////////
// arrayGraph is the data structure for the adjacency matrix representation of a graph.
type arrayGraph struct {
	numEdges int      // in the graph
	adjacent [][]bool // true at [v][w] iff {v,w} is an edge
}

// NewArrayGraph returns a pointer to a graph represented using an
// adjacency matrix.
// Pre: n > 0
// Pre violation: return a graph with 1 vertex.
// Normal return: return a graph with n vertices.
func NewArrayGraph(n int) *arrayGraph {
	result := new(arrayGraph)
	if n < 0 {
		n = 1
	}
	result.adjacent = make([][]bool, n)
	for i := 0; i < n; i++ {
		result.adjacent[i] = make([]bool, n)
	}
	return result
}

// Edges return the number of edges in the receiver graph.
func (g *arrayGraph) Edges() int {
	return g.numEdges
}

// Vertices return the number of vertices in the receiver graph.
func (g *arrayGraph) Vertices() int {
	return len(g.adjacent)
}

// AddEdge puts a new edge in the receiver graph; it does nothing if
// the edge is already there.
// Pre: v and w are in the graph.
// Pre violation: return false.
// Normal return: add the edge and return true.
func (g *arrayGraph) AddEdge(v, w int) error {
	if w == v {
		return errors.New("The edge vertices are not distinct")
	}
	if v < 0 || g.Vertices() <= v {
		return errors.New("The source vertex is not in the graph")
	}
	if w < 0 || g.Vertices() <= w {
		return errors.New("The target vertex is not in the graph")
	}
	if g.adjacent[v][w] {
		return nil
	}
	g.adjacent[v][w] = true
	g.adjacent[w][v] = true
	g.numEdges++
	return nil
}

// IsEdge determines whether the receiver graph contains edge {v,w}
func (g *arrayGraph) IsEdge(v, w int) bool {
	if w == v {
		return false
	}
	if v < 0 || g.Vertices() <= v {
		return false
	}
	if w < 0 || g.Vertices() <= w {
		return false
	}
	return g.adjacent[v][w]
}

// NewIterator returns an iterator over the vertices adjacent to v.
// Pre: 0 <= v <= g.Vertices()
// Pre violation: return nil and false.
// Normal return: return a new iterator and true.
func (g *arrayGraph) NewIterator(v int) (Iterator, error) {
	if v < 0 || g.Vertices() <= v {
		return nil, errors.New("The source vertex is not in the graph")
	}
	result := new(arrayGraphIterator)
	result.g = g
	result.v = v
	result.Reset()
	return result, nil
}

// String produces a string representation of a graph.
func (g *arrayGraph) String() string {
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

//////////////////////////////////////////////////////////////////////////////////////////
// arrayGraphVertexIterator holds data about iterating over vertices adjacent to vertex v.
type arrayGraphIterator struct {
	g *arrayGraph // the graph containing the iterator
	v int         // the source vertex
	w int         // the vertex we have reached so far
}

// Reset prepares for a new iteration.
func (iter *arrayGraphIterator) Reset() {
	for iter.w = 0; iter.w < iter.g.Vertices(); iter.w++ {
		if iter.g.adjacent[iter.v][iter.w] {
			break
		}
	}
}

// IsDone is true iff iteration is complete.
func (iter *arrayGraphIterator) IsDone() bool {
	return iter.g.Vertices() <= iter.w
}

// Next return the next adjacent vertex.
// Pre: Iteration is not complete.
// Pre violation: return 0 and false.
// Normal return: the next vertex and true.
func (iter *arrayGraphIterator) Next() (int, bool) {
	if iter.g.Vertices() <= iter.w {
		return 0, false
	}
	result := iter.w
	for iter.w++; iter.w < iter.g.Vertices(); iter.w++ {
		if iter.g.adjacent[iter.v][iter.w] {
			break
		}
	}
	return result, true
}

///////////////////////////////////////////////////////////////////////////////////////
// linkedGraph is the data structure for the adjacency lists representation of a graph.
type linkedGraph struct {
	numEdges int               // in the graph
	adjacent []containers.List // linked list of vertices adjacent to [v]
}

// NewLinkedGraph returns a pointer to a graph represented using adjacency lists.
// Pre: n > 0
// Pre violation: return a graph with 1 vertex.
// Normal return: return a graph with n vertices.
func NewLinkedGraph(n int) *linkedGraph {
	result := new(linkedGraph)
	if n < 0 {
		n = 1
	}
	result.adjacent = make([]containers.List, n)
	for i := 0; i < n; i++ {
		result.adjacent[i] = containers.NewLinkedList()
	}
	return result
}

// The containers only store interface types, so we must make one for vertices.
type Vertex int

// Edges return the number of edges in the receiver graph.
func (g *linkedGraph) Edges() int {
	return g.numEdges
}

// Vertices return the number of vertices in the receiver graph.
func (g *linkedGraph) Vertices() int {
	return len(g.adjacent)
}

// AddEdge puts a new edge in the receiver graph; it does nothing if
// the edge is already there.
// Pre: v and w are in the graph.
// Pre violation: return false.
// Normal return: add the edge and return true.
func (g *linkedGraph) AddEdge(v, w int) error {
	if w == v {
		return errors.New("The edge vertices are not distinct")
	}
	if v < 0 || g.Vertices() <= v {
		return errors.New("The source vertex is not in the graph")
	}
	if w < 0 || g.Vertices() <= w {
		return errors.New("The target vertex is not in the graph")
	}
	if g.IsEdge(v, w) {
		return nil
	}
	g.adjacent[v].Insert(0, Vertex(w))
	g.adjacent[w].Insert(0, Vertex(v))
	g.numEdges++
	return nil
}

// IsEdge determines whether the receiver graph contains edge {v,w}
func (g *linkedGraph) IsEdge(v, w int) bool {
	if w == v {
		return false
	}
	if v < 0 || g.Vertices() <= v {
		return false
	}
	if w < 0 || g.Vertices() <= w {
		return false
	}
	return g.adjacent[v].Contains(Vertex(w))
}

// NewIterator returns an iterator over the vertices adjacent to v.
// Pre: 0 <= v <= g.Vertices()
// Pre violation: return nil and false.
// Normal return: return a new iterator and true.
func (g *linkedGraph) NewIterator(v int) (Iterator, error) {
	if v < 0 || g.Vertices() <= v {
		return nil, errors.New("The source vertex is not in the graph")
	}
	result := new(linkedGraphIterator)
	result.g = g
	result.v = v
	result.Reset()
	return result, nil
}

///////////////////////////////////////////////////////////////////////////////////////////
// linkedGraphVertexIterator holds data about iterating over vertices adjacent to vertex v.
type linkedGraphIterator struct {
	g        *linkedGraph        // the graph containing the iterator
	v        int                 // the source vertex
	listIter containers.Iterator // to iterate through the list of vertices adjacent to v
}

// Reset prepares for a new iteration.
func (iter *linkedGraphIterator) Reset() {
	iter.listIter = iter.g.adjacent[iter.v].NewIterator()
}

// IsDone is true iff iteration is complete.
func (iter *linkedGraphIterator) IsDone() bool {
	return iter.listIter.IsDone()
}

// Next return the next adjacent vertex.
// Pre: Iteration is not complete.
// Pre violation: return 0 and false.
// Normal return: the next vertex and true.
func (iter *linkedGraphIterator) Next() (int, bool) {
	w, ok := iter.listIter.Next()
	if !ok {
		return 0, false
	}
	return int(w.(Vertex)), true
}

// String produces a string representation of a graph.
func (g *linkedGraph) String() string {
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
