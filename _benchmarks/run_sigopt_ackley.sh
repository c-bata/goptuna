#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}

go build -o ./tpe_solver ./_benchmarks/tpe_solver/main.go
go build -o ./cma_solver ./_benchmarks/cma_solver/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)
TPE_SOLVER=$($KUROBAKO solver command ./tpe_solver)
CMA_SOLVER=$($KUROBAKO solver command ./cma_solver)
PROBLEM=$($KUROBAKO problem sigopt --dim 5 ackley)

$KUROBAKO studies \
  --solvers $RANDOM_SOLVER $TPE_SOLVER $CMA_SOLVER \
  --problems $PROBLEM \
  --repeats 100 --budget 100 \
  | $KUROBAKO run --parallelism 4 > $1
