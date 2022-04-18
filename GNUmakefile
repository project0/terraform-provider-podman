# https://github.com/containers/podman/issues/12548
GOFLAGS ?= -tags=remote,exclude_graphdriver_btrfs,btrfs_noversion,exclude_graphdriver_devicemapper,containers_image_openpgp

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


.PHONY: generate lint
generate:
	go generate ./...

lint:
	golangci-lint run -v