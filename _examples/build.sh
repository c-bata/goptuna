#!/bin/sh

DIR=$(cd $(dirname $0); pwd)
BIN_DIR=$(cd $(dirname $(dirname $0)); pwd)/bin

mkdir -p ${BIN_DIR}

set -ex

cd ${DIR}
go build -o ${BIN_DIR}/cmaes ./cmaes/main.go
go build -o ${BIN_DIR}/cmaes_blackhole ./cmaes/blackhole/main.go
go build -o ${BIN_DIR}/enqueue_trial ./enqueue_trial/main.go
go build -o ${BIN_DIR}/trialnotify ./trialnotify/main.go
go build -o ${BIN_DIR}/signalhandling ./signalhandling/main.go
go build -o ${BIN_DIR}/simple_rdb ./simple_rdb/main.go
go build -o ${BIN_DIR}/simple_tpe ./simple_tpe/main.go
go build -o ${BIN_DIR}/sobol ./sobol/main.go
#go build -o ${BIN_DIR}/concurrency ${DIR}/concurrency/main.go
#go build -o ${BIN_DIR}/gorgonia_iris ${DIR}/gorgonia_iris/main.go
