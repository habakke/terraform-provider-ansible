NAME       := terraform-provider-ansible
ROOT_DIR   := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR  := $(ROOT_DIR)/dist
VERSION    := $(shell cat VERSION)
GITSHA     := $(shell git rev-parse --short HEAD)


BUILD_TIME         := $(shell date +'%Y-%m-%d_%T')
GO_OS              := $(if $(GOHOSTOS),$(GOHOSTOS),$(shell go env GOHOSTOS))
GO_ARCH            := $(if $(GOHOSTARCH),$(GOHOSTARCH),$(shell go env GOHOSTARCH))
OS_ARCH            := $(GO_OS)_$(GO_ARCH)
GIT_BRANCH         :=$(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION       :=$(shell git rev-list -1 HEAD)
GIT_REVISION_DIRTY :=$(shell (git diff-index --quiet HEAD -- . && git diff --staged --quiet -- .) || echo "-dirty")

.PHONY: prepare lint check sec build-dev build install test testacc fmt release-test release clean

all: testacc build

prepare:
	mkdir -p $(BUILD_DIR)

lint:
	go get -u golang.org/x/lint/golint
	$(shell go list -f {{.Target}} golang.org/x/lint/golint) ./...

check:
	go get -u honnef.co/go/tools/cmd/staticcheck
	$(shell go list -f {{.Target}} honnef.co/go/tools/cmd/staticcheck) ./...

sec:
	go get -u github.com/securego/gosec/v2/cmd/gosec
	$(shell go list -f {{.Target}} github.com/securego/gosec/v2/cmd/gosec) -fmt=golint ./...

build-dev:
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X main.buildTime=$(BUILD_TIME) -X main.version=$(VERSION) -X main.buildBy=${USER}" -o ~/.terraform.d/plugins/$(NAME)_$(VERSION) .

build: prepare
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X main.buildTime=$(BUILD_TIME) -X main.version=$(VERSION) -X main.buildBy=${USER}" -o $(BUILD_DIR)/$(NAME)_$(VERSION) .

install: build
	mkdir -p ~/.terraform.d/plugins/github.com/habakke/ansible/$(VERSION)/$(OS_ARCH)
	mv $(BUILD_DIR)/$(NAME)_$(VERSION) ~/.terraform.d/plugins/github.com/habakke/ansible/$(VERSION)/$(OS_ARCH)/$(NAME)_$(VERSION)

test: prepare
	go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: export TF_ACC=true
testacc: prepare
	go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

fmt:
	go fmt ./...

release-test: export GITHUB_SHA=$(GITSHA)
release-test: export GPG_FINGERPRINT=0000
release-test: testacc
	goreleaser release --skip-publish --snapshot --rm-dist --skip-sign

release: export GITHUB_SHA=$(GITSHA)
release: release-test
	git tag -a $(VERSION) -m "Release" && git push origin $(VERSION)

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ~/.terraform.d/plugins/github.com/habakke/ansible
