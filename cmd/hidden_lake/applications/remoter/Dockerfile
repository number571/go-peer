FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.16.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlr ./cmd/hidden_lake/applications/remoter/cmd/hlr

ENV SERVICE_PATH="/mounted"
CMD ./hlr -path="${SERVICE_PATH}"
