FROM golang:1.25-bookworm
WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
