FROM golang:1.20-alpine3.18

WORKDIR /build

COPY ./go.mod ./go.sum /build/
RUN go mod download

COPY . .

RUN go build main.go

EXPOSE 3000


# FROM golang:1.20 as build
# WORKDIR /build

# COPY go.mod /build/
# COPY go.sum /build/

# RUN go mod download
# RUN go mod tidy

# COPY . /build/

# FROM alpine:3.16.0
# WORKDIR /build



# EXPOSE 3000


