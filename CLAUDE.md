# calc — Go

A command-line arithmetic evaluator.

```
calc "2 + 3 * 4"   ->   14
```

`eval.Eval` in `eval/eval.go` is a stub: it parses a single bare number and
errors on everything else. Build it up over five steps, each harder than the
last. `./check.sh` (`.\check.ps1` on Windows) is the gate: it runs the steps in
order and stops at the first one that fails, so just keep making the next step
pass.

## Steps

1. **Core arithmetic** — `+ - * /` with standard precedence, parentheses, unary
   minus, integer and decimal numbers. Malformed input must return an error.
2. **Modulo** — the `%` operator, same precedence as `*` and `/`.
3. **Power** — the `^` operator, right-associative, binding tighter than `* /`
   (so `2 ^ 3 ^ 2 == 512` and `2 * 3 ^ 2 == 18`).
4. **Functions** — `sqrt`, `abs`, `min`, `max`, `pow`, called like `sqrt(16)`
   and `max(3, 7)`; arguments are themselves expressions.
5. **Error positions** — when input has an unexpected character, the returned
   error must name the 1-based column, e.g. `2 @ 3` → an error containing
   `column 3`.

## The spec and the gate

The tests in `eval/eval_test.go` (`TestStep1Core` … `TestStep5Errors`) are the
specification, one per step. The gate runs, in order: `gofmt`, `go build`, then
each step's unit tests plus an acceptance run of the built `calc` binary over
`cases-stepN.txt`. It prints `CHECK PASSED` only when all five steps are green.

You are not expected to finish all five — get as far as you can; each step
stands on its own.
