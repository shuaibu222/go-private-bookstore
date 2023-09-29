FROM golang:1.20

WORKDIR /app

# Copy go.mod, go.sum. Whith go.mod no need of go get -d -v
COPY go.mod . 
COPY go.sum .
COPY .env /app/.env

# Just download from go.mod not install
RUN go mod download

# Copy all source code then build
COPY . .

# Build it inside app directory
RUN go build -o bin .

# Run this command when running this image
ENTRYPOINT [ "/app/bin" ]