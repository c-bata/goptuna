#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
ROOTDIR=$(dirname $(dirname $DIR))

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
echo "2. Run Goptuna optimizations."
echo ""

go run ${DIR}/main.go mysql "goptuna:password@tcp(localhost:3306)/goptuna?parseTime=true"

##################################################################
echo ""
echo "3. View the optimization results on Goptuna bulit-in dashboard."
echo ""
go run ${ROOTDIR}/cmd/main.go dashboard --storage mysql+mysqldb://goptuna:password@127.0.0.1:3306/goptuna

##################################################################
echo ""
echo "4. Stop MYSQL Server"
echo ""

docker stop goptuna-mysql
