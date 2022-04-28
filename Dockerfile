FROM yedf/dtm:1.8.4

RUN apk add curl
RUN apk add jq

COPY .docker-tmp/consul /usr/bin/consul

RUN mkdir -p /usr/local/bin
COPY docker-entrypoint.sh /usr/local/bin
RUN cp /usr/local/bin/docker-entrypoint.sh /usr/local/bin/docker-entrypoint-inner.sh

RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]

