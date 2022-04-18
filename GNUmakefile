# https://github.com/containers/podman/issues/12548
GOFLAGS ?= -tags=remote,exclude_graphdriver_btrfs,btrfs_noversion,exclude_graphdriver_devicemapper,containers_image_openpgp

default: testacc

# Run acceptance tests
.PHONY: testacc test
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

test:
	go test -v ./...

.PHONY: build generate lint
build:
	go build

generate:
	go generate ./...

lint:
	golangci-lint run -v
