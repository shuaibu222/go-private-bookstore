FROM alpine:latest

RUN mkdir /app

COPY bookstoreApp /app

CMD [ "/app/bookstoreApp" ]