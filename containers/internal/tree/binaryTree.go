// binaryTree.go -- Implementation of a basic binary tree type
//
// author:  C. Fox
// version: 6/2017

// tree provides basic binary trees and binary search tees for use
// in implementing containers.
package tree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"containers"
	"containers/stack"
)

// A BinaryTree is a tree shose nodes have 0, 1, or 2 children.
// inv: BinaryTree.root == nil iff BinaryTree.count == 0
// inv: BinaryTree.count >= 0.
type BinaryTree struct {
	count int     // how many nodes in the tree
	root  *btNode // the root of the tree
}

// buildbinaryTree creates and returns a binary tree made from other binary trees.
// Note that when a new tree is created, a complete copy of the graphs of its
// constituent trees is created (i.e. this is a deep tree copy).
func buildBinaryTree(v interface{}, leftTree, rightTree BinaryTree) BinaryTree {
	var result BinaryTree
	result.count = 1 + rightTree.Size() + leftTree.Size()
	result.root = newBTNode(v, leftTree.root.clone(), rightTree.root.clone())
	return result
}

// RootValue returns the value at the root of the tree, if any.
// Precondition: The tree is not empty
// Precondition violation: Return nil and an error
// Normal return: The root value and nil
func (tree *BinaryTree) RootValue() (interface{}, error) {
	if tree.root == nil {
		return nil, errors.New("The empty tree has no root value")
	}
	return tree.root.value, nil
}

// LeftSubtree returns the left subtree of the root, if any.
// Precondition: The tree is not empty
// Precondition violation: Return nil and an error
// Normal return: A copy of the left subtree and nil
func (tree *BinaryTree) LeftSubtree() (BinaryTree, error) {
	var result BinaryTree
	if tree.root == nil {
		return result, errors.New("The empty tree has no left subtree")
	}
	result.count = tree.root.left.size()
	result.root = tree.root.left.clone()
	return result, nil
}

// RightSubtree returns the right subtree of the root, if any.
// Precondition: The tree is not empty.
// Precondition violation: Return nil and an error.
// Normal return: A copy of the right subtree and nil.
func (tree *BinaryTree) RightSubtree() (BinaryTree, error) {
	var result BinaryTree
	if tree.root == nil {
		return result, errors.New("The empty tree has no right subtree")
	}
	result.count = tree.root.right.size()
	result.root = tree.root.right.clone()
	return result, nil
}

// Size returns the number of nodes in a binary tree.
func (tree *BinaryTree) Size() int { return tree.count }

// Clear makes a binary tree empty.
func (tree *BinaryTree) Clear() {
	tree.count = 0
	tree.root = nil
}

// Empty returns true just in case the tree has no nodes.
func (tree *BinaryTree) Empty() bool { return tree.count == 0 }

// Height reports how many levels the tree has.
func (tree *BinaryTree) Height() int {
	return tree.root.getHeight()
}

// Contains determines whether a tree contains value e.
func (tree *BinaryTree) Contains(e interface{}) bool {
	if tree.root == nil {
		return false
	}
	return tree.root.contains(e)
}

// VisitPreorder is an internal iterator that applies a visit function f to every
// node in a binary tree in preorder (root, left subtree, then right subtree).
func (tree *BinaryTree) VisitPreorder(f func(e interface{})) {
	if tree.root == nil {
		return
	}
	tree.root.visitPreorder(f)
}

// VisitInorder is an internal iterator that applies a visit function f to every
// node in a binary tree inorder (left subtree, root, then right subtree).
func (tree *BinaryTree) VisitInorder(f func(e interface{})) {
	if tree.root == nil {
		return
	}
	tree.root.visitInorder(f)
}

// VisitPostorder is an internal iterator that applies a visit function f to every
// node in a binary tree in postorder (left subtree, right subtree, then root).
func (tree *BinaryTree) VisitPostorder(f func(e interface{})) {
	if tree.root == nil {
		return
	}
	tree.root.visitPostorder(f)
}

// NewPreorderIterator creates and returns a new preorder external iterator.
func (tree *BinaryTree) NewPreorderIterator() containers.Iterator {
	result := new(preorderIterator)
	result.stack = new(stack.LinkedStack)
	result.root = tree.root
	result.Reset()
	return result
}

// NewInorderIterator create and returns a new inorder external iterator.
func (tree *BinaryTree) NewInorderIterator() containers.Iterator {
	result := new(inorderIterator)
	result.stack = new(stack.LinkedStack)
	result.root = tree.root
	result.Reset()
	return result
}

// NewPostOrderIterator creates and returns a new postorder external iterator.
func (tree *BinaryTree) NewPostorderIterator() containers.Iterator {
	result := new(postorderIterator)
	result.stack = new(stack.LinkedStack)
	result.root = tree.root
	result.Reset()
	return result
}

// Display the tree as a string
func (tree *BinaryTree) String() string {
	result := "Tree size: " + strconv.Itoa(tree.count) + "\n"
	result += tree.root.toString(0)
	return result
}

////////////////////////////////////////////////////////////////////////////
// binary tree node type and helper functions //////////////////////////////

// This private struct is used for nodes in the graph of a binary tree.
type btNode struct {
	value  interface{} // what is stored at this node
	left   *btNode     // pointer to the root of the left subtree
	right  *btNode     // pointer to the root of the right subtree
	height int         // only used by AVL trees
}

// newBTNode allocates a new btNode from the heap and
// assigns values to its fields.
func newBTNode(v interface{}, leftTree, rightTree *btNode) *btNode {
	result := new(btNode)
	result.value, result.left, result.right = v, leftTree, rightTree
	return result
}

// clone makes a copy of the graph of a binary tree.
func (node *btNode) clone() *btNode {
	if node == nil {
		return nil
	}
	result := new(btNode)
	result.value = node.value
	result.left = node.left.clone()
	result.right = node.right.clone()
	return result
}

// getHeight figures out the height of tree recursively: the empty tree and
// the tree with a single node both have height 0; a tree with root r and subtrees
// t1 and t2 has height that is 1 more than the maximum of the heights of t1 and t2.
func (node *btNode) getHeight() int {
	if node == nil {
		return 0
	}
	if node.left == nil && node.right == nil {
		return 0
	}
	return 1 + max(node.left.getHeight(), node.right.getHeight())
}

// max finds the maximum of two ints.
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// size figures out how many nodes are in a tree recursively: the empty
// tree has size 0; a tree with root r and subtrees t1 and t2 has size 1 plus
// the sizes of t1 and t2.
func (node *btNode) size() int {
	if node == nil {
		return 0
	}
	return 1 + node.left.size() + node.right.size()
}

// Create a string representation of the tree rooted at this node.
func (node *btNode) toString(indent int) string {
	const tab = 3
	result := strings.Repeat(" ", indent)
	if node == nil {
		result += "-\n"
		return result
	}
	switch value := node.value.(type) {
	case fmt.Stringer:
		result += value.String()
	default:
		result += "?"
	}
	result += "\n"
	if node.left == nil && node.right == nil {
		return result
	}
	result += node.left.toString(indent + tab)
	result += node.right.toString(indent + tab)
	return result
}

// contains checks whether a root or its subtrees contain a value e.
func (node *btNode) contains(e interface{}) bool {
	if node.value == e {
		return true
	}
	if node.left != nil && node.left.contains(e) {
		return true
	}
	if node.right != nil && node.right.contains(e) {
		return true
	}
	return false
}

// visitPreorder applies a visit function f to root and its subtrees in preorder.
func (node *btNode) visitPreorder(f func(e interface{})) {
	f(node.value)
	if node.left != nil {
		node.left.visitPreorder(f)
	}
	if node.right != nil {
		node.right.visitPreorder(f)
	}
}

// visitInorder applies a visit function f to root and its subtrees in order.
func (node *btNode) visitInorder(f func(e interface{})) {
	if node.left != nil {
		node.left.visitInorder(f)
	}
	f(node.value)
	if node.right != nil {
		node.right.visitInorder(f)
	}
}

// visitPostorder applies a visit function f to root and its subtrees in postorder.
func (node *btNode) visitPostorder(f func(e interface{})) {
	if node.left != nil {
		node.left.visitPostorder(f)
	}
	if node.right != nil {
		node.right.visitPostorder(f)
	}
	f(node.value)
}

// Preorder Iterator implementation -----------------------------------------

// This private struct keeps track of the current state of preorder iteration.
// Invariant: current node is stack.Top()
type preorderIterator struct {
	stack stack.Stack // holds deferred nodes
	root  *btNode     // to reset to tree root
}

// Reset prepares for a new iteration.
func (iterator *preorderIterator) Reset() {
	iterator.stack.Clear()
	if iterator.root != nil {
		iterator.stack.Push(iterator.root)
	}
}

// Done indicates whether all elements have been accessed.
func (iterator *preorderIterator) Done() bool { return iterator.stack.Empty() }

// Next returns the next element and indication of whether there is one.
// Precondition: Iteration is not complete.
// Precondition violation: nil and false.
// Normal return: the next element and true.
func (iterator *preorderIterator) Next() (interface{}, bool) {
	e, err := iterator.stack.Pop()
	if err != nil {
		return nil, false
	}
	node := e.(*btNode)
	if node.right != nil {
		iterator.stack.Push(node.right)
	}
	if node.left != nil {
		iterator.stack.Push(node.left)
	}
	return node.value, true
}

// Inorder Iterator implementation -----------------------------------------

// This private struct keeps track of the current state of inorder iteration.
// Invariant: current node is stack.Top()
type inorderIterator struct {
	stack stack.Stack // holds deferred nodes
	root  *btNode     // to reset to tree root
}

// Reset prepares for a new iteration.
func (iterator *inorderIterator) Reset() {
	iterator.stack.Clear()
	node := iterator.root
	for node != nil {
		iterator.stack.Push(node)
		node = node.left
	}
}

// Done indicates whether all elements have been accessed.
func (iterator *inorderIterator) Done() bool {
	return iterator.stack.Empty()
}

// Next returns the next element and indication of whether there is one.
// Precondition: Iteration is not complete.
// Precondition violation: nil and false.
// Normal return: the next element and true.
func (iterator *inorderIterator) Next() (interface{}, bool) {
	e, err := iterator.stack.Pop()
	if err != nil {
		return nil, false
	}
	node := e.(*btNode)
	result := node.value
	node = node.right
	for node != nil {
		iterator.stack.Push(node)
		node = node.left
	}
	return result, true
}

// Postorder Iterator implementation ----------------------------------------

// This private struct keeps track of the current state of postorder iteration.
// Invariant: stack.Top() is the parent of nextNode (if any)
type postorderIterator struct {
	stack    stack.Stack // holds deferred nodes
	root     *btNode     // to reset to tree root
	nextNode *btNode     // node holding the value returned by Next()
}

// Reset prepares for a new iteration.
func (iterator *postorderIterator) Reset() {
	iterator.stack.Clear()
	node := iterator.root
	for node != nil {
		iterator.stack.Push(node)
		if node.left != nil {
			node = node.left
		} else {
			node = node.right
		}
	}
	if e, err := iterator.stack.Pop(); err == nil {
		iterator.nextNode = e.(*btNode)
	} else {
		iterator.nextNode = nil
	}
}

// Done indicates whether all elements have been accessed.
func (iterator *postorderIterator) Done() bool {
	return iterator.nextNode == nil
}

// Next returns the next element and indication of whether there is one.
// Precondition: Iteration is not complete.
// Precondition violation: nil and false.
// Normal return: the next element and true.
func (iterator *postorderIterator) Next() (interface{}, bool) {
	if iterator.nextNode == nil {
		return nil, false
	}
	result := iterator.nextNode.value
	if e, err := iterator.stack.Top(); err != nil {
		iterator.nextNode = nil
	} else {
		var node *btNode = e.(*btNode)
		if iterator.nextNode != node.right {
			node = node.right
			for node != nil {
				iterator.stack.Push(node)
				if node.left != nil {
					node = node.left
				} else {
					node = node.right
				}
			}
		}
		e, _ = iterator.stack.Pop()
		iterator.nextNode = e.(*btNode)
	}
	return result, true
}
