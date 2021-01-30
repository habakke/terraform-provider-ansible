NAME := terraform-provider-ansible
ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
BUILD_TIME := $(shell date +'%Y-%m-%d_%T')
GO_OS := $(if $(GOHOSTOS),$(GOHOSTOS),$(shell go env GOHOSTOS))
GO_ARCH := $(if $(GOHOSTARCH),$(GOHOSTARCH),$(shell go env GOHOSTARCH))
OS_ARCH = $(GO_OS)_$(GO_ARCH)

GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION=$(shell git rev-list -1 HEAD)
GIT_REVISION_DIRTY=$(shell (git diff-index --quiet HEAD -- . && git diff --staged --quiet -- .) || echo "-dirty")

.PHONY: build-dev

all: testacc build

prepare:
	mkdir -p $(BUILD_DIR)

build-dev:
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X main.buildTime=$(BUILD_TIME) -X main.version=${version} -X main.buildBy=${USER}" -o ~/.terraform.d/plugins/$(NAME)_${version} .

build: prepare
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X main.buildTime=$(BUILD_TIME) -X main.version=${version} -X main.buildBy=${USER}" -o $(BUILD_DIR)/$(NAME)_${version} .

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
