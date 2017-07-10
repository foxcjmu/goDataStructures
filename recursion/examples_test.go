package recursion

import "fmt"
import "testing"

func TestReverse(t *testing.T) {
	if "abçdef" != Reverse("fedçba") {
		t.Errorf("Reverse is broken")
	}
}
func TestRecursiveReverse(t *testing.T) {
	if "abçdef" != RecursiveReverse("fedçba") {
		t.Errorf("Reverse is broken")
	}
}

func TestHanoi(t *testing.T) {
	var s *HanoiState
	s = NewHanoiState(8)
	s.MoveTower(A, C, B, 8)
	if s.moveCount != 255 {
		t.Errorf(fmt.Sprintf("Recursive Hanoi is broken"))
	}
}

func TestHanoiStack(t *testing.T) {
	var s *HanoiState
	s = NewHanoiState(8)
	s.MoveTowerStack(A, C, B, 8)
	if s.moveCount != 255 {
		t.Errorf(fmt.Sprintf("Stack-based Hanoi is broken"))
	}
}

func testBalancedBracketsFunction(t *testing.T, isBalanced func(string) bool) {
	if !isBalanced("") {
		t.Errorf("%v fails on empty string", isBalanced)
	}
	if !isBalanced("[]") {
		t.Errorf("%vfails on []", isBalanced)
	}
	if !isBalanced("[][]") {
		t.Errorf("%v fails on [][]", isBalanced)
	}
	if !isBalanced("[][][]") {
		t.Errorf("%v fails on [][][]", isBalanced)
	}
	if !isBalanced("[[]]") {
		t.Errorf("%v fails on [[]]", isBalanced)
	}
	if !isBalanced("[[[]][[][]]]") {
		t.Errorf("%v fails on [[[]][[][]]]", isBalanced)
	}
	if isBalanced("[") {
		t.Errorf("%v fails on [", isBalanced)
	}
	if isBalanced("]") {
		t.Errorf("%v fails on ]", isBalanced)
	}
	if isBalanced("[[]") {
		t.Errorf("%v fails on [[]", isBalanced)
	}
	if isBalanced("[[]") {
		t.Errorf("%v fails on []]", isBalanced)
	}
	if isBalanced("[[[[][[]]]]") {
		t.Errorf("%v fails on [[[][[]]]]", isBalanced)
	}
}

func TestFactorail(t *testing.T) {
	if 1 != RecursiveFactorial(0) {
		t.Errorf("RecursiveFactorial is broken")
	}
	if 479001600 != RecursiveFactorial(12) {
		t.Errorf("RecursiveFactorial is broken")
	}
	if 1 != Factorial(0) {
		t.Errorf("Factorial is broken")
	}
	if 479001600 != Factorial(12) {
		t.Errorf("Factorial is broken")
	}
}

func TestBalancedBrackets(t *testing.T) {
	testBalancedBracketsFunction(t, IsBalancedRecursive)
	testBalancedBracketsFunction(t, IsBalancedStack)
}

func TestRecursiveSearch(t *testing.T) {
	slice := make([]int, 1000)
	for i := 0; i < 2000; i += 2 {
		slice[i/2] = i
	}
	if RecursiveSearch(nil, -1) {
		t.Errorf("Recursive search found 0 in nil")
	}
	if RecursiveSearch(slice, -1) {
		t.Errorf("Recursive search found -1")
	}
	if RecursiveSearch(slice, 1) {
		t.Errorf("Recursive search found 1")
	}
	if RecursiveSearch(slice, 2000) {
		t.Errorf("Recursive search found 2000")
	}
	if !RecursiveSearch(slice, 0) {
		t.Errorf("Recursive search did not find 0")
	}
	if !RecursiveSearch(slice, 1000) {
		t.Errorf("Recursive search did not find 1000")
	}
	if !RecursiveSearch(slice, 1998) {
		t.Errorf("Recursive search did not find 1998")
	}

	if Search(nil, -1) {
		t.Errorf("Recursive search found 0 in nil")
	}
	if Search(slice, -1) {
		t.Errorf("Recursive search found -1")
	}
	if Search(slice, 1) {
		t.Errorf("Recursive search found 1")
	}
	if Search(slice, 2000) {
		t.Errorf("Recursive search found 2000")
	}
	if !Search(slice, 0) {
		t.Errorf("Recursive search did not find 0")
	}
	if !Search(slice, 1000) {
		t.Errorf("Recursive search did not find 1000")
	}
	if !Search(slice, 1998) {
		t.Errorf("Recursive search did not find 1998")
	}
}
