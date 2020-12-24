NAME ?= federatedai/client
VERSION ?= v1.5.0
IMG ?= ${NAME}:${VERSION:v%=%-release}

docker-build:
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}

docker-save: docker-build
	docker save -o client-${VERSION}.docker ${IMG} 