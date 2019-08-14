#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

rm db.sqlite3

##################################################################
echo ""
echo "1. Create a study."
echo ""

go run ${REPOSITORY_ROOT}/cmd/main.go create-study --storage sqlite:///db.sqlite3 --study rdb

##################################################################
echo ""
echo "2. Check the tables of SQLite3."
echo ""

sqlite3 db.sqlite3 <<END_SQL
.header on
.mode column
.tables
select * from studies;
END_SQL

##################################################################
echo ""
echo "3. Run Goptuna optimization."
echo ""

go run ${DIR}/main.go sqlite3 db.sqlite3

##################################################################
echo ""
echo "4. View the optimization results on Optuna's dashboard."
echo ""

if [ -d ./venv ]; then
    source venv/bin/activate
else
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh
fi

optuna dashboard --storage sqlite:///db.sqlite3 --study rdb

##################################################################
echo ""
echo "5. Delete db.sqlite3"
echo ""

rm db.sqlite3
