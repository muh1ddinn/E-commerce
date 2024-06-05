#!/bin/bash

# Get the current directory from the first argument
CURRENT_DIR=$1

# Remove old generated files
rm -rf "${CURRENT_DIR}/genproto/*"

# Set paths to plugins
PROTOC_GEN_GO=$(which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(which protoc-gen-go-grpc)

# Check if plugins are found
if [[ -z "$PROTOC_GEN_GO" || -z "$PROTOC_GEN_GO_GRPC" ]]; then
  echo "protoc-gen-go or protoc-gen-go-grpc not found in PATH"
  exit 1
fi

# Generate code for each subdirectory in the protos folder
for x in $(find ${CURRENT_DIR}/protos/* -type d); do
  protoc --plugin="protoc-gen-go=${PROTOC_GEN_GO}" \
         --plugin="protoc-gen-go-grpc=${PROTOC_GEN_GO_GRPC}" \
         -I=${x} -I=${CURRENT_DIR}/protos -I /usr/local/include \
         --go_out=${CURRENT_DIR} \
         --go-grpc_out=${CURRENT_DIR} \
         ${x}/*.proto
done
