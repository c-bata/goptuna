#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}

go build -o ./kurobako-solver ./_benchmarks/kurobako-solver/main.go
go build -o ./quadratic-problem ./_benchmarks/quadratic-problem/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)
GOPTUNA_SOLVER=$($KUROBAKO solver command ./kurobako-solver)
PROBLEM=$($KUROBAKO problem sigopt --dim 5 ackley)

$KUROBAKO studies \
  --solvers $RANDOM_SOLVER $GOPTUNA_SOLVER \
  --problems $PROBLEM \
  --repeats 100 --budget 100 \
  | $KUROBAKO run --parallelism 4 > $1
