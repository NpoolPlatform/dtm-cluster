FROM yedf/dtm:1.8.4

RUN apk add curl
RUN apk add jq
RUN apk add mysql-client

COPY config.yaml /app/dtm/configs/


RUN chmod a+r /app/dtm/configs/config.yaml

COPY ./sqls/dtmcli.barrier.mysql.sql /
COPY ./sqls/dtmsvr.storage.mysql.sql /

COPY .docker-tmp/consul /usr/bin/consul

RUN mkdir -p /usr/local/bin
COPY docker-entrypoint.sh /usr/local/bin
RUN cp /usr/local/bin/docker-entrypoint.sh /usr/local/bin/docker-entrypoint-inner.sh

RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]

