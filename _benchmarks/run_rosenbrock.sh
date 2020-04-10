#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}

go build -o ./cma_solver ./_benchmarks/cma_solver/main.go
go build -o ./tpe_solver ./_benchmarks/tpe_solver/main.go
go build -o ./rosenbrock_problem ./_benchmarks/rosenbrock_problem/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)

CMA_SOLVER=$($KUROBAKO solver command ./cma_solver)
TPE_SOLVER=$($KUROBAKO solver command ./tpe_solver)
# It has a circular import problem.
# OPTUNA_SOLVER=$($KUROBAKO solver command python ./_benchmarks/optuna_solver/cmaes.py)
PROBLEM=$($KUROBAKO problem command ./rosenbrock_problem)

$KUROBAKO studies \
  --solvers $RANDOM_SOLVER $TPE_SOLVER $CMA_SOLVER \
  --problems $PROBLEM \
  --repeats 8 --budget 300 \
  | $KUROBAKO run --parallelism 1 > $1
