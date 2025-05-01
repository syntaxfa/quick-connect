ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

lint:
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8)
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...


PROTO_DIR ?= proto
OUT_DIR ?= gen

PROTOC_GEN_GO ?= $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC ?= $(shell which protoc)

PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

generate-proto:
	@if [ ! -x "$(PROTOC_GEN_GO)" ] || [ ! -x "$(PROTOC_GEN_GO_GRPC)" ]; then \
		echo "Error: protoc-gen-go and protoc-gen-go-grpc must be installed and in your PATH."; \
		exit 1; \
	fi

	@mkdir -p $(OUT_DIR)
	@for file in $(PROTO_FILES); do \
		protoc \
			--go_out=$(OUT_DIR) \
			--go-grpc_out=$(OUT_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			--proto_path=$(PROTO_DIR) \
			$$file; \
	done
