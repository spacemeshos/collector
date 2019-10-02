#!/bin/bash -e
./scripts/check-go-version.sh
./scripts/install-protobuf.sh

go mod download
protobuf_path=$(go list -m -f '{{.Dir}}' github.com/golang/protobuf)
echo "installing protoc-gen-go..."
go install $protobuf_path/protoc-gen-go

grpc_gateway_path=$(go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway)
echo "installing protoc-gen-grpc-gateway"
go install $grpc_gateway_path/protoc-gen-grpc-gateway

echo "installing protoc-gen-swagger"
go install $grpc_gateway_path/protoc-gen-swagger

echo "setup complete 🎉"
