#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

docker pull mysql:8.0
docker stop goptuna-mysql

set -e

cd ${DIR}
docker run \
  -d \
  --rm \
  -p 3306:3306 \
  --mount type=volume,src=mysql,dst=/etc/mysql/conf.d \
  -e MYSQL_USER=goptuna \
  -e MYSQL_DATABASE=goptuna \
  -e MYSQL_PASSWORD=password \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  --name goptuna-mysql \
  mysql:8.0
cd -

echo "Wait ready for mysql"
sleep 20

echo "DATABASES:"
mysql --host 127.0.0.1 --port 3306 --user goptuna -ppassword -e "SHOW DATABASES;"


go run ${REPOSITORY_ROOT}/cmd/main.go create-study --storage mysql://goptuna:password@localhost:3306/goptuna --study rdb

echo "TABLES:"
mysql --host 127.0.0.1 --port 3306 --user goptuna -ppassword -e "SHOW TABLES FROM goptuna;"

go run ${DIR}/main.go mysql "goptuna:password@tcp(localhost:3306)/goptuna"

if [ -d ./venv ]; then
    source venv/bin/activate
else
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh
fi

optuna dashboard --storage mysql://goptuna:password@localhost:3306/goptuna --study rdb

docker stop goptuna-mysql
