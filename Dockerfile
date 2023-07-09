FROM golang:1.20 as builder
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
