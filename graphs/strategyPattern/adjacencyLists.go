// adjacencyLists.go: This file contains an implementation of undirected graphs using
// adjacency lists. In particular, it contains the adjacencyLists and adjacencyListsIterator
// types and implementations of their method sets.
//
// author: C. Fox
// version: 8/2012

package graphs

import "containers" // use a linked list in the linked graph representation
import "errors"     // for illegal vertices and like errors

/////////////////////////////////////////////////////////////////////////////////
// adjacencyLists is the data structure for the linked representation of a graph.
type adjacencyLists struct {
	numEdges int               // in the graph
	adjacent []containers.List // linked list of vertices adjacent to [v]
}

// The containers only store interface types, so we must make one for vertices.
type vertex int

// edges returns the number of edges in the adjacency lists.
func (lists *adjacencyLists) edges() int {
	return lists.numEdges
}

// vertices return the number of vertices in the adjacency lists
func (lists *adjacencyLists) vertices() int {
	return len(lists.adjacent)
}

// addEdge puts a new edge in the adjacency lists between v and w; it does nothing
// if the edge is already there.
// Pre: v != w and 0 <= v, w < len(lists.adjacent)
// Pre violation: return an error indication.
// Normal return: add the edge and return nil.
func (lists *adjacencyLists) addEdge(v, w int) error {
	if w == v {
		return errors.New("The edge vertices are not distinct")
	}
	if v < 0 || lists.vertices() <= v {
		return errors.New("The source vertex is not in the graph")
	}
	if w < 0 || lists.vertices() <= w {
		return errors.New("The target vertex is not in the graph")
	}
	if lists.isEdge(v, w) {
		return nil
	}
	lists.adjacent[v].Insert(0, vertex(w))
	lists.adjacent[w].Insert(0, vertex(v))
	lists.numEdges++
	return nil
}

// isEdge determines whether the receiver adjacency lists contain edge {v,w}
func (lists *adjacencyLists) isEdge(v, w int) bool {
	if w == v {
		return false
	}
	if v < 0 || lists.vertices() <= v {
		return false
	}
	if w < 0 || lists.vertices() <= w {
		return false
	}
	return lists.adjacent[v].Contains(vertex(w))
}

// newIterator returns an iterator over the vertices adjacent to v.
// Pre: 0 <= v < lists.vertices()
// Pre violation: return nil and an error indication.
// Normal return: return a new iterator and nil.
func (lists *adjacencyLists) newIterator(v int) (Iterator, error) {
	if v < 0 || lists.vertices() <= v {
		return nil, errors.New("The source vertex is not in the graph")
	}
	result := new(adjacencyListsIterator)
	result.lists = lists
	result.v = v
	result.Reset()
	return result, nil
}

///////////////////////////////////////////////////////////////////////////////////////////
// adjacencyListsIterator holds data about iterating over vertices adjacent to vertex v.
type adjacencyListsIterator struct {
	lists    *adjacencyLists     // the graph containing the iterator
	v        int                 // the source vertex
	listIter containers.Iterator // to iterate through the list of vertices adjacent to v
}

// Reset prepares for a new iteration.
func (iter *adjacencyListsIterator) Reset() {
	iter.listIter = iter.lists.adjacent[iter.v].NewIterator()
}

// IsDone is true iff iteration is complete.
func (iter *adjacencyListsIterator) IsDone() bool {
	return iter.listIter.IsDone()
}

// Next return the next adjacent vertex.
// Pre: Iteration is not complete.
// Pre violation: return 0 and false.
// Normal return: the next vertex and true.
func (iter *adjacencyListsIterator) Next() (int, bool) {
	w, ok := iter.listIter.Next()
	if !ok {
		return 0, false
	}
	return int(w.(vertex)), true
}
