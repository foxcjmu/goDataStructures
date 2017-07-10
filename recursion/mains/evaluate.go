// evaluate.go: Ask users for an infix or prefix expression, translate it
// to postfix, and print the translation, then use a stack-based algorithm
// to evaluate the postfix expression and print the result. If an expression
// is ill-formed, an error messages is printed but no value. The program
// ends when the user types an empty line.
//
// Version: 10/13
// Author: C. Fox

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"recursion"
)

func main() {
	fmt.Printf("Infix or prefix expression evaluator.\n")
	console := bufio.NewReader(os.Stdin) // where we read from

	for {
		var (
			exp        []byte // expression typed by the user
			postfixExp string // expression translated to postfix
			err        error  // in case something goes wrong
			value      int    // postfix expression result
		)

		// prompt the user, collect the reply
		fmt.Printf("> ")
		if exp, err = console.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				fmt.Printf("IO Error: %v\n", err)
			}
			break
		}
		exp = exp[:len(exp)-1] // chop off the '\n'
		if 0 == len(exp) {
			break
		}

		// translate a prefix or infix expression to postfix
		switch exp[0] {
		case '+', '-', '*', '/', '%':
			postfixExp, err = recursion.Prefix2Postfix(string(exp))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '(':
			postfixExp, err = recursion.Infix2Postfix(string(exp))
		default:
			err = errors.New("Illegal character")
		}

		// either print an error message or evaluate the expression and print it
		if err == nil {
			value, err = recursion.EvalPostfixStack(postfixExp)
		}
		if err == nil {
			fmt.Printf("Evaluating %s -> %d\n", postfixExp, value)
		} else {
			fmt.Printf("%v\n", err)
		}
	}
	fmt.Printf("Bye\n")
}
