FROM yedf/dtm:1.8.4

RUN apk add --no-cache curl jq mysql-client

COPY config.yaml /app/dtm/configs/
COPY ./sqls /
COPY .docker-tmp/consul docker-entrypoint.sh /usr/local/bin/

RUN chmod a+r /app/dtm/configs/config.yaml && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
