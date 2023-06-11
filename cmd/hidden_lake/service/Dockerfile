FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.4.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hls ./cmd/hidden_lake/service/cmd/hls

ENV SERVICE_PATH="/mounted"
ENV SERVICE_KEY=""
CMD ./hls -path="${SERVICE_PATH}" -key="${SERVICE_KEY}"
