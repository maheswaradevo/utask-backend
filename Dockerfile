FROM golang:1.20-alpine3.18

WORKDIR /build

COPY ./go.mod ./go.sum /build/
RUN go mod download

COPY . .

RUN go build main.go

EXPOSE 3000