FROM golang:1.21 AS builder

RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/go/bin:${PATH}"

WORKDIR /app

ENV GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN protoc \
    --proto_path=./proto/ \
    --go_out=./pb --go_opt=paths=source_relative \
    --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
    user.proto

RUN go build -o grpc-tube ./cmd/main.go

RUN chmod +x grpc-tube

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/grpc-tube .

EXPOSE 8080

CMD ["./grpc-tube"]
