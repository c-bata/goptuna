#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

rm db.sqlite3

gtimeout 6 go run ${DIR}/main.go # brew install coreutils

sleep 0.5

echo ""
echo "*** check trials ***"
echo ""

sqlite3 db.sqlite3 <<END_SQL
.header on
.mode column
select * from trials;
END_SQL
