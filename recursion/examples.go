// examples.go: This file contains examples of recursive and stack-based
// functions in Go. These functions illustrate the "equivalence" of
// stacks and recursion in algorithms that need to use them.

package recursion

import (
	"containers"
	"fmt"
	"unicode/utf8"
)

//////////////////////////////////////////////////////////////////////////////
// String reversal function

// Reverse returns the reverse of its string argument using a stack.
func Reverse(s string) string {
	stack := containers.NewLinkedStack()
	for _, ch := range s {
		stack.Push(ch)
	}
	result := ""
	for v, err := stack.Pop(); err == nil; v, err = stack.Pop() {
		result += string(v.(rune))
	}
	return result
}

// RecrusiveReverse returns the reverse of its string argument using recursion.
func RecursiveReverse(s string) string {
	if len(s) == 0 {
		return s
	}
	_, size := utf8.DecodeRuneInString(s)
	return Reverse(s[size:]) + s[:size]
}

//////////////////////////////////////////////////////////////////////////////
// Towers of Hanoi
type HanoiState struct {
	moveCount int
	towerA    []byte
	towerB    []byte
	towerC    []byte
}

const ( // the towers
	A = iota
	B
	C
)

// NewHanoiState creates a HanoiState structure an initializes
// it to hold n disks on tower A and no disks on towers B and
// C. In other words, it returns a pointer to the initial
// HanoiState.
func NewHanoiState(n byte) *HanoiState {
	result := new(HanoiState)
	var i byte
	for i = 0; i < n; i++ {
		result.towerA = append(result.towerA, 'a'+i)
	}
	return result
}

// String creates a string representation of a HanoiState
func (s *HanoiState) String() string {
	return fmt.Sprintf("Tower A: %s\nTower B: %s\nTower C: %s\nMoves: %v\n",
		s.towerA, s.towerB, s.towerC, s.moveCount)
}

// Move disk moves the top disk from src tower to dst tower.
func (s *HanoiState) MoveDisk(src, dst int) {
	if src == A {
		if dst == B {
			s.towerB = append(s.towerB, s.towerA[len(s.towerA)-1])
		}
		if dst == C {
			s.towerC = append(s.towerC, s.towerA[len(s.towerA)-1])
		}
		s.towerA = s.towerA[:len(s.towerA)-1]
	} else if src == B {
		if dst == A {
			s.towerA = append(s.towerA, s.towerB[len(s.towerB)-1])
		}
		if dst == C {
			s.towerC = append(s.towerC, s.towerB[len(s.towerB)-1])
		}
		s.towerB = s.towerB[:len(s.towerB)-1]
	} else { // src == C
		if dst == A {
			s.towerA = append(s.towerA, s.towerC[len(s.towerC)-1])
		}
		if dst == B {
			s.towerB = append(s.towerB, s.towerC[len(s.towerC)-1])
		}
		s.towerC = s.towerC[:len(s.towerC)-1]
	}
	s.moveCount++
	//fmt.Printf("%s\n",s)			// uncomment this line to see all the moves
}

// MoveTower transfers n disks from src tower to dst tower using the aux
// tower as an auxiliary in accord with the rules of the game.
func (s *HanoiState) MoveTower(src, dst, aux, n int) {
	if n == 1 {
		s.MoveDisk(src, dst)
	} else {
		s.MoveTower(src, aux, dst, n-1)
		s.MoveDisk(src, dst)
		s.MoveTower(aux, dst, src, n-1)
	}
}

// moveTask stores tower move tasks pushed on a stack in the stack-based
// solution to the Towers of Hanoi problem.
type moveTask struct {
	src, dst, aux, n int
}

// newMoveTask creates and initializes a new moveTask instance.
func newMoveTask(src, dst, aux, n int) *moveTask {
	result := new(moveTask)
	result.src, result.dst, result.aux, result.n = src, dst, aux, n
	return result
}

// MoveTowerStack solves the Towers of Hanoi problem using a stack rather
// than recursion.
func (s *HanoiState) MoveTowerStack(src, dst, aux, n int) {
	stack := containers.NewLinkedStack()
	task := newMoveTask(src, dst, aux, n)
	stack.Push(task)
	for val, err := stack.Pop(); err == nil; val, err = stack.Pop() {
		task = val.(*moveTask)
		if task.n == 1 {
			s.MoveDisk(task.src, task.dst)
		} else {
			stack.Push(newMoveTask(task.aux, task.dst, task.src, task.n-1))
			stack.Push(newMoveTask(task.src, task.dst, task.src, 1))
			stack.Push(newMoveTask(task.src, task.aux, task.dst, task.n-1))
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
// Strings of balanced brackets

// IsBalancedRecursive tests whether s is a string of balanced
// brackets. This function is an interface function that sets things
// up and calls the recursive function isBalancedBrackets() to do the
// real work. Also handles the empty string as a valid argument.
func IsBalancedRecursive(s string) bool {
	if len(s) == 0 {
		return true
	}
	return isBalancedBrackets(NewTokenizer(s))
}

// isBalancedBrackets uses recursion to parse a string of brackets to
// see if they are balanced.
// Uses the grammar B -> [] | [B] | BB for recursive descent parsing
func isBalancedBrackets(current *Tokenizer) bool {
	if current.Char != '[' {
		return false
	}
	current.Next()
	if current.Char == '[' {
		if !isBalancedBrackets(current) {
			return false
		}
	}
	if current.Char != ']' {
		return false
	}
	current.Next()
	if current.Char == '[' {
		return isBalancedBrackets(current)
	}
	return true
}

// IsBalancedStack tests whether s is a string of balaned brackets.
// This function uses a stack to parse a string of brackets.
func IsBalancedStack(s string) bool {
	current := NewTokenizer(s)
	stack := containers.NewLinkedStack()
	for current.Char != '$' {
		switch current.Char {
		case '[':
			stack.Push(current.Char)
		case ']':
			if stack.IsEmpty() {
				return false
			}
			stack.Pop()
		default:
			return false
		}
		current.Next()
	}
	return stack.IsEmpty()
}

//////////////////////////////////////////////////////////////////////
// Factorial functions

// Recursvie factorial computes n! using recursion.
func RecursiveFactorial(n int) int {
	if n < 0 {
		panic("Factorial of a negative number is undefined")
	}
	if n <= 1 {
		return 1
	}
	return n * RecursiveFactorial(n-1)
}

// Factorial computes n! without recursion.
func Factorial(n int) int {
	if n < 0 {
		panic("Factorial of a negative number is undefined")
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

//////////////////////////////////////////////////////////////////////////////
// Slice search functions--tail recursion.

// RecursiveSearch looks through a slice for a value using recursion.
func RecursiveSearch(slice []int, value int) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	if slice[0] == value {
		return true
	}
	return RecursiveSearch(slice[1:], value)
}

// Search looks through a slice for a value without recursion.
// Note that since the recursive version is tail recursive, this function does
// not need a stack.
func Search(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
