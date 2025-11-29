ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

# Default valid values
GOLANGCI_LINT_VERSION ?= v2.5.0
IMAGE_NAME ?= quick-connect
PROTO_DIR ?= protobuf
OUT_DIR ?= .

# Protobuf tools detection.
PROTOC_GEN_GO ?= $(shell which protoc-gen-go)
# FIX: Pointing to the actual grpc plugin, not protoc itself.
PROTOC_GEN_GO_GRPC ?= $(shell which protoc-gen-go-grpc)
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

# Define all non-file targets as PHONY to avoid conflicts with files of the same name.
# Added: all, test, clean
.PHONY: all test clean lint generate-proto update-proto-tools \
	chat-swag-init manager-swag-init notification-swag-init example-micro1-swag-init \
	test-general chat-test manager-test notification-test admin-test all-in-one-test \
	chat-build manager-build notification-build admin-build all-in-one-build \
	generate-example-proto

# 1. Standard 'all' target: Runs linting and all tests by default
all: lint test

# 2. Standard 'test' target: Aggregates all your sub-tests
test: test-general chat-test manager-test notification-test admin-test

# 3. Standard 'clean' target: Cleans Go cache (or remove binaries if you had any)
clean:
	go clean
	@echo "Cleaned up."

lint:
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION))
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

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
			--go_opt=paths=import \
			--go-grpc_opt=paths=import \
			--proto_path=$(PROTO_DIR) \
			$$file; \
	done

update-proto-tools:
	@echo "Updating Go protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Go protoc plugins updated successfully."

# Swagger Generators
chat-swag-init:
	swag init -g cmd/chat/main.go -o app/chatapp/docs/ --tags=Websocket,Chat --instanceName chat

manager-swag-init:
	swag init -g cmd/manager/main.go -o app/managerapp/docs/ --tags=Manager,User,Token,Internal,Guest --instanceName manager

notification-swag-init:
	swag init -g cmd/notification/main.go -o app/notificationapp/docs/ --tags=Notification,NotificationClient,NotificationAdmin --instanceName notification

example-micro1-swag-init:
	swag init -g example/observability/microservice1/main.go -o example/observability/internal/microservice1/docs --tags=Micro1

# Specific Tests
test-general:
	go test ./pkg/...
	go test ./adapter/...
	go test ./config/...
	go test ./outbox/...

chat-test:
	go test ./app/chatapp/...

manager-test:
	go test ./app/managerapp/...

notification-test:
	go test ./app/notificationapp/...

admin-test:
	go test ./app/adminapp/...

all-in-one-test:
	go test ./app/managerapp/...
	go test ./app/chatapp/...
	go test ./app/notificationapp/...
	go test ./app/adminapp/...

# Builds (Ensure IMAGE_NAME is set or passed as argument).

build_tag = $(if $(filter quick-connect,$(IMAGE_NAME)),$(IMAGE_NAME):$(1),$(IMAGE_NAME))

chat-build:
	docker build -t $(call build_tag,chat) -f deploy/chat/deploy/Dockerfile .

manager-build:
	docker build -t $(call build_tag,manager) -f deploy/manager/deploy/Dockerfile .

notification-build:
	docker build -t $(call build_tag,notification) -f deploy/notification/deploy/Dockerfile .

admin-build:
	docker build -t $(call build_tag,admin) -f deploy/admin/deploy/Dockerfile .

all-in-one-build:
	docker build -t $(call build_tag,aio) -f deploy/all-in-one/deploy/Dockerfile .

generate-example-proto:
	@protoc \
		--proto_path=protobuf "protobuf/example/proto/example.proto" \
		--go_out=. \
		--go-grpc_out=.