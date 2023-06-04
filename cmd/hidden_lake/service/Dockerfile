FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://dl.google.com/go/go1.20.3.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.20.3.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hls ./cmd/hidden_lake/service/cmd/hls

ENTRYPOINT ["./hls"]
CMD ["-path", "/mounted"]
