// avlTree.go: Implementation of a balanced binary search tree.
// This type stores values implementing Comparer, which has Equal()
// and Less() operations for navigating the tree.
//
// author:  C. Fox
// version: 6/2017

package tree

import (
	"fmt"
	"strconv"
	"strings"

	"containers"
)

// An AVL tree is a BinarySearchTree that stays balanced.
type AVLTree struct {
	BinarySearchTree
}

//////////////////////////////////////////////////////////////////
// Override BinaryTree and BinarySearchTree methods for AVL trees.

// Override BinaryTree.Height to use the height field in the root.
func (t *AVLTree) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.height
}

// Create a new node holding value v and put it at the right spot
// at the bottom of the binary search tree. If v is already in the tree,
// replace the value at the node with v.
func (tree *AVLTree) Add(v containers.Comparer) {
	if tree.root == nil {
		tree.root = newAVLNode(v, nil, nil)
		tree.count = 1
		return
	}
	if !tree.root.contains(v) {
		tree.count++
	}
	tree.root.add(v)
}

// Take a node with value v out of the tree. If v is not in the tree, do nothing.
func (tree *AVLTree) Remove(v containers.Comparer) {
	if !tree.Contains(v) {
		return
	}
	tree.count--
	tree.root = tree.root.remove(v)
}

// Override the BinaryTree String method to display the balance factors
// at each node.
func (tree *AVLTree) String() string {
	result := "Tree size: " + strconv.Itoa(tree.count) + "\n"
	if tree.root == nil {
		return result + "empty\n"
	}
	result += tree.root.toStringAVL(0)
	return result
}

////////////////////////////////////////////////////////////////
// Add methods to and for the btNode type for AVL trees.

// This method sets the height field when creating an AVL tree node.
func newAVLNode(v interface{}, left, right *btNode) *btNode {
	result := newBTNode(v, left, right)
	result.setHeight()
	return result
}

// Adjust the height of node based on its children's heights.
func (node *btNode) setHeight() {
	leftHeight := -1
	if node.left != nil {
		leftHeight = node.left.height
	}
	rightHeight := -1
	if node.right != nil {
		rightHeight = node.right.height
	}
	if leftHeight < rightHeight {
		node.height = rightHeight + 1
	} else {
		node.height = leftHeight + 1
	}
}

// Compute the AVL tree balance factor for a node, which is
// the left sub-tree height minus the right sub-tree height.
func (node *btNode) balance() int {
	leftHeight := -1
	if node.left != nil {
		leftHeight = node.left.height
	}
	rightHeight := -1
	if node.right != nil {
		rightHeight = node.right.height
	}
	return leftHeight - rightHeight
}

// Check the balance factor and it is 2 or -2, then rebalance the sub-tree from this node.
// There are two kinds of rebalancing depending on the shape of the sub-trees. If the balance
// factor is 2, then an R rotation is done when the left child has a balance factor of 1 or 0;
// otherwise (that is, a balance factor of -1), an LR rotation is done. The rotations done
// with a balance factor of -2 are symmetric.
// The height is always reset because even if no rotations occur, this method is only called
// when the tree has been altered, so height must be reset anyway.
func (node *btNode) rebalance() {
	switch node.balance() {
	case 2:
		if node.left.balance() == -1 {
			node.rotateLR()
		} else {
			node.rotateR()
		}
	case -2:
		if node.right.balance() == 1 {
			node.rotateRL()
		} else {
			node.rotateL()
		}
	}
	node.setHeight()
}

// Insert a value in the tree rooted at node.
func (node *btNode) add(v containers.Comparer) {
	switch {
	case v.Equal(node.value):
		node.value = v
		return
	case v.Less(node.value):
		if node.left == nil {
			node.left = newAVLNode(v, nil, nil)
			node.setHeight()
			return
		} else {
			node.left.add(v)
		}
	default:
		if node.right == nil {
			node.right = newAVLNode(v, nil, nil)
			node.setHeight()
			return
		} else {
			node.right.add(v)
		}
	}
	node.rebalance()
}

// Delete a value in the tree rooted at node.
func (node *btNode) remove(v containers.Comparer) *btNode {
	switch {
	case v.Equal(node.value):
		if node.left == nil {
			return node.right
		}
		if node.right == nil {
			return node.left
		}
		node.value, node.right = node.right.removeSuccessor()
	case v.Less(node.value):
		if node.left == nil {
			return node
		} else {
			node.left = node.left.remove(v)
		}
	default:
		if node.right == nil {
			return node
		} else {
			node.right = node.right.remove(v)
		}
	}
	node.rebalance()
	return node
}

// Find the successor of a deleted node, remove it, and return its value along with the
// nodes along the path to the successor.
func (node *btNode) removeSuccessor() (interface{}, *btNode) {
	if node.left == nil {
		return node.value, node.right
	}
	resultValue, resultNode := node.left.removeSuccessor()
	node.left = resultNode
	node.rebalance()
	return resultValue, node
}

// Rotate left, meaning the right child becomes the root of the sub-tree.
func (node *btNode) rotateL() {
	newLeft := newAVLNode(node.value, node.left, node.right.left)
	node.value = node.right.value
	node.right = node.right.right
	node.left = newLeft
}

// Rotate right, meaning the left child becomes the root of the sub-tree.
func (node *btNode) rotateR() {
	newRight := newAVLNode(node.value, node.left.right, node.right)
	node.value = node.left.value
	node.left = node.left.left
	node.right = newRight
}

// Rotate left-right, meaning the right child of the left child of this node
// becomes the root of the sub-tree.
func (node *btNode) rotateLR() {
	newRight := newAVLNode(node.value, node.left.right.right, node.right)
	node.value = node.left.right.value
	node.left.right = node.left.right.left
	node.left.setHeight()
	node.right = newRight
}

// Rotate right-left, meaning the left child of the right child of this node
// becomes the root of the sub-tree.
func (node *btNode) rotateRL() {
	newLeft := newAVLNode(node.value, node.left, node.right.left.left)
	node.value = node.right.left.value
	node.right.left = node.right.left.right
	node.right.setHeight()
	node.left = newLeft
}

// Make a string representation of the tree rooted at node.
func (node *btNode) toStringAVL(indent int) string {
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
	result += " (" + strconv.Itoa(node.balance()) + ")\n"
	if node.left == nil && node.right == nil {
		return result
	}
	result += node.left.toStringAVL(indent + tab)
	result += node.right.toStringAVL(indent + tab)
	return result
}
