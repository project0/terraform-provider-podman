# https://github.com/containers/podman/issues/12548
GOFLAGS ?= -tags=remote,exclude_graphdriver_btrfs,btrfs_noversion,exclude_graphdriver_devicemapper,containers_image_openpgp

TEST_DOCKER_COMPOSE ?= tcp://localhost:10888
TF_ACC_TEST_PROVIDER_PODMAN_URI ?= $(TEST_DOCKER_COMPOSE)

export PODMAN_VERSION ?= latest

default: testacc

# Run acceptance tests (optional in rootless docker)
.PHONY: testacc test
testacc:
ifeq ($(TEST_DOCKER_COMPOSE),$(TF_ACC_TEST_PROVIDER_PODMAN_URI))
	docker-compose up -d
endif
	TF_ACC_TEST_PROVIDER_PODMAN_URI=$(TF_ACC_TEST_PROVIDER_PODMAN_URI) TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
ifeq ($(TEST_DOCKER_COMPOSE),$(TF_ACC_TEST_PROVIDER_PODMAN_URI))
	docker-compose down
endif

test:
	go test -v ./...

.PHONY: build generate lint
build:
	go build

generate:
	go generate ./...

lint:
	golangci-lint run -v
