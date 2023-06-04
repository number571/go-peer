FROM ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://dl.google.com/go/go1.20.3.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.20.3.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlm ./cmd/hidden_lake/messenger/cmd/hlm

ENTRYPOINT ["./hlm"]
CMD ["-path", "/mounted"]
