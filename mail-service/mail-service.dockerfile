FROM alpine:latest

RUN mkdir /app

COPY mailServiceApp /app

ADD templates /app/templates

WORKDIR /app

CMD [ "/app/mailServiceApp" ]

