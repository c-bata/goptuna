#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}

go build -o ./himmelblau_problem ./_benchmarks/himmelblau_problem/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)
OPTUNA_SOLVER=$($KUROBAKO solver command python ./_benchmarks/optuna_solver/cmaes.py)
PROBLEM=$($KUROBAKO problem command ./himmelblau_problem)

$KUROBAKO studies \
  --solvers $RANDOM_SOLVER $OPTUNA_SOLVER \
  --problems $PROBLEM \
  --repeats 25 --budget 50 \
  | $KUROBAKO run --parallelism 5 > $1
