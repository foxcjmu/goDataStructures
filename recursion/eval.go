// eval.go: This file contains recursive and stack-based algorithms for evaluating simple
// prefix, infix, and postfix expressions held in strings. These expressions must have
// only the operators +, -, *, /, and % (with their usual meanings in integer arithemetic),
// and operands that are one digit long. There are no negative operands.  Infix
// expressions may have parentheses. No white space is allowed in expressions.

package recursion

import (
	"containers"
	"errors"
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////
// Expression evaluation utility functions.

// Determine whether a character is a digit
func isDigit(ch byte) bool {
	return strings.ContainsRune("0123456789", rune(ch))
}

// Determine whether a character is an operator
func isOperator(ch byte) bool {
	return strings.ContainsRune("+-*/%", rune(ch))
}

// Apply an operator designated by op to two arguments
func applyOperator(op byte, leftArg, rightArg int) (int, error) {
	switch op {
	case '+':
		return leftArg + rightArg, nil
	case '-':
		return leftArg - rightArg, nil
	case '*':
		return leftArg * rightArg, nil
	case '/':
		return leftArg / rightArg, nil
	case '%':
		return leftArg % rightArg, nil
	default:
		return 0, errors.New(fmt.Sprintf("Bad character %c", op))
	}
	panic("Reached impossible spot")
}

//////////////////////////////////////////////////////////////////////////
// Prefix: These functions evaluate a prefix expression held in a string.

// EvalPrefixRecursive uses recursion to parse and evaluate a prefix
// expression.
// Pre: The expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
func EvalPrefixRecursive(s string) (int, error) {
	current := NewTokenizer(s)
	result, err := evalPrefix(current)
	if err == nil && current.Char != '$' {
		return 0, errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// evalPrefix is a private function that recursively parses and evaluates
// a prefix expression provided by a Tokenizer.
// Strategy: Handle the case of a single digit as a special case. Otherwise,
// remember the operator and call evalPrefix recursively twice to evaluate the
// two operand expressions.
func evalPrefix(current *Tokenizer) (int, error) {
	if current.Char == '$' {
		return 0, errors.New("Missing argument")
	}

	// handle the case of a single digit
	if isDigit(current.Char) {
		result := int(current.Char - '0')
		current.Next()
		return result, nil
	}

	// handle the case of an operator followed by two expressions
	op := current.Char
	current.Next()
	leftArg, err := evalPrefix(current)
	if err != nil {
		return 0, err
	}
	rightArg, err := evalPrefix(current)
	if err != nil {
		return 0, err
	}
	return applyOperator(op, leftArg, rightArg)
}

// EvalPrefixStack uses a stack to parse and evaluate a prefix expression.
// Pre: The expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
// Strategy: Push all operators on the opStack. If a digit is encountered and
// the top of the opStack is an operator, then push the number on the valStack
// and a special v marker on the opStack. Otherwise, as long as there is a v
// marker on the opStack, use the current value as the rightArg, pop the top
// value from the valStack as the leftArg, and pop both the v marker and the
// operator under it and apply it to the leftArg and rightArg, leaving the
// result in the rightArg. When there are no more v markers on the opStack,
// push the rightArg on the valStack and a v on the opStack. At the end of
// the expression, there should be just a v on the opStack and the final value
// in the valStack.
func EvalPrefixStack(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	opStack := containers.NewLinkedStack()
	valStack := containers.NewLinkedStack()

	// process the entire string unless an error is encountered
	for current.Char != '$' {
		switch {
		case isOperator(current.Char):
			opStack.Push(current.Char)
		case isDigit(current.Char):
			rightArg := int(current.Char - '0')
			op, err := opStack.Top()
			for err == nil && op == 'v' {
				opStack.Pop()
				if op, err = opStack.Pop(); err != nil {
					return 0, errors.New("Missing operator")
				}
				var leftArg int // argument from the stack
				if elem, err := valStack.Pop(); err != nil {
					return 0, errors.New("Missing left argument")
				} else {
					leftArg = elem.(int)
				}
				if rightArg, err = applyOperator(op.(byte), leftArg, rightArg); err != nil {
					return 0, err
				}
				op, err = opStack.Top()
			}
			valStack.Push(rightArg)
			opStack.Push('v')
		default:
			return 0, errors.New(fmt.Sprintf("Illegal character %v", current.Char))
		}
		current.Next()
	}

	// if all is well, v should be on the opStack and the result on the valStack
	if op, err := opStack.Pop(); err != nil || op != 'v' {
		return 0, errors.New("Missing argument")
	}
	if !opStack.IsEmpty() {
		return 0, errors.New("Missing argument")
	}
	result, err := valStack.Pop()
	if err != nil {
		return 0, errors.New("Missing argument")
	}
	if !valStack.IsEmpty() {
		return 0, errors.New("Too many arguments")
	}
	return result.(int), nil
}

//////////////////////////////////////////////////////////////////////////
// Infix: These functions evaluate an infix expression held in a string.

// EvalInfixRecursive uses recursion to parse and evaluate an infix
// expression.
// Pre: Expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
func EvalInfixRecursive(s string) (int, error) {
	current := NewTokenizer(s)
	result, err := evalInfix(current)
	if err == nil && current.Char != '$' {
		return 0, errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// evalInfix is a private function to parse and evaluate an infix expression
// using recursion. The strategy is to appy operators to operands from left
// to right, with recursive calls to handle parenthesized sub-expressions.
func evalInfix(current *Tokenizer) (result int, err error) {
	// get the left argument first
	var leftArg int
	if current.Char == '(' {
		current.Next()
		leftArg, err = evalInfix(current)
		if err != nil {
			return 0, err
		}
		if current.Char != ')' {
			return 0, errors.New("Missing right parenthesis")
		}
	} else if isDigit(current.Char) {
		leftArg = int(current.Char - '0')
	} else {
		return 0, errors.New("Missing left argument")
	}
	current.Next()

	// apply the next operator to the following operand as long as there is one
	for isOperator(current.Char) {
		op, rightArg := current.Char, 0
		current.Next()
		if current.Char == '(' {
			current.Next()
			rightArg, err = evalInfix(current)
			if err != nil {
				return 0, err
			}
			if current.Char != ')' {
				return 0, errors.New("Missing right parenthesis")
			}
		} else if isDigit(current.Char) {
			rightArg = int(current.Char - '0')
		} else {
			return 0, errors.New("Missing right argument")
		}
		current.Next()
		leftArg, err = applyOperator(op, leftArg, rightArg)
		if err != nil {
			return 0, err
		}
	}

	// when we are out of operators, we are done and the result is the leftArg
	return leftArg, nil
}

// EvalInfixStack parses and evaluates an infix expression using a stack.
// Pre: Expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
// Strategy: Push all operators and left parens on the opStack, digits on the valueStack,
// and check right parens against left parens on the top of the opStack. After pushing
// a digit or checking a right parens, apply the top operator on the opStack to the top
// two operands on the valueStack, and push the result on the valueStack, as long as the
// opStack has an operator on it. The result should be in the valueStack at the end.
func EvalInfixStack(s string) (int, error) {
	current := NewTokenizer(s)
	opStack := containers.NewLinkedStack()
	valueStack := containers.NewLinkedStack()
	for current.Char != '$' {
		if isOperator(current.Char) || current.Char == '(' {
			opStack.Push(current.Char)
		} else {
			if isDigit(current.Char) {
				valueStack.Push(int(current.Char - '0'))
			} else if current.Char == ')' {
				if op, err := opStack.Top(); err != nil || op.(byte) != '(' {
					return 0, errors.New("Missing left parenthesis")
				}
				opStack.Pop()
			} else {
				return 0, errors.New("Illegal character in expression")
			}
			op, err := opStack.Top()
			if err == nil && isOperator(op.(byte)) {
				opStack.Pop()
				rightArg, err := valueStack.Pop()
				if err != nil {
					return 0, errors.New("Missing right argument")
				}
				leftArg, err := valueStack.Pop()
				if err != nil {
					return 0, errors.New("Missing left argument")
				}
				if value, err := applyOperator(op.(byte), leftArg.(int), rightArg.(int)); err == nil {
					valueStack.Push(value)
				} else {
					return 0, err
				}
			}
		}
		current.Next()
	}
	if !opStack.IsEmpty() {
		return 0, errors.New("Missing argument")
	}
	result, err := valueStack.Pop()
	if err != nil {
		return 0, errors.New("Missing expression")
	}
	if !valueStack.IsEmpty() {
		return 0, errors.New("Too many arguments")
	}
	return result.(int), nil
}

//////////////////////////////////////////////////////////////////////////
// Postfix: These functions evaluate a postfix expression held in a string.

// EvalPostfixRecursive uses recursion to parse and evaluate a postfix
// expression.
// Pre: The expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
func EvalPostfixRecursive(s string) (int, error) {
	current := NewTokenizer(s)
	result, err := evalPostfix(current)
	if err == nil && current.Char != '$' {
		return 0, errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// evalPostfix is a private function that recursively parses and evaluates
// a postfix expression provided by a Tokenizer.
// Strategy: The first character must be a digit, so remember it as the leftArg.
// Assume the next digit is the rightArg. If the character after that is not an
// operator, then back up one character and call evalPostfix recursively to
// evaluate the right operand expression. When the operator is finally found,
// apply it to leftArg and rightArg and leave the result in leftArg. Look for
// another digit as the start of a possible following expression and repeat.
func evalPostfix(current *Tokenizer) (resul int, err error) {
	if !isDigit(current.Char) {
		return 0, errors.New("Missing argument")
	}
	leftArg := int(current.Char - '0')
	current.Next()
	for isDigit(current.Char) {
		rightArg := int(current.Char - '0')
		current.Next()
		if current.Char == '$' {
			return 0, errors.New("Missing operator")
		}
		if isDigit(current.Char) {
			current.Last()
			for {
				rightArg, err = evalPostfix(current)
				if err != nil {
					return 0, err
				}
				if !isDigit(current.Char) {
					break
				}
			}
		}
		leftArg, err = applyOperator(current.Char, leftArg, rightArg)
		if err != nil {
			return 0, err
		}
		current.Next()
	}
	return leftArg, nil
}

// EvalPostfixStack uses a stack to parse and evaluate a postfix expression.
// Pre: The expression in s is well formed
// Pre violation: return 0 and an error indication
// Normal return: the expression value and nil
// Strategy: Put all operands in the stack, and whenever an operator is
// encountered, apply it to the top two values in the stack and push the
// result back on the stack. At the end, the stack should contain the result.
func EvalPostfixStack(s string) (int, error) {
	current := NewTokenizer(s)
	stack := containers.NewLinkedStack()
	for current.Char != '$' {
		if isDigit(current.Char) {
			stack.Push(int(current.Char - '0'))
		} else {
			rightArg, err := stack.Pop()
			if err != nil {
				return 0, errors.New("Missing right argument")
			}
			leftArg, err := stack.Pop()
			if err != nil {
				return 0, errors.New("Missing left argument")
			}
			value, err := applyOperator(current.Char, leftArg.(int), rightArg.(int))
			if err == nil {
				stack.Push(value)
			} else {
				return 0, err
			}
		}
		current.Next()
	}
	result, err := stack.Pop()
	if err != nil {
		return 0, errors.New("Missing expression")
	}
	if !stack.IsEmpty() {
		return 0, errors.New("Too many arguments")
	}
	return result.(int), nil
}
