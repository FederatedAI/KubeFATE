NAME = federatedai/kubefate
VERSION = 2.0-alpha


build all: compile zip docker-build

compile:
	go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

zip:
	zip -q -r kubefate.zip *

docker-build:
	docker build -t ${NAME}:${VERSION} .

clean:
	rm kubefate.zip kubefate
