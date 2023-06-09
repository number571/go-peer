FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.4.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlt ./cmd/hidden_lake/traffic/cmd/hlt

ENTRYPOINT ["./hlt"]
CMD ["-path", "/mounted"]
