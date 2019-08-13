#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

rm db.sqlite3

set -e

go run ${REPOSITORY_ROOT}/cmd/main.go create-study --storage sqlite:///db.sqlite3 --study rdb
go run ${DIR}/main.go sqlite3 db.sqlite3

if [ -d ./venv ]; then
    source venv/bin/activate
else
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh
fi

optuna dashboard --storage sqlite:///db.sqlite3 --study rdb
