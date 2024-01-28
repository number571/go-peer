FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.16.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlm ./cmd/hidden_lake/applications/filer/cmd/hlf

ENV SERVICE_PATH="/mounted"
CMD ./hlf -path="${SERVICE_PATH}"
