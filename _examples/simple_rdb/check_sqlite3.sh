#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

set -x

go run ${DIR}/main.go sqlite3 db.sqlite3

set +x

echo "Create virtualenv"

if [ ! -d ./venv ]; then
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh
else
    source venv/bin/activate
fi

set -x

pip install -U git+https://github.com/optuna/optuna

optuna dashboard --storage sqlite:///db.sqlite3 --study rdb
