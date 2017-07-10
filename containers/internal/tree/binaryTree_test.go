// Test BinaryTree interface and the Binary Tree implementation
// author: C. Fox
// version: 7/2012

package tree

import (
	"testing"
	//"fmt"

	"containers"
)

func TestEmptyBinaryTree(t *testing.T) {
	var r *BinaryTree = new(BinaryTree)

	// make sure a new BinaryTree is empty
	if !r.Empty() || r.Size() != 0 {
		t.Error("BinaryTree should be empty when new")
	}

	// test operations
	if r.Contains(6) {
		t.Error("Contains fails on empty BinaryTree")
	}
	if r.Height() != 0 {
		t.Error("Empty BinaryTree should have height 0")
	}
	if _, err := r.LeftSubtree(); err == nil {
		t.Error("Got left subtree of empty BinaryTree")
	}
	if _, err := r.RightSubtree(); err == nil {
		t.Error("Got right subtree of empty BinaryTree")
	}
	if _, err := r.RootValue(); err == nil {
		t.Error("Got root value of empty BinaryTree")
	}
	visitor := func(v interface{}) {
		t.Error("Visiting a non-empty tree")
	}
	r.VisitInorder(visitor)
	r.VisitPostorder(visitor)
	r.VisitPreorder(visitor)
	iter := r.NewInorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		visitor(e)
	}
	iter = r.NewPostorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		visitor(e)
	}
	iter = r.NewPreorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		visitor(e)
	}
}

func TestNonEmptyBinaryTree(t *testing.T) {

	// make a tree from other trees
	var empty BinaryTree
	r := buildBinaryTree(8, empty, empty)
	r = buildBinaryTree(12, r, buildBinaryTree(6, r, empty))

	if r.Empty() {
		t.Error("BinaryTree should not be empty")
	}
	if r.Size() != 4 {
		t.Errorf("BinaryTree size should be 4 but is %v", r.Size())
	}
	if r.Height() != 2 {
		t.Errorf("BinaryTree height should be 2 but is %v", r.Height())
	}
	if v, err := r.RootValue(); err != nil {
		t.Error("BinaryTree root should exist")
	} else {
		if v != 12 {
			t.Errorf("BinaryTree root should be 12 but is %v", v)
		}
	}
	if lTree, err := r.LeftSubtree(); err != nil {
		t.Error("BinaryTree left subtree exist")
	} else {
		if v, _ := lTree.RootValue(); v != 8 {
			t.Errorf("BinaryTree left subtree root should be 8 but is %v", v)
		}
	}
	if rTree, err := r.RightSubtree(); err != nil {
		t.Error("BinaryTree right subtree exist")
	} else {
		if v, _ := rTree.RootValue(); v != 6 {
			t.Errorf("BinaryTree left subtree root should be 6 but is %v", v)
		}
	}
	if !r.Contains(12) {
		t.Error("BinaryTree should contain 12")
	}
	if !r.Contains(8) {
		t.Error("BinaryTree should contain 8")
	}
	if !r.Contains(6) {
		t.Error("BinaryTree should contain 6")
	}
	if r.Contains(2) {
		t.Error("BinaryTree should not contain 2")
	}

	// check preorder iteration
	preorder := []int{12, 8, 6, 8}
	var iter containers.Iterator
	i := 0
	r.VisitPreorder(func(e interface{}) {
		if e != preorder[i] {
			t.Errorf("Preoder value is %v should be %v", e, preorder[i])
		}
		i++
	})
	i = 0
	iter = r.NewPreorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if e != preorder[i] {
			t.Errorf("Preorder external iterator value is %v should be %v", e, preorder[i])
		}
		i++
	}
	if i != 4 {
		t.Errorf("Preorder external iterator did not complete; stopped with i at %v", i)
	}
	if !iter.Done() {
		t.Error("Preorder external iterator shoult be done")
	}

	// check inorder iteration
	inorder := []int{8, 12, 8, 6}
	i = 0
	r.VisitInorder(func(e interface{}) {
		if e != inorder[i] {
			t.Errorf("Inorder internal iterator value is %v should be %v", e, inorder[i])
		}
		i++
	})
	i = 0
	iter = r.NewInorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if e != inorder[i] {
			t.Errorf("Inorder external iterator value is %v should be %v", e, inorder[i])
		}
		i++
	}
	if i != 4 {
		t.Errorf("Inorder external iterator did not complete; stopped with i at %v", i)
	}
	if !iter.Done() {
		t.Error("Inorder external iterator shoult be done")
	}

	// check postorder iteration
	postorder := []int{8, 8, 6, 12}
	i = 0
	r.VisitPostorder(func(e interface{}) {
		if e != postorder[i] {
			t.Errorf("Postoder value is %v should be %v", e, postorder[i])
		}
		i++
	})
	i = 0
	iter = r.NewPostorderIterator()
	for e, ok := iter.Next(); ok; e, ok = iter.Next() {
		if e != postorder[i] {
			t.Errorf("Postorder external iterator value is %v should be %v", e, postorder[i])
		}
		i++
	}
	if i != 4 {
		t.Errorf("Postorder external iterator did not complete; stopped with i at %v", i)
	}
	if !iter.Done() {
		t.Error("Postorder external iterator shoult be done")
	}

	// make sure a cleared BinaryTree is empty
	r.Clear()
	if !r.Empty() || r.Size() != 0 {
		t.Error("BinaryTree should be empty after Clear()")
	}
}
