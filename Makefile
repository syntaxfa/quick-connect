ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

GOLANGCI_LINT_VERSION ?= v2.5.0

lint:
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION))
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

PROTO_DIR ?= protobuf
OUT_DIR ?= .

PROTOC_GEN_GO ?= $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC ?= $(shell which protoc)

PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

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

chat-swag-init:
	swag init -g cmd/chat/main.go -o app/chatapp/docs/ --tags=Websocket,Chat

manager-swag-init:
	swag init -g cmd/manager/main.go -o app/managerapp/docs/ --tags=Manager,User,Token

notification-swag-init:
	swag init -g cmd/notification/main.go -o app/notificationapp/docs/ --tags=Notification,NotificationClient,NotificationAdmin

example-micro1-swag-init:
	swag init -g example/observability/microservice1/main.go -o example/observability/internal/microservice1/docs --tags=Micro1

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

notification-test:
	go test ./app/notificationapp/...

chat-build:
	docker build -t $(IMAGE_NAME) -f deploy/chat/deploy/Dockerfile .

manager-build:
	docker build -t $(IMAGE_NAME) -f deploy/manager/deploy/Dockerfile .

notification-build:
	docker build -t $(IMAGE_NAME) -f deploy/notification/deploy/Dockerfile .

generate-example-proto:
	@protoc \
		--proto_path=protobuf "protobuf/example/proto/example.proto" \
		--go_out=. \
  		--go-grpc_out=.
