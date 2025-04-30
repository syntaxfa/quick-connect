ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

lint:
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8)
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...
