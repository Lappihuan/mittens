#!/bin/sh
set -eu

# Developer tool versions (tracked by Renovate)
# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION=v2.7.2
# renovate: datasource=github-tags depName=gotestyourself/gotestsum
GOTESTSUM_VERSION=v1.13.0
# renovate: datasource=github-tags depName=kubernetes-sigs/kind
KIND_VERSION=v0.31.0
# renovate: datasource=github-tags depName=helm/helm
HELM_VERSION=v4.0.4
# renovate: datasource=github-tags depName=mvdan/sh
GOFUMPT_VERSION=v3.12.0

# canonical GOBIN fallback
GOBIN_DIR=${GOBIN:-"$(go env GOPATH 2>/dev/null || echo $HOME/go)/bin"}
mkdir -p "$GOBIN_DIR"

ensure_kubectl() {
  if ! [ -x "$(command -v kubectl)" ]; then
    echo "kubectl is not installed; please install kubectl or run this script in CI where kubectl is preinstalled."
    return 1
  fi
}

ensure_kubectl || true

# golangci-lint
if ! [ -x "$(command -v golangci-lint)" ]; then
  echo "Installing golangci-lint ${GOLANGCI_LINT_VERSION}"
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
fi

# gotestsum: try go install with a couple of patterns, fall back to warning
if ! [ -x "$(command -v gotestsum)" ]; then
  echo "Installing gotestsum"
  if ! go install gotest.tools/gotestsum@${GOTESTSUM_VERSION} 2>/dev/null; then
    if ! go install github.com/gotestyourself/gotestsum@${GOTESTSUM_VERSION} 2>/dev/null; then
      if ! go install gotest.tools/gotestsum@latest 2>/dev/null; then
        if ! go install github.com/gotestyourself/gotestsum@latest 2>/dev/null; then
          echo "WARNING: gotestsum could not be installed via 'go install'. Please install gotestsum manually or ensure the environment allows its installation."
        fi
      fi
    fi
  fi
fi

# helm
if ! [ -x "$(command -v helm)" ]; then
  echo "Installing helm ${HELM_VERSION}"
  go install helm.sh/helm/v4/cmd/helm@${HELM_VERSION}
fi

# kind
if ! [ -x "$(command -v kind)" ]; then
  echo "Installing kind ${KIND_VERSION}"
  go install sigs.k8s.io/kind@${KIND_VERSION}
fi

# gofumpt (preferred formatter)
if ! [ -x "$(command -v gofumpt)" ]; then
  echo "Installing gofumpt ${GOFUMPT_VERSION}"
  if ! go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} 2>/dev/null; then
    go install mvdan.cc/gofumpt@latest || true
  fi
fi

echo "Developer tools setup complete. Ensure ${GOBIN_DIR} is on your PATH."
