FROM scratch

ADD ./bin/wormhole /wormhole/

EXPOSE 80
VOLUME ["/wormhole/config", "/wormhole/logs"]

ENV APP_ENV production
ENV APP_CONFIG /wormhole/config

CMD ["/wormhole/wormhole", "-log_dir=\"/wormhole/logs\"", "-v=1"]
