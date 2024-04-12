FROM alpine:latest

RUN mkdir /app

COPY listenerApp /app

WORKDIR  /app

CMD [ "/app/listenerApp" ]