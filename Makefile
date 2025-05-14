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

chat-swag-init:
	swag init -g cmd/chat/main.go -o app/chatapp/docs/ --tags=Websocket,Chat

manager-swag-init:
	swag init -g cmd/manager/main.go -o app/managerapp/docs/ --tags=Manager,User,Token

test-general:
	go test ./pkg/...
	go test ./adapter/...
	#go test ./cli/...
	go test ./config/...
	go test ./outbox/...

chat-test:
	go test ./app/chatapp/...

manager-test:
	go test ./app/managerapp/...

chat-build:
	docker build -t $(IMAGE_NAME) -f deploy/chat/deploy/Dockerfile .

manager-build:
	docker build -t $(IMAGE_NAME) -f deploy/manager/deploy/Dockerfile .
