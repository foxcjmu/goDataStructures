// twoThreeTree.go: Impelementation of 2-3 trees.
// author: C. Fox
// versioin: 6/2017

package tree

import (
	"fmt"
	"strings"

	"containers"
	"containers/stack"
)

///////////////////////////////////////////////////////////////////
// A 2-3 tree is a perfectly balanced search tree whose nodes
// have zero, 2 or 3 children. All operations are O(lg n).
type TwoThreeTree struct {
	count int           // how many nodes in the tree
	root  *twoThreeNode // the root of the tree
}

// Determine whether a tree is empty.
func (t *TwoThreeTree) Empty() bool { return t.count == 0 }

// Determine how many nodes are in this tree.
func (t *TwoThreeTree) Size() int { return t.count }

// Make this tree empty.
func (t *TwoThreeTree) Clear() {
	t.count = 0
	t.root = nil
}

// Determine the height of this tree.
func (t *TwoThreeTree) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.height()
}

// Return a value from the tree and true, or nil and false if it is missing.
func (t *TwoThreeTree) Get(v containers.Comparer) (interface{}, bool) {
	if t.root == nil {
		return nil, false
	}
	return t.root.get(v)
}

// Determine whether a value is in the tree.
func (t *TwoThreeTree) Contains(v containers.Comparer) bool {
	if t.root == nil {
		return false
	}
	_, present := t.root.get(v)
	return present
}

// Add one value to the tree, or replace it if it is already present.
func (t *TwoThreeTree) Add(v containers.Comparer) {
	if t.root == nil {
		t.root = newTwoThreeNode(v, nil, nil)
		t.count = 1
	} else {
		root, addition := t.root.add(v)
		t.root = root
		if addition {
			t.count++
		}
	}
}

// Remove a value from this 2-3 tree, or do nothing if it is not present.
func (t *TwoThreeTree) Remove(v containers.Comparer) {
	if t.root == nil {
		return
	}
	deletion := t.root.remove(v, nil, 0)
	if deletion {
		t.count--
	}
	if t.root.nodeType == 1 {
		t.root = t.root.left
	}
}

// Apply a visitor method to every value in the tree in order.
// Note that in a 2-3 tree, the leftmost subtree is visited first, followed
// by the leftmost value, followed by the middle subtree, followed (if the
// node is a 3-node) by the rightmost value, followed by the rightmost
// subtree.
func (t *TwoThreeTree) Visit(visitor func(interface{})) {
	if t.root == nil {
		return
	}
	t.root.visitInorder(visitor)
}

// Apply a visitor method to every value in the tree in preorder.
// Note that in a 2-3 tree, the values at the node are visited first,
// followed by the leftmost subtree, followed by the middle subtree,
// followed (if the node is a 3-node) by the rightmost subtree.
func (t *TwoThreeTree) VisitPreorder(visitor func(interface{})) {
	if t.root == nil {
		return
	}
	t.root.visitPreorder(visitor)
}

// Apply a visitor method to every value in the tree in post order.
// Note that in a 2-3 tree, the values in the leftmost subtree are visited
// first, followed by the middle subtree, followed (if the node is a 3-node)
// by the rightmost subtree, followed by the value(s) at this node.
func (t *TwoThreeTree) VisitPostorder(visitor func(interface{})) {
	if t.root == nil {
		return
	}
	t.root.visitPostorder(visitor)
}

// Create and return in inorder iterator over this tree.
func (t *TwoThreeTree) NewIterator() containers.Iterator {
	result := new(twoThreeTreeIterator)
	result.stack = new(stack.LinkedStack)
	result.root = t.root
	result.Reset()
	return result
}

// Make a string representation of a tree.
const INDENT = "   " // how far to indent at each tree level
func (t *TwoThreeTree) String() string {
	if t.root == nil {
		return "Empty 2-3 tree\n"
	}
	result := fmt.Sprintf("2-3 tree size %v\n", t.count)
	result += t.root.toString(0)
	return result
}

///////////////////////////////////////////////////////////////////////////
// twoThreeNode declarations and receiver functions ///////////////////////

// twoThreeTree nodes actually do most of the work implementing
// 2-3 trees using // recursive algorithms on node receivers.
type twoThreeNode struct {
	nodeType int           // 1, 2, or 3 node
	value1   interface{}   // leftmost value
	value2   interface{}   // rightmost value
	left     *twoThreeNode // leftmost child
	mid      *twoThreeNode // right (middle) child in a 2-node (3-node)
	right    *twoThreeNode // rightmost child in a 3-node
}

// Create a return a new 2-node.
func newTwoThreeNode(v interface{}, leftChild, midChild *twoThreeNode) *twoThreeNode {
	r := new(twoThreeNode)
	r.nodeType, r.value1, r.left, r.mid = 2, v, leftChild, midChild
	return r
}

// Return true iff r is a leaf node.
func (r *twoThreeNode) isLeaf() bool {
	return r.left == nil
}

// Determine the height of the tree rooted at this node.
// Since the tree is balanced, we need only look down one path to the leaves.
func (r *twoThreeNode) height() int {
	if r.isLeaf() {
		return 0
	}
	return 1 + r.left.height()
}

// Return a value matching a given value in the tree rooted at this node
// and true, or nil and false if it is not present.
func (r *twoThreeNode) get(v containers.Comparer) (interface{}, bool) {
	switch {
	case v.Less(r.value1):
		if r.isLeaf() {
			return nil, false
		}
		return r.left.get(v)
	case v.Equal(r.value1):
		return r.value1, true
	case r.nodeType == 2 || v.Less(r.value2):
		if r.isLeaf() {
			return nil, false
		}
		return r.mid.get(v)
	case v.Equal(r.value2):
		return r.value2, true
	default:
		if r.isLeaf() {
			return nil, false
		}
		return r.right.get(v)
	}
}

// Put a new value into the subtree rooted at this node.
// Return a (possibly new) subtree root and whether a replacement was made.
// This is a somewhat complex operation:
// - insertion always occurs at a leaf, so we recurse down the tree
// - if the leaf is a 2-node, the value is added and it becomes a 3-node
// - if the leaf is a 3-node, it is split to create a sub-tree with three
//   2-nodes, and the parent node is returned
// - after insertion in a subtree, if a split node is returned, then if
//   this node is a 2-node, it is made into a 3-node and the subtrees are
//   incorporated into it
// - if this is a 3-node and a split node from a subtree insertion, then this
//   node is split to make a new-subtree that is in turn returned
func (r *twoThreeNode) add(v containers.Comparer) (*twoThreeNode, bool) {
	switch {
	case v.Less(r.value1): // v is in the left sub-tree
		if r.isLeaf() {
			if r.nodeType == 3 {
				newLeft := newTwoThreeNode(v, nil, nil)
				newRight := newTwoThreeNode(r.value2, nil, nil)
				return newTwoThreeNode(r.value1, newLeft, newRight), true
			}
			r.nodeType, r.value1, r.value2 = 3, v, r.value1
			return r, true
		}
		newLeft, addition := r.left.add(v)
		if newLeft == r.left {
			return r, addition
		}
		if r.nodeType == 3 {
			newRight := newTwoThreeNode(r.value2, r.mid, r.right)
			return newTwoThreeNode(r.value1, newLeft, newRight), true
		}
		r.nodeType, r.value1, r.value2, r.left, r.mid, r.right =
			3, newLeft.value1, r.value1, newLeft.left, newLeft.mid, r.mid
		return r, true

	case v.Equal(r.value1): // v is already in the tree--replace it
		r.value1 = v
		return r, false

	case r.nodeType == 2: // v is in the right subtree and this is a 2-node
		if r.isLeaf() {
			r.nodeType, r.value2 = 3, v
			return r, true
		}
		newMid, addition := r.mid.add(v)
		if newMid != r.mid {
			r.nodeType, r.value2, r.mid, r.right =
				3, newMid.value1, newMid.left, newMid.mid
		}
		return r, addition

	case v.Less(r.value2): // v is in the mid subtree of a 3-node
		if r.isLeaf() {
			newLeft := newTwoThreeNode(r.value1, nil, nil)
			newRight := newTwoThreeNode(r.value2, nil, nil)
			return newTwoThreeNode(v, newLeft, newRight), true
		}
		newMid, addition := r.mid.add(v)
		if newMid != r.mid {
			newLeft := newTwoThreeNode(r.value1, r.left, newMid.left)
			newRight := newTwoThreeNode(r.value2, newMid.mid, r.right)
			return newTwoThreeNode(newMid.value1, newLeft, newRight), true
		}
		return r, addition

	case v.Equal(r.value2): // v is already in the tree--replace it
		r.value2 = v
		return r, false

	default: // v is in the right subtree of a 3-node
		if r.isLeaf() {
			newLeft := newTwoThreeNode(r.value1, nil, nil)
			newRight := newTwoThreeNode(v, nil, nil)
			return newTwoThreeNode(r.value2, newLeft, newRight), true
		}
		newRight, addition := r.right.add(v)
		if newRight != r.right {
			newLeft := newTwoThreeNode(r.value1, r.left, r.mid)
			return newTwoThreeNode(r.value2, newLeft, newRight), true
		}
		return r, addition
	}
} // add

// Recursively delete a value from this tree.
// Deletion always occurs starting from a leaf, so if the deleted value is in an internal
// node, its successor is copied in to the node holding the delete value (the target node),
// and the successor is deleted (just like in a BST or AVL tree).
// If the deleted value is in a leaf, then we delete it. This may change the node from a
// 3-node to a 2-node, which is ok, or from a 2-node to a 1-node, which propagates a problem
// to the parent (hence the usefulness of recursion).
// If a node has a child that is a 1-node then
// - if the child has a sibling that is a 3-node, a value and a subtree are borrowed from
//   the sibling fix the 1-node
// - if this node (the parent) is a 3-node, then a value and a subtree are borrowed from
//   the parent to fix the 1-node
// - otherwise, this is a 2-node with a 2-node child and a 1-node child: this node becomes a
//   1-node with a 3-node child holding the 2-node value, the parent value, and the 1-node and
//   2-node children. The new 1-node is returned, propagating the problem up the tree
func (r *twoThreeNode) remove(v containers.Comparer, target *twoThreeNode, which int) bool {

	if r.isLeaf() {
		if target != nil {
			switch which {
			case 1:
				target.value1 = r.value1
			case 2:
				target.value2 = r.value1
			}
			r.value1 = r.value2
		} else {
			switch {
			case v.Less(r.value1):
				return false
			case v.Equal(r.value1):
				r.value1 = r.value2
			case r.nodeType == 2 || v.Less(r.value2):
				return false
			case !v.Equal(r.value2):
				return false
			}
		}
		r.nodeType--
		return true
	}

	// do the deletion
	var deletion = true
	if target != nil {
		r.left.remove(v, target, which)
	} else {
		switch {
		case v.Less(r.value1):
			deletion = r.left.remove(v, nil, 0)
		case v.Equal(r.value1):
			deletion = r.mid.remove(v, r, 1)
		case r.nodeType == 2 || v.Less(r.value2):
			deletion = r.mid.remove(v, nil, 0)
		case v.Equal(r.value2):
			deletion = r.right.remove(v, r, 2)
		default:
			deletion = r.right.remove(v, nil, 0)
		}
	}

	// if any child is a 1-node, fix it
	if r.left.nodeType == 1 {
		if r.mid.nodeType == 3 || (r.nodeType == 3 && r.right.nodeType == 3) {
			r.leftBorrowsFromMid()
			if r.mid.nodeType == 1 {
				r.midBorrowsFromRight()
			}
		} else if r.nodeType == 3 {
			r.foldLeftIntoMid()
		} else {
			r.pushLeftIntoMid()
		}
	}
	if r.mid.nodeType == 1 {
		if r.left.nodeType == 3 {
			r.midBorrowsFromLeft()
		} else if r.nodeType == 3 {
			if r.right.nodeType == 3 {
				r.midBorrowsFromRight()
			} else {
				r.foldMidIntoRight()
			}
		} else {
			r.pushMidIntoLeft()
		}
	}
	if r.nodeType == 3 && r.right.nodeType == 1 {
		if r.left.nodeType == 3 || r.mid.nodeType == 3 {
			r.rightBorrowsFromMid()
			if r.mid.nodeType == 1 {
				r.midBorrowsFromLeft()
			}
		} else {
			r.foldRightIntoMid()
		}
	}
	return deletion

} // remove

// Move the rightmost portions of this node to the left.
func (r *twoThreeNode) shiftLeft() {
	r.value1, r.left, r.mid = r.value2, r.mid, r.right
}

//  Move the leftmost portions of this node to the right.
func (r *twoThreeNode) shiftRight() {
	r.value2, r.right, r.mid = r.value1, r.mid, r.left
}

// During a deletion, a leftmost subtree is a 1-node that is made
// into a 2-node by borrowing from the subtree to its right.
// Pre: left.nodeType == 1 mid.nodeType == 2..3
func (r *twoThreeNode) leftBorrowsFromMid() {
	r.left.value1 = r.value1
	r.value1 = r.mid.value1
	r.left.mid = r.mid.left
	r.mid.shiftLeft()
	r.left.nodeType = 2
	r.mid.nodeType--
}

// During a deletion, a middle subtree is a 1-node that is made
// into a 2-node by borrowing from the subtree to its left.
// Pre: mid.nodeType == 1 left.nodeType == 3
func (r *twoThreeNode) midBorrowsFromLeft() {
	r.mid.value1 = r.value1
	r.value1 = r.left.value2
	r.mid.mid = r.mid.left
	r.mid.left = r.left.right
	r.mid.nodeType = 2
	r.left.nodeType = 2
}

// During a deletion, a middle subtree is a 1-node that is made
// into a 2-node by borrowing from the subtree to its right.
// Pre: mid.nodeType == 1 right.nodeType == 3
func (r *twoThreeNode) midBorrowsFromRight() {
	r.mid.value1 = r.value2
	r.value2 = r.right.value1
	r.mid.mid = r.right.left
	r.right.shiftLeft()
	r.mid.nodeType = 2
	r.right.nodeType = 2
}

// During a deletion, a rightmost subtree in a 3-node is a 1-node that
// is made into a 2-node by borrowing from the subtree to its left.
// pre: right.nodeType == 1 mid.nodeType == 2..3
func (r *twoThreeNode) rightBorrowsFromMid() {
	r.right.value1 = r.value2
	r.right.mid = r.right.left
	switch r.mid.nodeType {
	case 2:
		r.value2 = r.mid.value1
		r.right.left = r.mid.mid
	case 3:
		r.value2 = r.mid.value2
		r.right.left = r.mid.right
	default:
		panic("Unreachable code")
	}
	r.right.nodeType = 2
	r.mid.nodeType--
}

// During a deletion, a leftmost subtree that is a 1-node is incorporated
// into the 2-node to its right, borrowing a value from this node and thus
// making it into a 2-node.
// Pre: nodeType == 3 left.nodeType == 1 mid.nodeType == 2
func (r *twoThreeNode) foldLeftIntoMid() {
	r.mid.shiftRight()
	r.mid.value1 = r.value1
	r.mid.left = r.left.left
	r.mid.nodeType = 3
	r.shiftLeft()
	r.nodeType = 2
}

// During a deletion, a middle subtree that is a 1-node is incorporated
// into the 2-node to its right, borrowing a value from this node and thus
// making it into a 2-node.
// pre: nodeType == 3 mid.nodeType == 1 right.nodeType == 2
func (r *twoThreeNode) foldMidIntoRight() {
	r.right.shiftRight()
	r.right.value1 = r.value2
	r.right.left = r.mid.left
	r.mid = r.right
	r.right.nodeType = 3
	r.nodeType = 2
}

// During a deletion, a rightmost subtree in a 3-node that is a 1-node is
// incorporated into the 2-node to its left, borrowing a value from this
// node and thus making it into a 2-node.
// Pre: nodeType == 3 mid.nodeType == 2 right.nodeType == 1
func (r *twoThreeNode) foldRightIntoMid() {
	r.mid.value2 = r.value2
	r.mid.right = r.right.left
	r.mid.nodeType = 3
	r.nodeType = 2
}

// During a deletion, this 2-node has a left child that is a 1-node and
// a right child that is a 2-node, forcing the left child to be combined
// with the right, borrowing a value from this, to form a 3-node that is
// the only child of this, which then becomes a 1-node.
// Pre: nodeType == 2 left.nodeType == 1 mid.nodeType == 2
func (r *twoThreeNode) pushLeftIntoMid() {
	r.mid.shiftRight()
	r.mid.value1 = r.value1
	r.mid.left = r.left.left
	r.mid.nodeType = 3
	r.left = r.mid
	r.nodeType = 1
}

// During a deletion, this 2-node has a right child that is a 1-node and
// a left child that is a 2-node, forcing the right child to be combined
// with the left, borrowing a value from this, to form a 3-node that is
// the only child of this, which then becomes a 1-node.
// Pre: nodeType == 2 left.nodeType == 2 mid.nodeType == 1
func (r *twoThreeNode) pushMidIntoLeft() {
	r.left.value2 = r.value1
	r.left.right = r.mid.left
	r.left.nodeType = 3
	r.nodeType = 1
}

// Recursively apply a visitor function in order.
func (r *twoThreeNode) visitInorder(visitor func(interface{})) {
	if r.isLeaf() {
		visitor(r.value1)
		if r.nodeType == 3 {
			visitor(r.value2)
		}
	} else {
		r.left.visitInorder(visitor)
		visitor(r.value1)
		r.mid.visitInorder(visitor)
		if r.nodeType == 3 {
			visitor(r.value2)
			r.right.visitInorder(visitor)
		}
	}
}

// Recursively apply a visitor function in preorder.
func (r *twoThreeNode) visitPreorder(visitor func(interface{})) {
	if r.isLeaf() {
		visitor(r.value1)
		if r.nodeType == 3 {
			visitor(r.value2)
		}
	} else {
		visitor(r.value1)
		if r.nodeType == 3 {
			visitor(r.value2)
		}
		r.left.visitPreorder(visitor)
		r.mid.visitPreorder(visitor)
		if r.nodeType == 3 {
			r.right.visitPreorder(visitor)
		}
	}
}

// Recursively apply a visitor function in preorder.
func (r *twoThreeNode) visitPostorder(visitor func(interface{})) {
	if r.isLeaf() {
		visitor(r.value1)
		if r.nodeType == 3 {
			visitor(r.value2)
		}
	} else {
		r.left.visitPostorder(visitor)
		r.mid.visitPostorder(visitor)
		if r.nodeType == 3 {
			r.right.visitPostorder(visitor)
		}
		visitor(r.value1)
		if r.nodeType == 3 {
			visitor(r.value2)
		}
	}
}

// Recursively make a copy of the tree rooted at this node.
func (r *twoThreeNode) clone() *twoThreeNode {
	var (
		newLeft  *twoThreeNode
		newMid   *twoThreeNode
		newRight *twoThreeNode
	)
	if r.left != nil {
		newLeft = r.left.clone()
	}
	if r.mid != nil {
		newMid = r.mid.clone()
	}
	result := newTwoThreeNode(r.value1, newLeft, newMid)
	if r.nodeType == 3 {
		if r.right != nil {
			newRight = r.right.clone()
		}
		result.nodeType, result.value2, result.right = 3, r.value2, newRight
	}
	return result
}

///////////////////////////////////////////////////////////////////
// A 2-3 tree in order iterator ///////////////////////////////////

// Keep track of where we are iterating over a 2-3 tree.
// The stack holds tree nodes, and the top node in the stack always
// contains the current value as value1.
type twoThreeTreeIterator struct {
	stack stack.Stack
	root  *twoThreeNode
}

// Get ready for another iteration.
func (iter *twoThreeTreeIterator) Reset() {
	iter.stack.Clear()
	if iter.root != nil {
		node := iter.root
		for node != nil {
			iter.stack.Push(node)
			node = node.left
		}
	}
}

// Determine whether iteration is complete.
func (iter *twoThreeTreeIterator) Done() bool {
	return iter.stack.Empty()
}

// This algorithm is like the binary tree in order iteration algorithm except that
// it converts 3-nodes to 2-nodes and pushes them back on the stack after they are
// visited the first time.
func (iter *twoThreeTreeIterator) Next() (interface{}, bool) {
	if top, err := iter.stack.Pop(); err != nil {
		return nil, false
	} else {
		node := top.(*twoThreeNode)
		result := node.value1
		if node.nodeType == 3 {
			iter.stack.Push(newTwoThreeNode(node.value2, node.mid, node.right))
		}
		if !node.isLeaf() {
			node := node.mid
			for node != nil {
				iter.stack.Push(node)
				node = node.left
			}
		}
		return result, true
	}
}

// Make a string representation of the tree rooted at this node.
func (r *twoThreeNode) toString(level int) string {
	result := strings.Repeat(INDENT, level) + fmt.Sprint(r.value1)
	if r.nodeType == 3 {
		result += fmt.Sprint(" / ", r.value2)
	}
	result += "\n"
	if r.isLeaf() {
		return result
	}
	result += r.left.toString(level + 1)
	result += r.mid.toString(level + 1)
	if r.nodeType == 3 {
		result += r.right.toString(level + 1)
	}
	return result
}
