# Build
FROM golang:1.18-alpine as builder

RUN apk add --no-cache gcc musl-dev

COPY go.mod /polygon-pools/
COPY go.sum /polygon-pools/
RUN cd /polygon-pools && go mod download

ADD . /polygon-pools
RUN cd /polygon-pools && go build main.go

# Pull
FROM alpine:latest

COPY --from=builder /polygon-pools/main /usr/local/bin/

EXPOSE 8035
ENTRYPOINT ["main"]