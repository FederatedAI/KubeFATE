NAME = federatedai/kubefate
VERSION = 2.0-alpha

build-linux-binary:
	go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

build-docker-image:
	docker build -t ${NAME}:${VERSION} .

clean:
	rm kubefate
