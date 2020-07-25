#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}
DIR=$(cd $(dirname $0); pwd)
BINDIR=$(dirname $DIR)/bin

usage() {
    cat <<EOF
$(basename ${0}) is an entrypoint to run benchmarkers.
Usage:
    $ $(basename ${0}) <problem> <json-path>
Problem:
    rosenbrock     : https://www.sfu.ca/~ssurjano/rosen.html
    himmelblau     : https://en.wikipedia.org/wiki/Himmelblau%27s_function
    ackley         : Ackley function in https://github.com/sigopt/evalset
    weierstrass    : Weierstrass function in https://github.com/sigopt/evalset
    schwefel20     : Schwefel20 function in https://github.com/sigopt/evalset
    schwefel36     : Schwefel36 function in https://github.com/sigopt/evalset
Options:
    --help, -h         print this
Example:
    $ $(basename ${0}) rosenbrock ./tmp/kurobako.json
    $ cat ./tmp/kurobako.json | kurobako plot curve --errorbar -o ./tmp
EOF
}

mkdir -p $BINDIR

go build -o ${BINDIR}/goptuna_solver ${DIR}/goptuna_solver/main.go
go build -o ${BINDIR}/himmelblau_problem ${DIR}/himmelblau_problem/main.go
go build -o ${BINDIR}/rosenbrock_problem ${DIR}/rosenbrock_problem/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)
CMA_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver cmaes)
IPOP_CMA_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver ipop-cmaes)
TPE_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver tpe)
OPTUNA_CMA_SOLVER=$($KUROBAKO solver command python ${DIR}/optuna_solver.py cmaes)
OPTUNA_TPE_SOLVER=$($KUROBAKO solver command python ${DIR}/optuna_solver.py tpe)

case "$1" in
    himmelblau)
        PROBLEM=$($KUROBAKO problem command ${BINDIR}/himmelblau_problem)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $CMA_SOLVER $IPOP_CMA_SOLVER $OPTUNA_CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 5 --budget 100 \
          | $KUROBAKO run --parallelism 1 > $2
        ;;
    rosenbrock)
        PROBLEM=$($KUROBAKO problem command ${BINDIR}/rosenbrock_problem)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $CMA_SOLVER $IPOP_CMA_SOLVER $OPTUNA_CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 5 --budget 100 \
          | $KUROBAKO run --parallelism 1 > $2
        ;;
    ackley)
        PROBLEM=$($KUROBAKO problem sigopt --dim 10 ackley)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $CMA_SOLVER $IPOP_CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 5 --budget 1000 \
          | $KUROBAKO run --parallelism 5 > $2
        ;;
    weierstrass)
        PROBLEM=$($KUROBAKO problem sigopt --dim 10 weierstrass)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $IPOP_CMA_SOLVER $CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 5 --budget 1000 \
          | $KUROBAKO run --parallelism 5 > $2
        ;;
    schwefel20)
        PROBLEM=$($KUROBAKO problem sigopt --dim 2 schwefel20)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $IPOP_CMA_SOLVER $CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 10 --budget 100 \
          | $KUROBAKO run --parallelism 5 > $2
        ;;
    schwefel36)
        PROBLEM=$($KUROBAKO problem sigopt --dim 2 schwefel36)
        $KUROBAKO studies \
          --solvers $RANDOM_SOLVER $IPOP_CMA_SOLVER $CMA_SOLVER $TPE_SOLVER $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed 1 --repeats 10 --budget 100 \
          | $KUROBAKO run --parallelism 5 > $2
        ;;
    help|--help|-h)
        usage
        exit 0
        ;;
    *)
        echo "[Error] Invalid problem '${1}'"
        usage
        exit 1
        ;;
esac
