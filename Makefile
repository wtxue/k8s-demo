VERSION ?= v0.0.4
# Image URL to use all building/pushing image targets
IMG_REG ?= symcn.tencentcloudcr.com/symcn

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# This repo's root import path (under GOPATH).
ROOT := github.com/wtxue/k8s-demo

GO_VERSION := 1.14.4
ARCH     ?= $(shell go env GOARCH)
BUILD_DATE = $(shell date +'%Y-%m-%dT%H:%M:%SZ')
COMMIT    = $(shell git rev-parse --short HEAD)
GOENV    := CGO_ENABLED=0 GOOS=$(shell uname -s | tr A-Z a-z) GOARCH=$(ARCH) GOPROXY=https://goproxy.io,direct
#GO       := $(GOENV) go build -mod=vendor
GO       := $(GOENV) go build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: docker-build-provider

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

proto:
	protoc --proto_path=pkg/api/echo --go_out=plugins=grpc:pkg/api/echo echo.proto

# Build the docker image
docker-build-consumer:
	docker run --rm -v "$$PWD":/go/src/${ROOT} -v ${GOPATH}/pkg/mod:/go/pkg/mod -w /go/src/${ROOT} golang:${GO_VERSION} make build-consumer

docker-build-provider:
	docker run --rm -v "$$PWD":/go/src/${ROOT} -v ${GOPATH}/pkg/mod:/go/pkg/mod -w /go/src/${ROOT} golang:${GO_VERSION} make build-provider

docker-build-grpc:
	docker run --rm -v "$$PWD":/go/src/${ROOT} -v ${GOPATH}/pkg/mod:/go/pkg/mod -w /go/src/${ROOT} golang:${GO_VERSION} make build-grpc

build-consumer:
	$(GO) -v -o bin/consumer -ldflags "-s -w -X $(ROOT)/pkg/version.Release=$(VERSION) -X  $(ROOT)/pkg/version.Commit=$(COMMIT)   \
	-X  $(ROOT)/pkg/version.BuildDate=$(BUILD_DATE)" cmd/consumer/main.go

build-grpc:
	$(GO) -v -o bin/grpc-server -ldflags "-s -w -X  $(ROOT)/pkg/version.Release=$(VERSION) -X  $(ROOT)/pkg/version.Commit=$(COMMIT)   \
	-X  $(ROOT)/pkg/version.BuildDate=$(BUILD_DATE)" cmd/grpc-server/main.go

build-provider:
	$(GO) -v -o bin/user-provider -ldflags "-s -w -X  $(ROOT)/pkg/version.Release=$(VERSION) -X  $(ROOT)/pkg/version.Commit=$(COMMIT)   \
	-X  $(ROOT)/pkg/version.BuildDate=$(BUILD_DATE)" cmd/user-provider/main.go
	$(GO) -v -o bin/order-provider -ldflags "-s -w -X  $(ROOT)/pkg/version.Release=$(VERSION) -X  $(ROOT)/pkg/version.Commit=$(COMMIT)   \
	-X  $(ROOT)/pkg/version.BuildDate=$(BUILD_DATE)" cmd/order-provider/main.go

push-provider: docker-build-provider
	docker build -t $(IMG_REG)/user-provider:${VERSION} -f ./misc/user-provider/Dockerfile .
	docker push $(IMG_REG)/user-provider:${VERSION}
	docker build -t $(IMG_REG)/order-provider:${VERSION} -f ./misc/order-provider/Dockerfile .
	docker push $(IMG_REG)/order-provider:${VERSION}

push-consumer: docker-build-consumer
	docker build -t $(IMG_REG)/consumer:${VERSION} -f ./misc/consumer/Dockerfile .
	docker push $(IMG_REG)/consumer:${VERSION}

push-grpc: docker-build-grpc
	docker build -t $(IMG_REG)/grpc-server:${VERSION} -f ./misc/grpc-server/Dockerfile .
	docker push $(IMG_REG)/grpc-server:${VERSION}



