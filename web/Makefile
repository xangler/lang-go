GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

TARGETS := go-web-server go-proxy-server go-leader-server
FULL_BUILD := $(shell git rev-parse HEAD)
BUILD := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(BRANCH)-$(BUILD)
IMAGE_NAME := go-web-demo

REGISTRY_ADDRESS ?= registry.demo.com
REPO_USERNAME ?= demo
REPO_PASSWORD ?= demo
REPO_PROJECT ?= demo
IMAGE_FULLNAME := $(REGISTRY_ADDRESS)/$(REPO_PROJECT)/$(IMAGE_NAME):$(VERSION)
BUILD_IMAGE := protoc:v1.0.0

project=github.com/learn-go/web
LDFLAGS += -X "$(project)/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(project)/version.GitHash=$(FULL_BUILD)"
LDFLAGS += -X "$(project)/version.Version=$(VERSION)"
LDFLAGS += -X "$(project)/version.GitBranch=$(BRANCH)"
CodeRoot=/go/src/$(project)

all: check build

build: $(TARGETS)

$(TARGETS): $(SRC)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o $@.out $(project)/cmd/$@ 

ibuild:
	$(call docker_env, make build)

image-fullname:
	@echo $(IMAGE_FULLNAME)

image-version:
	@echo $(IMAGE_NAME):$(VERSION)

image:
	docker build -t $(IMAGE_FULLNAME) .

image-repo:
	docker login -u '${REPO_USERNAME}' -p '${REPO_PASSWORD}' ${REGISTRY_ADDRESS}

image-push:
	docker push $(IMAGE_FULLNAME)

lint:
	@golangci-lint run -v --deadline=5m

packages = $(shell go list ./...|grep -v /vendor/|grep -v /mock/)

mock:
	@rm -rf mock
	@go generate -x ./...

test: check
	$(GO) test ${packages}

test-no-dao: check
	$(GO) test $(shell go list ./...|grep -v /vendor/|grep -v /dao)

cov: check
	gocov test $(packages) | gocov-html > coverage.html

check:
	@echo skip go vet

clean:
	rm -f $(TARGETS)

.PHONY: clean all build check image image-fullname image-version lint mock test cov check clean

define docker_env
	docker run --rm \
		-v ${PWD}:${CodeRoot} \
		-w ${CodeRoot} \
		${BUILD_IMAGE} bash -c "$(1)"
endef
