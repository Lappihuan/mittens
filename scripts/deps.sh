#!/bin/sh

GOLANGCI_LINT_VERSION=v1.35.2
GOTESTSUM_VERSION=v0.6.0
KIND_VERSION=v0.9.0
HELM_VERSION=v3.5.0
GOFUMPT_VERSION=v0.3.0
GOFUMPORTS_VERSION=v0.3.0

cd

if ! [ -x "$(command -v kubectl)" ]; then
  echo "kubectl is not installed"
  exit 1
fi


if ! [ -x "$(command -v golangci-lint)" ]; then
  # Use 'go install' with a version tag (works outside modules)
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
fi

if ! [ -x "$(command -v gotestsum)" ]; then
  go install gotest.tools/gotestsum@${GOTESTSUM_VERSION}
fi

if ! [ -x "$(command -v helm)" ]; then
  # official module path for Helm v3
  go install helm.sh/helm/v3/cmd/helm@${HELM_VERSION}
fi


if ! [ -x "$(command -v kind)" ]; then
  go install sigs.k8s.io/kind@${KIND_VERSION}
fi

if ! [ -x "$(command -v gofumpt)" ]; then
  go install mvdan.cc/gofumpt@${GOFUMPT_VERSION}
fi

if ! [ -x "$(command -v gofumports)" ]; then
  go install mvdan.cc/gofumpt/gofumports@${GOFUMPORTS_VERSION}
fi

cd -
