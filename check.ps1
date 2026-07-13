#!/usr/bin/env pwsh
# The verification gate. Runs fmt and build, then each lab STEP in order,
# stopping at the first step that fails. PowerShell equivalent of check.sh.
$ErrorActionPreference = 'Continue'
Set-Location -Path $PSScriptRoot

Write-Host "== fmt (gofmt) =="
$unformatted = & gofmt -l .
if ($unformatted) {
    Write-Host "CHECK FAILED (fmt): not gofmt-clean:"
    $unformatted | ForEach-Object { Write-Host $_ }
    exit 1
}

Write-Host "== build =="
go build ./...
if ($LASTEXITCODE -ne 0) { Write-Host "CHECK FAILED (build)"; exit 1 }
$bin = Join-Path $PSScriptRoot 'calc.exe'
go build -o $bin .
if ($LASTEXITCODE -ne 0) { Write-Host "CHECK FAILED (build)"; exit 1 }

function Test-Cases($file) {
    if (-not (Test-Path $file)) { return $true }
    $ok = $true
    foreach ($line in Get-Content $file) {
        if (-not $line.Trim()) { continue }
        $parts = $line -split '\|', 2
        $expr = $parts[0]; $expected = $parts[1]
        $got = (& $bin $expr | Out-String).Trim()
        if ($got -ne $expected) {
            Write-Host "    FAIL: '$expr' => got '$got', want '$expected'"
            $ok = $false
        }
    }
    return $ok
}

$steps = @(
    @{ num = 1; title = 'Core arithmetic (+ - * /, parens, unary minus, decimals)'; pattern = 'TestStep1Core'; cases = 'cases-step1.txt' },
    @{ num = 2; title = 'Modulo operator %'; pattern = 'TestStep2Modulo'; cases = 'cases-step2.txt' },
    @{ num = 3; title = 'Power operator ^ (right-associative, binds tighter than * /)'; pattern = 'TestStep3Power'; cases = 'cases-step3.txt' },
    @{ num = 4; title = 'Built-in functions: sqrt abs min max pow'; pattern = 'TestStep4Functions'; cases = 'cases-step4.txt' },
    @{ num = 5; title = 'Error messages carry the column of the bad character'; pattern = 'TestStep5Errors'; cases = '' }
)

foreach ($s in $steps) {
    Write-Host ("== STEP {0}/5: {1} ==" -f $s.num, $s.title)
    $out = & go test -run $s.pattern ./eval/ 2>&1
    if ($LASTEXITCODE -ne 0) {
        $out | ForEach-Object { Write-Host "    $_" }
        Write-Host ("CHECK FAILED (step {0})" -f $s.num); exit 1
    }
    if ($s.cases) {
        if (-not (Test-Cases $s.cases)) {
            Write-Host ("CHECK FAILED (step {0})" -f $s.num); exit 1
        }
    }
    Write-Host ("  STEP {0} PASSED" -f $s.num)
}

Write-Host "CHECK PASSED"
exit 0
