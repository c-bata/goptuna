#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

rm db.sqlite3

gtimeout 6 go run ${DIR}/main.go sqlite3 db.sqlite3  # brew install coreutils

sqlite3 db.sqlite3 <<END_SQL
.header on
.mode column
.tables
select * from trials;
END_SQL
