FROM alpine

RUN apk add --update ca-certificates

COPY bin/sicily /usr/bin/sicily

EXPOSE 3000

ENTRYPOINT sicily \
    -citizens-host=$CITIZENS_HOST \
    -citizens-port=$CITIZENS_PORT \
    -palermo-host=$PALERMO_HOST \
    -palermo-port=$PALERMO_PORT \
    -plato-host=$PLATO_HOST \
    -plato-port=$PLATO_PORT \
    -helenia-host=$HELENIA_HOST \
    -helenia-port=$HELENIA_PORT
