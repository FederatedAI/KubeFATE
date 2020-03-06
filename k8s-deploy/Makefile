NAME = federatedai/kubefate
VERSION = 2.0-alpha

build-linux-binary:
	go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

build-docker-image:
	docker build -t ${NAME}:${VERSION} .

release: build-linux-binary
	tar -czvf kubefate-k8s-${RELEASE_VERSION}.tar.gz kubefate cluster.yaml config.yaml kubefate.yaml rbac-config.yaml

clean:
	rm kubefate
