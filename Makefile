BUILD-PATH = ${shell pwd}
K8S-DEPLOY="./k8s-deploy"

define sub_make
 cd $1 && make $2 && cd ${BUILD-PATH}
endef

all: build-linux-binary zip build-docker-image

build-linux-binary:
	$(call sub_make, ${K8S-DEPLOY}, build-linux-binary)

build-docker-image:
	$(call sub_make, ${K8S-DEPLOY}, build-docker-image)

zip:
	tar -czvf kubefate-docker-compose.tar.gz ./docker-deploy

release: zip
	${call sub_make, ${K8S-DEPLOY}, release RELEASE_VERSION=v${RELEASE_VERSION}} && mv ${K8S-DEPLOY}/kubefate-k8s-v${RELEASE_VERSION}.tar.gz ./ && curl -LJO https://federatedai.github.io/KubeFATE/fate-${RELEASE_VERSION}.tgz
clean:
	rm kubefate-docker-compose.tar.gz kubefate-k8s-*.tar.gz fate-*.tgz && $(call sub_make, ${K8S-DEPLOY}, clean)