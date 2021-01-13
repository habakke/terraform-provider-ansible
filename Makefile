ROOT_DIR := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR = $(ROOT_DIR)/build
NAME := terraform-provider-ansible

.PHONY: build-dev

prepare:
	mkdir -p $(BUILD_DIR)

build-dev:
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -o ~/.terraform.d/plugins/$(NAME)_${version} .

build: prepare
	@[ "${version}" ] || ( echo ">> please provide version=vX.Y.Z"; exit 1 )
	go build -o $(BUILD_DIR)/$(NAME)_${version} .

install: build
	mv $(NAME) ~/.terraform.d/plugins/github.com/habakke/$(NAME)/${version}/darwin_amd64

test: prepare
	go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: prepare
	TF_ACC=true go test -v -coverprofile=$(BUILD_DIR)/cover.out ./...

clean:
	rm -rf $(BUILD_DIR)
