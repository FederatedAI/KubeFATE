module github.com/FederatedAI/KubeFATE/k8s-deploy

go 1.13

require (
	github.com/Masterminds/sprig/v3 v3.1.0
	github.com/appleboy/gin-jwt/v2 v2.6.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/logger v0.0.2
	github.com/gin-gonic/gin v1.6.3
	github.com/gofrs/flock v0.7.1
	github.com/gosuri/uitable v0.0.4
	github.com/jinzhu/gorm v1.9.14
	github.com/json-iterator/go v1.1.10
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.19.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.0
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.0
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)
