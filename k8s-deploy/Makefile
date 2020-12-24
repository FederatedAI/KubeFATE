NAME ?= federatedai/kubefate
VERSION ?= v1.2.0
IMG ?= ${NAME}:${VERSION}

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: kubefate

# Run tests
test: fmt vet 
	go test ./... -coverprofile cover.out

# Build manager binary
kubefate: fmt vet
	go build -o bin/kubefate kubefate.go

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

package: kubefate
	mkdir -p package;
	cp -r ./{bin/kubefate,cluster.yaml,cluster-serving.yaml,cluster-spark.yaml,config.yaml,kubefate.yaml,rbac-config.yaml,examples} package/;
	tar -czvf kubefate-k8s-${RELEASE_VERSION}.tar.gz -C package/ .;
	rm -r package;

release: package docker-save
	mkdir -p release;
	mv ./{kubefate-k8s-${RELEASE_VERSION}.tar.gz,kubefate-${VERSION}.docker} release/;

clean:
	rm -r release

