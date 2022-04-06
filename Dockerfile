FROM yedf/dtm:1.8.4

RUN apk add curl
RUN apk add jq
RUN apk add mysql-client

COPY config.yaml /app/dtm/configs/

RUN ls /app/dtm/configs/

COPY ./sqls/dtmcli.barrier.mysql.sql /
COPY ./sqls/dtmcli.barrier.mysql.sql /

COPY .docker-tmp/consul /usr/bin/consul

RUN mkdir -p /usr/local/bin
COPY docker-entrypoint.sh /usr/local/bin
RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

