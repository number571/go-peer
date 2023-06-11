FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.4.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlm ./cmd/hidden_lake/messenger/cmd/hlm

ENV SERVICE_PATH="/mounted"
CMD ./hlm -path="${SERVICE_PATH}"
