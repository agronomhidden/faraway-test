FROM golang:1.21.11 as builder

WORKDIR /go/src/app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

COPY . /go/src/app
RUN go build -o app ./main.go

FROM alpine:latest

COPY --from=builder /go/src/app/app app

ENTRYPOINT ./app
