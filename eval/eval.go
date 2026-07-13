// Package eval evaluates arithmetic expressions.
//
// LAB 3 — the feature to build:
// Right now Eval is a stub. It parses a single bare number but does not
// understand operators, parentheses, or unary minus. Your task is to make
// Eval a full expression evaluator so that ./check.sh passes.
//
// Target grammar (standard precedence, left-associative):
//
//	expr    = term   (('+' | '-') term)*
//	term    = factor (('*' | '/') factor)*
//	factor  = '-' factor | primary
//	primary = number | '(' expr ')'
//
// Numbers may be integers or decimals. Whitespace is insignificant.
// Malformed input (empty, unbalanced parens, dangling operators, garbage)
// must return a non-nil error rather than a bogus value.
package eval

import (
	"fmt"
	"strconv"
	"strings"
)

// Eval evaluates an arithmetic expression and returns its value.
//
// TODO(lab3): implement operator precedence, parentheses, and unary minus.
func Eval(expr string) (float64, error) {
	s := strings.TrimSpace(expr)
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("eval: not implemented for %q", expr)
}
