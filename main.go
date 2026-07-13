// Command calc evaluates an arithmetic expression given as its argument.
//
//	calc "2 + 3 * 4"   ->   14
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"calc/eval"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `usage: calc "<expression>"`)
		os.Exit(2)
	}
	expr := strings.Join(os.Args[1:], " ")
	v, err := eval.Eval(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(strconv.FormatFloat(v, 'f', -1, 64))
}
