#!/bin/bash

PROTO_PATH="./protocol/"
OUTPUT_DIR="./protos"
mkdir -p "$OUTPUT_DIR"

# 遍历所有 .proto 文件并编译
for proto_file in $(find "$PROTO_PATH" -name "*.proto"); do
    echo "Compiling: $proto_file"
    protoc  -I $PROTO_PATH \
            -I ~/code/protocol-buffer/googleapis \
            --go_out=./protos \
            --go_opt=paths=source_relative \
            --go-grpc_out=./protos \
            --go-grpc_opt=paths=source_relative \
            "$proto_file"
done

mv "$OUTPUT_DIR/core/contract"/* "$OUTPUT_DIR/core/"
rmdir "$OUTPUT_DIR/core/contract"