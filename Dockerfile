
FROM --platform=amd64 golang:1.22.2 AS builder

RUN apt-get install git gcc libc-dev curl -y
RUN mkdir /builder
WORKDIR /builder
COPY . .
RUN go build -o kongo .

FROM --platform=amd64 kong/kong-gateway:3.4.2.0 as kong
COPY --from=builder /builder/kongo /usr/local/kong/
USER kong