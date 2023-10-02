FROM golang:1.20

WORKDIR /app

# Copy go.mod, go.sum. Whith go.mod no need of go get -d -v
COPY go.mod . 
COPY go.sum .

# Copy all source code then build
COPY . .
# Just download from go.mod not install
RUN go mod download


# Build it inside app directory
RUN go build -o bin .

# Run this command when running this image
ENTRYPOINT [ "/app/bin" ]