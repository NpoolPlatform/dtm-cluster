FROM yedf/dtm:1.8.4

RUN apk add curl jq mysql-client

COPY config.yaml /app/dtm/configs/

RUN chmod a+r /app/dtm/configs/config.yaml

COPY ./sqls /
COPY .docker-tmp/consul docker-entrypoint.sh /usr/local/bin/

RUN cp /usr/local/bin/docker-entrypoint.sh /usr/local/bin/docker-entrypoint-inner.sh

RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
