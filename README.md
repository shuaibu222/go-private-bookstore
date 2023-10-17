### private and public bookstore golang REST API


version: "3"

services:

  bookstore:
    build:
      context: .
      dockerfile: Dockerfile
    restart: always
    ports:
      - "9000:7000"
    depends_on:
      - postgres
    deploy:
      mode: replicated
      replicas: 1
    environment: 
      - JWT_SECRET=${JWT_SECRET}
      - WEB_PORT=${WEB_PORT}
    env_file:
      - .env

  postgres:
    image: "postgres"
    ports:
      - "5433:5432"
    restart: always
    deploy:
      mode: replicated
      replicas: 1
    environment:
      - POSTGRES_HOST=${POSTGRES_HOST}
      - POSTGRES_PORT=${POSTGRES_PORT}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    env_file:
      - .env
    volumes:
      - postgres:/var/lib/postgresql/data/

volumes:
  postgres: {}