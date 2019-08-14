#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
REPOSITORY_ROOT=$(cd $(dirname $(dirname $(dirname $0))); pwd)

##################################################################
echo ""
echo "1. Prepare MYSQL 8.0 Server using Docker."
echo ""

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

##################################################################
echo ""
echo "2. Create a study."
echo ""

go run ${REPOSITORY_ROOT}/cmd/main.go create-study --storage mysql+mysqldb://goptuna:password@127.0.0.1:3306/goptuna --study rdb

##################################################################
echo ""
echo "3. Check the tables of MYSQL."
echo ""

echo "DATABASES:"
mysql --host 127.0.0.1 --port 3306 --user goptuna -ppassword -e "SHOW DATABASES;"

echo "TABLES:"
mysql --host 127.0.0.1 --port 3306 --user goptuna -ppassword -e "SHOW TABLES FROM goptuna;"

##################################################################
echo ""
echo "4. Run Goptuna optimizations."
echo ""

go run ${DIR}/main.go mysql "goptuna:password@tcp(localhost:3306)/goptuna?parseTime=true"

##################################################################
echo ""
echo "5. View the optimization results on Optuna's dashboard."
echo ""

if [ -d ./venv ]; then
    source venv/bin/activate
    pip install mysqlclient
else
    python3.7 -m venv venv
    source venv/bin/activate
    pip install optuna bokeh mysqlclient
fi

optuna dashboard --storage mysql+mysqldb://goptuna:password@127.0.0.1:3306/goptuna --study rdb

##################################################################
echo ""
echo "6. Stop MYSQL Server"
echo ""

docker stop goptuna-mysql
