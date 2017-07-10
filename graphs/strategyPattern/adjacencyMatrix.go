// adjacencyMatrix.go: This file contains an implementation of undirected graphs using
// the adjacency matrix representation. In particular, this file includes the
// adjacencyMatrix and adjacencyMatrixIterator types and implementations of their
// method sets.
//
// author: C. Fox
// version: 8/2012

package graphs

import "errors" // for illegal vertices and like errors

//////////////////////////////////////////////////////////////////////////////////////
// adjacencyMatrix is the data structure for the contiguous representation of a graph.
type adjacencyMatrix struct {
	numEdges int      // in the graph
	adjacent [][]bool // true at [v][w] iff {v,w} is an edge
}

// edges returns the number of edges in the adjacency matrix.
func (matrix *adjacencyMatrix) edges() int {
	return matrix.numEdges
}

// vertices return the number of vertices in the adjacency matrix.
func (matrix *adjacencyMatrix) vertices() int {
	return len(matrix.adjacent)
}

// addEdge puts a new edge in the adjacency matrix between v and w; it does nothing
// if the edge is already there.
// Pre: v != w and 0 <= v, w < len(matrix.adjacent)
// Pre violation: return an error indication.
// Normal return: add the edge and return nil.
func (matrix *adjacencyMatrix) addEdge(v, w int) error {
	if w == v {
		return errors.New("The edge vertices are not distinct")
	}
	if v < 0 || matrix.vertices() <= v {
		return errors.New("The source vertex is not in the graph")
	}
	if w < 0 || matrix.vertices() <= w {
		return errors.New("The target vertex is not in the graph")
	}
	if !matrix.adjacent[v][w] {
		matrix.adjacent[v][w] = true
		matrix.adjacent[w][v] = true
		matrix.numEdges++
	}
	return nil
}

// isEdge determines whether the adjacency matrix contains edge {v,w}
func (matrix *adjacencyMatrix) isEdge(v, w int) bool {
	if w == v {
		return false
	}
	if v < 0 || matrix.vertices() <= v {
		return false
	}
	if w < 0 || matrix.vertices() <= w {
		return false
	}
	return matrix.adjacent[v][w]
}

// newIterator returns an iterator over the vertices adjacent to v.
// Pre: 0 <= v < matrix.vertices()
// Pre violation: return nil and an error indication.
// Normal return: return a new iterator and nil.
func (matrix *adjacencyMatrix) newIterator(v int) (Iterator, error) {
	if v < 0 || matrix.vertices() <= v {
		return nil, errors.New("The source vertex is not in the graph")
	}
	result := new(adjacencyMatrixIterator)
	result.matrix = matrix
	result.v = v
	result.Reset()
	return result, nil
}

/////////////////////////////////////////////////////////////////////////////////////////
// adjacencyMatrixIterator holds data about iterating over vertices adjacent to vertex v.
type adjacencyMatrixIterator struct {
	matrix *adjacencyMatrix // the matrix containing the vertices and edges
	v      int              // the source vertex
	w      int              // the vertex we have reached so far
}

// Reset prepares for a new iteration.
func (iter *adjacencyMatrixIterator) Reset() {
	for iter.w = 0; iter.w < iter.matrix.vertices(); iter.w++ {
		if iter.matrix.adjacent[iter.v][iter.w] {
			break
		}
	}
}

// IsDone is true iff iteration is complete.
func (iter *adjacencyMatrixIterator) IsDone() bool {
	return iter.matrix.vertices() <= iter.w
}

// Next return the next adjacent vertex.
// Pre: Iteration is not complete.
// Pre violation: return 0 and false.
// Normal return: the next vertex and true.
func (iter *adjacencyMatrixIterator) Next() (int, bool) {
	if iter.matrix.vertices() <= iter.w {
		return 0, false
	}
	result := iter.w
	for iter.w++; iter.w < iter.matrix.vertices(); iter.w++ {
		if iter.matrix.adjacent[iter.v][iter.w] {
			break
		}
	}
	return result, true
}
