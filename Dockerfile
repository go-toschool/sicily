FROM alpine

RUN apk add --update ca-certificates

COPY bin/sicily /usr/bin/sicily

ENV CITIZENS_HOST "syracyse"
ENV CITIZENS_PORT 8001
ENV PALERMO_HOST "palermo"
ENV PALERMO_PORT 8003

EXPOSE 3000

ENTRYPOINT sicily -citizens-host=$CITIZENS_HOST -citizens-port=$CITIZENS_PORT -palermo-host=$PALERMO_HOST -palermo-port=$PALERMO_PORT