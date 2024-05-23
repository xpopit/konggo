FROM golang:1.19-alpine3.15 AS plugin-builder

WORKDIR /builder

COPY . .

RUN apk add make
RUN make build