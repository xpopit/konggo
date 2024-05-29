
FROM --platform=amd64 golang:1.22.2 AS builder

RUN apt-get install git gcc libc-dev curl -y
RUN mkdir /builder
WORKDIR /builder
COPY . . 
# RUN curl https://raw.githubusercontent.com/Kong/go-plugins/master/go-hello.go -o /builder/main.go
# RUN go mod init kongo && \
#     go get github.com/Kong/go-pdk && \
#     go mod tidy
# RUN go build -o kongo .
RUN go get github.com/Kong/go-pdk && \
    go mod tidy
RUN go build -o kongo .

FROM --platform=amd64 kong/kong:3.6.1 as kong
# COPY ./config.yaml /kong/
# COPY --from=builder /builder/kongo /usr/local/kong/
COPY --from=builder /builder/kongo  ./kong/
USER kong