package eval

import (
	"math"
	"testing"
)

// TestEval is the spec: these expressions must evaluate to these values.
func TestEval(t *testing.T) {
	cases := []struct {
		expr string
		want float64
	}{
		{"42", 42},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"10 / 4", 2.5},
		{"2 * -4", -8},
		{"-3 + 5", 2},
		{"1 + 2 + 3 + 4", 10},
		{"2 * (3 + 4) - 1", 13},
		{"100 / 5 / 2", 10},
		{"3.5 * 2", 7},
		{"-(2 + 3)", -5},
		{"((1))", 1},
	}
	for _, c := range cases {
		got, err := Eval(c.expr)
		if err != nil {
			t.Errorf("Eval(%q) returned error: %v", c.expr, err)
			continue
		}
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", c.expr, got, c.want)
		}
	}
}

// TestEvalErrors requires malformed input to be rejected, not silently accepted.
func TestEvalErrors(t *testing.T) {
	bad := []string{"", "1 +", "(1 + 2", "2 * * 3", "abc"}
	for _, expr := range bad {
		if _, err := Eval(expr); err == nil {
			t.Errorf("Eval(%q) expected an error, got nil", expr)
		}
	}
}
