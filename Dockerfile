FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY bookstoreApp /app

CMD [ "/app/bookstoreApp" ]