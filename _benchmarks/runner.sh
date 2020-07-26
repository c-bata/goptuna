#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}
DIR=$(cd $(dirname $0); pwd)
BINDIR=$(dirname $DIR)/bin
REPEATS=${REPEATS:-5}
BUDGET=${BUDGET:-300}
SEED=${BUDGET:-1}
DIM=${DIM:-2}

usage() {
    cat <<EOF
$(basename ${0}) is an entrypoint to run benchmarkers.
Usage:
    $ $(basename ${0}) <problem> <json-path>
Problem:
    rosenbrock     : https://www.sfu.ca/~ssurjano/rosen.html
    himmelblau     : https://en.wikipedia.org/wiki/Himmelblau%27s_function
    ackley         : https://www.sfu.ca/~ssurjano/ackley.html
    rastrigin      : https://www.sfu.ca/~ssurjano/rastr.html
    weierstrass    : Weierstrass function in https://github.com/sigopt/evalset
    schwefel20     : https://www.sfu.ca/~ssurjano/schwef.html
    schwefel36     : https://www.sfu.ca/~ssurjano/schwef.html
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
go build -o ${BINDIR}/rastrigin_problem ${DIR}/rastrigin_problem/main.go

RANDOM_SOLVER=$($KUROBAKO solver random)
CMA_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver cmaes)
IPOP_CMA_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver ipop-cmaes)
BIPOP_CMA_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver bipop-cmaes)
TPE_SOLVER=$($KUROBAKO solver command ${BINDIR}/goptuna_solver tpe)
OPTUNA_CMA_SOLVER=$($KUROBAKO solver command python ${DIR}/optuna_solver.py cmaes)
OPTUNA_TPE_SOLVER=$($KUROBAKO solver command python ${DIR}/optuna_solver.py tpe)

case "$1" in
    himmelblau)
        PROBLEM=$($KUROBAKO problem command ${BINDIR}/himmelblau_problem)
        ;;
    rosenbrock)
        PROBLEM=$($KUROBAKO problem command ${BINDIR}/rosenbrock_problem)
        ;;
    ackley)
        PROBLEM=$($KUROBAKO problem sigopt --dim $DIM ackley)
        ;;
    rastrigin)
        # "kurobako problem sigopt --dim 8 rastrigin" only accepts 8-dim.
        PROBLEM=$($KUROBAKO problem command ${BIDIR}/rastrigin_problem $DIM)
        ;;
    weierstrass)
        PROBLEM=$($KUROBAKO problem sigopt --dim $DIM weierstrass)
        ;;
    schwefel20)
        PROBLEM=$($KUROBAKO problem sigopt --dim 2 schwefel20)
        ;;
    schwefel36)
        PROBLEM=$($KUROBAKO problem sigopt --dim 2 schwefel36)
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

$KUROBAKO studies \
  --solvers \
    $RANDOM_SOLVER \
    $CMA_SOLVER \
    $IPOP_CMA_SOLVER \
    $BIPOP_CMA_SOLVER \
    $TPE_SOLVER \
    $OPTUNA_CMA_SOLVER \
    $OPTUNA_TPE_SOLVER \
  --problems $PROBLEM \
  --seed $SEED --repeats $REPEATS --budget $BUDGET \
  | $KUROBAKO run --parallelism 1 > $2
