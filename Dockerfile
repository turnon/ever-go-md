FROM golang:alpine as builder

WORKDIR $GOPATH/src/github.com/turnon/ever-go-md
COPY . ./

RUN apk add --no-cache git \
    && export GO111MODULE=on \
    && go get ./... \
    && go build -o /ever-go-md \
    && apk del git

FROM alpine:latest

WORKDIR /usr/bin
COPY --from=builder /ever-go-md .

ENTRYPOINT ["/usr/bin/ever-go-md", "-from", "/from", "-to", "/to", "-clean"]