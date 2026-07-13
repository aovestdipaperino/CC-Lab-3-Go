package eval

import (
	"math"
	"strings"
	"testing"
)

func mustEval(t *testing.T, expr string) float64 {
	t.Helper()
	got, err := Eval(expr)
	if err != nil {
		t.Fatalf("Eval(%q) returned error: %v", expr, err)
	}
	return got
}

func checkTable(t *testing.T, cases map[string]float64) {
	t.Helper()
	for expr, want := range cases {
		got, err := Eval(expr)
		if err != nil {
			t.Errorf("Eval(%q) returned error: %v", expr, err)
			continue
		}
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", expr, got, want)
		}
	}
}

// ---- Step 1: core arithmetic --------------------------------------------

func TestStep1Core(t *testing.T) {
	checkTable(t, map[string]float64{
		"42":              42,
		"2 + 3 * 4":       14,
		"(2 + 3) * 4":     20,
		"10 / 4":          2.5,
		"2 * -4":          -8,
		"-3 + 5":          2,
		"1 + 2 + 3 + 4":   10,
		"2 * (3 + 4) - 1": 13,
		"100 / 5 / 2":     10,
		"3.5 * 2":         7,
		"-(2 + 3)":        -5,
		"((1))":           1,
	})
	// Malformed input must be rejected.
	for _, bad := range []string{"", "1 +", "(1 + 2", "2 * * 3", "abc"} {
		if _, err := Eval(bad); err == nil {
			t.Errorf("Eval(%q) expected an error, got nil", bad)
		}
	}
}

// ---- Step 2: modulo operator --------------------------------------------

func TestStep2Modulo(t *testing.T) {
	checkTable(t, map[string]float64{
		"10 % 3":    1,
		"17 % 5":    2,
		"20 % 4":    0,
		"2 + 9 % 4": 3, // % binds like * and /
		"9 % 4 * 3": 3,
	})
}

// ---- Step 3: power operator (right-associative, binds tighter than * /) --

func TestStep3Power(t *testing.T) {
	checkTable(t, map[string]float64{
		"2 ^ 10":      1024,
		"2 ^ 3 ^ 2":   512, // right associative: 2^(3^2)
		"2 * 3 ^ 2":   18,  // ^ before *
		"(1 + 1) ^ 8": 256,
		"9 ^ 0.5":     3,
	})
}

// ---- Step 4: built-in functions -----------------------------------------

func TestStep4Functions(t *testing.T) {
	checkTable(t, map[string]float64{
		"sqrt(16)":      4,
		"abs(-5)":       5,
		"max(3, 7)":     7,
		"min(4, 2)":     2,
		"pow(2, 10)":    1024,
		"sqrt(2 * 8)":   4,
		"max(1 + 1, 3)": 3,
		"abs(2 - 9)":    7,
	})
}

// ---- Step 5: error messages carry the column of the offending character --

func TestStep5Errors(t *testing.T) {
	cases := map[string]string{
		"2 @ 3":  "column 3",
		"1 # 2":  "column 3",
		"12 ~ 5": "column 4",
	}
	for expr, want := range cases {
		_, err := Eval(expr)
		if err == nil {
			t.Errorf("Eval(%q) expected an error, got nil", expr)
			continue
		}
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Eval(%q) error = %q, want it to contain %q", expr, err.Error(), want)
		}
	}
}
