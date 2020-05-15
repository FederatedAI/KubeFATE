#!/usr/bin/env sh
kubefateVersion="v1.0.3"

gitCommit=$(git rev-parse "HEAD" 2>/dev/null);
gitVersion=$(git describe --tags --match='v*' --abbrev=14 "${gitCommit}" 2>/dev/null)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ' 2>/dev/null)
git_version="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.gitVersion=${gitVersion}'"
git_commit="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.gitCommit=${gitCommit}'"
build_data="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.buildDate=${buildDate}'"
kubefate_version="-X 'github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service.kubefateVersion=${kubefateVersion}'"

echo "${git_version} ${git_commit} ${build_data} ${kubefate_version}"