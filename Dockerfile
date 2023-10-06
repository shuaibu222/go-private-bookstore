#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v .

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app

# Set execute permissions
RUN chmod +x /app

# Copy the .env file into the container
COPY .env /app

ENTRYPOINT /app
LABEL Name=gobookstore Version=0.0.1
EXPOSE 9000
