#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)

set -ex

go run ${DIR}/main.go sqlite3 db.sqlite3
go run ${DIR}/cmd/main.go dashboard --storage sqlite:///db.sqlite3
