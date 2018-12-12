FROM ubuntu:16.04

WORKDIR /build
COPY . /build/


ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOVERSION 1.11.2
ENV GOROOT /opt/go
ENV GOPATH /root/.go

RUN sed -i -e 's/archive\.ubuntu/us.archive.ubuntu/' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends --fix-missing ca-certificates wget git gcc build-essential && \
    cd /opt && \
    wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH && \
    cd /build/cmd/fdb-clusterfile && \
    go build .


