#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

set -e

echo ""
echo "1. Run Goptuna optimization."
echo ""

go run ${DIR}/main.go sqlite3 db.sqlite3

echo ""
echo "2. View the optimization results on Optuna's dashboard."
echo ""

if [ -d ./venv ]; then
    source venv/bin/activate
else
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh
fi

optuna dashboard --storage sqlite:///db.sqlite3 --study rdb
