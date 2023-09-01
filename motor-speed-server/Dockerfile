FROM alpine:latest

RUN mkdir /app

COPY ./cmd/api/templates /templates

COPY ./motorSpeedServerApp /app/motorSpeedServerApp

CMD [ "/app/motorSpeedServerApp" ]