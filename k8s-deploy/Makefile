NAME ?= federatedai/kubefate
VERSION ?= v1.4.1
IMG ?= ${NAME}:${VERSION}

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


ifeq ($(OS),Windows_NT)
BUILD_MODE = -buildmode=exe
OUTPUT_FILE = bin/kubefate.exe
else
BUILD_MODE = 
OUTPUT_FILE = bin/kubefate
endif


all: kubefate

# Run tests
test: fmt vet 
	go test ./... -coverprofile cover.out

# Build manager binary
kubefate: fmt vet swag
	CGO_ENABLED=0 go build -a --ldflags '-extldflags "-static"' -o ${OUTPUT_FILE} ${BUILD_MODE} kubefate.go

kubefate-without-swag: fmt vet
	go build -o ${OUTPUT_FILE} ${BUILD_MODE} kubefate.go

run: fmt vet 
	go run ./kubefate.go service

install: 
	kubectl apply -f rbac-config.yaml
	kubectl apply -f kubefate.yaml

uninstall: 
	kubectl delete -f kubefate.yaml
	kubectl delete -f rbac-config.yaml

docker-build: test
	docker build . -t ${IMG}

docker-build-without-test:
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}

docker-save: docker-build-without-test
	docker save -o kubefate-${VERSION}.docker ${IMG} 

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...
# Generate swag API file
swag: swag-bin
	$(SWAG_BIN) init --parseDependency --parseInternal -g pkg/api/service.go

package: kubefate
	mkdir -p package;
	cp -r bin/kubefate cluster.yaml cluster-serving.yaml cluster-spark.yaml config.yaml kubefate.yaml rbac-config.yaml examples package/;
	tar -czvf kubefate-k8s-${RELEASE_VERSION}.tar.gz -C package/ .;
	rm -r package;

release: package docker-save
	mkdir -p release;
	mv kubefate-k8s-${RELEASE_VERSION}.tar.gz kubefate-${VERSION}.docker release/;

clean:
	rm -r release

swag-bin:
ifeq (, $(shell which swag))
	@{ \
	set -e ;\
	SWAG_BIN_TMP_DIR=$$(mktemp -d) ;\
	cd $$SWAG_BIN_TMP_DIR ;\
	go mod init tmp ;\
	go get -u github.com/swaggo/swag/cmd/swag ;\
	rm -rf $$SWAG_BIN_TMP_DIR ;\
	}
SWAG_BIN=$(GOBIN)/swag
else
SWAG_BIN=$(shell which swag)
endif
