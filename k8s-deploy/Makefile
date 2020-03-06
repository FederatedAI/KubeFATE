NAME = federatedai/kubefate
VERSION = v1.0.1

build-linux-binary:
	go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

build-docker-image:
	docker build -t ${NAME}:${VERSION} .

clean:
	rm kubefate
