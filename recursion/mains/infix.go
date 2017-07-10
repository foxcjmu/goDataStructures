// infix.go: Ask users for an infix expression, evaluate it using a recursive evaluation
// algorithm and then a stack-based algorithm, and print the results, or error messages
// if there are any. The program ends when the user types an empty line.
//
// Version: 9/12
// Author: C. Fox

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"recursion"
)

func main() {
	fmt.Printf("Infix expression evaluator.\n")

	console := bufio.NewReader(os.Stdin) // where we read from

	for {
		var (
			exp []byte // infix expression typed by the user
			err error  // in case something goes wrong on input
		)

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

		if v, err := recursion.EvalInfixRecursive(string(exp)); err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%d\n", v)
		}
		if v, err := recursion.EvalInfixStack(string(exp)); err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%d\n", v)
		}
	}
	fmt.Printf("Bye\n")
}
