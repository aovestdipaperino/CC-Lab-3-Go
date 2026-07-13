#!/usr/bin/env bash
# The verification gate. Runs fmt and build, then each lab STEP in order,
# stopping at the first step that fails. Each step is checked by its unit
# tests and (where present) an acceptance run of the built binary over
# cases-stepN.txt. The gate passes only when every step is green.
set -uo pipefail
cd "$(dirname "$0")"

log="${TMPDIR:-/tmp}/calc-step.log"

echo "== fmt (gofmt) =="
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "CHECK FAILED (fmt): not gofmt-clean:"
  echo "$unformatted"
  exit 1
fi

echo "== build =="
if ! go build ./...; then
  echo "CHECK FAILED (build)"
  exit 1
fi
bin="./calc"
if ! go build -o "$bin" .; then
  echo "CHECK FAILED (build)"
  exit 1
fi

run_cases() {
  local f="$1" fail=0 got
  [ -f "$f" ] || return 0
  while IFS='|' read -r expr expected; do
    [ -z "${expr// }" ] && continue
    got=$("$bin" "$expr")
    if [ "$got" != "$expected" ]; then
      echo "    FAIL: '$expr' => got '$got', want '$expected'"
      fail=1
    fi
  done < "$f"
  return $fail
}

# num | title | go test -run pattern | cases file (optional)
steps=(
  "1|Core arithmetic (+ - * /, parens, unary minus, decimals)|TestStep1Core|cases-step1.txt"
  "2|Modulo operator %|TestStep2Modulo|cases-step2.txt"
  "3|Power operator ^ (right-associative, binds tighter than * /)|TestStep3Power|cases-step3.txt"
  "4|Built-in functions: sqrt abs min max pow|TestStep4Functions|cases-step4.txt"
  "5|Error messages carry the column of the bad character|TestStep5Errors|"
)

for entry in "${steps[@]}"; do
  IFS='|' read -r num title pattern casefile <<< "$entry"
  echo "== STEP $num/5: $title =="
  if ! go test -run "$pattern" ./eval/ > "$log" 2>&1; then
    sed 's/^/    /' "$log"
    echo "CHECK FAILED (step $num)"
    exit 1
  fi
  if [ -n "$casefile" ]; then
    if ! run_cases "$casefile"; then
      echo "CHECK FAILED (step $num)"
      exit 1
    fi
  fi
  echo "  STEP $num PASSED"
done

echo "CHECK PASSED"
exit 0
