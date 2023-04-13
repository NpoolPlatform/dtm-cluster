FROM uhub.service.ucloud.cn/entropypool_public/dtm:1.17.1

RUN apk add --no-cache curl jq mysql-client

COPY config.yaml /app/dtm/configs/
COPY ./sqls /
COPY .docker-tmp/consul docker-entrypoint.sh /usr/local/bin/

RUN chmod a+r /app/dtm/configs/config.yaml && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
