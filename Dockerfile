FROM golang:1.20 as builder
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

FROM postgres:14 as postgres
COPY database/010_schema.sql /docker-entrypoint-initdb.d/010_schema.sql
COPY database/020_postsetup.sql /docker-entrypoint-initdb.d/020_postsetup.sql
