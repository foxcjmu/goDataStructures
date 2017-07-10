// translate.go: This file contains recursive and stack-based algorithms for translating
// between simple prefix, infix, and postfix expressions held in strings. These expressions
// must have only the operators +, -, *, /, and % (with their usual meanings in integer
// arithemetic), and operands that are one digit long. There are no negative operands.
// Infix expressions may have parentheses. No white space is allowed in expressions.

package recursion

import (
	"containers"
	"errors"
	"fmt"
	//"strings"
)

//////////////////////////////////////////////////////////////////////////
// Prefix2: Translate prefix expressions to infix or postfix.

// Prefix2InfixRecursive uses recursion to parse and translate a prefix
// expression to a fuly parenthesized infix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in infix form and nil
func Prefix2InfixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := prefix2otherfix(current, "infix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// Prefix2PostfixRecursive uses recursion to parse and translate a prefix
// expression to a postfix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in postfix form and nil
func Prefix2PostfixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := prefix2otherfix(current, "postfix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// prefix2otherfix is a private function that recursively parses and translates
// a prefix expression provided by a Tokenizer into a either a postfix or
// infix expression as indicated by the fixity argument.
// pre: fixity must be "prefix" or "infix" and s is well-formed
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: Handle the case of a single digit as a special case. Otherwise,
// remember the operator and call prefix2otherfix recursively twice to evaluate
// the two operand expressions.
func prefix2otherfix(current *Tokenizer, fixity string) (string, error) {
	if current.Char == '$' {
		return "", errors.New("Missing argument")
	}

	// handle the case of a single digit
	if isDigit(current.Char) {
		result := string(current.Char)
		current.Next()
		return result, nil
	}

	// handle the case of an operator followed by two expressions
	op := current.Char
	current.Next()
	leftArg, err := prefix2otherfix(current, fixity)
	if err != nil {
		return "", err
	}
	rightArg, err := prefix2otherfix(current, fixity)
	if err != nil {
		return "", err
	}
	if fixity == "infix" {
		return "(" + leftArg + string(op) + rightArg + ")", nil
	}
	if fixity == "postfix" {
		return leftArg + rightArg + string(op), nil
	}
	panic("Bad fixity value")
}

// Prefix2Infix uses as stack to parse and translate a prefix
// expression to a fuly parenthesized infix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in infix form and nil
func Prefix2Infix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return prefix2other(current, "infix")
}

// Prefix2Postfix uses as stack to parse and translate a prefix
// expression to a postfix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in postifx form and nil
func Prefix2Postfix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return prefix2other(current, "postfix")
}

// prefix2other uses a stack to parse and translate a prefix expression
// into a infix or postfix expression.
// pre: Expression s is well formed; fixity is either "infix" or "postfix"
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: Push all operators on the opStack. If a digit is encountered
// and the top of the opStack is an operator, then push the digit on the
// expStack and a special 'e' marker on the opStack to indicate that the
// operator has a left argument ex[pression on the expStack. Otherwise,
// as long as there is an 'e' marker on the opStack, use the current
// expression as the rightArg, pop the top expression from the expStack
// as the leftArg, pop both the 'e' marker and the operator under it, form
// a translated expression using the leftArg, the rightArg, and the op (based
// on the fixity), and save the resulting expression in the rightArg.
// When there are no more 'e' markers on the opStack, push the rightArg on
// the expStack and a 'e' on the opStack. At the end of the input
// expression, there should be just an 'e' on the opStack and the final
// expression on the expStack.
func prefix2other(current *Tokenizer, fixity string) (string, error) {
	opStack := containers.NewLinkedStack()
	expStack := containers.NewLinkedStack()

	// process the entire string unless an error is encountered
	for current.Char != '$' {
		if isOperator(current.Char) {
			opStack.Push(current.Char)
		} else if isDigit(current.Char) {
			rightArg := string(current.Char)
			op, err := opStack.Top()
			for err == nil && op == 'e' {
				opStack.Pop()
				if op, err = opStack.Pop(); err != nil {
					return "", errors.New("Missing operator")
				}
				var leftArg string // argument from the stack
				if exp, err := expStack.Pop(); err != nil {
					return "", errors.New("Missing left argument")
				} else {
					leftArg = exp.(string)
				}
				switch {
				case fixity == "infix":
					rightArg = "(" + leftArg + string(op.(byte)) + rightArg + ")"
				case fixity == "postfix":
					rightArg = leftArg + rightArg + string(op.(byte))
				default:
					panic("Bad fixity value")
				}
				op, err = opStack.Top()
			}
			expStack.Push(rightArg)
			opStack.Push('e')
		} else {
			return "", errors.New(fmt.Sprintf("Illegal character %v", current.Char))
		}
		current.Next()
	}

	// if all is well, v should be on the opStack and the result on the expStack
	if op, err := opStack.Pop(); err != nil || op != 'e' {
		return "", errors.New("Missing argument")
	}
	if !opStack.IsEmpty() {
		return "", errors.New("Missing argument")
	}
	result, err := expStack.Pop()
	if err != nil {
		return "", errors.New("Missing argument")
	}
	if !expStack.IsEmpty() {
		return "", errors.New("Too many arguments")
	}
	return result.(string), nil
}

//////////////////////////////////////////////////////////////////////////
// Infix2: Translate infix expressions to prefix or postfix.

// Infix2PrefixRecursive uses recursion to translate an infix expression
// to preifx.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in prefix form and nil
func Infix2PrefixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := infix2otherfix(current, "prefix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// Infix2PostfixRecursive uses recursion to translate an infix expression
// to postifx.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in postfix form and nil
func Infix2PostfixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := infix2otherfix(current, "postfix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// infix2otherfix is a private function to translate an infix expression
// to postfix or prefix (as indicated by the fixity parameter) using
// recursion.
// pre: fixity must be "prefix" or "infix" and s is well-formed
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: Transform the infix expression from left to right, with recursive
// calls to handle parenthesized sub-expressions.
func infix2otherfix(current *Tokenizer, fixity string) (result string, err error) {
	if current.Char == '(' {
		current.Next()
		if result, err = infix2otherfix(current, fixity); err != nil {
			return
		}
		if current.Char != ')' {
			return "", errors.New("Missing right parenthesis")
		}
	} else if isDigit(current.Char) {
		result = string(current.Char)
	} else {
		return "", errors.New("Missing left argument")
	}
	current.Next()

	// apply the next operator to the following operand as long as there is one
	for isOperator(current.Char) {
		var rightArg string
		op := string(current.Char)
		current.Next()
		if current.Char == '(' {
			current.Next()
			if rightArg, err = infix2otherfix(current, fixity); err != nil {
				return
			}
			if current.Char != ')' {
				return "", errors.New("Missing right parenthesis")
			}
		} else if isDigit(current.Char) {
			rightArg = string(current.Char)
		} else {
			return "", errors.New("Missing right argument")
		}
		current.Next()
		switch {
		case fixity == "prefix":
			result = op + result + rightArg
		case fixity == "postfix":
			result += rightArg + op
		default:
			panic("Bad fixity value")
		}
	}

	// when we are out of operators, we are done
	return result, nil
}

// Infix2Prefix translates an infix expression to prefix using a stack.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the prefix expression value and nil
func Infix2Prefix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return infix2other(current, "prefix")
}

// Infix2Postfix translates an infix expression to postfix using a stack.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the postfix expression value and nil
func Infix2Postfix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return infix2other(current, "postfix")
}

// infix2other is a private function to translate an infix expression
// to postfix or prefix (as indicated by the fixity parameter) using
// a stack.
// pre: fixity must be "prefix" or "infix" and s is well-formed
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: Push all operators and left parens on the opStack, digits on
// the expStack, and check right parens against left parens on the top of
// the opStack. After pushing a digit or checking a right parens, apply the
// top operator on the opStack to the top two operands on the expStack,
// and push the result on the expStack, as long as the opStack has an
// operator on it. The result should be in the expStack at the end.
func infix2other(current *Tokenizer, fixity string) (string, error) {
	opStack := containers.NewLinkedStack()
	expStack := containers.NewLinkedStack()
	for current.Char != '$' {
		if isOperator(current.Char) || current.Char == '(' {
			opStack.Push(current.Char)
		} else {
			if isDigit(current.Char) {
				expStack.Push(string(current.Char))
			} else if current.Char == ')' {
				if op, err := opStack.Top(); err != nil || op.(byte) != '(' {
					return "", errors.New("Missing left parenthesis")
				}
				opStack.Pop()
			} else {
				return "", errors.New("Illegal character in expression")
			}
			op, err := opStack.Top()
			if err == nil && isOperator(op.(byte)) {
				opStack.Pop()
				rightArg, err := expStack.Pop()
				if err != nil {
					return "", errors.New("Missing right argument")
				}
				leftArg, err := expStack.Pop()
				if err != nil {
					return "", errors.New("Missing left argument")
				}
				switch {
				case fixity == "prefix":
					expStack.Push(string(op.(byte)) + leftArg.(string) + rightArg.(string))
				case fixity == "postfix":
					expStack.Push(leftArg.(string) + rightArg.(string) + string(op.(byte)))
				default:
					panic("Bad fixity value")
				}
			}
		}
		current.Next()
	}
	if !opStack.IsEmpty() {
		return "", errors.New("Missing argument")
	}
	result, err := expStack.Pop()
	if err != nil {
		return "", errors.New("Missing expression")
	}
	if !expStack.IsEmpty() {
		return "", errors.New("Too many arguments")
	}
	return result.(string), nil
}

//////////////////////////////////////////////////////////////////////////
// Postfix2: Translate a postfix expression to prefix or infix.

// Postfix2PrefixRecursive translates a postfix expression to prefix
// using recursion.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in prefix form and nil
func Postfix2PrefixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := postfix2otherfix(current, "prefix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// Postfix2InfixRecursive translates a postfix expression to infix
// using recursion.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in infix form and nil
func Postfix2InfixRecursive(s string) (string, error) {
	current := NewTokenizer(s)
	result, err := postfix2otherfix(current, "infix")
	if err == nil && current.Char != '$' {
		return "", errors.New("Extra characters at the end of the expression")
	}
	return result, err
}

// infix2otherfix is a private function to translate an infix expression
// to postfix or prefix (as indicated by the fixity parameter) using
// recursion.
// pre: fixity must be "prefix" pr "postfix" and s is well-formed
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: Transform the infix expression from left to right, with recursive
// calls to handle parenthesized sub-expressions.
// strategy: The first character must be a digit, so remember it as the leftArg.
// Assume the next digit is the rightArg. If the character after that is not an
// operator, then back up one character and call postfix2otherfix recursively to
// evaluate the right operand expression. When the operator is finally found,
// apply it to leftArg and rightArg and leave the result in leftArg. Look for
// another digit as the start of a possible following expression and repeat.
func postfix2otherfix(current *Tokenizer, fixity string) (result string, err error) {
	if !isDigit(current.Char) {
		return "", errors.New("Missing argument")
	}
	leftArg := string(current.Char)
	current.Next()
	for isDigit(current.Char) {
		rightArg := string(current.Char)
		current.Next()
		if isDigit(current.Char) {
			current.Last()
			for {
				if rightArg, err = postfix2otherfix(current, fixity); err != nil {
					return "", err
				}
				if !isDigit(current.Char) {
					break
				}
			}
		}
		if current.Char == '$' {
			return "", errors.New("Missing operator")
		}
		switch {
		case fixity == "prefix":
			leftArg = string(current.Char) + leftArg + rightArg
		case fixity == "infix":
			leftArg = "(" + leftArg + string(current.Char) + rightArg + ")"
		default:
			panic("Bad fixity value")
		}
		current.Next()
	}
	return leftArg, nil
}

// Postfix2Prefix uses as stack to translate a postfix expression to
// a prefix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in prefix form and nil
func Postfix2Prefix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return postfix2other(current, "prefix")
}

// Postfix2Infix uses as stack to translate a postfix expression
// to a fuly parenthesized infix expression.
// pre: Expression s is well formed
// pre violation: return "" and an error indication
// normal return: the expression in infix form and nil
func Postfix2Infix(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("Missing argument")
	}
	current := NewTokenizer(s)
	return postfix2other(current, "infix")
}

// postfix2otherfix is a private function to translate a postfix expression
// to prefix or infix (as indicated by the fixity parameter) using a stack.
// pre: fixity must be "prefix" or "infix" and s is well-formed
// pre violation: if s is not well-formed: the empty string and an error
//		if fixity is unrecognized: panic
// normal return: translated expression and nil
// strategy: put all operands in the stack, and whenever an operator is
// encountered, apply it to the top two values in the stack and push the
// result back on the stack. At the end, the stack should contain the result.
func postfix2other(current *Tokenizer, fixity string) (string, error) {
	stack := containers.NewLinkedStack() // expressions during evaluation
	for current.Char != '$' {
		if isDigit(current.Char) {
			stack.Push(string(current.Char))
		} else {
			rightArg, err := stack.Pop()
			if err != nil {
				return "", errors.New("Missing right argument")
			}
			leftArg, err := stack.Pop()
			if err != nil {
				return "", errors.New("Missing left argument")
			}
			switch {
			case fixity == "prefix":
				stack.Push(string(current.Char) + leftArg.(string) + rightArg.(string))
			case fixity == "infix":
				stack.Push("(" + leftArg.(string) + string(current.Char) + rightArg.(string) + ")")
			default:
				panic("Bad fixity value")
			}
		}
		current.Next()
	}
	result, err := stack.Pop()
	if err != nil {
		return "", errors.New("Missing expression")
	}
	if !stack.IsEmpty() {
		return "", errors.New("Too many arguments")
	}
	return result.(string), nil
}
