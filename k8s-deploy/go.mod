module github.com/FederatedAI/KubeFATE/k8s-deploy

go 1.13

require (
	github.com/DataDog/zstd v1.4.4 // indirect
	github.com/Masterminds/sprig/v3 v3.0.2
	github.com/appleboy/gin-jwt/v2 v2.6.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/logger v0.0.2
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/gofrs/flock v0.7.1
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gosuri/uitable v0.0.1
	github.com/json-iterator/go v1.1.9
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.17.2
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.3.2
	github.com/urfave/cli/v2 v2.1.1
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.2.0
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	helm.sh/helm/v3 v3.0.2
	k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	sigs.k8s.io/yaml v1.1.0
)

replace (
	// github.com/Azure/go-autorest/autorest has different versions for the Go
	// modules than it does for releases on the repository. Note the correct
	// version when updating.
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.9.0
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

	// Kubernetes imports github.com/miekg/dns at a newer version but it is used
	// by a package Helm does not need. Go modules resolves all packages rather
	// than just those in use (like Glide and dep do). This sets the version
	// to the one oras needs. If oras is updated the version should be updated
	// as well.
	github.com/miekg/dns => github.com/miekg/dns v0.0.0-20181005163659-0d29b283ac0f
	gopkg.in/inf.v0 v0.9.1 => github.com/go-inf/inf v0.9.1
	gopkg.in/square/go-jose.v2 v2.3.0 => github.com/square/go-jose v2.3.0+incompatible

	rsc.io/letsencrypt => github.com/dmcgowan/letsencrypt v0.0.0-20160928181947-1847a81d2087
)
