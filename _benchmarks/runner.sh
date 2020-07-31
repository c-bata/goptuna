#!/bin/sh

set -e

KUROBAKO=${KUROBAKO:-kurobako}
DIR=$(cd $(dirname $0); pwd)
BINDIR=$(dirname $DIR)/bin
TMPDIR=$(dirname $DIR)/tmp
REPEATS=${REPEATS:-5}
BUDGET=${BUDGET:-300}
SEED=${SEED:-1}
DIM=${DIM:-2}
SOLVERS=${SOLVERS:-all}
LOGLEVEL=${LOGLEVEL:-error}

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
    hpobench-naval
    hpobench-parkinson
    hpobench-protein
    hpobench-slice
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

OPTUNA_CMA_SOLVER=$($KUROBAKO solver --name Optuna-CMAES optuna --loglevel ${LOGLEVEL} --sampler CmaEsSampler)
OPTUNA_TPE_SOLVER=$($KUROBAKO solver --name Optuna-TPE optuna --loglevel ${LOGLEVEL} --sampler TPESampler)
OPTUNA_RANDOM_MEDIAN_SOLVER=$(kurobako solver --name Optuna-RANDOM-MEDIAN optuna --loglevel ${LOGLEVEL} --sampler RandomSampler --pruner MedianPruner)
OPTUNA_RANDOM_ASHA_SOLVER=$(kurobako solver --name Optuna-RANDOM-ASHA optuna --loglevel ${LOGLEVEL} --sampler RandomSampler --pruner SuccessiveHalvingPruner)
OPTUNA_TPE_MEDIAN_SOLVER=$(kurobako solver --name Optuna-TPE-MEDIAN optuna --loglevel  ${LOGLEVEL} --sampler TPESampler --pruner MedianPruner)
OPTUNA_TPE_ASHA_SOLVER=$(kurobako solver --name Optuna-TPE-ASHA optuna --loglevel ${LOGLEVEL} --sampler TPESampler --pruner SuccessiveHalvingPruner)

if [[ $1 == "hpobench-"* ]] ; then
  if [ ! -d "$TMPDIR/fcnet_tabular_benchmarks" ] ; then
    if [ ! -f "$TMPDIR/fcnet_tabular_benchmarks.tar.gz" ] ; then
      wget -O $TMPDIR/fcnet_tabular_benchmarks.tar.gz http://ml4aad.org/wp-content/uploads/2019/01/fcnet_tabular_benchmarks.tar.gz
    fi
    tar -xf $TMPDIR/fcnet_tabular_benchmarks.tar.gz -C $TMPDIR
  else
    echo "HPOBench dataset has already downloaded."
  fi
fi

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
        PROBLEM=$($KUROBAKO problem command ${BINDIR}/rastrigin_problem $DIM)
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
    hpobench-naval)
        PROBLEM=$($KUROBAKO problem hpobench "${TMPDIR}/fcnet_tabular_benchmarks/fcnet_naval_propulsion_data.hdf5")
        ;;
    hpobench-parkinson)
        PROBLEM=$($KUROBAKO problem hpobench "${TMPDIR}/fcnet_tabular_benchmarks/fcnet_parkinsons_telemonitoring_data.hdf5")
        ;;
    hpobench-protein)
        PROBLEM=$($KUROBAKO problem hpobench "${TMPDIR}/fcnet_tabular_benchmarks/fcnet_protein_structure_data.hdf5")
        ;;
    hpobench-slice)
        PROBLEM=$($KUROBAKO problem hpobench "${TMPDIR}/fcnet_tabular_benchmarks/fcnet_slice_localization_data.hdf5")
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
case $SOLVERS in
    all)
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
          | $KUROBAKO run --parallelism 7 > $2
        ;;
    cmaes)
        $KUROBAKO studies \
          --solvers \
            $RANDOM_SOLVER \
            $CMA_SOLVER \
            $IPOP_CMA_SOLVER \
            $BIPOP_CMA_SOLVER \
            $OPTUNA_CMA_SOLVER \
          --problems $PROBLEM \
          --seed $SEED --repeats $REPEATS --budget $BUDGET \
          | $KUROBAKO run --parallelism 5 > $2
        ;;
    tpe)
        $KUROBAKO studies \
          --solvers \
            $RANDOM_SOLVER \
            $TPE_SOLVER \
            $OPTUNA_TPE_SOLVER \
          --problems $PROBLEM \
          --seed $SEED --repeats $REPEATS --budget $BUDGET \
          | $KUROBAKO run --parallelism 3 > $2
        ;;
    pruner)
        $KUROBAKO studies \
          --solvers \
            $RANDOM_SOLVER \
            $OPTUNA_RANDOM_MEDIAN_SOLVER \
            $OPTUNA_RANDOM_ASHA_SOLVER \
            $OPTUNA_TPE_MEDIAN_SOLVER \
            $OPTUNA_TPE_ASHA_SOLVER \
          --problems $PROBLEM \
          --seed $SEED --repeats $REPEATS --budget $BUDGET \
          | $KUROBAKO run --parallelism 3 > $2
        ;;
    *)
        echo "[Error] Invalid solver '${SOLVERS}'"
        usage
        exit 1
        ;;
esac
