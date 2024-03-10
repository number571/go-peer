FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.16.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.16.linux-amd64.tar.gz

WORKDIR /go-peer
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o consumer ./cmd/hidden_lake/helpers/adapters/common/consumer

ENV SERVICE_ADDR=""
ENV TRAFFIC_ADDR=""
ENV LOG_ENABLED=""
CMD ./consumer "${SERVICE_ADDR}" "${TRAFFIC_ADDR}" "${LOG_ENABLED}"
