NAME ?= federatedai/kubefate
VERSION ?= v1.0.2

# Image URL to use all building/pushing image targets
IMG ?= federatedai/kubefate:latest

all: kubefate

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Build manager binary
kubefate: fmt vet
	go build -o bin/kubefate kubefate.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./kubefate.go service

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

build-linux-binary:
	go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

build-docker-image:
	docker build -t ${NAME}:${VERSION} .

release: build-linux-binary
	tar -czvf kubefate-k8s-${RELEASE_VERSION}.tar.gz kubefate cluster.yaml cluster-serving.yaml config.yaml kubefate.yaml rbac-config.yaml

clean:
	rm kubefate
