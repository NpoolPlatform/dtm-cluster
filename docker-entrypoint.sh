#!/bin/sh
echo "start*****************"
echo $ENV_CLUSTER_NAMESPACE
CONSUL_HTTP_ADDR=${ENV_CONSUL_HOST}:${ENV_CONSUL_PORT} consul services register -address=dtm.${ENV_CLUSTER_NAMESPACE}.svc.cluster.local -name=dtm.npool.top -port=36790
if [ ! $? -eq 0 ]; then
  echo "FAIL TO REGISTER CONFIGSERVICE TO CONSUL"
  exit 1
fi

MYSQL_HOST=`curl http://${ENV_CONSUL_HOST}:${ENV_CONSUL_PORT}/v1/agent/health/service/name/mysql.npool.top | jq '.[0] | .Service | .Address'`
if [ ! $? -eq 0 ]; then
  echo "FAIL TO GET MYSQL HOST"
  exit 1
fi

MYSQL_PORT=`curl http://${ENV_CONSUL_HOST}:${ENV_CONSUL_PORT}/v1/agent/health/service/name/mysql.npool.top | jq '.[0] | .Service | .Port'`
if [ ! $? -eq 0 ]; then
  echo "FAIL TO GET MYSQL PORT"
  exit 1
fi

MYSQL_HOST=`echo $MYSQL_HOST | sed 's/"//g'`


mysql -uroot -p$MYSQL_PASSWORD -h $MYSQL_HOST < /dtmcli.barrier.mysql.sql
mysql -uroot -p$MYSQL_PASSWORD -h $MYSQL_HOST < /dtmsvr.storage.mysql.sql

if [ ! $? -eq 0 ]; then
  echo "FAIL TO IMPORT SQL FILE with options $MYSQL_HOST $MYSQL_PORT"
fi
ls /app/dtm/configs/
sed -i "s/HOST/$MYSQL_HOST/g" /app/dtm/configs/config.yaml
sed -i "s/PORT/$MYSQL_PORT/g" /app/dtm/configs/config.yaml
sed -i "s/PWD/$MYSQL_PASSWORD/g" /app/dtm/configs/config.yaml

/usr/local/bin/docker-entrypoint.sh $@