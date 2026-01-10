FROM golang:1.25.5-alpine3.23


COPY go.mod /workspace/go.mod
COPY go.sum /workspace/go.sum
COPY entrypoint.go /workspace/entrypoint.go

WORKDIR /workspace

RUN go mod download
RUN go build -o /workspace/entrypoint entrypoint.go

ENTRYPOINT ["/workspace/entrypoint"]
CMD []