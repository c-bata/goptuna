#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
BIN_DIR=$(cd $(dirname $(dirname $0)); pwd)/bin

mkdir -p ${BIN_DIR}
go build -o ${BIN_DIR}/simple_tpe ${DIR}/simple_tpe/main.go
go build -o ${BIN_DIR}/simple_random_search ${DIR}/simple_random_search/main.go
