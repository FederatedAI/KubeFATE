#!/usr/bin/env bash
kubefateVersion="v1.0.2"
gitCommit=$(git rev-parse "HEAD" 2>/dev/null);
gitVersion=$(git describe --tags --match='v*' --abbrev=14 "${gitCommit}" 2>/dev/null)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ' 2>/dev/null)
ldflags+="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.gitVersion=${gitVersion}'"

ldflags+="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.gitCommit=${gitCommit}'"

ldflags+="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.buildDate=${buildDate}'"

ldflags+="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.kubefateVersion=${kubefateVersion}'"

echo ldflags:${ldflags}
go build -o kubefate.exe -i -v -gcflags='-N -l' -ldflags="${ldflags}"