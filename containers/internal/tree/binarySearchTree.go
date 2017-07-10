// binarySearchTree.go: Implementation of an unbalanced binary search tree.
// This type stores values implementing Comparer, which has Equal()
// and Less() operations for navigating the tree.
//
// author:  C. Fox
// version: 6/2017

package tree

import (
	"containers"
)

// A BinarySearchTree is BinaryTree whose nodes are in order when travered in order.
type BinarySearchTree struct {
	BinaryTree
}

// Return true iff element e is in the tree. This is a reimplementation of
// the version from binaryTree that takes advantage of the structure of
// a binary search tree to get better performance.
func (tree *BinarySearchTree) Contains(e interface{}) bool {
	value := e.(containers.Comparer)
	for node := tree.root; node != nil; {
		switch {
		case value.Equal(node.value):
			return true
		case value.Less(node.value):
			node = node.left
		default:
			node = node.right
		}
	}
	return false
}

// Create a new node holding value v and put it at the right spot
// at the bottom of the binary search tree. If v is already in the tree,
// replace the value at the node with v.
func (tree *BinarySearchTree) Add(v containers.Comparer) {
	if tree.root == nil {
		tree.root = new(btNode)
		tree.root.value = v
		tree.count++
		return
	}
	for node := tree.root; ; {
		switch {
		case v.Equal(node.value):
			node.value = v
			return
		case v.Less(node.value):
			if node.left == nil {
				node.left = newBTNode(v, nil, nil)
				tree.count++
				return
			} else {
				node = node.left
			}
		default:
			if node.right == nil {
				node.right = newBTNode(v, nil, nil)
				tree.count++
				return
			} else {
				node = node.right
			}
		}
	}
}

// Return the value in the tree matching argument v, if any.
// Precondition: Value v is in the tree.
// Preconditon gviolation: return nil and false.
// Normal return: the nodeValue and true.
func (tree *BinarySearchTree) Get(v containers.Comparer) (interface{}, bool) {
	for node := tree.root; node != nil; {
		switch {
		case v.Equal(node.value):
			return node.value, true
		case v.Less(node.value):
			node = node.left
		default:
			node = node.right
		}
	}
	return nil, false
}

// Take a node with value v out of the tree. If v is not in the tree, do nothing.
func (tree *BinarySearchTree) Remove(v containers.Comparer) {
	var (
		parent *btNode // the parent of target node in the tree
		target *btNode
	) // ultimately, the node that contains v

	// find the node containing v; return if not found
	for target = tree.root; target != nil; {
		targetValue := target.value
		if v.Equal(targetValue) {
			break
		}
		parent = target
		if v.Less(targetValue) {
			target = target.left
		} else {
			target = target.right
		}
	}
	if target == nil {
		return
	}
	tree.deleteNode(target, parent)
}

// Helper functions ------------------------------------------------------

// Remove a node from a binary search tree.
// If the deleted node has one child, attach the child to the node's parent.
// Otherwise, find the node's successor, swap values, and remove the
// successor node (which has no left child).
func (tree *BinarySearchTree) deleteNode(node, parent *btNode) {
	tree.count--
	switch {
	case node.left == nil:
		tree.attach(parent, node, node.right)
	case node.right == nil:
		tree.attach(parent, node, node.left)
	default:
		successor := node.right
		parent = node
		for successor.left != nil {
			parent = successor
			successor = successor.left
		}
		node.value, successor.value = successor.value, node.value
		if parent == node {
			parent.right = successor.right
		} else {
			parent.left = successor.right
		}
	}
}

// Attach a node to it's parent's parent
func (tree *BinarySearchTree) attach(parent, child, attached *btNode) {
	if parent == nil {
		tree.root = attached
	} else {
		if child == parent.left {
			parent.left = attached
		} else {
			parent.right = attached
		}
	}
}
