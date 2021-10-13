FROM golang:1.17 as builder
WORKDIR /root/api
COPY api api
COPY cli cli
COPY providers providers
COPY go.mod go.sum /root/api/
RUN GOOS=linux go build -o build/api cli/api/main.go

FROM debian:11
WORKDIR /root/
COPY --from=builder /root/api/build/api .
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
CMD ["/root/api"]
