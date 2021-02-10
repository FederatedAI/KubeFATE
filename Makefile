BUILD-PATH = ${shell pwd}
K8S-DEPLOY = ./k8s-deploy
CHART = ./helm-charts
FML = ./fml_manager
RELEASE_VERSION ?= ${shell git describe --tags}

define sub_make
 cd $1 && make $2 && cd ${BUILD-PATH}
endef

all: release

release: k8s-release docker-compose-release helm-chart-release
	mkdir -p ${BUILD-PATH}/release
	mv ${K8S-DEPLOY}/release/* release/
	mv ${BUILD-PATH}/kubefate-docker-compose-${RELEASE_VERSION}.tar.gz release/
	mv ${CHART}/fate-*.tgz release/

clean:
	rm -r release

k8s-release:
	${call sub_make, ${K8S-DEPLOY}, release RELEASE_VERSION=${RELEASE_VERSION}}
docker-compose-release:
	tar -czvf kubefate-docker-compose-${RELEASE_VERSION}.tar.gz ./docker-deploy/* ./docker-deploy/.env
helm-chart-release:
	${call sub_make, ${CHART}, release  RELEASE_VERSION=${RELEASE_VERSION}}

fml-release:
	${call sub_make, ${FML}, docker-save VERSION=${RELEASE_VERSION}}
