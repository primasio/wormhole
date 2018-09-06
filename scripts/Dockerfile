FROM scratch

ADD ./bin/wormhole /wormhole/
ADD ./certs/* /etc/ssl/certs/

EXPOSE 80
VOLUME ["/wormhole/config", "/wormhole/logs"]

ENV APP_ENV production
ENV APP_CONFIG /wormhole/config

ENTRYPOINT ["/wormhole/wormhole"]
