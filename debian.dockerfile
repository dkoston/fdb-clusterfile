FROM golang:1.11.2-stretch

WORKDIR /build
COPY . /build/

RUN apt-get update && \
    apt-get install -y --no-install-recommends git build-essential gcc && \
    cd /build/cmd/fdb-clusterfile && \
    go build .