ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
BUILD_TIME := $(shell date +'%Y-%m-%d_%T')
BUILD_COMMIT := $(shell git rev-parse HEAD)
NAME := terraform-provider-ansible
GO_OS := $(if $(GOHOSTOS),$(GOHOSTOS),$(shell go env GOHOSTOS))
GO_ARCH := $(if $(GOHOSTARCH),$(GOHOSTARCH),$(shell go env GOHOSTARCH))
OS_ARCH = $(GO_OS)_$(GO_ARCH)

.PHONY: build-dev

all: testacc build

prepare:
	mkdir -p $(BUILD_DIR)

build-dev:
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -ldflags "-X util.commit=$(BUILD_COMMIT) -X util.buildTime=$(BUILD_TIME) -X util.version=${version} -X util.buildBy=${USER}" -o ~/.terraform.d/plugins/$(NAME)_${version} .

build: prepare
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -ldflags "-X util.commit=$(BUILD_COMMIT) -X util.buildTime=$(BUILD_TIME) -X util.version=${version} -X util.buildBy=${USER}" -o $(BUILD_DIR)/$(NAME)_${version} .

install: build
	mkdir -p ~/.terraform.d/plugins/github.com/habakke/ansible/${version}/$(OS_ARCH)
	mv $(BUILD_DIR)/$(NAME)_${version} ~/.terraform.d/plugins/github.com/habakke/ansible/${version}/$(OS_ARCH)/$(NAME)_${version}

test: prepare
	go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: prepare
	TF_ACC=true go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

release: testacc
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	git tag ${version}
	git push --tags

clean:
	rm -rf $(BUILD_DIR)
